package controllers

import (
	"encoding/json"
	"log"
	"nano-service/config"
	"nano-service/models"
	"nano-service/repository"
	"nano-service/repository/implRepo"

	"github.com/gin-gonic/gin"
)

type NanoRepo struct {
	repo repository.Nano
}

func NewNanoHandler(db *config.DB) *NanoRepo {
	return &NanoRepo{
		repo: implRepo.NewSQLNano(db.SQL),
	}
}

func (p *NanoRepo) GetPelayanan(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Content-type")
	c.Header("Access-Control-Allow-Method", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Origin", "*")
	var responses models.ResponsePelayanan
	payload, err := p.repo.GetPelayanan()
	if err != nil {
		log.Panicln(err)
	}
	responses.Status = 200
	responses.Message = "Success"
	responses.Data = payload
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(c.Writer).Encode(responses)
}