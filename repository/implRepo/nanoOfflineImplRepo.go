package implRepo

import (
	"errors"
	"fmt"
	"log"
	"nano-service/models"
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
	start1 := "08:00"
	startParse1, _ := time.Parse(layoutJam, start1)

	end1 := "09:00"
	endParse1, _ := time.Parse(layoutJam, end1)

	jam1 := inTimeSpan(startParse1, endParse1, datesParse)

	// ======================== jam ke 2 ========================
	start2 := "09:00"
	startParse2, _ := time.Parse(layoutJam, start2)

	end2 := "10:00"
	endParse2, _ := time.Parse(layoutJam, end2)

	jam2 := inTimeSpan(startParse2, endParse2, datesParse)

	// ======================== jam ke 3 ========================
	start3 := "10:00"
	startParse3, _ := time.Parse(layoutJam, start3)

	end3 := "11:00"
	endParse3, _ := time.Parse(layoutJam, end3)

	jam3 := inTimeSpan(startParse3, endParse3, datesParse)

	// ======================== jam ke 4 ========================
	start4 := "11:00"
	startParse4, _ := time.Parse(layoutJam, start4)

	end4 := "12:00"
	endParse4, _ := time.Parse(layoutJam, end4)

	jam4 := inTimeSpan(startParse4, endParse4, datesParse)

	// ======================== jam ke 5 ========================
	start5 := "13:00"
	startParse5, _ := time.Parse(layoutJam, start5)

	end5 := "14:00"
	endParse5, _ := time.Parse(layoutJam, end5)

	jam5 := inTimeSpan(startParse5, endParse5, datesParse)

	// ======================== jam ke 6 ========================
	start6 := "14:00"
	startParse6, _ := time.Parse(layoutJam, start6)

	end6 := "15:00"
	endParse6, _ := time.Parse(layoutJam, end6)
	jam6 := inTimeSpan(startParse6, endParse6, datesParse)

	start7 := "15:00"
	startParse7, _ := time.Parse(layoutJam, start7)

	end7 := "23:00"
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
	} else if jam5 == true {
		idJam = 5
	} else if jam6 == true {
		idJam = 6
	} else if jam7 == true {
		idJam = 7
	}

	log.Println("INI DIA ID JAM NYA ", idJam)

	return idJam
}

func (m *mySQLNano) CreateAntrianOffline(f models.FormIsian) (models.GetAntrian, error) {
	// defer m.Conn.Close()
	// dt := time.Now()
	// dates := dt.Format("2006.01.02 15:04:05")

	// ca := `INSERT INTO tran_form_isian (nama_lengkap, no_identitas, jenis_kelamin, alamat, email, no_hp, tanggal_kedatangan, jam_kedatangan, id_pelayanan) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	// err := m.Conn.MustExec(ca, f.Nama_lengkap, f.No_identitas, f.Jenis_kelamin, f.Alamat, f.Email, f.No_hp, dates, f.Jam_kedatangan, f.Id_pelayanan)
	// log.Println("ID return ", r)
	// var id int

	var rm models.GetAntrian
	dt := time.Now()
	currentDate := dt.Format("2006-01-02")
	idJam := getJamKedatanganID()
	noAntrain, errAnt := m.GenerateNoAntrianOffline(f.Id_pelayanan, currentDate, idJam)
	if errAnt != nil {
		return rm, errAnt
	}
	log.Println("DATE ", currentDate)

	row, err := m.Conn.NamedQuery(`INSERT INTO tran_form_isian (nama_lengkap, no_identitas, jenis_kelamin, alamat, email, no_hp, tanggal_kedatangan, jam_kedatangan, id_pelayanan, no_antrian, status, metode) 
	VALUES(:nl, :ni, :jk, :al, :em, :nh, :tk, :jkk, :idp, :na, :st, :mt) RETURNING id, nama_lengkap, no_identitas, jenis_kelamin, alamat, email, no_hp, tanggal_kedatangan, jam_kedatangan, id_pelayanan, no_antrian`, map[string]interface{}{
		"nl":  "userOffline",
		"ni":  nil,
		"jk":  "-",
		"al":  "-",
		"em":  "-",
		"nh":  "-",
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
		row.Scan(&rm.ID, &rm.Nama_lengkap, &rm.No_identitas, &rm.Jenis_kelamin, &rm.Alamat, &rm.Email, &rm.No_hp, &rm.Tanggal_kedatangan, &rm.Jam_kedatangan, &rm.Id_pelayanan, &rm.No_Antrian)
	}
	// // jadwal := rm.Tanggal_kedatangan.Format("2006-01-02")
	// var jamKdtng string
	// switch *rm.Jam_kedatangan {
	// case 1:
	// 	jamKdtng = "08.00 WIB"
	// case 2:
	// 	jamKdtng = "09.00 WIB"
	// case 3:
	// 	jamKdtng = "10.00 WIB"
	// case 4:
	// 	jamKdtng = "11.00 WIB"
	// case 5:
	// 	jamKdtng = "13.00 WIB"
	// case 6:
	// 	jamKdtng = "14.00 WIB"
	// }
	// var loket string
	// errPl := m.Conn.Get(&loket, `SELECT nama FROM mst_pelayanan WHERE id =$1`, rm.Id_pelayanan)
	// if errPl != nil {
	// 	log.Println("ID PELAYANAN TIDAK TERSEDIA")
	// }
	// t := strconv.Itoa(rm.ID)
	// log.Println("ID ", t)
	// body := map[string]interface{}{
	// 	"id":      t,
	// 	"jadwal":  rm.Tanggal_kedatangan,
	// 	"antrian": rm.No_Antrian,
	// 	"loket":   loket,
	// 	"email":   rm.Email,
	// 	"name":    rm.Nama_lengkap,
	// 	"waktu":   jamKdtng,
	// }
	// fmt.Println("body ", body)
	// br, errBr := json.Marshal(body)
	// if errBr != nil {
	// 	log.Panicln(errBr)
	// }
	// request, errReq := http.NewRequest("POST", "http://43.229.254.22:8081/generate", bytes.NewBuffer(br))
	// request.Header.Set("Content-type", "application/json")
	// timeout := time.Duration(30 * time.Second)
	// client := http.Client{
	// 	Timeout: timeout,
	// }
	// if errReq != nil {
	// 	log.Panicln(errReq.Error())
	// }
	// resp, errResp := client.Do(request)
	// log.Println("LOG BODY RESPONSE ", resp)
	// if errResp != nil {
	// 	log.Panicln(errResp.Error())
	// }
	// defer resp.Body.Close()
	// bd, errBody := ioutil.ReadAll(resp.Body)

	// if errBody != nil {
	// 	log.Panicln(errBody.Error())
	// }
	// var dataErrorRes ErrorBody
	// json.Unmarshal(bd, &dataErrorRes)
	// log.Println("LOG REQUEST EMAIL", dataErrorRes)

	return rm, nil
}

func (m *mySQLNano) GenerateNoAntrianOffline(idp int, tgl_kedatangan string, jk int) (string, error) {
	var jamK int
	var noAtrian string

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
		} else if jk == 5 {
			jamKD := 20
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "A ", i)
		} else if jk == 6 {
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
		} else if jk == 5 {
			jamKD := 20
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "B ", i)
		} else if jk == 6 {
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
		} else if jk == 5 {
			jamKD := 20
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "C ", i)
		} else if jk == 6 {
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
		} else if jk == 5 {
			jamKD := 20
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "D ", i)
		} else if jk == 6 {
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
		} else if jk == 5 {
			jamKD := 20
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "E ", i)
		} else if jk == 6 {
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
		} else if jk == 5 {
			jamKD := 20
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "F ", i)
		} else if jk == 6 {
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
		} else if jk == 5 {
			jamKD := 20
			i := jamKD + jamK + 1
			noAtrian = fmt.Sprintf("%s%d", "G ", i)
		} else if jk == 6 {
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
