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
	GetDomain() string
}

type ImageCmdRepo interface {
	Create(ctx context.Context, entity *Image) error
}

var (
	ErrCannotUploadImage = errors.New("can not upload image")
	ErrCannotFindImage   = errors.New("can not find image")
	ErrCannotUploadFile  = errors.New("can not upload file")
	ErrCannotOpenFile    = errors.New("can not open file")
	ErrCannotReadFile    = errors.New("can not read file")
)
