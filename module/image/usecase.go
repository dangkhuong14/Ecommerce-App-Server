package image

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type imageUseCase struct {
	imageSaver ImageSaver
	repo       ImageCmdRepo
}

func NewImageUseCase(imageSaver ImageSaver, repo ImageCmdRepo) imageUseCase {
	return imageUseCase{
		imageSaver: imageSaver,
		repo:       repo,
	}
}

func (uc imageUseCase) UploadImage(ctx context.Context, dto UploadImageDTO) (*Image, error) {
	// 1. Create destination file name: folder name (if specified)/ + Unix time in nano second + file name
	dstFileName := fmt.Sprintf("%d_%s", time.Now().UTC().UnixNano(), dto.FileName)

	// 2. Upload image to S3
	if err := uc.imageSaver.SaveFileUploaded(ctx, dto.FileData, dstFileName); err != nil {
		return nil, err
	}

	// 3. Create new record of uploaded image in "images" table
	newImage := NewImage(WithTitle(dto.Name),
		WithFileName(dstFileName),
		WithFileSize(dto.FileSize),
		WithFileType(dto.FileType),
		WithStorageProvider(uc.imageSaver.GetName()),
		WithStatus(StatusUploaded))

	if err := uc.repo.Create(ctx, newImage); err != nil {
		return nil, err
	}

	return newImage, nil
}

type ImageSaver interface {
	SaveFileUploaded(ctx context.Context, data []byte, dst string) error
	GetName() string
}

type ImageCmdRepo interface {
	Create(ctx context.Context, entity *Image) error
}

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
	}
}

var (
	ErrCannotUploadImage = errors.New("can not upload image")
	ErrCannotFindImage   = errors.New("can not find image")
	ErrCannotUploadFile  = errors.New("can not upload file")
	ErrCannotOpenFile    = errors.New("can not open file")
	ErrCannotReadFile    = errors.New("can not read file")
)
