package repository

import (
	"nano-service/models"
)
type Nano interface {
	GetPelayanan()([]models.Pelayanan, error)
	GenerateNoAntrian(idp int, tgl_kedatangan string, jk int)(string, error)
	CreateAntrian(models.FormIsian)(models.GetAntrian, error)
	GetAntrian(id int)(models.GetAntrian, error)
	CekAntrian(tk string, jkd int, idp int) bool
	GetAvailJam(tk string, idp int) []int
	GetPDF(id int)(models.GetAntrian, error)
}