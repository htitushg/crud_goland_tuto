package config

import "text/template"

/*type dbase struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	sslmode  string
}*/

/*const (
	Port = ":8090"
)*/

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
type Livre struct {
	Livre_Id    string
	Titre       string
	Editeur_Id  string
	Isbn        string
	Description string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}
type Livre_A struct {
	Livre_Id    string
	Titre       string
	Editeur_Id  string
	Isbn        string
	Description string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
	Auteurs     []Auteur2
}
type Auteur2 struct {
	LMS_Id    string
	AMS_Id    string
	AuteurId  string
	Nom       string
	Prenom    string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}
type UnLivre struct {
	TitrePage string
	Titre     string
	Livre     Livre
}
type ListeLivres struct {
	TitrePage string
	Nombre    int
	Titre     string
	Livres    []Livre
}
type ListeAuteurs struct {
	TitrePage string
	Nombre    int
	Titre     string
	Auteurs   []Auteur
}
type LivreAuteur struct {
	Livre_Id   string
	Titre      string
	Editeur_Id string
	Isbn       string
	Auteur_Id  string
	Nom        string
	Prenom     string
}

type LivresEtAuteurs2 struct {
	TitrePage string
	NombreL   int
	NombreA   int
	LivresA   []Livre_A
}
type LivresEtAuteurs struct {
	TitrePage string
	NombreL   int
	NombreA   int
	Livres    []Livre
	Auteurs   []LivreAuteur
}

type Auteur struct {
	Auteur_Id string
	Nom       string
	Prenom    string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}
type Editeur struct {
	Editeur_Id string
	Nom        string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}
type UnEditeur struct {
	TitrePage string
	Titre     string
	Editeur   Editeur
}
type UnAuteur struct {
	TitrePage string
	Titre     string
	Auteur    Auteur
}
type ListeEditeurs struct {
	TitrePage string
	Nombre    int
	Nom       string
	Editeurs  []Editeur
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
