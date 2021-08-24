package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"nano-service/config"
	handler "nano-service/handlers"
	"nano-service/models"
	"nano-service/repository"
	"nano-service/repository/implRepo"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/knq/escpos"
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
	var responses models.ResponseAntrian
	var form models.FormIsian
	errBind := c.BindJSON(&form)
	if errBind != nil {
		c.AbortWithStatusJSON(c.Writer.Status(), handler.ErrorHandler(c.Writer.Status(), 404, errBind.Error()))
		log.Panicln(errBind.Error())
		return
	}

	idAnt, err := n.repo.CreateAntrian(form)
	if err != nil {
		c.AbortWithStatusJSON(400, handler.ErrorHandler(400, 404, err.Error()))
		log.Panicln(err)
		return
	}
	responses.Status = 200
	responses.Message = "Success"
	responses.Data = idAnt
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(c.Writer).Encode(responses)

}

func (n *NanoRepo) CreateAntrianOffline(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Content-type")
	c.Header("Access-Control-Allow-Method", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Origin", "*")
	var responses models.ResponseAntrianOff
	var form models.FormIsian
	errBind := c.BindJSON(&form)
	if errBind != nil {
		c.AbortWithStatusJSON(c.Writer.Status(), handler.ErrorHandler(c.Writer.Status(), 404, errBind.Error()))
		log.Panicln(errBind.Error())
		return
	}

	idAnt, err := n.repo.CreateAntrianOffline(form)
	if err != nil {
		c.AbortWithStatusJSON(400, handler.ErrorHandler(400, 404, err.Error()))
		log.Panicln(err)
		return
	}

	PrintTicket(idAnt.No_Antrian, idAnt.Pelayanan, idAnt.Tanggal_kedatangan, idAnt.Jam_kedatangan)

	// fmt.Printf("TES %s", p)

	responses.Status = 200
	responses.Message = "Success"
	responses.Data = idAnt
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(c.Writer).Encode(responses)

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
	var responses models.ResponseCekAntrian
	var cekD models.CekAntrian

	cek := n.repo.CekAntrian(tk, j, i)

	cekD.IsAvailable = cek
	cekD.AvailableTime = n.repo.GetAvailJam(tk, i)
	responses.Status = 200
	responses.Message = "Success"
	responses.Data = cekD
	c.Header("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(c.Writer).Encode(responses)

}

// var (
// 	lpDev = flag.String("p", os.Getenv("PRINTER_ADDRESSS"), "Printer dev file")
// 	// imgPath   = flag.String("i", "/home/septiansah/Pictures/Screenshot from 2021-07-22 15-49-25.png", "Input image")
// 	// threshold = flag.Float64("t", 0.5, "Black/white threshold")
// 	align = flag.String("a", "center", "Alignment (left, center, right)")
// 	doCut = flag.Bool("c", false, "Cut after print")
// 	// maxWidth  = flag.Int("printer-max-width", 512, "Printer max width in pixels")
// )

func PrintTicket(noAntrian, pelayanan, tgl, jam string) {
	

	// flag.Parse()
	// log.Println("INI PAtH ", *lpDev)
	f, err := os.OpenFile(os.Getenv("PRINTER_ADDRESSS"), os.O_RDWR, 0)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()
	// log.Print(*lpDev, " open.")

	ep := escpos.New(f)

	ep.Init()

	// ep.SetAlign(*align)

	ep.SetAlign("center")

	ep.SetSmooth(1)
	ep.SetFontSize(2, 3)
	ep.SetFont("A")
	ep.Write("Pengadilan Negeri Jakarta Selatan")
	ep.Formfeed()

	ep.SetFontSize(1, 2)
	ep.SetFont("A")
	ep.Write("Jl. Ampera Raya No.133, Ragunan, Kec. Ps. Minggu,")
	ep.Write("Kota Jakarta Selatan, Daerah Khusus Ibukota Jakarta 12940")
	ep.Formfeed()

	ep.SetFontSize(3, 4)
	ep.SetFont("B")
	ep.Write("Antrian")
	ep.Formfeed()

	ep.SetFontSize(1, 2)
	ep.Write("Loket " + pelayanan)
	ep.Formfeed()

	ep.SetFontSize(4, 5)
	ep.Write(noAntrian)
	ep.Formfeed()

	ep.SetFontSize(1, 2)
	ep.Write("Tanggal " + tgl + " | Waktu: " + jam)
	ep.FormfeedN(5)
	// if *doCut {
		ep.Cut()
	// }
	ep.End()
}

// func (n *NanoRepo) DownloadPdf(c *gin.Context) {
// 	c.Header("Access-Control-Allow-Headers", "Content-type")
// 	c.Header("Access-Control-Allow-Method", "POST, GET, OPTIONS, PUT, DELETE")
// 	c.Header("Access-Control-Allow-Origin", "*")
// 	var responses models.ResponseGA
// 	id := c.Query("id")
// 	i, _ := strconv.Atoi(id)

// 	payload, err := n.repo.GetAntrian(i)

// 	if err != nil {
// 		c.AbortWithStatusJSON(400, handler.ErrorHandler(400, 404, err.Error()))
// 		log.Panicln(err)
// 		return
// 	}
// 	errPdf := GeneratePDF(payload)
// 	if errPdf != nil {
// 		c.AbortWithStatusJSON(400, handler.ErrorHandler(400, 404, err.Error()))
// 		log.Panicln(err)
// 		return
// 	}

// 	responses.Status = 200
// 	responses.Message = "Success"
// 	responses.Data = payload
// 	c.Header("Content-Type", "application/json")
// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 	json.NewEncoder(c.Writer).Encode(responses)
// }

// func GeneratePDF(m models.GetAntrian) error{

// 	pdf := gofpdf.New("L", "mm", "A4", "")
// 	pdf.AddPage()
// 	pdf.SetFont("Arial", "B", 16)
// 	// pdf.Text(90, 30, m.Pelayanan)
// 	pdf.CellFormat(250, 7, "judul "+ m.Pelayanan, "0", 0, "CM", false, 0, "")

// 	errPdf := pdf.OutputFileAndClose("./documents/file2.pdf")
// 	if errPdf != nil {
// 		log.Println("ERROR", errPdf.Error())
// 	}
// 	return nil
// }
