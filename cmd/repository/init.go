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
	TagRepository            TagRepository
	CategoryRepository       CategoryRepository
	BlogRepository           BlogRepository
	DemoRepository           DemoRepository
	FileRepository           FileRepository
}

func Init(db *database.Database) {
	once.Do(func() {
		RepoPool = &PoolRepository{
			UserRepository:           NewUserRepository(db),
			AuthenticationRepository: NewAuthRepository(db),
			TagRepository:            NewTagRepository(db),
			CategoryRepository:       NewCatRepository(db),
			BlogRepository:           NewBlogRepository(db),
			DemoRepository:           NewDemoRepository(db),
			FileRepository:           NewFileRepository(db),
		}
	})
}
