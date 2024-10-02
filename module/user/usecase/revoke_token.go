package usecase

import (
	"context"
	"ecommerce/common"
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
		return err
	}
	return nil
}