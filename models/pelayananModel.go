package models

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