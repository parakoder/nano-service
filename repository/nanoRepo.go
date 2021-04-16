package repository

import (
	"nano-service/models"
)
type Nano interface {
	GetPelayanan()([]models.Pelayanan, error)
	CreateAntrian(models.FormIsian)(error)
	GetAntrian(id int)(models.GetAntrian, error)
	CekAntrian(tk string, jkd int, idp int) bool
	GetAvailJam(tk string, idp int) ([]int, error)
	GetPDF(id int)(models.GetAntrian, error)
}