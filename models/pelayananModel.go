package models

import "time"

type ResponsePelayanan struct {
	Status  int             `json:"status"`
	Message string          `json:"messages"`
	Data    []Pelayanan `json:"data"`
}

type Pelayanan struct {
	ID int `json:"id"`
	Nama string `json:"pelayanan"`
	Description []string `json:"description"`
}


type DetailPelayanan struct {
	ID int `json:"id"`
	Value_detail string `json:"detail"`
	Id_pelayanan int `json:"pelayananID"`
}

type FormIsian struct{
	Nama_lengkap string `json:"namaLengkap"`
	No_identitas string `json:"noIdentitas"`
	Jenis_kelamin string `json:"jenisKelamin"`
	Alamat string `json:"alamat"`
	Email string `json:"email"`
	No_hp string `json:"noHp"`
	Tanggal_kedatangan time.Time `json:"tanggalKedatangan"`
	Jam_kedatangan string `json:"jamKedatangan"`
	Id_pelayanan int `json:"idPelayanan"`
}