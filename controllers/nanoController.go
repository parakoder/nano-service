package controllers

import (
	"encoding/json"
	"log"
	"nano-service/config"
	handler "nano-service/handlers"
	"nano-service/models"
	"nano-service/repository"
	"nano-service/repository/implRepo"
	"strconv"

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

func (n *NanoRepo) CreateAntrian(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Content-type")
	c.Header("Access-Control-Allow-Method", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Origin", "*")
	var form models.FormIsian
	errBind := c.BindJSON(&form)
	if errBind != nil {
		c.AbortWithStatusJSON(c.Writer.Status(), handler.ErrorHandler(c.Writer.Status(), 404, errBind.Error()))
		log.Panicln(errBind.Error())
		return
	}

	err := n.repo.CreateAntrian(form)
	if err != nil {
		c.AbortWithStatusJSON(400, handler.ErrorHandler(400, 404, err.Error()))
		log.Panicln(err)
		return
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json")
	c.JSON(200, gin.H{
		"status":     200,
		"message_id": "Suskes membuat antrian baru",
	})

}


func (n *NanoRepo) GetAntrian(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Content-type")
	c.Header("Access-Control-Allow-Method", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Origin", "*")
	var responses models.ResponseGA
	id := c.Query("id")
	i, _ := strconv.Atoi(id)

	payload, err := n.repo.GetAntrian(i)
	if err != nil {
		c.AbortWithStatusJSON(400, handler.ErrorHandler(400, 404, err.Error()))
		log.Panicln(err)
		return
	}
	responses.Status = 200
	responses.Message = "Success"
	responses.Data = payload
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(c.Writer).Encode(responses)

}