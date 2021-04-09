package repository

import (
	"nano-service/models"
)
type Nano interface {
	GetPelayanan()([]models.Pelayanan, error)
}