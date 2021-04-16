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
	"github.com/jung-kurt/gofpdf"
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

func (n *NanoRepo) CekAntrian(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Content-type")
	c.Header("Access-Control-Allow-Method", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Origin", "*")
	tk := c.Query("tanggalKedatangan")
	jk := c.Query("jamKedatangan")
	idp := c.Query("idPelayanan")
	j, _ := strconv.Atoi(jk)
	i, _ := strconv.Atoi(idp)

	cek := n.repo.CekAntrian(tk, j, i)
	if cek == true {
		c.AbortWithStatusJSON(400, handler.ErrorHandler(400, 404, "maaf antrian penuh"))
		// log.Panicln(err)
		return
	}
	getJam, _ := n.repo.GetAvailJam(tk,i)

	c.JSON(200, gin.H{
		"status":     200,
		"message_id": "Suskes antrian masih kosong",
		"antrian" : getJam,
	})

}

func (n *NanoRepo) DownloadPdf(c *gin.Context) {
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
	errPdf := GeneratePDF(payload)
	if errPdf != nil {
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


func GeneratePDF(m models.GetAntrian) error{

	
	
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	// pdf.Text(90, 30, m.Pelayanan)
	pdf.CellFormat(250, 7, "judul "+ m.Pelayanan, "0", 0, "CM", false, 0, "")


	errPdf := pdf.OutputFileAndClose("./documents/file2.pdf")
	if errPdf != nil {
		log.Println("ERROR", errPdf.Error())
	}
	return nil
}