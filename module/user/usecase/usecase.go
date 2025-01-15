package usecase

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"
)

type UseCase interface {
	Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error
	LoginEmailPassword(ctx context.Context, dto EmailPasswordLoginDTO) (*TokenResponseDTO, error)
	RevokeToken(ctx context.Context, sessionID common.UUID) error
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponseDTO, error)
}

type useCase struct {
	//Embed
	*loginEmailPasswordUC
	*registerUC
	*revokeTokenUC
	*refreshTokendUC
}

type Builder interface {
	BuildUserQueryRepo() UserQueryRepository
	BuildUserCmdRepo() UserCommandRepository
	BuildHasher() Hasher
	BuildTokenProvider() common.TokenProvider
	BuildSessionRepo() SessionRepository
	BuildSessionQueryRepo() SessionQueryRepository
	BuildSessionCmdRepo() SessionCommandRepository
}

func NewUseCaseWithBuilder(b Builder) *useCase {
	return &useCase{
		loginEmailPasswordUC: NewLoginEmailPasswordUC(b.BuildUserQueryRepo(), b.BuildSessionCmdRepo(), b.BuildHasher(), b.BuildTokenProvider()),
		registerUC:           NewRegisterUC(b.BuildUserQueryRepo(), b.BuildUserCmdRepo(), b.BuildHasher()),
		revokeTokenUC:        NewRevokeTokenUC(b.BuildSessionCmdRepo()),
		refreshTokendUC:      NewRefreshTokenUC(b.BuildUserQueryRepo(), b.BuildSessionRepo(), b.BuildHasher(), b.BuildTokenProvider()),
	}
}

func NewUseCase(repo UserRepository, sessionRepo SessionRepository, hasher Hasher, tokenProvider common.TokenProvider) *useCase {
	return &useCase{
		loginEmailPasswordUC: NewLoginEmailPasswordUC(repo, sessionRepo, hasher, tokenProvider),
		registerUC:           NewRegisterUC(repo, repo, hasher),
	}
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
	Find(ctx context.Context, userID string) (*domain.User, error)
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
	Delete(ctx context.Context, sessionID common.UUID) error
}

type SessionQueryRepository interface {
	Find(ctx context.Context, sessionID string) (*domain.Session, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
}
