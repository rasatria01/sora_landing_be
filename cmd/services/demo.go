package services

import (
	"context"
	"database/sql"
	"sora_landing_be/cmd/dto"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/dto/response"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/pkg/database"

	"github.com/uptrace/bun"
)

type DemoService interface {
	CreateDemo(ctx context.Context, payload requests.Demo) error
	GetDemoByID(ctx context.Context, id string) (response.Demo, error)
	ListDemos(ctx context.Context, payload requests.ListDemo) (dto.PaginationResponse[response.Demo], error)
	UpdateDemo(ctx context.Context, id string, payload requests.Demo) error
	DeleteDemo(ctx context.Context, id string) error
}

type demoService struct {
	demoRepo repository.DemoRepository
}

func NewDemoService(demoRepo repository.DemoRepository) DemoService {
	return &demoService{
		demoRepo: demoRepo,
	}
}

func (s *demoService) CreateDemo(ctx context.Context, payload requests.Demo) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		data := payload.ToDomain()

		err := s.demoRepo.CreateDemo(ctx, &data)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *demoService) GetDemoByID(ctx context.Context, id string) (response.Demo, error) {
	var res response.Demo
	data, err := s.demoRepo.GetDemoByID(ctx, id)
	if err != nil {
		return res, err
	}

	res = response.NewDemo(data)
	return res, nil
}

func (s *demoService) ListDemos(ctx context.Context, payload requests.ListDemo) (dto.PaginationResponse[response.Demo], error) {
	var paginateRes dto.PaginationResponse[response.Demo]
	res, count, err := s.demoRepo.ListDemos(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewListDemo(res))
	return paginateRes, nil
}

func (s *demoService) UpdateDemo(ctx context.Context, id string, payload requests.Demo) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		data := payload.ToDomain()
		data.ID = id

		err := s.demoRepo.UpdateDemo(ctx, &data)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *demoService) DeleteDemo(ctx context.Context, id string) error {
	err := s.demoRepo.DeleteDemo(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
