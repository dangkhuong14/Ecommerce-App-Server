package usecase

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"time"
)

type refreshTokendUC struct {
	userRepo      UserQueryRepository
	sessionRepo   SessionRepository
	tokenProvider common.TokenProvider
	hasher        Hasher
}

func NewRefreshTokenUC(
	userRepo UserQueryRepository, sessionRepo SessionRepository,
	hasher Hasher, tokenProvider common.TokenProvider,
) *refreshTokendUC {
	return &refreshTokendUC{
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
	}
}

func (uc *refreshTokendUC) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponseDTO, error) {
	// 1. Find session by refresh token and check if it is expired
	session, err := uc.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// if refresh token's expire time is before now: this means this token is expired
	if session.GetRefreshExpAt().Before(time.Now().UTC()) {
		return nil, domain.ErrRefreshTokenExpired
	}

	// 2. Find user by id and check user's status
	user, err := uc.userRepo.Find(ctx, session.GetUserID().String())
	if err != nil {
		return nil, err
	}

	if user.GetStatus() == "banned" {
		return nil, domain.ErrUserBanned
	}

	// 3. Generate new JWT payload: session id, user id
	// Pre generate session id
	newSessionID := common.GenNewUUID()

	token, err := uc.tokenProvider.IssueToken(ctx, newSessionID.String(), user.GetID().String())

	if err != nil {
		return nil, err
	}

	// Generate random string for refresh token
	newRefreshToken, _ := uc.hasher.RandomStr(16)

	// 4. Create new session
	accessExpAt := time.Now().UTC().Add(time.Second * time.Duration((uc.tokenProvider.TokenExpireInSeconds())))
	refreshExpAt := time.Now().UTC().Add(time.Second * time.Duration((uc.tokenProvider.RefreshExpireInSeconds())))

	newSession := domain.NewSession(newSessionID, user.GetID(), newRefreshToken, accessExpAt, refreshExpAt)

	if err := uc.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, err
	}

	tokenResponseDTO := TokenResponseDTO{
		AccessToken:       token,
		AccessTokenExpIn:  uc.tokenProvider.TokenExpireInSeconds(),
		RefreshToken:      newRefreshToken,
		RefreshTokenExpIn: uc.tokenProvider.RefreshExpireInSeconds(),
	}

	// 5. Delete old session
	// run async since we don't need result
	go func() {
		uc.sessionRepo.Delete(ctx, session.GetID())
	}()

	return &tokenResponseDTO, nil
}
