package usecase

import(
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"errors"
	"log"
)

type registerUC struct{
	userQueryRepo UserQueryRepository
	userCmdRepo UserCommandRepository
	hasher Hasher
}

func NewRegisterUC(userQueryRepo UserQueryRepository, userCmdRepo UserCommandRepository, hasher Hasher) *registerUC {
	return &registerUC{
		userQueryRepo: userQueryRepo,
		userCmdRepo: userCmdRepo,
		hasher: hasher,
	}
}

func (uc *registerUC) Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error {
	// Find user with provided email
	user, err := uc.userQueryRepo.FindByEmail(ctx, dto.Email)
	if user != nil {
		return domain.ErrEmailExists
	}
	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		return err
	}

	// Generate salt
	salt, err := uc.hasher.RandomStr(30)

	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	// Create hashed password
	hashedPassword, err := uc.hasher.HashPassword(salt, dto.Password)

	if err != nil {
		log.Printf("error: %v", err)
		return err
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
		log.Printf("error: %v", err)
		return err
	}

	// Create new User
	if err := uc.userCmdRepo.Create(ctx, newUser); err != nil {
		return err
	}
	return nil
}