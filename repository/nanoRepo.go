package repository

import (
	"nano-service/models"
	"time"
)
type Nano interface {
	GetPelayanan()([]models.Pelayanan, error)
	GenerateNoAntrian(idp int, tgl_kedatangan time.Time)string
	CreateAntrian(models.FormIsian)(int, error)
	GetAntrian(id int)(models.GetAntrian, error)
	CekAntrian(tk string, jkd int, idp int) bool
	GetAvailJam(tk string, idp int) []int
	GetPDF(id int)(models.GetAntrian, error)
}