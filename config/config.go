package config

import "text/template"

type dbase struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	sslmode  string
}

const (
	Port = ":8090"
)

type User struct {
	Id        string
	Nom       string
	Prenom    string
	Email     string
	Password  string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

var Arr User

type ListeUtilisateurs struct {
	Titre  string
	Nombre int
	Name   string
	Users  []User
}
type Utilisateur struct {
	Titre string
	Name  string
	User  User
}
type Config struct {
	TemplateCache map[string]*template.Template
	Port          string
}

var AppConfig Config

//type Car struct {
//	Id           int       `json:"id"`
//	Manufacturer string    `json:"manufacturer"`
//	Design       string    `json:"design"`
//	Style        string    `json:"style"`
//	Doors        uint8     `json:"doors"`
//	CreatedAt    time.Time `json:"created_at"`
//	UpdatedAt    time.Time `json:"updated_at"`
//}
//
//type Cars []Car

/*
Pour les besoins d'inscription dans une base de données MySql:

Conversion de date : string vers []Byte
DatedeCreation = []byte(user.CreatedAt)
DatedeMaj = []byte(time.Now().Format("2006-01-02"))

Conversion de date: []Byte vers string
Arr.UpdatedAt = string(DatedeMaj)
Arr.CreatedAt = string(DatedeCreation)
*/
