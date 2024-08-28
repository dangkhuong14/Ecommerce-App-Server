package usecase

import (
	"context"
	"ecommerce/module/user/domain"
)

type UseCase interface {
	Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error
	LoginEmailPassword(ctx context.Context, dto EmailPasswordLoginDTO) (*TokenResponseDTO, error)
}

type useCase struct {
	//Embed
	*loginEmailPasswordUC
	*registerUC
}

func NewUseCase(repo UserRepository, sessionRepo SessionRepository, hasher Hasher, tokenProvider TokenProvider) *useCase {
	return &useCase{
		loginEmailPasswordUC: NewLoginEmailPasswordUC(repo, sessionRepo, hasher, tokenProvider),
		registerUC: NewRegisterUC(repo, repo, hasher),
	}
}

type TokenProvider interface {
	IssueToken(ctx context.Context, id, sub string) (token string, err error)
	TokenExpireInSeconds() int
	RefreshExpireInSeconds() int
}

type Hasher interface {
	RandomStr(length int) (string, error)
	HashPassword(salt, password string) (string, error)
	CompareHashPassword(hashedPassword, salt, password string) bool
}


type UserRepository interface {
	UserQueryRepository
	UserCommandRepository
}

type UserQueryRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserCommandRepository interface {
	Create(ctx context.Context, data *domain.User) error
}

type SessionRepository interface {
	SessionQueryRepository
	SessionCommandRepository
}

type SessionCommandRepository interface {
	Create(ctx context.Context, data *domain.Session) error
}

type SessionQueryRepository interface{}
