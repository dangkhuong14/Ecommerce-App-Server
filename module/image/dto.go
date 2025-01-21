package image

import (
	"fmt"
	"time"
)

type UploadImageDTO struct {
	Name     string
	FileName string
	FileType string
	FileSize int
	FileData []byte
}

// Replace ID from UUID to string
type ImageResponseDTO struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	FileName        string    `json:"file_name"`
	FileURL         string    `json:"file_url"`
	FileSize        int       `json:"file_size"`
	FileType        string    `json:"file_type"`
	StorageProvider string    `json:"storage_provider"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func toImageResponseDTO(image *Image) *ImageResponseDTO {
	return &ImageResponseDTO{
		ID:              image.ID.String(), // UUID -> string
		Title:           image.Title,
		FileName:        image.FileName,
		FileSize:        image.FileSize,
		FileType:        image.FileType,
		StorageProvider: image.StorageProvider,
		Status:          image.Status,
		CreatedAt:       image.CreatedAt,
		UpdatedAt:       image.UpdatedAt,
		FileURL:         "", // Need to set by call setFileUrl method
	}
}

func (imgResponseDTO *ImageResponseDTO) setFileUrl(domain string) {
	fileUrl := fmt.Sprintf("%s/%s", domain, imgResponseDTO.FileName)
	imgResponseDTO.FileURL =fileUrl
}
