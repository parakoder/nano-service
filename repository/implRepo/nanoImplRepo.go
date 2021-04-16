package implRepo

import (
	"errors"
	"log"
	"nano-service/models"
	repo "nano-service/repository"
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
	
	q, err := m.Conn.Queryx(`SELECT * FROM mst_pelayanan`)
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
	err1 := m.Conn.Get(&totalJam, `select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and jam_kedatangan = $2 and id_pelayanan=$3`,tk, jk, idp)
	// log.Println("data ", totalJam)
	if err1 != nil {
		log.Panicln(err1)
		return true
	}
// log.Println("DATA ", totalJam)
	if totalJam >= 5 {
		return false
	}
	return true

}

func (m *mySQLNano) CreateAntrian(f models.FormIsian) error {
	dt := time.Now()
	dates := dt.Format("2006.01.02 15:04:05")
	// cek := m.CekAntrian(f.Jam_kedatangan)
	// if cek == true {
	// 	return errors.New("Antrian Sudah full")
	// }

	ca := `INSERT INTO tran_form_isian (nama_lengkap, no_identitas, jenis_kelamin, alamat, email, no_hp, tanggal_kedatangan, jam_kedatangan, id_pelayanan) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	err := m.Conn.MustExec(ca, f.Nama_lengkap, f.No_identitas, f.Jenis_kelamin, f.Alamat, f.Email, f.No_hp, dates, f.Jam_kedatangan, f.Id_pelayanan)

	if err == nil {
		return errors.New("Error ketika create antrian")
	}
	return nil
}

func (m mySQLNano)GetAvailJam(tk string, idp int) []int{
	var (
		jam1 int
		jam2 int
		jam3 int
		jam4 int
		jam5 int
	)
	var arrJam []int
	err := m.Conn.Get(&jam1,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and jam_kedatangan =1 and id_pelayanan = $2`, tk, idp)
	if err != nil {
		log.Panicln(err)
	}
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
	err5 := m.Conn.Get(&jam5,`select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and jam_kedatangan =5 and id_pelayanan = $2`, tk, idp)
	if err5 != nil {
		log.Panicln(err5)
	}
	if jam5 < 5 {
		arrJam = append(arrJam, 5)
	}
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