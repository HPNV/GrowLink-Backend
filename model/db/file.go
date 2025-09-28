package db

type File struct {
	UUID         string `db:"uuid"`
	OriginalName string `db:"original_name"`
	FileName     string `db:"file_name"`
	FilePath     string `db:"file_path"`
	FileSize     int64  `db:"file_size"`
	MimeType     string `db:"mime_type"`
	UploadedBy   string `db:"uploaded_by"`
	CreatedAt    string `db:"created_at"`
}
