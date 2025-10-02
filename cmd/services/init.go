package services

import (
	"sora_landing_be/cmd/repository"
	"sync"
)

var once = &sync.Once{}
var ServicePool *PoolService

type PoolService struct {
	AuthService     AuthService
	UserService     UserService
	TagService      TagService
	CategoryService CategoryService
	BlogService     BlogService
	DemoService     DemoService
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
			TagService:      NewTagService(repo.TagRepository),
			CategoryService: NewCatService(repo.CategoryRepository),
			BlogService:     NewBlogService(repo.BlogRepository, repo.TagRepository, repo.CategoryRepository),
			DemoService:     NewDemoService(repo.DemoRepository),
		}
	})
}
