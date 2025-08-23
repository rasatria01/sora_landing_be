package services

import (
	"sora_landing_be/cmd/repository"
	"sync"
)

var once = &sync.Once{}
var ServicePool *PoolService

type PoolService struct {
	AuthService AuthService
	UserService UserService
}

func Init() {
	once.Do(func() {
		repo := repository.RepoPool
		ServicePool = &PoolService{
			AuthService: NewAuthSrv(repo.AuthenticationRepository),
			UserService: NewUserSrv(
				repo.UserRepository,
				repo.AuthenticationRepository,
			),
		}
	})
}
