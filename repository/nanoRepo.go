package repository

import (
	"nano-service/models"
)
type Nano interface {
	GetPelayanan()([]models.Pelayanan, error)
	CreateAntrian(models.FormIsian)(error)
	GetAntrian(id int)(models.GetAntrian, error)

	GetPDF(id int)(models.GetAntrian, error)
}