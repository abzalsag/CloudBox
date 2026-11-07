package trash

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"CloudBox/internal/entity"
)

type Usecase struct {
	repo    Repository
	folders FolderReader
	users   UserRepository
	storage Storage
	tx      TxManager
	log     *zap.Logger
}

func NewUsecase(repo Repository, folders FolderReader, users UserRepository, storage Storage, tx TxManager, log *zap.Logger) *Usecase {
	return &Usecase{repo: repo, folders: folders, users: users, storage: storage, tx: tx, log: log}
}

type Listing struct {
	Folders []entity.Folder
	Files   []entity.File
}

func (u *Usecase) List(ctx context.Context, ownerID uuid.UUID) (*Listing, error) {
	folders, err := u.repo.ListTrashedFolders(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	files, err := u.repo.ListTrashedFiles(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return &Listing{Folders: folders, Files: files}, nil
}

func (u *Usecase) RestoreFile(ctx context.Context, ownerID, id uuid.UUID) (*entity.File, error) {
	f, err := u.repo.GetFileAny(ctx, ownerID, id)
	if err != nil {
		return nil, err
	}
	if f.DeletedAt == nil {
		return nil, entity.ErrNotFound
	}
	if err := u.ensureParentActive(ctx, ownerID, f.FolderID); err != nil {
		return nil, err
	}
	return u.repo.RestoreFile(ctx, ownerID, id)
}

func (u *Usecase) RestoreFolder(ctx context.Context, ownerID, id uuid.UUID) (*entity.Folder, error) {
	f, err := u.repo.GetFolderAny(ctx, ownerID, id)
	if err != nil {
		return nil, err
	}
	if f.DeletedAt == nil {
		return nil, entity.ErrNotFound
	}
	if err := u.ensureParentActive(ctx, ownerID, f.ParentID); err != nil {
		return nil, err
	}

	ids, err := u.repo.FolderTreeIDs(ctx, ownerID, id)
	if err != nil {
		return nil, err
	}
	if err := u.repo.RestoreFolderTree(ctx, ownerID, ids); err != nil {
		return nil, err
	}
	return u.folders.GetByID(ctx, ownerID, id)
}

func (u *Usecase) PurgeFile(ctx context.Context, ownerID, id uuid.UUID) error {
	f, err := u.repo.GetFileAny(ctx, ownerID, id)
	if err != nil {
		return err
	}
	if f.DeletedAt == nil {
		return entity.ErrNotFound
	}
	return u.purgeFileRef(ctx, entity.FileRef{ID: f.ID, OwnerID: f.OwnerID, Size: f.Size, StorageKey: f.StorageKey})
}

func (u *Usecase) PurgeFolder(ctx context.Context, ownerID, id uuid.UUID) error {
	f, err := u.repo.GetFolderAny(ctx, ownerID, id)
	if err != nil {
		return err
	}
	if f.DeletedAt == nil {
		return entity.ErrNotFound
	}

	ids, err := u.repo.FolderTreeIDs(ctx, ownerID, id)
	if err != nil {
		return err
	}
	files, err := u.repo.FilesInFolders(ctx, ownerID, ids)
	if err != nil {
		return err
	}

	err = u.tx.WithTx(ctx, func(ctx context.Context) error {
		var total int64
		for _, ref := range files {
			if err := u.repo.HardDeleteFile(ctx, ownerID, ref.ID); err != nil {
				return err
			}
			total += ref.Size
		}
		if err := u.repo.HardDeleteFolders(ctx, ownerID, ids); err != nil {
			return err
		}
		return u.users.DecreaseStorageUsed(ctx, ownerID, total)
	})
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(files))
	for _, ref := range files {
		keys = append(keys, ref.StorageKey)
	}
	u.removeObjects(ctx, keys)
	return nil
}

func (u *Usecase) PurgeExpired(ctx context.Context, batch int) (purgedFiles, purgedFolders int, err error) {
	files, err := u.repo.ListExpiredFiles(ctx, batch)
	if err != nil {
		return 0, 0, err
	}
	for _, ref := range files {
		if err := u.purgeFileRef(ctx, ref); err != nil {
			u.log.Error("trash: purge expired file", zap.String("file_id", ref.ID.String()), zap.Error(err))
			continue
		}
		purgedFiles++
	}

	folders, err := u.repo.ListExpiredFolders(ctx, batch)
	if err != nil {
		return purgedFiles, 0, err
	}
	for _, ref := range folders {
		err := u.PurgeFolder(ctx, ref.OwnerID, ref.ID)
		if errors.Is(err, entity.ErrNotFound) {
			continue
		}
		if err != nil {
			u.log.Error("trash: purge expired folder", zap.String("folder_id", ref.ID.String()), zap.Error(err))
			continue
		}
		purgedFolders++
	}
	return purgedFiles, purgedFolders, nil
}

func (u *Usecase) purgeFileRef(ctx context.Context, ref entity.FileRef) error {
	err := u.tx.WithTx(ctx, func(ctx context.Context) error {
		if err := u.repo.HardDeleteFile(ctx, ref.OwnerID, ref.ID); err != nil {
			return err
		}
		return u.users.DecreaseStorageUsed(ctx, ref.OwnerID, ref.Size)
	})
	if err != nil {
		return err
	}
	u.removeObjects(ctx, []string{ref.StorageKey})
	return nil
}

func (u *Usecase) removeObjects(ctx context.Context, keys []string) {
	for _, key := range keys {
		if err := u.storage.Delete(context.WithoutCancel(ctx), key); err != nil {
			u.log.Warn("trash: delete object failed", zap.String("key", key), zap.Error(err))
		}
	}
}

func (u *Usecase) ensureParentActive(ctx context.Context, ownerID uuid.UUID, parentID *uuid.UUID) error {
	if parentID == nil {
		return nil
	}
	_, err := u.folders.GetByID(ctx, ownerID, *parentID)
	if errors.Is(err, entity.ErrNotFound) {
		return fmt.Errorf("%w: restore the parent folder first", entity.ErrInvalidInput)
	}
	return err
}
