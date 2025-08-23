package repository

import (
	"sora_landing_be/pkg/database"
	"sync"
)

var once = &sync.Once{}
var RepoPool *PoolRepository

type PoolRepository struct {
	UserRepository           UserRepository
	AuthenticationRepository AuthRepository
}

func Init(db *database.Database) {
	once.Do(func() {
		RepoPool = &PoolRepository{
			UserRepository:           NewUserRepository(db),
			AuthenticationRepository: NewAuthRepository(db),
		}
	})
}
