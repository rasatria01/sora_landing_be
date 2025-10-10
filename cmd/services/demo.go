package services

import (
	"context"
	"database/sql"
	"fmt"
	"sora_landing_be/cmd/dto"
	"sora_landing_be/cmd/dto/requests"
	"sora_landing_be/cmd/dto/response"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/pkg/database"

	"github.com/uptrace/bun"
	"github.com/xuri/excelize/v2"
)

type DemoService interface {
	CreateDemo(ctx context.Context, payload requests.Demo) error
	GetDemoByID(ctx context.Context, id string) (response.Demo, error)
	ListDemos(ctx context.Context, payload requests.ListDemo) (dto.PaginationResponse[response.Demo], error)
	UpdateDemo(ctx context.Context, id string, payload requests.Demo) error
	DeleteDemo(ctx context.Context, id string) error
	ExportDemo(ctx context.Context, payload requests.ExportDemo) (*excelize.File, error)
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

func (s *demoService) ExportDemo(ctx context.Context, payload requests.ExportDemo) (*excelize.File, error) {
	entries, err := s.demoRepo.ExportDemo(ctx, payload)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	mainSheet := "Client"
	f.SetSheetName("Sheet1", mainSheet)

	header := []string{"ID", "Nama", "Brand", "No HP", "Email", "Waktu", "Tanggal Janjian", "Tanggal Input"}

	for i, h := range header {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(mainSheet, cell, h)
	}

	for i, entry := range entries {
		row := i + 2
		f.SetCellValue(mainSheet, fmt.Sprintf("A%d", row), entry.ID)
		f.SetCellValue(mainSheet, fmt.Sprintf("B%d", row), entry.Nama)
		f.SetCellValue(mainSheet, fmt.Sprintf("C%d", row), entry.Brand)
		f.SetCellValue(mainSheet, fmt.Sprintf("D%d", row), entry.NoHP)
		f.SetCellValue(mainSheet, fmt.Sprintf("E%d", row), entry.Email)
		f.SetCellValue(mainSheet, fmt.Sprintf("F%d", row), entry.Waktu)
		f.SetCellValue(mainSheet, fmt.Sprintf("G%d", row), entry.Tanggal.Format("15-03-2006"))
		f.SetCellValue(mainSheet, fmt.Sprintf("H%d", row), entry.BaseEntity.CreatedAt.Format("15-03-2006"))
	}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
			Size:  12,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4F81BD"}, // blue-gray
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	f.SetCellStyle(mainSheet, "A1", "H1", headerStyle)

	optionsSheet := "Options"
	f.NewSheet(optionsSheet)
	options := []string{"pagi", "siang", "malam"}
	for i, val := range options {
		cell := fmt.Sprintf("A%d", i+1)
		f.SetCellValue(optionsSheet, cell, val)
	}

	dv := excelize.NewDataValidation(true)
	dv.Type = "list"
	dv.AllowBlank = true
	dv.ShowDropDown = true
	dv.Formula1 = fmt.Sprintf("%s!$A$1:$A$%d", optionsSheet, len(options))
	dv.Sqref = "F2:F1000"
	if err := f.AddDataValidation(mainSheet, dv); err != nil {

		return nil, err
	}

	bodyStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "DDDDDD", Style: 1},
			{Type: "top", Color: "DDDDDD", Style: 1},
			{Type: "bottom", Color: "DDDDDD", Style: 1},
			{Type: "right", Color: "DDDDDD", Style: 1},
		},
	})
	lastRow := len(entries) + 1
	f.SetCellStyle(mainSheet, "A2", fmt.Sprintf("G%d", lastRow), bodyStyle)

	for i := range header {
		col := string(rune('A' + i))
		f.SetColWidth(mainSheet, col, col, 20)
	}
	f.SetPanes(mainSheet, &excelize.Panes{
		Freeze:      true,
		Split:       true,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	return f, nil
}
