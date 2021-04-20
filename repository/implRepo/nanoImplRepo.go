package implRepo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"nano-service/models"
	repo "nano-service/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type mySQLNano struct {
	Conn *sqlx.DB
}

func NewSQLNano(Conn *sqlx.DB) repo.Nano {
	return &mySQLNano{
		Conn: Conn,
	}
}

func (m *mySQLNano) GetPelayanan() ([]models.Pelayanan, error) {
	var arrP []models.Pelayanan
	
	q, err := m.Conn.Queryx(`SELECT * FROM mst_pelayanan ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	
	for q.Next() {
		var p models.Pelayanan
		errScan := q.StructScan(&p)
		if errScan != nil {
			return nil, err
		}

		qD, errD := m.Conn.Queryx(`SELECT value_detail FROM mst_detail_pelayanan WHERE id_pelayanan =$1`, p.ID)
		if errD != nil {
			return nil, errD
		}
		for qD.Next(){
			var d models.DetailPelayanan
			errScanD := qD.Scan(&d.Value_detail)
			if errScanD != nil {
				return nil, errScanD
			}
			p.Description = append(p.Description, d.Value_detail)
		}
		
		arrP = append(arrP, p)
	}
	return arrP, nil
}

func (m *mySQLNano) CekAntrian(tk string,jk int, idp int) bool {
	var totalJam int
	err1 := m.Conn.Get(&totalJam, `select COUNT(jam_kedatangan) from tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan=$2`,tk, idp)
	if err1 != nil {
		log.Panicln(err1)
		return true
	}
	if totalJam >= 20 {
		return false
	}
	return true

}

func (m *mySQLNano) GenerateNoAntrian(idp int, tgl_kedatangan string, jk int)(string, error) {
	var jamK int
	var noAtrian string
	log.Println("JAM ", jamK)
	
	switch idp {
	case 1 :
		err := m.Conn.Get(&jamK,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A",  i)
		} else if jk == 2{
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A",  i)
		} else if jk == 3{
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A",  i)
		}else if jk == 4{
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A",  i)
		}
		
		log.Println("JAM CASE ", jamK)
		// noAtrian = fmt.Sprintf("%s%d", "A", jamK +1)

	case 2 :
		err := m.Conn.Get(&jamK,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B",  i)
		} else if jk == 2{
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B",  i)
		} else if jk == 3{
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B",  i)
		}else if jk == 4{
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B",  i)
		}
		// noAtrian = fmt.Sprintf("%s%d", "B", jamK +1)

	case 3 :
		err := m.Conn.Get(&jamK,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C",  i)
		} else if jk == 2{
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C",  i)
		} else if jk == 3{
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C",  i)
		}else if jk == 4{
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C",  i)
		}

		// noAtrian = fmt.Sprintf("%s%d", "C", jamK +1)

	case 4 :
		err := m.Conn.Get(&jamK,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D",  i)
		} else if jk == 2{
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D",  i)
		} else if jk == 3{
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D",  i)
		}else if jk == 4{
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D",  i)
		}

		// noAtrian = fmt.Sprintf("%s%d", "D", jamK +1)

	

	}
	// log.Println("JAM ", jamK)
	if jamK > 19 {
		return "", errors.New("Antrian Sudah Penuh untuk hari ini")
	}
	return noAtrian, nil
}

type ErrorBody struct {
    Status              int             `json:“status”`
    DetailStatus        int             `json:“detail_status”`
    MessageID           string          `json:“message_id”`
    MessageEN           string          `json:“message_en”`
    Error               string          `json:“error”`
}

func (m *mySQLNano) CreateAntrian(f models.FormIsian) (int, error) {
	// defer m.Conn.Close()
	// dt := time.Now()
	// dates := dt.Format("2006.01.02 15:04:05")

	// ca := `INSERT INTO tran_form_isian (nama_lengkap, no_identitas, jenis_kelamin, alamat, email, no_hp, tanggal_kedatangan, jam_kedatangan, id_pelayanan) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
 	// err := m.Conn.MustExec(ca, f.Nama_lengkap, f.No_identitas, f.Jenis_kelamin, f.Alamat, f.Email, f.No_hp, dates, f.Jam_kedatangan, f.Id_pelayanan)
	// log.Println("ID return ", r)
	var id int
	var rm models.GetAntrian
	noAntrain, errAnt := m.GenerateNoAntrian(f.Id_pelayanan, f.Tanggal_kedatangan, f.Jam_kedatangan)
	if errAnt != nil {
		return 0, errAnt
	}
	fmt.Println("NO ANTRIAN ", noAntrain)

	row, err := m.Conn.NamedQuery(`INSERT INTO tran_form_isian (nama_lengkap, no_identitas, jenis_kelamin, alamat, email, no_hp, tanggal_kedatangan, jam_kedatangan, id_pelayanan, no_antrian) 
	VALUES(:nl, :ni, :jk, :al, :em, :nh, :tk, :jkk, :idp, :na) RETURNING id, tanggal_kedatangan, no_antrian, id_pelayanan, email, nama_lengkap, jam_kedatangan`, map[string]interface{}{
		"nl" : f.Nama_lengkap,
		"ni" : f.No_identitas,
		"jk" : f.Jenis_kelamin,
		"al" : f.Alamat,
		"em" : f.Email,
		"nh" : f.No_hp,
		"tk" : f.Tanggal_kedatangan,
		"jkk" : f.Jam_kedatangan,
		"idp" : f.Id_pelayanan,
		"na" : noAntrain,
	})
		
	if err != nil {
		return 0, err
	}
	for row.Next(){
		row.Scan(&id, &rm.Tanggal_kedatangan, &rm.No_Pelayanan, &rm.Id_pelayanan, &rm.Email, &rm.Nama_lengkap, &rm.Jam_kedatangan)
	}
	jadwal := rm.Tanggal_kedatangan.Format("2006-01-02")
	log.Println("TANGGAL ", jadwal)
	var jamKdtng string
	switch *rm.Jam_kedatangan {
	case 1 : 
		jamKdtng = "09.00 WIB"
	case 2 :
		jamKdtng = "10.00 WIB"
	case 3 : 
		jamKdtng = "11.00 WIB"
	case 4 : 
		jamKdtng = "12.00 WIB"
	}
	var loket string
	errPl := m.Conn.Get(&loket, `SELECT nama FROM mst_pelayanan WHERE id =$1`, rm.Id_pelayanan)
	if errPl != nil {
		log.Println("ID PELAYANAN TIDAK TERSEDIA")
	} 
	t := strconv.Itoa(id)
	log.Println("ID ", t)
	body := map[string]interface{}{
		"id": t,
		"jadwal": jadwal,
		"antrian": noAntrain,
		"loket": loket,
		"email": rm.Email,
		"name": rm.Nama_lengkap,
		"waktu": jamKdtng,
	}
	fmt.Println("body ", body)
	br, errBr := json.Marshal(body)
	if errBr != nil {
		log.Panicln(errBr)
	}
		request, errReq := http.NewRequest("POST", "http://43.229.254.22:8081/generate", bytes.NewBuffer(br))
			request.Header.Set("Content-type", "application/json")
			timeout := time.Duration(5 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			if errReq != nil {
				log.Panicln(errReq.Error())
			}
			resp, errResp := client.Do(request)
			log.Println("LOG BODY RESPONSE ", resp)
			if errResp != nil {
				log.Panicln(errResp.Error())
			}
			defer resp.Body.Close()
			bd, errBody := ioutil.ReadAll(resp.Body)
			
			if errBody != nil {
				log.Panicln(errBody.Error())
			}
			var dataErrorRes ErrorBody
			json.Unmarshal(bd, &dataErrorRes)
			log.Println("LOG REQUEST EMAIL", dataErrorRes)


	return id, nil
}

func (m mySQLNano)GetAvailJam(tk string, idp int) []int{
	var (
		jam1 int
		jam2 int
		jam3 int
		jam4 int
		// jam5 int
	)
	log.Println("PARAMS ", tk, idp)
	var arrJam []int
	err := m.Conn.Get(&jam1,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and jam_kedatangan =1 and id_pelayanan = $2`, tk, idp)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("JAM AJA ", jam1)
	if jam1 < 5 {
		arrJam = append(arrJam, 1)
	}

	err2 := m.Conn.Get(&jam2,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and jam_kedatangan =2 and id_pelayanan = $2`, tk, idp)
	if err != nil {
		log.Panicln(err2)
	}
	if jam2 < 5 {
		arrJam = append(arrJam, 2)
	}

	err3 := m.Conn.Get(&jam3,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and jam_kedatangan =3 and id_pelayanan = $2`, tk, idp)
	if err3 != nil {
		log.Panicln(err3)
	}
	
	if jam3 < 5 {
		arrJam = append(arrJam, 3)
	}


	err4 := m.Conn.Get(&jam4,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and jam_kedatangan =4 and id_pelayanan = $2`, tk, idp)
	if err4 != nil {
		log.Panicln(err4)
	}
	if jam4 < 5 {
		arrJam = append(arrJam, 4)
	}
	// err5 := m.Conn.Get(&jam5,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and jam_kedatangan =5 and id_pelayanan = $2`, tk, idp)
	// if err5 != nil {
	// 	log.Panicln(err5)
	// }
	// if jam5 < 5 {
	// 	arrJam = append(arrJam, 5)
	// }
	log.Println("JAM AVAIL ", arrJam)
	return arrJam
	
}

func (m *mySQLNano) GetAntrian(id int) (models.GetAntrian, error) {
	var f models.GetAntrian
	err := m.Conn.Get(&f, `SELECT t.*, p.nama as pelayanan FROM tran_form_isian t left join mst_pelayanan p on p.id = t.id_pelayanan WHERE t.id = $1`, id)
	if err != nil {
		log.Panicln(err)
		return f, err
	}
	return f, nil
}

func (m *mySQLNano) GetPDF(id int)(models.GetAntrian, error) {
	var f models.GetAntrian
	err := m.Conn.Get(&f, `SELECT t.*, p.nama as pelayanan FROM tran_form_isian t left join mst_pelayanan p on p.id = t.id_pelayanan WHERE t.id = $1`, id)
	if err != nil {
		log.Panicln(err)
		return f, err
	}
	return f, nil
}