package usecase

import (
	"context"
	"ecommerce/common"

	"github.com/viettranx/service-context/core"
)

type revokeTokenUC struct {
	// Find token by user's id
	sessionCmdRepo SessionCommandRepository

}

func NewRevokeTokenUC(sessionCmdRepo SessionCommandRepository) *revokeTokenUC{
	return &revokeTokenUC{
		sessionCmdRepo: sessionCmdRepo,
	}
}

func (u *revokeTokenUC) RevokeToken(ctx context.Context, sessionID common.UUID) error {
	// Delete token
	err := u.sessionCmdRepo.Delete(ctx, sessionID)
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error()).WithError("Can not log out right now, something went wrong")
	}
	return nil
}