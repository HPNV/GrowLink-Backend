package file

const (
	CreateQuery = `
		INSERT INTO files (uuid, original_name, file_name, file_path, file_size, mime_type, uploaded_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	GetByUUIDQuery = `
		SELECT uuid, original_name, file_name, file_path, file_size, mime_type, uploaded_by, created_at
		FROM files WHERE uuid = $1
	`

	DeleteQuery = `DELETE FROM files WHERE uuid = $1`

	GetByUploadedByQuery = `
		SELECT uuid, original_name, file_name, file_path, file_size, mime_type, uploaded_by, created_at
		FROM files WHERE uploaded_by = $1 ORDER BY created_at DESC
	`
)
