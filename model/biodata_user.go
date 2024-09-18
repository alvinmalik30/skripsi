package model

import "time"

type BiodataUser struct {
	Id             string
	UserCredential UserCredential
	NamaLengkap    string
	Nik            string
	NomorTelepon   string
	Pekerjaan      string
	TempatLahir    string
	TanggalLahir   time.Time
	KodePos        string
	IsAglible      bool
	StatusUpdate   bool
	Information    string
}
