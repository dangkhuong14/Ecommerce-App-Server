package builder

import (
	"ecommerce/common"
	"ecommerce/module/user/domain"
	"ecommerce/module/user/infras/repository"
	"ecommerce/module/user/usecase"

	"context"

	"gorm.io/gorm"
)

type simpleBuilder struct {
	db *gorm.DB
	tp usecase.TokenProvider
}

func NewSimpleBuilder(db *gorm.DB, tp usecase.TokenProvider) simpleBuilder {
	return simpleBuilder{
		db: db,
		tp: tp,
	}
}

func (s simpleBuilder) BuildUserQueryRepo() usecase.UserQueryRepository {
	return repository.NewMysqlUser(s.db)
}

func (s simpleBuilder) BuildUserCmdRepo() usecase.UserCommandRepository {
	return repository.NewMysqlUser(s.db)
}

func (s simpleBuilder) BuildHasher() usecase.Hasher {
	return &common.Hasher{}
}

func (s simpleBuilder) BuildTokenProvider() usecase.TokenProvider {
	return s.tp
}

func (s simpleBuilder) BuildSessionQueryRepo() usecase.SessionQueryRepository {
	return repository.NewMysqlSession(s.db)
}

func (s simpleBuilder) BuildSessionCmdRepo() usecase.SessionCommandRepository {
	return repository.NewMysqlSession(s.db)
}

// Use this complex builder to overwrite BuildUserQueryRepo method of simpleBuilder
type cmplxBuilder struct {
	simpleBuilder
}

func NewCmplxBuilder(s simpleBuilder) cmplxBuilder {
	return cmplxBuilder{
		simpleBuilder: s,
	}
}

func (c cmplxBuilder) BuildUserQueryRepo() usecase.UserQueryRepository {
	return NewUserCacheRepo(repository.NewMysqlUser(c.db))
}

type userCacheRepo struct {
	usecase.UserQueryRepository
	cache    map[string]*domain.User
}

func NewUserCacheRepo(realRepo usecase.UserQueryRepository) *userCacheRepo {
	return &userCacheRepo{
		UserQueryRepository: realRepo,
		cache:    make(map[string]*domain.User),
	}
}

func (u *userCacheRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Check if email already in the cache
	if user, ok := u.cache[email]; ok {
		return user, nil
	}

	// If not in the cache use the mysql repo
	user, err := u.UserQueryRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	} else {
		// Save user to the cache
		u.cache[email] = user
		return user, nil
	}
}
