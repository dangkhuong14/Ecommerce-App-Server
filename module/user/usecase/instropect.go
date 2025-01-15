package usecase

import (
	"context"
	"ecommerce/common"
	"ecommerce/module/user/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/viettranx/service-context/core"
)

type TokenParser interface {
	ParseToken(ctx context.Context, tokenString string) (claims *jwt.RegisteredClaims, err error)
}

type IntrospectTokenUC struct {
	userQueryRepo UserQueryRepository
	sessionQueryRepo SessionQueryRepository
	tokenParser  TokenParser
}

func NewIntrospectTokenUC(userQueryRepo UserQueryRepository, sessionQueryRepo SessionQueryRepository, tokenParser TokenParser) *IntrospectTokenUC {
	return &IntrospectTokenUC{
		userQueryRepo: userQueryRepo,
		sessionQueryRepo: sessionQueryRepo,
		tokenParser: tokenParser,
	}
}

func (uc *IntrospectTokenUC) IntrospectToken(ctx context.Context, token string) (common.Requester, error){
	// 1. Parse token, get user id, token id
	claims, err := uc.tokenParser.ParseToken(ctx, token)
	if err != nil {
		return nil, core.ErrBadRequest.WithWrap(err).WithDebug(err.Error()).WithError("Can not parse this token")
	}

	// 2. Find session by session id
	if _, err := uc.sessionQueryRepo.Find(ctx, claims.ID); err != nil {
		return nil, core.ErrBadRequest.WithWrap(err).WithDebug(err.Error()).WithError("Can not find session")
	}
	// 3. Find user by user id
	user, err := uc.userQueryRepo.Find(ctx, claims.Subject)

	if err != nil {
		return nil, core.ErrBadRequest.WithWrap(err).WithDebug(err.Error()).WithError("Can not find user")
	}

	if user.GetStatus() == "banned"{
		return nil, core.ErrForbidden.WithError(domain.ErrUserBanned.Error())
	}

	return common.NewRequester(
		common.UUID(uuid.MustParse(claims.Subject)),
		common.UUID(uuid.MustParse(claims.ID)),
		user.GetFirstName(),
		user.GetLastName(),
		user.GetRole().String(),
		user.GetStatus(),
	), nil
}


