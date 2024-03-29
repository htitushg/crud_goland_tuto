package controllers

import (
	"bytes"
	"context"
	"crud_golang_tuto/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"net/http"
	"text/template"
)

func renderTemplate(w http.ResponseWriter, tmplName string, td any) {

	templateCache := appConfig.TemplateCache

	tmpl, ok := templateCache[tmplName+".html"]
	if !ok {
		http.Error(w, "Le template n'existe pas!", http.StatusInternalServerError)
		return
	}
	buffer := new(bytes.Buffer)
	tmpl.Execute(buffer, td)
	buffer.WriteTo(w)
}
func CreateTemplateCache() (map[string]*template.Template, error) {
	path, err2 := os.Getwd()
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(path)
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(path + "/views/*.html")
	if err != nil {
		return cache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		tmpl := template.Must(template.ParseFiles(page))
		layouts, err := filepath.Glob(path + "/views/layouts/*.layout.html")
		if err != nil {
			return cache, err
		}
		if len(layouts) > 0 {
			tmpl.ParseGlob(path + "/views/layouts/*.layout.html")
		}
		cache[name] = tmpl
	}
	return cache, nil
}

var appConfig *config.Config

func CreateTemplates(app *config.Config) {
	appConfig = app
}

func Home(w http.ResponseWriter, r *http.Request) {
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	rows, err := Db.Query("SELECT * FROM users where DeletedAt < ?", "1900-01-02")
	config.CheckError(err)
	defer rows.Close()
	UnUser := config.User{}
	myUsers := make([]config.User, 0)
	i := 0
	for rows.Next() {
		err = rows.Scan(&UnUser.Nom, &UnUser.Prenom, &UnUser.Email, &UnUser.Password, &UnUser.Id, &UnUser.CreatedAt, &UnUser.UpdatedAt, &UnUser.DeletedAt)
		config.CheckError(err)
		UnUser.Id = strings.Join(strings.Fields(UnUser.Id), "")
		UnUser.Email = strings.Join(strings.Fields(UnUser.Email), "")
		myUsers = append(myUsers, UnUser)
		i++

		//fmt.Printf("Id : %s, Nom: %s, Prenom: %s, Courriel: %s, Pass: %s\n", UnUser.Id, UnUser.Nom, UnUser.Prenom, UnUser.Email, UnUser.Password)

	}
	data := config.ListeUtilisateurs{
		Nombre: i,
		Name:   "Liste des utilisateurs inscrits dans la base mysql",
		Users:  myUsers,
	}
	//fmt.Printf("data.Name : %s, data.Nombre: %d", data.Name, data.Nombre)

	renderTemplate(w, "home", &data)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "AddUser", nil)
}

func AddUserPost(w http.ResponseWriter, r *http.Request) {

	// Ajouté le 8/11/2023 à 12h30

	user := config.User{
		Id:        r.FormValue("id"),
		Nom:       r.FormValue("nom"),
		Prenom:    r.FormValue("prenom"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		CreatedAt: r.FormValue("CreatedAt"),
		UpdatedAt: r.FormValue("UpdatedAt"),
		DeletedAt: r.FormValue("DeletedAt"),
	}

	user.Email = strings.Join(strings.Fields(user.Email), "")
	var DatedeCreation, DatedeMaj []uint8

	DatedeCreation = []byte(time.Now().Format("2006-01-02"))
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))

	//fmt.Printf("Avant création dans la table = Id= %s, Nom: %s, Prenom: %s, Courriel: %s, Pass: %s, CreatedAt: %s, UpdatedAt: %s\n", user.Id, user.Nom, user.Prenom, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	// INSERT INTO users
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)

	query := "INSERT INTO users (Nom, Prenom, Email,Password,CreatedAt, UpdatedAt, DeletedAt) VALUES (?, ?, ?, ?, ?, ?, ?)"
	insertResult, err := Db.ExecContext(context.Background(), query, user.Nom, user.Prenom, user.Email, user.Password, DatedeCreation, DatedeMaj, "1900-01-01")
	if err != nil {
		log.Fatalf("impossible insert user: %s", err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}
	log.Printf("inserted id: %d", id)
	rows, err := Db.Query("SELECT * FROM users where DeletedAt < ?", "1900-01-02")
	config.CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.Nom, &user.Prenom, &user.Email, &user.Password, &user.Id, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		config.CheckError(err)
		user.UpdatedAt = string(DatedeMaj)
		user.CreatedAt = string(DatedeCreation)
		user.Email = strings.Join(strings.Fields(user.Email), "")
	}
	data := config.Utilisateur{
		Name: "Création d'un Utilisateur :" + user.Nom,
		User: user,
	}
	// Fin de l'ajout du 8/11/2023 à 12h30

	renderTemplate(w, "AddUserPost", &data)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/UpdateUser/")
	index, err := strconv.Atoi(id)
	config.CheckError(err)

	var DatedeMaj []uint8 = []byte(time.Now().Format("2006-01-02"))
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	unuser := config.User{}
	rows, err := Db.Query("SELECT * FROM users WHERE Id= ?", index)
	config.CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&unuser.Nom, &unuser.Prenom, &unuser.Email, &unuser.Password, &index, &unuser.CreatedAt, &DatedeMaj, &unuser.DeletedAt) //DatedeCreation, &DatedeMaj, &index)
		config.CheckError(err)
		unuser.Email = strings.Join(strings.Fields(unuser.Email), "")
		unuser.Id = id
		unuser.UpdatedAt = string(DatedeMaj)
		//unuser.CreatedAt = string(DatedeCreation)
		//unuser.DeletedAt = string(DatedeCreation)
	}

	//fmt.Printf("CreatedAt: %#v \nUpdatedAt: %#v \n", DatedeCreation, DatedeMaj)
	//fmt.Printf("UpdateUser: Id : %s, Nom: %s, Prenom: %s, Courriel: %s, Pass: %s, CreatedDat: %s, UpdatedAt: %s \n", unuser.Id, unuser.Nom, unuser.Prenom, unuser.Email, unuser.Password, unuser.CreatedAt, DatedeMaj)

	data := config.Utilisateur{
		Name: "Mise à jour de l'utilisateur :" + unuser.Prenom + " " + unuser.Nom,
		User: unuser,
	}

	renderTemplate(w, "UpdateUser", &data)

}

func UpdateUserPost(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/UpdateUserPost/")
	//fmt.Printf("UpdateUserPost avant conversion id : %s ", id)
	index, err := strconv.Atoi(id)
	config.CheckError(err)

	user := config.User{
		Nom:       r.FormValue("nom"),
		Prenom:    r.FormValue("prenom"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		CreatedAt: r.FormValue("CreatedAt"),
		UpdatedAt: r.FormValue("UpdatedAt"),
		DeletedAt: r.FormValue("DeletedAt"),
	}
	user.Email = strings.Join(strings.Fields(user.Email), "")
	var DatedeCreation, DatedeMaj []uint8
	DatedeCreation = []byte(user.CreatedAt)
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))
	fmt.Printf("UpdateUserPost avant calcul DateEffacement : %s ", user.DeletedAt)
	// UPDATE INTO users
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	config.PingDB(Db)

	//Update db
	stmt, err := Db.Prepare("UPDATE users SET nom=? , prenom=?, email=?, password=?, CreatedAt=?, UpdatedAt=?, DeletedAt=? where id =? ")
	config.CheckError(err)

	// execute
	res, err := stmt.Exec(user.Nom, user.Prenom, user.Email, user.Password, DatedeCreation, DatedeMaj, "1900-01-01", index)
	config.CheckError(err)

	a, err := res.RowsAffected()
	config.CheckError(err)

	log.Printf("updated id: %d, %v", index, a)

	// Fin de la mise à jour du 8/11/2023 à 12h30
	user.UpdatedAt = string(DatedeMaj)
	user.CreatedAt = string(DatedeCreation)

	data := config.Utilisateur{
		Name: "Mise à jour Utilisateur :" + user.Nom + user.Prenom,
		User: user,
	}

	renderTemplate(w, "AddUserPost", &data)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/DeleteUser/")
	index, err := strconv.Atoi(id)
	config.CheckError(err)
	//Arr := config.User{}
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	unuser := config.User{}
	rows, err := Db.Query("SELECT * FROM users WHERE Id= ?", index)
	config.CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&unuser.Nom, &unuser.Prenom, &unuser.Email, &unuser.Password, &index, &unuser.CreatedAt, &unuser.UpdatedAt, &unuser.DeletedAt)
		config.CheckError(err)
	}
	unuser.Email = strings.Join(strings.Fields(unuser.Email), "")
	unuser.Id = id
	//fmt.Printf("DeleteUser= Id : %s, Nom: %s, Prenom: %s, Courriel: %s, Pass: %s, CreatedAt: %s, UpdatedAt: %s\n", unuser.Id, unuser.Nom, unuser.Prenom, unuser.Email, unuser.Password, unuser.CreatedAt, unuser.UpdatedAt)
	data := config.Utilisateur{
		Name: "Effacer un Utilisateur :" + unuser.Prenom + " " + unuser.Nom,
		User: unuser,
	}

	renderTemplate(w, "DeleteUser", &data)
}

func DeleteUserPost(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/DeleteUserPost/")
	//fmt.Printf("DeleteUserPost avant conversion id : %s ", id)
	index, err := strconv.Atoi(id)
	config.CheckError(err)
	// Ajouté le 8/11/2023 à 12h30
	Arr := config.User{
		Nom:       r.FormValue("nom"),
		Prenom:    r.FormValue("prenom"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		CreatedAt: r.FormValue("CreatedAt"),
		UpdatedAt: r.FormValue("UpdatedAt"),
		DeletedAt: r.FormValue("DeletedAt"),
	}
	Arr.Email = strings.Join(strings.Fields(Arr.Email), "")
	var DateEffacement []uint8 = []byte(time.Now().Format("2006-01-02"))
	// Delete utilisateur dans users
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	config.PingDB(Db)

	//Update db
	//stmt, err := Db.Prepare("DELETE FROM users WHERE id =?")
	stmt, err := Db.Prepare("UPDATE users SET nom=? , prenom=?, email=?, password=?, CreatedAt=?, UpdatedAt=?, DeletedAt=? where id =? ")
	config.CheckError(err)

	// execute
	res, err := stmt.Exec(Arr.Nom, Arr.Prenom, Arr.Email, Arr.Password, Arr.CreatedAt, Arr.UpdatedAt, DateEffacement, index)
	config.CheckError(err)

	a, err := res.RowsAffected()
	config.CheckError(err)

	log.Printf("deleted id: %v, %v", id, a)

	// Fin de la mise à jour du 8/11/2023 à 12h30
	data := config.Utilisateur{
		Name: "l'utilisateur " + Arr.Nom + " " + Arr.Prenom + " a été effacé",
	}

	renderTemplate(w, "DeleteUserPost", &data)

}
func Contact(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Contact", nil)
}
func AddLivre(w http.ResponseWriter, r *http.Request) {
	//Ouverture de la base et execution requete pour aller chercher les editeurs
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	rows, err := Db.Query("SELECT * FROM Editeurs where DeletedAt < ?", "1900-01-02")
	config.CheckError(err)
	defer rows.Close()
	var unEditeur config.Editeur
	Editeurs := make([]config.Editeur, 0)
	i := 0
	for rows.Next() {
		err = rows.Scan(&unEditeur.Editeur_Id, &unEditeur.Nom, &unEditeur.CreatedAt, &unEditeur.UpdatedAt, &unEditeur.DeletedAt)
		config.CheckError(err)
		Editeurs = append(Editeurs, unEditeur)
		fmt.Printf("ListeDesEditeurs 324 Editeur_Id : %s, Nom: %s\n", unEditeur.Editeur_Id, unEditeur.Nom)
		i++
	}
	fmt.Printf("ListeDesEditeurs 327 Editeurs : %#v", Editeurs)
	renderTemplate(w, "AddLivre", &Editeurs)
}

func AddLivrePost(w http.ResponseWriter, r *http.Request) {

	livre := config.Livre{
		Livre_Id:    r.FormValue("Livre_Id"),
		Titre:       r.FormValue("Titre"),
		Editeur_Id:  r.FormValue("Editeur_Id"),
		Isbn:        r.FormValue("Isbn"),
		Description: r.FormValue("Description"),
		CreatedAt:   r.FormValue("CreatedAt"),
		UpdatedAt:   r.FormValue("UpdatedAt"),
		DeletedAt:   r.FormValue("DeletedAt"),
	}
	var DatedeCreation, DatedeMaj []uint8
	DatedeCreation = []byte(time.Now().Format("2006-01-02"))
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))

	//fmt.Printf("Avant création du livre = Id= %s, Titre: %s, Auteur_Id: %s, Editeur: %s, Isbn: %s, CreatedAt: %s, UpdatedAt: %s\n", livre.Id, livre.Titre, livre.Auteur_Id, livre.Editeur, livre.Isbn, livre.CreatedAt, livre.UpdatedAt)

	// INSERT INTO livres
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	fmt.Printf("AddLivrePost 352 Livre_Id: %s, Titre: %s, Editeur_Id : %s, Isbn: %s, Description: %s", livre.Livre_Id, livre.Titre, livre.Editeur_Id, livre.Isbn, livre.Description)
	id_editeur, err := strconv.Atoi(livre.Editeur_Id)
	if err != nil {
		log.Fatalf("impossible de convertir livre.Editeur_Id: %s", err)
	}
	fmt.Printf("AddLivrePost 357 id_editeur: %d\n", id_editeur)
	query := "INSERT INTO Livres (Titre, Editeur_Id, Isbn, Description, CreatedAt, UpdatedAt, DeletedAt) VALUES (?, ?, ?, ?, ?, ?, ?)"
	insertResult, err := Db.ExecContext(context.Background(), query, livre.Titre, id_editeur, livre.Isbn, livre.Description, DatedeCreation, DatedeMaj, "1900-01-01")
	if err != nil {
		log.Fatalf("Impossible d'inserer le livre: %s\n", err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}
	log.Printf("inserted id: %d", id)
	data := config.UnLivre{
		TitrePage: "Création d'un Livre :" + livre.Titre,
		Livre:     livre,
	}
	renderTemplate(w, "AddLivrePost", &data)
}
func ListeDesLivres(w http.ResponseWriter, r *http.Request) {
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	rows, err := Db.Query("SELECT * FROM Livres where DeletedAt < ?", "1900-01-02")
	config.CheckError(err)
	defer rows.Close()
	unlivre := config.Livre{}
	mylivres := make([]config.Livre, 0)
	i := 0
	for rows.Next() {
		err = rows.Scan(&unlivre.Livre_Id, &unlivre.Titre, &unlivre.Editeur_Id, &unlivre.Isbn, &unlivre.Description, &unlivre.CreatedAt, &unlivre.UpdatedAt, &unlivre.DeletedAt)
		config.CheckError(err)
		mylivres = append(mylivres, unlivre)
		fmt.Printf("ListeDesLivres 387 Id : %s, Titre: %s, Editeur: %s, Isbn: %s\n", unlivre.Livre_Id, unlivre.Titre, unlivre.Editeur_Id, unlivre.Isbn)
		i++
	}
	data := config.ListeLivres{
		TitrePage: "Liste des livres inscrits dans la base mysql",
		Nombre:    i,
		Titre:     "Liste des Livres",
		Livres:    mylivres,
	}
	fmt.Printf("ListeDesLivres 396 data.TitrePage : %s, data.Nombre: %d, data.Titre: %s", data.TitrePage, data.Nombre, data.Titre)

	renderTemplate(w, "ListeDesLivres", &data)
}
func UpdateLivre(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/UpdateLivre/")
	index, err := strconv.Atoi(id)
	config.CheckError(err)

	var DatedeMaj []uint8 = []byte(time.Now().Format("2006-01-02"))
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	mylivre := config.Livre{}
	rows, err := Db.Query("SELECT * FROM Livres WHERE Livre_Id= ?", index)
	config.CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&index, &mylivre.Titre, &mylivre.Editeur_Id, &mylivre.Isbn, &mylivre.Description, &mylivre.CreatedAt, &DatedeMaj, &mylivre.DeletedAt) //DatedeCreation, &DatedeMaj, &index)
		config.CheckError(err)
		mylivre.Livre_Id = id
		mylivre.UpdatedAt = string(DatedeMaj)
	}

	//fmt.Printf("CreatedAt: %#v \nUpdatedAt: %#v \n", DatedeCreation, DatedeMaj)
	//fmt.Printf("UpdateUser: Id : %s, Nom: %s, Prenom: %s, Courriel: %s, Pass: %s, CreatedDat: %s, UpdatedAt: %s \n", unuser.Id, unuser.Nom, unuser.Prenom, unuser.Email, unuser.Password, unuser.CreatedAt, DatedeMaj)
	data := config.UnLivre{
		TitrePage: "Mise à jour du livre: " + mylivre.Titre,
		Livre:     mylivre,
	}
	fmt.Printf("UpdateLivre 426 Id : %s, Titre: %s, Editeur: %s, Isbn: %s, CreatedAt: %s, UpdatedAt: %s\n", mylivre.Livre_Id, mylivre.Titre, mylivre.Editeur_Id, mylivre.Isbn, mylivre.CreatedAt, mylivre.UpdatedAt)

	renderTemplate(w, "UpdateLivre", &data)
}

func UpdateLivrePost(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/UpdateLivrePost/")
	fmt.Printf("UpdateLivrePost 433 avant conversion id : %s ", id)
	index, err := strconv.Atoi(id)
	config.CheckError(err)
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	config.PingDB(Db)
	mylivre := config.Livre{
		Livre_Id:    r.FormValue("Livre_Id"),
		Titre:       r.FormValue("Titre"),
		Editeur_Id:  r.FormValue("Editeur_Id"),
		Isbn:        r.FormValue("Isbn"),
		Description: r.FormValue("Description"),
		CreatedAt:   r.FormValue("CreatedAt"),
		UpdatedAt:   r.FormValue("UpdatedAt"),
		DeletedAt:   r.FormValue("DeletedAt"),
	}
	mylivre.Livre_Id = id
	fmt.Printf("UpdateLivrePost 450 Livre_Id: %s, Titre: %s, Editeur_Id : %s, Isbn: %s, Description: %s\n", mylivre.Livre_Id, mylivre.Titre, mylivre.Editeur_Id, mylivre.Isbn, mylivre.Description)
	var DatedeCreation, DatedeMaj []uint8
	mylivre.UpdatedAt = string(DatedeMaj)

	DatedeCreation = []byte(mylivre.CreatedAt)
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))
	id_editeur, err := strconv.Atoi(mylivre.Editeur_Id)
	config.CheckError(err)
	fmt.Printf("UpdateLivrePost 458 Livre_Id: %s, index: %d, Titre: %s, Editeur_Id : %s, id_editeur: %d, Isbn: %s, Description: %s\n", mylivre.Livre_Id, index, mylivre.Titre, mylivre.Editeur_Id, id_editeur, mylivre.Isbn, mylivre.Description)
	// UPDATE INTO livres
	stmt, err := Db.Prepare("UPDATE Livres SET Titre=? , Editeur_Id=?, Isbn=?, Description=?, CreatedAt=?, UpdatedAt=?, DeletedAt=? where Livre_Id =? ")
	config.CheckError(err)

	// execute
	res, err := stmt.Exec(mylivre.Titre, id_editeur, mylivre.Isbn, mylivre.Description, DatedeCreation, DatedeMaj, "1900-01-01", index)
	config.CheckError(err)

	a, err := res.RowsAffected()
	config.CheckError(err)

	log.Printf("updated id: %d, %v", index, a)

	// Fin de la mise à jour du 8/11/2023 à 12h30
	mylivre.UpdatedAt = string(DatedeMaj)
	mylivre.CreatedAt = string(DatedeCreation)
	mylivre.Livre_Id = id
	data := config.UnLivre{
		TitrePage: "Mise à jour Utilisateur :" + mylivre.Titre,
		Livre:     mylivre,
	}
	fmt.Printf("UpdateLivre 480 Id : %s, Titre: %s, Editeur: %s, Isbn: %s, Description: %s, CreatedAt: %s, UpdatedAt: %s\n", mylivre.Livre_Id, mylivre.Titre, mylivre.Editeur_Id, mylivre.Isbn, mylivre.Description, mylivre.CreatedAt, mylivre.UpdatedAt)

	renderTemplate(w, "UpdateLivrePost", &data)
}
func DeleteLivre(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/DeleteLivre/")
	index, err := strconv.Atoi(id)
	config.CheckError(err)
	//Arr := config.User{}
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	unlivre := config.Livre{}
	rows, err := Db.Query("SELECT * FROM Livres WHERE Livre_Id= ?", index)
	config.CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&index, &unlivre.Titre, &unlivre.Editeur_Id, &unlivre.Isbn, &unlivre.Description, &unlivre.CreatedAt, &unlivre.UpdatedAt, &unlivre.DeletedAt)
		config.CheckError(err)
	}
	unlivre.Livre_Id = id
	fmt.Printf("DeleteLivre 496  Id : %s, Titre: %s, Editeur: %s, Isbn: %s, Description: %s, CreatedAt: %s, UpdatedAt: %s\n", unlivre.Livre_Id, unlivre.Titre, unlivre.Editeur_Id, unlivre.Isbn, unlivre.Description, unlivre.CreatedAt, unlivre.UpdatedAt)
	data := config.UnLivre{
		TitrePage: "Effacer un Utilisateur :" + unlivre.Titre,
		Livre:     unlivre,
	}
	renderTemplate(w, "DeleteLivre", &data)
}
func DeleteLivrePost(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/DeleteLivrePost/")
	//fmt.Printf("DeleteUserPost avant conversion id : %s ", id)
	index, err := strconv.Atoi(id)
	config.CheckError(err)
	// Ajouté le 8/11/2023 à 12h30
	unlivre := config.Livre{
		Titre:       r.FormValue("Titre"),
		Editeur_Id:  r.FormValue("Editeur"),
		Isbn:        r.FormValue("Isbn"),
		Description: r.FormValue("Description"),
		CreatedAt:   r.FormValue("CreatedAt"),
		UpdatedAt:   r.FormValue("UpdatedAt"),
		DeletedAt:   r.FormValue("DeletedAt"),
	}
	var DateEffacement []uint8 = []byte(time.Now().Format("2006-01-02"))
	id_editeur, err := strconv.Atoi(unlivre.Editeur_Id)
	config.CheckError(err)
	// Delete utilisateur dans livres
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	config.PingDB(Db)

	//Update db
	//stmt, err := Db.Prepare("DELETE FROM livres WHERE Livre_Id =?")
	stmt, err := Db.Prepare("UPDATE Livres SET Titre=? ,Editeur_Id=?, Isbn=?, Description=?, CreatedAt=?, UpdatedAt=?, DeletedAt=? where Livre_Id =? ")
	config.CheckError(err)

	// execute
	res, err := stmt.Exec(unlivre.Titre, id_editeur, unlivre.Isbn, unlivre.Description, unlivre.CreatedAt, unlivre.UpdatedAt, DateEffacement, index)
	config.CheckError(err)

	a, err := res.RowsAffected()
	config.CheckError(err)

	log.Printf("deleted id: %v, %v", id, a)

	// Fin de la mise à jour du 8/11/2023 à 12h30
	data := config.UnLivre{
		TitrePage: "l'utilisateur " + unlivre.Titre + " a été effacé",
		Livre:     unlivre,
	}
	fmt.Printf("DeleteLivrePost= Id : %s, Titre: %s, Editeur: %s, Isbn: %s, Description: %s, CreatedAt: %s, UpdatedAt: %s, DateEffacement: %s\n", unlivre.Livre_Id, unlivre.Titre, unlivre.Editeur_Id, unlivre.Isbn, unlivre.Description, unlivre.CreatedAt, unlivre.UpdatedAt, DateEffacement)

	renderTemplate(w, "DeleteLivrePost", &data)

}

func AfficheLivresEtAuteurs(w http.ResponseWriter, r *http.Request) {

	// Modif du 23/11
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	Unlivre := config.Livre_A{}
	//unlivreauteur := config.LivreAuteur{}
	var desLivres []config.Livre_A
	//desauteurs := make([]config.Livre_A, 0)
	rowsL, err := Db.Query("SELECT Livres.Livre_Id, Livres.Titre, Livres.Editeur_Id, Livres.Isbn, Livres.CreatedAt, Livres.UpdatedAt, Livres.DeletedAt FROM Livres")
	config.CheckError(err)
	i := 0
	j := 0
	for rowsL.Next() {

		err = rowsL.Scan(&Unlivre.Livre_Id, &Unlivre.Titre, &Unlivre.Editeur_Id, &Unlivre.Isbn, &Unlivre.CreatedAt, &Unlivre.UpdatedAt, &Unlivre.DeletedAt)
		config.CheckError(err)
		desLivres = append(desLivres, Unlivre)
		fmt.Printf("LLLL-Livre_Id : %s, Titre: %s, Editeur_Id: %s, Isbn: %s\n",
			desLivres[i].Livre_Id, desLivres[i].Titre, desLivres[i].Editeur_Id, desLivres[i].Isbn)
		rowsA, err := Db.Query("SELECT Auteurs.Auteur_Id, Auteurs.Nom,Auteurs.Prenom FROM Livres JOIN  LivreMembership ON Livres.Livre_Id = LivreMembership.Livre_Id JOIN Auteurs	ON LivreMembership.Auteur_Id = Auteurs.Auteur_Id WHERE Livres.Livre_Id=?", desLivres[i].Livre_Id)
		//rowsa, err := Db.Query("SELECT Auteurs.Auteur_Id, Auteurs.Nom, Auteurs.Prenom FROM LivreMembership JOIN Auteurs ON LivreMembership.Auteur_Id = Auteurs.Auteur_Id WHERE Livre_Id = ?", desLivres[i].Livre_Id)
		config.CheckError(err)
		for rowsA.Next() {
			fmt.Printf("LLLL 582 i: %d, j: %d\n", i, j)
			err = rowsA.Scan(&desLivres[i].Auteurs[j].AuteurId, &desLivres[i].Auteurs[j].Nom, &desLivres[i].Auteurs[j].Prenom)
			config.CheckError(err)
			//fmt.Printf("46-Livre_Id : %s, Auteur_Id : %s, Nom: %s, Prenom : %s, \n", unlivreauteur.Livre_Id, unlivreauteur.Auteur_Id, unlivreauteur.Nom, unlivreauteur.Prenom)
			//desauteurs = append(desauteurs, unlivreauteur)
			fmt.Printf("48-Auteur_Id : %s, Nom: %s, Prenom: %s\n", desLivres[j].Auteurs[i].AuteurId, desLivres[j].Auteurs[i].Nom, desLivres[j].Auteurs[i].Prenom)
			j++
		}
		i++
		defer rowsA.Close()
	}
	defer rowsL.Close()
	data := config.LivresEtAuteurs2{
		TitrePage: "Liste des Livres et de leurs auteurs",
		NombreL:   i,
		NombreA:   j,
		LivresA:   desLivres,
	}
	renderTemplate(w, "AfficheLivresEtAuteurs", &data)

}

func AddEditeur(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "AddEditeur", nil)
}

func AddEditeurPost(w http.ResponseWriter, r *http.Request) {

	// Ajouté le 8/11/2023 à 12h30

	editeur := config.Editeur{

		Editeur_Id: r.FormValue("Editeur_Id"),
		Nom:        r.FormValue("Nom"),
		CreatedAt:  r.FormValue("CreatedAt"),
		UpdatedAt:  r.FormValue("UpdatedAt"),
		DeletedAt:  r.FormValue("DeletedAt"),
	}

	var DatedeCreation, DatedeMaj []uint8
	DatedeCreation = []byte(time.Now().Format("2006-01-02"))
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))

	//fmt.Printf("Avant création du livre = Id= %s, Titre: %s, Auteur_Id: %s, Editeur: %s, Isbn: %s, CreatedAt: %s, UpdatedAt: %s\n", livre.Id, livre.Titre, livre.Auteur_Id, livre.Editeur, livre.Isbn, livre.CreatedAt, livre.UpdatedAt)

	// INSERT INTO livres
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	fmt.Printf("AddEditeurPost 601 Editeur_Id: %s, Nom : %s\n", editeur.Editeur_Id, editeur.Nom)
	//id_editeur, err := strconv.Atoi(editeur.Editeur_Id)
	if err != nil {
		log.Fatalf("impossible de convertir editeur.Editeur_Id: %s", err)
	}
	//fmt.Printf("AddEditeurPost 606 id_editeur: %d\n", id_editeur)
	query := "INSERT INTO Editeurs (Nom, CreatedAt, UpdatedAt, DeletedAt) VALUES (?, ?, ?, ?)"
	insertResult, err := Db.ExecContext(context.Background(), query, editeur.Nom, DatedeCreation, DatedeMaj, "1900-01-01")
	if err != nil {
		log.Fatalf("Impossible d'inserer le livre: %s\n", err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}
	log.Printf("inserted id: %d", id)

	data := config.UnEditeur{
		TitrePage: "Création d'un Editeur :" + editeur.Nom,
		Editeur:   editeur,
	}
	// Fin de l'ajout du 8/11/2023 à 12h30

	renderTemplate(w, "AddEditeurPost", &data)
}
func ListeDesEditeurs(w http.ResponseWriter, r *http.Request) {
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	rows, err := Db.Query("SELECT * FROM Editeurs where DeletedAt < ?", "1900-01-02")
	config.CheckError(err)
	defer rows.Close()
	unEditeur := config.Editeur{}
	myEditeurs := make([]config.Editeur, 0)
	i := 0
	for rows.Next() {
		err = rows.Scan(&unEditeur.Editeur_Id, &unEditeur.Nom, &unEditeur.CreatedAt, &unEditeur.UpdatedAt, &unEditeur.DeletedAt)
		config.CheckError(err)
		myEditeurs = append(myEditeurs, unEditeur)
		fmt.Printf("ListeDesEditeurs 664 Editeur_Id : %s, Nom: %v, CreatedAt: %s, UpdatedAt: %s, DeletedAt: %s\n", unEditeur.Editeur_Id, &unEditeur.Nom, unEditeur.CreatedAt, unEditeur.UpdatedAt, unEditeur.DeletedAt)
		i++
	}
	data := config.ListeEditeurs{
		TitrePage: "Liste des Editeurs inscrits dans la base mysql",
		Nombre:    i,
		Editeurs:  myEditeurs,
	}
	fmt.Printf("ListeDesEditeurs 672 data.TitrePage : %s, data.Nombre: %d\n", data.TitrePage, data.Nombre)

	renderTemplate(w, "ListeDesEditeurs", &data)
}
func UpdateEditeur(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/UpdateEditeur/")
	index, err := strconv.Atoi(id)
	config.CheckError(err)

	var DatedeMaj []uint8 = []byte(time.Now().Format("2006-01-02"))
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	myEditeur := config.Editeur{}
	rows, err := Db.Query("SELECT * FROM Editeurs WHERE Editeur_Id= ?", index)
	config.CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&myEditeur.Editeur_Id, &myEditeur.Nom, &myEditeur.CreatedAt, &DatedeMaj, &myEditeur.DeletedAt, &index) //DatedeCreation, &DatedeMaj, &index)
		config.CheckError(err)
		myEditeur.Editeur_Id = id
		myEditeur.UpdatedAt = string(DatedeMaj)
	}

	//fmt.Printf("CreatedAt: %#v \nUpdatedAt: %#v \n", DatedeCreation, DatedeMaj)
	//fmt.Printf("UpdateUser: Id : %s, Nom: %s, Prenom: %s, Courriel: %s, Pass: %s, CreatedDat: %s, UpdatedAt: %s \n", unuser.Id, unuser.Nom, unuser.Prenom, unuser.Email, unuser.Password, unuser.CreatedAt, DatedeMaj)
	data := config.UnEditeur{
		TitrePage: "Mise à jour du livre: " + myEditeur.Nom,
		Editeur:   myEditeur,
	}
	fmt.Printf("UpdateEditeur 716 Editeur_Id : %s, Nom: %s, CreatedAt: %s, UpdatedAt: %s\n", myEditeur.Editeur_Id, myEditeur.Nom, myEditeur.CreatedAt, myEditeur.UpdatedAt)

	renderTemplate(w, "UpdateEditeur", &data)
}

func UpdateEditeurPost(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/UpdateEditeurPost/")
	fmt.Printf("UpdateEditeurPost 723 avant conversion id : %s ", id)
	index, err := strconv.Atoi(id)
	config.CheckError(err)
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	config.PingDB(Db)
	myEditeur := config.Editeur{
		Editeur_Id: r.FormValue("Editeur_Id"),
		Nom:        r.FormValue("Nom"),
		CreatedAt:  r.FormValue("CreatedAt"),
		UpdatedAt:  r.FormValue("UpdatedAt"),
		DeletedAt:  r.FormValue("DeletedAt"),
	}
	myEditeur.Editeur_Id = id
	fmt.Printf("UpdateEditeurPost 737 Editeur_Id: %s, Nom: %v, CreatedAt: %s, UpdatedAt: %s\n", myEditeur.Editeur_Id, myEditeur.Nom, myEditeur.CreatedAt, myEditeur.UpdatedAt)
	var DatedeCreation, DatedeMaj []uint8
	myEditeur.UpdatedAt = string(DatedeMaj)

	DatedeCreation = []byte(myEditeur.CreatedAt)
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))
	id_Editeur, err := strconv.Atoi(myEditeur.Editeur_Id)
	config.CheckError(err)
	fmt.Printf("UpdateEditeurPost 745 Editeur_Id: %s, index: %d, Nom: %s, CreatedAt : %s, UpdatedAt: %s\n", myEditeur.Editeur_Id, id_Editeur, myEditeur.Nom, myEditeur.CreatedAt, myEditeur.UpdatedAt)
	// UPDATE INTO livres
	stmt, err := Db.Prepare("UPDATE Editeurs SET Nom=? , Editeur_Id=?, CreatedAt=?, UpdatedAt=?, DeletedAt=? where Editeur_Id =? ")
	config.CheckError(err)

	// execute
	res, err := stmt.Exec(myEditeur.Nom, DatedeCreation, DatedeMaj, "1900-01-01", id_Editeur)
	config.CheckError(err)

	a, err := res.RowsAffected()
	config.CheckError(err)

	log.Printf("updated id: %d, %v", index, a)

	// Fin de la mise à jour du 8/11/2023 à 12h30
	myEditeur.UpdatedAt = string(DatedeMaj)
	myEditeur.CreatedAt = string(DatedeCreation)
	myEditeur.Editeur_Id = id
	data := config.UnEditeur{
		TitrePage: "Mise à jour Editeur :" + myEditeur.Nom,
		Editeur:   myEditeur,
	}
	fmt.Printf("UpdateEditeurPost 767 Editeur_Id: %s, index: %d, Nom: %s, CreatedAt : %s, UpdatedAt: %s\n", myEditeur.Editeur_Id, id_Editeur, myEditeur.Nom, myEditeur.CreatedAt, myEditeur.UpdatedAt)

	renderTemplate(w, "UpdateEditeurPost", &data)
}

func ListeDesAuteurs(w http.ResponseWriter, r *http.Request) {
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)

	/*SELECT Livres.Livre_Id, Livres.Titre, Livres.Editeur_Id, Livres.Isbn, LivreMembership.Livre_Id, LivreMembership.Auteur_Id, Auteurs.Auteur_Id, Auteurs.Nom,Auteurs.Prenom
	FROM Livres
	JOIN  LivreMembership
	ON   Livres.Livre_Id = LivreMembership.Livre_Id
	 JOIN Auteurs
	ON LivreMembership.Auteur_Id = Auteurs.Auteur_Id WHERE Livres.Livre_Id= ?, Livre_Id*/
	rows, err := Db.Query("SELECT LivreMembership.Livre_Id, LivreMembership.Auteur_Id, Auteurs.Auteur_Id, Auteurs.Nom,Auteurs.Prenom FROM Auteurs JOIN LivreMembership ON LivreMembership.Auteur_Id = Auteurs.Auteur_Id where Auteurs.DeletedAt < ?", "1900-01-02")
	config.CheckError(err)
	defer rows.Close()
	var unAuteur config.Auteur2
	var myAuteurs []config.Auteur2
	i := 0
	for rows.Next() {
		err = rows.Scan(&unAuteur.LMS_Id, &unAuteur.AMS_Id, &unAuteur.AuteurId, &unAuteur.Nom, &unAuteur.Prenom)
		config.CheckError(err)
		myAuteurs = append(myAuteurs, unAuteur)
		fmt.Printf("ListeDesAuteurs 787 LMS_Id: %s, AMS_Id: %s, Auteur_Id : %s, Nom: %v, Prenom: %v\n", unAuteur.LMS_Id, unAuteur.AMS_Id, unAuteur.AuteurId, unAuteur.Nom, unAuteur.Prenom)
		i++
	}
	data := struct {
		TitrePage string
		Nombre    int
		Auteurs   []config.Auteur2
	}{
		TitrePage: "Liste des Auteurs inscrits dans la base mysql",
		Nombre:    i,
		Auteurs:   myAuteurs,
	}
	fmt.Printf("ListeDesAuteurs 801 data.TitrePage : %s, data.Nombre: %d\n", data.TitrePage, data.Nombre)

	renderTemplate(w, "ListeDesAuteurs", &data)
}
func AfficheAuteursduLivre(w http.ResponseWriter, r *http.Request) {
	// Ajouté le 26/11
	id := strings.TrimPrefix(r.URL.Path, "/AfficheAuteursduLivre/")
	index, err := strconv.Atoi(id)
	config.CheckError(err)
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)

	rows, err := Db.Query("SELECT Livres.Livre_Id, Livres.Titre, Livres.Isbn, LivreMembership.Livre_Id, LivreMembership.Auteur_Id, Auteurs.Auteur_Id, Auteurs.Nom,Auteurs.Prenom FROM Livres JOIN  LivreMembership ON Livres.Livre_Id = LivreMembership.Livre_Id JOIN Auteurs ON LivreMembership.Auteur_Id = Auteurs.Auteur_Id WHERE Livres.Livre_Id= ?", index)
	//var DatedeMaj []uint8 = []byte(time.Now().Format("2006-01-02"))
	config.CheckError(err)
	defer rows.Close()
	//Fin affichage 26/11
	defer rows.Close()
	var unlivre config.Livre_A

	i := 0
	for rows.Next() {
		unAuteur := config.Auteur2{}
		err = rows.Scan(&unlivre.Livre_Id, &unlivre.Titre, &unlivre.Isbn, &unAuteur.LMS_Id, &unAuteur.AMS_Id, &unAuteur.AuteurId, &unAuteur.Nom, &unAuteur.Prenom)
		config.CheckError(err)
		unlivre.Auteurs = append(unlivre.Auteurs, unAuteur)
		fmt.Printf("AfficheAuteursduLivre 822 Auteur_Id : %s, Nom: %v, Prenom: %v\n", unAuteur.AuteurId, unAuteur.Nom, unAuteur.Prenom)
		i++
	}
	data := struct {
		TitrePage string
		Nombre    int
		Livre     config.Livre_A
	}{
		TitrePage: "Liste des Auteurs du livre : n°" + unlivre.Livre_Id + " " + unlivre.Titre,
		Nombre:    i,
		Livre:     unlivre,
	}
	fmt.Printf("AfficheAuteursduLivre 834 data.TitrePage : %s, data.Nombre: %d\n", data.TitrePage, data.Nombre)

	renderTemplate(w, "ListeDesAuteursDuLivre", &data)
}
func AddAuteur(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "AddAuteur", nil)
}

func AddAuteurPost(w http.ResponseWriter, r *http.Request) {

	// Ajouté le 8/11/2023 à 12h30

	auteur := config.Auteur{

		Auteur_Id: r.FormValue("Auteur_Id"),
		Nom:       r.FormValue("Nom"),
		Prenom:    r.FormValue("Prenom"),
		CreatedAt: r.FormValue("CreatedAt"),
		UpdatedAt: r.FormValue("UpdatedAt"),
		DeletedAt: r.FormValue("DeletedAt"),
	}

	var DatedeCreation, DatedeMaj []uint8
	DatedeCreation = []byte(time.Now().Format("2006-01-02"))
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))

	//fmt.Printf("Avant création du livre = Id= %s, Titre: %s, Auteur_Id: %s, Editeur: %s, Isbn: %s, CreatedAt: %s, UpdatedAt: %s\n", livre.Id, livre.Titre, livre.Auteur_Id, livre.Editeur, livre.Isbn, livre.CreatedAt, livre.UpdatedAt)

	// INSERT INTO livres
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	fmt.Printf("AddAuteurPost 822 Nom : %s, Prenom : %s\n", auteur.Nom, auteur.Prenom)
	//id_editeur, err := strconv.Atoi(editeur.Editeur_Id)
	if err != nil {
		log.Fatalf("impossible de convertir auteur.Auteur_Id: %s", err)
	}
	//fmt.Printf("AddEditeurPost 606 id_editeur: %d\n", id_editeur)
	query := "INSERT INTO Auteurs (Nom, Prenom, CreatedAt, UpdatedAt, DeletedAt) VALUES (?, ?, ?, ?, ?)"
	insertResult, err := Db.ExecContext(context.Background(), query, auteur.Nom, auteur.Prenom, DatedeCreation, DatedeMaj, "1900-01-01")
	if err != nil {
		log.Fatalf("Impossible d'inserer l'Auteur: %s\n", err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}
	log.Printf("inserted id: %d", id)

	data := config.UnAuteur{
		TitrePage: "Création d'un Auteur :" + auteur.Nom,
		Auteur:    auteur,
	}
	// Fin de l'ajout du 8/11/2023 à 12h30

	renderTemplate(w, "AddAuteurPost", &data)
}
func UpdateAuteur(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/UpdateAuteur/")
	index, err := strconv.Atoi(id)
	config.CheckError(err)

	var DatedeMaj []uint8 = []byte(time.Now().Format("2006-01-02"))
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	myauteur := config.Auteur{}
	rows, err := Db.Query("SELECT Auteurs.Auteur_Id, Auteurs.Nom, Auteurs.Prenom, Auteurs.CreatedAt, Auteurs.UpdatedAt, Auteurs.DeletedAt FROM Auteurs WHERE Auteur_Id= ?", index)
	config.CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&myauteur.Auteur_Id, &myauteur.Nom, &myauteur.Prenom, &myauteur.CreatedAt, &DatedeMaj, &myauteur.DeletedAt) //DatedeCreation, &DatedeMaj, &index)
		config.CheckError(err)
		myauteur.Auteur_Id = id
		myauteur.UpdatedAt = string(DatedeMaj)
	}

	//fmt.Printf("CreatedAt: %#v \nUpdatedAt: %#v \n", DatedeCreation, DatedeMaj)
	//fmt.Printf("UpdateUser: Id : %s, Nom: %s, Prenom: %s, Courriel: %s, Pass: %s, CreatedDat: %s, UpdatedAt: %s \n", unuser.Id, unuser.Nom, unuser.Prenom, unuser.Email, unuser.Password, unuser.CreatedAt, DatedeMaj)
	data := config.UnAuteur{
		TitrePage: "Mise à jour de l'Auteur: " + myauteur.Nom,
		Titre:     myauteur.Nom,
		Auteur:    myauteur,
	}
	fmt.Printf("UpdateAuteur 874 Auteur_Id : %s, Nom: %s, Prenom: %v, CreatedAt: %s, UpdatedAt: %s\n", myauteur.Auteur_Id, myauteur.Nom, &myauteur.Prenom, myauteur.CreatedAt, myauteur.UpdatedAt)

	renderTemplate(w, "UpdateAuteur", &data)
}

func UpdateAuteurPost(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/UpdateAuteurPost/")
	fmt.Printf("UpdateAuteurPost 881 avant conversion id : %s ", id)
	index, err := strconv.Atoi(id)
	config.CheckError(err)
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/librairie")
	config.CheckError(err)
	config.PingDB(Db)
	myauteur := config.Auteur{
		Auteur_Id: r.FormValue("Auteur_Id"),
		Nom:       r.FormValue("Nom"),
		Prenom:    r.FormValue("Prenom"),
		CreatedAt: r.FormValue("CreatedAt"),
		UpdatedAt: r.FormValue("UpdatedAt"),
		DeletedAt: r.FormValue("DeletedAt"),
	}
	myauteur.Auteur_Id = id
	fmt.Printf("UpdateAuteurPost 895 Auteur_Id: %s, Nom: %v, Prenom: %v, CreatedAt: %s, UpdatedAt: %s\n", myauteur.Auteur_Id, myauteur.Nom, myauteur.Prenom, myauteur.CreatedAt, myauteur.UpdatedAt)
	var DatedeCreation, DatedeMaj []uint8
	myauteur.UpdatedAt = string(DatedeMaj)

	DatedeCreation = []byte(myauteur.CreatedAt)
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))
	id_Auteur, err := strconv.Atoi(myauteur.Auteur_Id)
	config.CheckError(err)
	fmt.Printf("UpdateAuteurPost 904 Auteur_Id: %s, index: %d, Nom: %s, Prenom: %s, CreatedAt : %s, UpdatedAt: %s\n", myauteur.Auteur_Id, id_Auteur, myauteur.Nom, myauteur.Prenom, myauteur.CreatedAt, myauteur.UpdatedAt)
	// UPDATE INTO livres
	stmt, err := Db.Prepare("UPDATE Auteurs SET Nom=? , Prenom=?, CreatedAt=?, UpdatedAt=?, DeletedAt=? where Auteur_Id =? ")
	config.CheckError(err)

	// execute
	res, err := stmt.Exec(myauteur.Nom, myauteur.Prenom, DatedeCreation, DatedeMaj, "1900-01-01", id_Auteur)
	config.CheckError(err)

	a, err := res.RowsAffected()
	config.CheckError(err)

	log.Printf("updated id: %d, %v", index, a)

	// Fin de la mise à jour du 8/11/2023 à 12h30
	myauteur.UpdatedAt = string(DatedeMaj)
	myauteur.CreatedAt = string(DatedeCreation)
	myauteur.Auteur_Id = id
	data := config.UnAuteur{
		TitrePage: "Mise à jour Editeur :" + myauteur.Nom,
		Auteur:    myauteur,
	}
	fmt.Printf("UpdateAuteurPost 926 Auteur_Id: %s, index: %d, Nom: %s, Prenom: %s, CreatedAt : %s, UpdatedAt: %s\n", myauteur.Auteur_Id, id_Auteur, myauteur.Nom, myauteur.Prenom, myauteur.CreatedAt, myauteur.UpdatedAt)

	renderTemplate(w, "UpdateAuteurPost", &data)
}
