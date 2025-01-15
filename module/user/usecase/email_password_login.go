package usecase

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"time"

	"github.com/viettranx/service-context/core"
)

type loginEmailPasswordUC struct {
	userRepo      UserQueryRepository
	sessionRepo   SessionCommandRepository
	tokenProvider common.TokenProvider
	hasher        Hasher
}

func NewLoginEmailPasswordUC(
	userRepo UserQueryRepository, sessionRepo SessionCommandRepository,
	hasher Hasher, tokenProvider common.TokenProvider,
) *loginEmailPasswordUC {
	return &loginEmailPasswordUC{
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
	}
}

func (uc *loginEmailPasswordUC) LoginEmailPassword(ctx context.Context, dto EmailPasswordLoginDTO) (*TokenResponseDTO, error) {
	// 1. Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return nil, core.ErrBadRequest.WithWrap(err).WithDebug(err.Error()).WithError(domain.ErrInvalidEmailPassword.Error())
	}
	// 2. Compare user's salt and password
	if ok := uc.hasher.CompareHashPassword(user.GetPassword(), user.GetSalt(), dto.Password); !ok {
		return nil, core.ErrBadRequest.WithDebug("User's password is incorrect").WithError(domain.ErrInvalidEmailPassword.Error())
	}
	// 3. Generate JWT payload: session id, user id
	// Pre generate session id
	newSessionID := common.GenNewUUID()

	token, err := uc.tokenProvider.IssueToken(ctx, newSessionID.String(), user.GetID().String())

	if err != nil {
		return nil, core.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	// Generate random string for refresh token
	refreshToken, _ := uc.hasher.RandomStr(16)

	// 4. Create new session
	accessExpAt := time.Now().UTC().Add(time.Second * time.Duration((uc.tokenProvider.TokenExpireInSeconds())))
	refreshExpAt := time.Now().UTC().Add(time.Second * time.Duration((uc.tokenProvider.RefreshExpireInSeconds())))

	newSession := domain.NewSession(newSessionID, user.GetID(), refreshToken, accessExpAt, refreshExpAt)

	if err := uc.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, core.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	tokenResponseDTO := TokenResponseDTO{
		AccessToken:       token,
		AccessTokenExpIn:  uc.tokenProvider.TokenExpireInSeconds(),
		RefreshToken:      refreshToken,
		RefreshTokenExpIn: uc.tokenProvider.RefreshExpireInSeconds(),
	}
	return &tokenResponseDTO, nil
}
