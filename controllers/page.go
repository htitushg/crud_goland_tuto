package controllers

import (
	"bytes"
	"context"
	"crud_golang_tuto/config"
	"database/sql"
	"log"
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
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./views/*.html")
	if err != nil {
		return cache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		tmpl := template.Must(template.ParseFiles(page))
		layouts, err := filepath.Glob("./views/layouts/*.layout.html")
		if err != nil {
			return cache, err
		}
		if len(layouts) > 0 {
			tmpl.ParseGlob("./views/layouts/*.layout.html")
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
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/base1")
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
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/base1")
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
		Name: "Mise à jour Utilisateur :" + user.Nom,
		User: user,
	}
	// Fin de l'ajout du 8/11/2023 à 12h30

	renderTemplate(w, "AddUserPost", &data)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/UpdateUser/")
	index, err := strconv.Atoi(id)
	config.CheckError(err)

	var DatedeCreation, DatedeMaj []uint8
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/base1")
	config.CheckError(err)
	unuser := config.User{}
	rows, err := Db.Query("SELECT * FROM users WHERE Id= ?", index)
	config.CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&unuser.Nom, &unuser.Prenom, &unuser.Email, &unuser.Password, &index, &DatedeCreation, &DatedeMaj, &unuser.DeletedAt) //DatedeCreation, &DatedeMaj, &index)
		config.CheckError(err)
		unuser.Email = strings.Join(strings.Fields(unuser.Email), "")
		unuser.Id = id
		unuser.UpdatedAt = string(DatedeMaj)
		unuser.CreatedAt = string(DatedeCreation)
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

	// UPDATE INTO users
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/base1")
	config.CheckError(err)
	config.PingDB(Db)

	//Update db
	stmt, err := Db.Prepare("UPDATE users SET nom=? , prenom=?, email=?, password=?, CreatedAt=?, UpdatedAt=?, DeletedAt=? where id =? ")
	config.CheckError(err)

	// execute
	res, err := stmt.Exec(user.Nom, user.Prenom, user.Email, user.Password, DatedeCreation, DatedeMaj, user.DeletedAt, index)
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
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/base1")
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
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/base1")
	config.CheckError(err)
	config.PingDB(Db)
	//Arr := config.User{}
	//rows, err := Db.Prepare("SELECT * FROM users WHERE Id= ?", index)
	//config.CheckError(err)
	//defer rows.Close()
	//res, err = rows.Exec(&Arr.Nom, &Arr.Prenom, &Arr.Email, &Arr.Password, &index, &Arr.CreatedAt, &Arr.UpdatedAt)
	//config.CheckError(err)
	//fmt.Printf("DeleteUserPost avant effacement = Id= %s, Nom: %s, Prenom: %s, Courriel: %s, Pass: %s, CreatedDat: %s, UpdatedAt: %s\n", id, Arr.Nom, Arr.Prenom, Arr.Email, Arr.Password, Arr.CreatedAt, Arr.UpdatedAt)

	//Update db
	//stmt, err := Db.Prepare("DELETE FROM users WHERE id =?")
	stmt, err := Db.Prepare("UPDATE users SET nom=? , prenom=?, email=?, password=?, CreatedAt=?, UpdatedAt=?, DeletedAt=? where id =? ")
	config.CheckError(err)

	// execute
	//res, err := stmt.Exec(index)
	//config.CheckError(err)
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
