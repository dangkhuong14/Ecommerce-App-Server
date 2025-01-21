package common

import (
	"time"
)

const (
	TbName          = "images"
	ProviderS3      = "aws_s3"
	StatusUploaded  = "uploaded"
	StatusActivated = "activated"
	StatusDeleted   = "deleted"
)

type Image struct {
	ID              UUID      `json:"id" gorm:"id"`
	Title           string    `json:"title" gorm:"title"`
	FileName        string    `json:"file_name" gorm:"file_name"`
	FileSize        int       `json:"file_size" gorm:"file_size"`
	FileType        string    `json:"file_type" gorm:"file_type"`
	StorageProvider string    `json:"storage_provider" gorm:"storage_provider"`
	Status          string    `json:"status" gorm:"status"`
	CreatedAt       time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"updated_at"`
}

type Option func(*Image)

func NewImage(opts ...Option) *Image {
	now := time.Now().UTC()

	newImage := &Image{
		// Create ID, created - updated time fields
		ID:        GenNewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, opt := range opts {
		opt(newImage)
	}
	return newImage
}

func WithTitle(title string) Option {
	return func(img *Image) {
		img.Title = title
	}
}
func WithFileName(fileName string) Option {
	return func(img *Image) {
		img.FileName = fileName
	}
}
func WithFileSize(fileSize int) Option {
	return func(img *Image) {
		img.FileSize = fileSize
	}
}
func WithFileType(fileType string) Option {
	return func(img *Image) {
		img.FileType = fileType
	}
}
func WithStorageProvider(storageProvider string) Option {
	return func(img *Image) {
		img.StorageProvider = storageProvider
	}
}
func WithStatus(status string) Option {
	return func(img *Image) {
		img.Status = status
	}
}
