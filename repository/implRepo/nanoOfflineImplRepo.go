package implRepo

import (
	"errors"
	"fmt"
	"log"
	"nano-service/models"
	"strings"
	"time"
)

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func getJamKedatanganID() int {
	// tx := m.Conn.MustBegin()
	// var jam2 bool
	dt := time.Now()
	var idJam int
	layoutJam := "15:04"
	dates := dt.Format("15:04")

	datesParse, _ := time.Parse(layoutJam, dates)

	// ======================== jam ke 1 ========================
	start1 := "09:00"
	startParse1, _ := time.Parse(layoutJam, start1)

	end1 := "10:00"
	endParse1, _ := time.Parse(layoutJam, end1)

	jam1 := inTimeSpan(startParse1, endParse1, datesParse)

	// ======================== jam ke 2 ========================
	start2 := "10:00"
	startParse2, _ := time.Parse(layoutJam, start2)

	end2 := "11:00"
	endParse2, _ := time.Parse(layoutJam, end2)

	jam2 := inTimeSpan(startParse2, endParse2, datesParse)

	// ======================== jam ke 3 ========================
	start3 := "11:00"
	startParse3, _ := time.Parse(layoutJam, start3)

	end3 := "12:00"
	endParse3, _ := time.Parse(layoutJam, end3)

	jam3 := inTimeSpan(startParse3, endParse3, datesParse)

	// ======================== jam ke 4 ========================
	start4 := "13:00"
	startParse4, _ := time.Parse(layoutJam, start4)

	end4 := "14:00"
	endParse4, _ := time.Parse(layoutJam, end4)

	jam4 := inTimeSpan(startParse4, endParse4, datesParse)

	start7 := "15:00"
	startParse7, _ := time.Parse(layoutJam, start7)

	end7 := "06:00"
	endParse7, _ := time.Parse(layoutJam, end7)

	jam7 := inTimeSpan(startParse7, endParse7, datesParse)

	if jam1 == true {
		idJam = 1
	} else if jam2 == true {
		idJam = 2
	} else if jam3 == true {
		idJam = 3
	} else if jam4 == true {
		idJam = 4
	} else if jam7 == true {
		idJam = 7
	}

	log.Println("INI DIA ID JAM NYA ", idJam)

	return idJam
}

var newStr string

func (m *mySQLNano) CreateAntrianOffline(f models.FormIsian) (models.AntrianOffline, error) {
	// defer m.Conn.Close()
	// dt := time.Now()
	// dates := dt.Format("2006.01.02 15:04:05")

	// ca := `INSERT INTO tran_form_isian (nama_lengkap, no_identitas, jenis_kelamin, alamat, email, no_hp, tanggal_kedatangan, jam_kedatangan, id_pelayanan) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	// err := m.Conn.MustExec(ca, f.Nama_lengkap, f.No_identitas, f.Jenis_kelamin, f.Alamat, f.Email, f.No_hp, dates, f.Jam_kedatangan, f.Id_pelayanan)
	// log.Println("ID return ", r)
	// var id int

	var rm models.AntrianOffline
	var id int
	var idPlyn int
	var tgl *time.Time
	var pelayanan string
	var jam string
	dt := time.Now()
	currentDate := dt.Format("2006-01-02")
	idJam := getJamKedatanganID()
	noAntrain, errAnt := m.GenerateNoAntrianOffline(f.Id_pelayanan, currentDate, idJam)
	if errAnt != nil {
		log.Panic(errAnt)
		return rm, errAnt
	}
	log.Println("DATE ", noAntrain)

	row, err := m.Conn.NamedQuery(`INSERT INTO tran_form_isian (nama_lengkap, tanggal_kedatangan, jam_kedatangan, id_pelayanan, no_antrian, status, metode) 
	VALUES(:nl, :tk, :jkk, :idp, :na, :st, :mt) RETURNING id, no_antrian, tanggal_kedatangan, jam_kedatangan, id_pelayanan`, map[string]interface{}{
		"nl":  "userOffline",
		"tk":  currentDate,
		"jkk": idJam,
		"idp": f.Id_pelayanan,
		"na":  noAntrain,
		"st":  "Waiting",
		"mt":  "offline",
	})

	if err != nil {
		return rm, err
	}
	for row.Next() {
		row.Scan(&id, &rm.No_Antrian, &tgl, &idJam, &idPlyn)
	}
	errPlyn := m.Conn.Get(&pelayanan, `SELECT nama FROM mst_pelayanan where id = $1`, idPlyn)
	if errPlyn != nil {
		log.Panic(errPlyn)
	}

	errJam := m.Conn.Get(&jam, `SELECT keterangan FROM ref_jam_kedatangan where id = $1`, idJam)
	if errJam != nil {
		log.Panic(errJam)
	}
	rm.Pelayanan = pelayanan
	newJam := strings.Split(jam, "-")
	rm.Jam_kedatangan = newJam[0]

	tes := tgl.Format("Monday 02, January 2006")
	tes1 := strings.Split(tes, " ")

	switch tes1[0] {
	case "Monday":
		newStr = "Senin"
	case "Tuesday":
		newStr = "Selasa"
	case "Wednesday":
		newStr = "Rabu"
	case "Thursday":
		newStr = "Kamis"
	case "Friday":
		newStr = "Jumaat"
	case "Saturday":
		newStr = "Sabtu"
	}

	newDate := strings.Replace(tes, tes1[0], newStr, 3)

	log.Println("MANTAP ", rm.Jam_kedatangan)

	rm.Tanggal_kedatangan = newDate
	return rm, nil
}

func (m *mySQLNano) GenerateNoAntrianOffline(idp int, tgl_kedatangan string, jk int) (string, error) {
	var jamK int
	var noAtrian string

	log.Println("PARAMS ", idp, tgl_kedatangan, jk)

	switch idp {
	case 1:
		err := m.Conn.Get(&jamK, `select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3 AND metode = 'offline'`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A ", i)
		} else if jk == 2 {
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A ", i)
		} else if jk == 3 {
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A ", i)
		} else if jk == 4 {
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A ", i)
		} else if jk == 7 {
			jamKD := 25
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A ", i)
		}

		// noAtrian = fmt.Sprintf("%s%d", "A", jamK +1)

	case 2:
		err := m.Conn.Get(&jamK, `select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B ", i)
		} else if jk == 2 {
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B ", i)
		} else if jk == 3 {
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B ", i)
		} else if jk == 4 {
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B ", i)
		} else if jk == 7 {
			jamKD := 25
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B ", i)
		}
		// noAtrian = fmt.Sprintf("%s%d", "B", jamK +1)

	case 3:
		err := m.Conn.Get(&jamK, `select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C ", i)
		} else if jk == 2 {
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C ", i)
		} else if jk == 3 {
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C ", i)
		} else if jk == 4 {
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C ", i)
		} else if jk == 7 {
			jamKD := 25
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C ", i)
		}

		// noAtrian = fmt.Sprintf("%s%d", "C", jamK +1)

	case 4:
		err := m.Conn.Get(&jamK, `select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D ", i)
		} else if jk == 2 {
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D ", i)
		} else if jk == 3 {
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D ", i)
		} else if jk == 4 {
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D ", i)
		} else if jk == 7 {
			jamKD := 25
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D ", i)
		}

	case 5:
		err := m.Conn.Get(&jamK, `select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "E ", i)
		} else if jk == 2 {
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "E ", i)
		} else if jk == 3 {
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "E ", i)
		} else if jk == 4 {
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "E ", i)
		} else if jk == 7 {
			jamKD := 25
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "E ", i)
		}

	case 6:
		err := m.Conn.Get(&jamK, `select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "F ", i)
		} else if jk == 2 {
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "F ", i)
		} else if jk == 3 {
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "F ", i)
		} else if jk == 4 {
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "F ", i)
		} else if jk == 7 {
			jamKD := 25
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "F ", i)
		}

	case 7:
		err := m.Conn.Get(&jamK, `select COUNT(jam_kedatangan) from  tran_form_isian where tanggal_kedatangan::date = $1 and id_pelayanan = $2 and jam_kedatangan = $3`, tgl_kedatangan, idp, jk)
		if err != nil {
			log.Panicln(err)
		}
		if jk == 1 {
			jamKD := 0
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "G ", i)
		} else if jk == 2 {
			jamKD := 5
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "G ", i)
		} else if jk == 3 {
			jamKD := 10
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "G ", i)
		} else if jk == 4 {
			jamKD := 15
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "G ", i)
		} else if jk == 7 {
			jamKD := 25
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "G ", i)
		}

		// noAtrian = fmt.Sprintf("%s%d", "D", jamK +1)
	}
	if jamK > 19 {
		return "", errors.New("Antrian Sudah Penuh untuk hari ini")
	}
	return noAtrian, nil
}
