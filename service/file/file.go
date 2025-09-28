package file

import (
	"fmt"
	"mime/multipart"

	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/repository"
	"github.com/jmoiron/sqlx"
)

type IFile interface {
	UploadImage(file multipart.File, header *multipart.FileHeader, uploadedBy string) (*dto.FileUploadResponse, error)
	GetByUUID(uuid string) (*dto.FileUploadResponse, error)
	Delete(uuid string) error
	GetByUploadedBy(uploadedBy string) ([]*dto.FileUploadResponse, error)
}

type File struct {
	repo repository.IRegistry
}

func NewFile(repo repository.IRegistry) IFile {
	return &File{
		repo: repo,
	}
}

func (f *File) UploadImage(file multipart.File, header *multipart.FileHeader, uploadedBy string) (*dto.FileUploadResponse, error) {
	var result *dto.FileUploadResponse

	err := f.repo.WithTransaction(func(tx *sqlx.Tx) error {
		fileRecord, err := f.repo.GetFile().UploadImage(tx, file, header, uploadedBy)
		if err != nil {
			return err
		}

		result = &dto.FileUploadResponse{
			UUID:         fileRecord.UUID,
			OriginalName: fileRecord.OriginalName,
			FileName:     fileRecord.FileName,
			FilePath:     fileRecord.FilePath,
			FileSize:     fileRecord.FileSize,
			MimeType:     fileRecord.MimeType,
			URL:          fmt.Sprintf("/static/images/%s", fileRecord.FileName),
			CreatedAt:    fileRecord.CreatedAt,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f *File) GetByUUID(uuid string) (*dto.FileUploadResponse, error) {
	fileRecord, err := f.repo.GetFile().GetByUUID(uuid)
	if err != nil {
		return nil, err
	}

	return &dto.FileUploadResponse{
		UUID:         fileRecord.UUID,
		OriginalName: fileRecord.OriginalName,
		FileName:     fileRecord.FileName,
		FilePath:     fileRecord.FilePath,
		FileSize:     fileRecord.FileSize,
		MimeType:     fileRecord.MimeType,
		URL:          fmt.Sprintf("/static/images/%s", fileRecord.FileName),
		CreatedAt:    fileRecord.CreatedAt,
	}, nil
}

func (f *File) Delete(uuid string) error {
	return f.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return f.repo.GetFile().Delete(tx, uuid)
	})
}

func (f *File) GetByUploadedBy(uploadedBy string) ([]*dto.FileUploadResponse, error) {
	files, err := f.repo.GetFile().GetByUploadedBy(uploadedBy)
	if err != nil {
		return nil, err
	}

	var responses []*dto.FileUploadResponse
	for _, file := range files {
		responses = append(responses, &dto.FileUploadResponse{
			UUID:         file.UUID,
			OriginalName: file.OriginalName,
			FileName:     file.FileName,
			FilePath:     file.FilePath,
			FileSize:     file.FileSize,
			MimeType:     file.MimeType,
			URL:          fmt.Sprintf("/static/images/%s", file.FileName),
			CreatedAt:    file.CreatedAt,
		})
	}

	return responses, nil
}
