package usecase

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"fmt"
	"log"

	"ecommerce/common/pubsub"

	"github.com/viettranx/service-context/core"
)

type changeAvatarUC struct {
	userQueryRepo UserQueryRepository
	userCmdRepo   UserCommandRepository
	imageRepo     ImageRepository
}

func NewChangeAvatarUC(userQueryRepo UserQueryRepository, userCmdRepo UserCommandRepository,
	imgRepo ImageRepository) *changeAvatarUC {
	return &changeAvatarUC{
		userQueryRepo: userQueryRepo,
		userCmdRepo:   userCmdRepo,
		imageRepo:     imgRepo,
	}
}

func (uc *changeAvatarUC) ChangeAvatar(ctx context.Context, dto SingleAvatarChangeDTO) error {
	// 1. Get user by requester
	userID := dto.Requester.UserId()
	user, err := uc.userQueryRepo.Find(ctx, userID.String())
	if err != nil {
		return core.ErrBadRequest.WithWrap(err).WithDebug(err.Error()).WithError("can not find user")
	}
	// 2. Find image by image id
	imageID, err := common.ParseUUID(dto.ImageID)
	if err != nil {
		return core.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	image, err := uc.imageRepo.Find(ctx, imageID)
	if err != nil {
		return core.ErrBadRequest.WithWrap(err).WithDebug(err.Error()).WithError("can not find image")
	}
	// 3. Update user's avatar
	// Create file url: domain/file name
	fileUrl := fmt.Sprintf("%s/%s", dto.CDNDomain, image.FileName)
	// Create new avatar
	avatar, err := domain.NewAvatar(image.ID, image.FileName, fileUrl)

	if err != nil {
		return core.ErrInternalServerError.WithWrap(err).WithDebug(err.Error()).WithError("can not update user's avatar")
	}
	if err := uc.userCmdRepo.UpdateAvatar(ctx, user, avatar); err != nil {
		return core.ErrInternalServerError.WithWrap(err).WithDebug(err.Error()).WithError("can not update user's avatar")
	}
	
	// 4. Update image's status using pubsub pattern
	// Publish event
	// Get pub sub component from context
	ps := ctx.Value("pubsub").(pubsub.PubSub)
	
	// Create message
	
	if err := ps.Publish(ctx, common.TopicUserChangedAvatar, pubsub.NewMessage(map[string]any{
		"img_id" : imageID.String(),
		"user_id" : dto.Requester.UserId().String(), 
	})); err != nil {
		log.Println(err)
	}
	return nil
}

type ImageRepository interface {
	Find(ctx context.Context, id common.UUID) (*common.Image, error)
	SetImageStatusActivated(ctx context.Context, id common.UUID) error
}

