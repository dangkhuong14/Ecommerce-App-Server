package image

import (
	"ecommerce/common"
	"net/http"

	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
)

type ImageService struct {
	serviceContext sctx.ServiceContext
}

func NewImageService(serviceContext sctx.ServiceContext) *ImageService {
	return &ImageService{serviceContext: serviceContext}
}


// Use to upload image to provider (AWS S3, local...)
func (s *ImageService) handleUploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get file from request multipart/form-data
		f, err := c.FormFile("file")
		if err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(ErrCannotUploadFile.Error()).WithDebug(err.Error()))
			return
		}

		// Open file
		file, err := f.Open()
		if err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(ErrCannotOpenFile.Error()).WithDebug(err.Error()))
			return
		}
		defer file.Close()

		// Create byte slice to store the file
		fileData := make([]byte, f.Size)

		// Read file data to slice
		if _, err := file.Read(fileData); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(ErrCannotReadFile.Error()).WithDebug(err.Error()))
			return
		}

		// Create UploadImageDTO to response and save image
		dto := UploadImageDTO{
			Name: f.Filename,
			FileName: f.Filename,
			FileType: http.DetectContentType(fileData),
			FileSize: int(f.Size),
			FileData: fileData,
		}
		

		// Upload file to AWS S3 bucket
		// Create image use case
		s3Uploader, ok := s.serviceContext.Get(common.KeyAwsS3Component)
		if !ok {
			common.WriteErrorResponse(c,core.ErrInternalServerError.WithDebug("can not get AWS S3 component from service context"))
			return
		}
		imageSaver, ok := s3Uploader.(ImageSaver)
		if !ok {
			common.WriteErrorResponse(c,core.ErrInternalServerError.WithDebug("can not assert s3Uploader to ImageSaver interface"))
			return
		}

		// Image repository
		dbComponent := s.serviceContext.MustGet(common.KeyGormComponent).(common.GormCompContext)
		imageRepo := NewImageRepo(dbComponent.GetDB())

		imgUseCase := NewImageUseCase(imageSaver, imageRepo)

		image, err := imgUseCase.UploadImage(c, dto)
		if err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(ErrCannotUploadFile.Error()).WithDebug(err.Error()))
			return
		}
		imgResponseDTO := toImageResponseDTO(image)
		// Set file url to image response
		imgResponseDTO.setFileUrl(imageSaver.GetDomain())
		c.JSON(http.StatusOK, core.ResponseData(imgResponseDTO))
	}
}

func (s *ImageService) Routes(group *gin.RouterGroup) {
	group.POST("/upload-image", s.handleUploadImage())
}

