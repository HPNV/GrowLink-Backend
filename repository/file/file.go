package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/HPNV/growlink-backend/model/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IFile interface {
	UploadImage(tx *sqlx.Tx, file multipart.File, header *multipart.FileHeader, uploadedBy string) (*db.File, error)
	GetByUUID(uuid string) (*db.File, error)
	Delete(tx *sqlx.Tx, uuid string) error
	GetByUploadedBy(uploadedBy string) ([]*db.File, error)
}

type File struct {
	db *sqlx.DB
}

func NewFile(db *sqlx.DB) IFile {
	return &File{
		db: db,
	}
}

func (f *File) UploadImage(tx *sqlx.Tx, file multipart.File, header *multipart.FileHeader, uploadedBy string) (*db.File, error) {
	defer file.Close()

	// Validate file type
	if !f.isValidImageType(header.Filename) {
		return nil, fmt.Errorf("invalid file type. Only jpg, jpeg, png, gif are allowed")
	}

	id := uuid.New()
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".png"
	}

	filename := fmt.Sprintf("%s%s", id.String(), ext)
	uploadDir := "./static/images"
	filePath := filepath.Join(uploadDir, filename)

	// Create upload directory
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Create file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy file data
	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %v", err)
	}

	// Create file record
	fileRecord := &db.File{
		UUID:         id.String(),
		OriginalName: header.Filename,
		FileName:     filename,
		FilePath:     filePath,
		FileSize:     size,
		MimeType:     header.Header.Get("Content-Type"),
		UploadedBy:   uploadedBy,
	}

	// Save to database
	err = tx.QueryRow(CreateQuery,
		fileRecord.UUID,
		fileRecord.OriginalName,
		fileRecord.FileName,
		fileRecord.FilePath,
		fileRecord.FileSize,
		fileRecord.MimeType,
		fileRecord.UploadedBy,
	).Scan(&fileRecord.CreatedAt)

	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file record: %v", err)
	}

	return fileRecord, nil
}

func (f *File) GetByUUID(uuid string) (*db.File, error) {
	file := &db.File{}
	err := f.db.Get(file, GetByUUIDQuery, uuid)
	return file, err
}

func (f *File) Delete(tx *sqlx.Tx, uuid string) error {
	// Get file info first
	file, err := f.GetByUUID(uuid)
	if err != nil {
		return err
	}

	// Delete from database
	_, err = tx.Exec(DeleteQuery, uuid)
	if err != nil {
		return err
	}

	// Delete physical file
	err = os.Remove(file.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete physical file: %v", err)
	}

	return nil
}

func (f *File) GetByUploadedBy(uploadedBy string) ([]*db.File, error) {
	var files []*db.File
	err := f.db.Select(&files, GetByUploadedByQuery, uploadedBy)
	return files, err
}

func (f *File) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}
