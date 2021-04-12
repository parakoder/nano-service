package implRepo

import (
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
		qD, errD := m.Conn.Queryx(`SELECT value_detail FROM mst_detail_pelayanan WHERE id_pelayanan =?`, p.ID)
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

func (m *mySQLNano) CreateAntrian(f models.FormIsian)error {
	dt := time.Now()
	dates := dt.Format("2006.01.02 15:04:05")

	log.Println("tanggal ", dt )
	_, err := m.Conn.NamedQuery(`INSERT INTO tran_form_isian
	(nama_lengkap, no_identitas, jenis_kelamin, alamat, email, no_hp, tanggal_kedatangan, jam_kedatangan, id_pelayanan)
	VALUES(:nl, :ni, :jk, :almt, :email, :nh, :tk, :jkd, :idp)`, map[string]interface{}{
		"nl" : f.Nama_lengkap,
		"ni" : f.No_identitas,
		"jk" : f.Jenis_kelamin,
		"almt" : f.Alamat,
		"email" : f.Email,
		"nh" : f.No_hp,
		"tk" : dates,
		"jkd" : f.Jam_kedatangan,
		"idp" : f.Id_pelayanan,
	})
	if err != nil {
		log.Panicln(err)
		return err
	}
	return nil
}

func (m *mySQLNano) GetAntrian(id int) (models.GetAntrian, error) {
	var f models.GetAntrian
	err := m.Conn.Get(&f, `SELECT t.*, p.nama as pelayanan FROM tran_form_isian t left join mst_pelayanan p on p.id = t.id_pelayanan WHERE t.id = ?`, id)
	if err != nil {
		log.Panicln(err)
		return f, err
	}
	return f, nil
}