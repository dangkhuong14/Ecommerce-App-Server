package usecase

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"errors"
	"log"
)

type UseCase interface {
	Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error
	LoginEmailPassword(ctx context.Context, dto EmailPasswordLoginDTO) (*TokenResponseDTO, error)

}

type TokenProvider interface {
	IssueToken(ctx context.Context, id, sub string) (token string, err error)
	TokenExpireInSeconds() int
	RefreshExpireInSeconds() int
}

type useCase struct {
	repo          UserRepository
	sessionRepo   SessionRepository
	hasher        Hasher
	tokenProvider TokenProvider
}

func NewUseCase(repo UserRepository, sessionRepo SessionRepository, hasher Hasher, tokenProvider TokenProvider) *useCase {
	return &useCase{
		repo:          repo,
		sessionRepo:   sessionRepo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
	}
}

type Hasher interface {
	RandomStr(length int) (string, error)
	HashPassword(salt, password string) (string, error)
	CompareHashPassword(hashedPassword, salt, password string) bool
}

func (uc *useCase) Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error {
	// Find user with provided email
	user, err := uc.repo.FindByEmail(ctx, dto.Email)
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
		domain.RoleUser,
	)

	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	// Create new User
	if err := uc.repo.Create(ctx, newUser); err != nil {
		return err
	}
	return nil
}

type UserRepository interface {
	// Find(ctx context.Context, id common.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, data *domain.User) error
	// Update(ctx context.Context, data *domain.User) error
	// Delete(ctx context.Context, data *domain.User) error
}

type SessionRepository interface {
	Create(ctx context.Context, data *domain.Session) error
}
