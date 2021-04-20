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

type ResponseCekAntrian struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data CekAntrian `json:"data"`
}

type CekAntrian struct {
	IsAvailable bool `json:"isAvailable"`
	AvailableTime []int `json:"availableTime"`
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
	Tanggal_kedatangan string `json:"tanggalKedatangan"`
	Jam_kedatangan int `json:"jamKedatangan"`
	Id_pelayanan int `json:"idPelayanan"`
}


type ResponseGA struct {
	Status  int             `json:"status"`
	Message string          `json:"messages"`
	Data    GetAntrian `json:"data"`
}

type GetAntrian struct {
	ID int `json:"ID"`
	Nama_lengkap string `json:"namaLengkap"`
	No_identitas string `json:"noIdentitas"`
	Jenis_kelamin string `json:"jenisKelamin"`
	Alamat string `json:"alamat"`
	Email string `json:"email"`
	No_hp string `json:"noHp"`
	Tanggal_kedatangan *time.Time `json:"tanggalKedatangan"`
	Jam_kedatangan *int `json:"jamKedatangan"`
	Id_pelayanan int `json:"idPelayanan"`
	Pelayanan string `json:"pelayanan"`
	No_Pelayanan string `json:"noPelayanan"`
}