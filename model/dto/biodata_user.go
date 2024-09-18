package dto

type BiodataRequest struct {
	NamaLengkap  string `json:"full name"`
	Nik          string `json:"nik"`
	NomorTelepon string `json:"phone"`
	Pekerjaan    string `json:"job"`
	TempatLahir  string `json:"place of birth"`
	TanggalLahir string `json:"date of birth"`
	KodePos      string `json:"postal code"`
}
type UpdateBioRequest struct {
	UserCredentialId string `json:"user credential id"`
	IsAglible        bool   `json:"eglibility"`
	Information      string `json:"information"`
}

type BiodataResponse struct {
	Id             string          `json:"id"`
	NamaLengkap    string          `json:"full name"`
	UserCredential GetAuthResponse `json:"user credential"`
	Nik            string          `json:"nik"`
	NomorTelepon   string          `json:"phone"`
	Pekerjaan      string          `json:"job"`
	TempatLahir    string          `json:"place of birth"`
	TanggalLahir   string          `json:"date of birth"`
	KodePos        string          `json:"postal code"`
	IsAglible      bool            `json:"is eglible"`
	StatusUpdate   bool            `json:"status update"`
	Information    string          `json:"information"`
}
