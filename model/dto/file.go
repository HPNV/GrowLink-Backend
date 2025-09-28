package dto

type FileUploadResponse struct {
	UUID         string `json:"uuid"`
	OriginalName string `json:"original_name"`
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`
	URL          string `json:"url"`
	CreatedAt    string `json:"created_at"`
}
