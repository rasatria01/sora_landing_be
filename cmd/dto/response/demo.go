package response

import (
	"sora_landing_be/cmd/domain"
	"time"
)

type (
	Demo struct {
		ID      string    `json:"id"`
		Name    string    `json:"name"`
		Brand   string    `json:"brand"`
		NoHP    string    `json:"no_hp"`
		Email   string    `json:"email"`
		Waktu   string    `json:"waktu"`
		Tanggal time.Time `json:"tanggal"`
	}
	
)

func NewListDemo(demos []domain.DemoEntry) []Demo {
	var res []Demo
	for _, demo := range demos {
		res = append(res, Demo{
			ID:      demo.ID,
			Name:    demo.Nama,
			Brand:   demo.Brand,
			NoHP:    demo.NoHP,
			Email:   demo.Email,
			Waktu:   demo.Waktu,
			Tanggal: demo.Tanggal,
		})
	}

	return res
}

func NewDemo(demo domain.DemoEntry) Demo {
	return Demo{
		ID:      demo.ID,
		Name:    demo.Nama,
		Brand:   demo.Brand,
		NoHP:    demo.NoHP,
		Email:   demo.Email,
		Waktu:   demo.Waktu,
		Tanggal: demo.Tanggal,
	}

}
