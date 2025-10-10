package requests

import (
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/dto"
	"time"
)

type Demo struct {
	Nama    string    `json:"name"`
	Brand   string    `json:"brand"`
	NoHP    string    `json:"no_hp"`
	Email   string    `json:"email"`
	Waktu   string    `json:"waktu"`
	Tanggal time.Time `json:"tanggal"`
}

type ListDemo struct {
	dto.PaginationRequest
	Search string `form:"search,omitempty"`
}

type ExportDemo struct {
	StartDate *time.Time `form:"start_date,omitempty"`
	EndDate   *time.Time `form:"end_date,omitempty"`
}

func (r *Demo) ToDomain() domain.DemoEntry {

	return domain.DemoEntry{
		Nama:    r.Nama,
		Brand:   r.Brand,
		NoHP:    r.NoHP,
		Email:   r.Email,
		Waktu:   r.Waktu,
		Tanggal: r.Tanggal,
	}
}
