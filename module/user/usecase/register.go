package usecase

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"errors"

	"github.com/viettranx/service-context/core"
)

type registerUC struct {
	userQueryRepo UserQueryRepository
	userCmdRepo   UserCommandRepository
	hasher        Hasher
}

func NewRegisterUC(userQueryRepo UserQueryRepository, userCmdRepo UserCommandRepository, hasher Hasher) *registerUC {
	return &registerUC{
		userQueryRepo: userQueryRepo,
		userCmdRepo:   userCmdRepo,
		hasher:        hasher,
	}
}

func (uc *registerUC) Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error {
	// Find user with provided email
	user, err := uc.userQueryRepo.FindByEmail(ctx, dto.Email)
	if user != nil {
		return core.ErrBadRequest.WithError(domain.ErrEmailExists.Error())
	}
	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		return core.ErrInternalServerError.WithError("Registration is not available right now").WithWrap(err).WithDebug(err.Error())
	}

	// Generate salt
	salt, err := uc.hasher.RandomStr(30)

	if err != nil {
		return core.ErrInternalServerError.WithError("Registration is not available right now").WithWrap(err).WithDebug(err.Error())
	}
	// Create hashed password
	hashedPassword, err := uc.hasher.HashPassword(salt, dto.Password)

	if err != nil {
		return core.ErrInternalServerError.WithError("Registration is not available right now").WithWrap(err).WithDebug(err.Error())
	}

	// Transform dto entity to User entity
	newUser, err := domain.NewUser(
		common.GenNewUUID(),
		dto.FirstName,
		dto.LastName,
		dto.Email,
		hashedPassword,
		salt,
		"active",
		domain.RoleUser,
	)

	if err != nil {
		return core.ErrInternalServerError.WithError("Registration is not available right now").WithWrap(err).WithDebug(err.Error())
	}

	// Create new User
	if err := uc.userCmdRepo.Create(ctx, newUser); err != nil {
		return core.ErrInternalServerError.WithError("Registration is not available right now").WithWrap(err).WithDebug(err.Error())
	}
	return nil
}
