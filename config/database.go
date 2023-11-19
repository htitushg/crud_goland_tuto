package config

import (
	"context"
	"database/sql"
	"log"
	"time"
)

var (
	ctx context.Context
	Db  *sql.DB
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
func PingDB(db *sql.DB) {
	err := db.Ping()
	CheckError(err)
}
func ConnectBase() bool {
	Db, err := sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/base1")
	CheckError(err)
	defer Db.Close()
	err = Db.Ping()
	CheckError(err)
	return true
}
func CreateTable(nomtable string) {
	_, err := Db.Exec("CREATE TABLE IF NOT EXISTS ? (id serial,manufacturer varchar(20), design varchar(20), style varchar(20), doors int, created_at timestamp default NULL, updated_at timestamp default NULL, constraint pk primary key(id))", nomtable)
	CheckError(err)
}

// Getter for db var
//
//	func Db() *sql.DB {
//		return Db
//	}
func NewUser(user *User) {
	if user == nil {
		log.Fatal(user)
	}
	currentTime := time.Now()
	DatedeCreation := currentTime.Format("2006-01-02")
	DatedeMaj := currentTime.Format("2006-01-02")

	query := "INSERT INTO users (Nom, Prenom, Email,Password) VALUES (?, ?, ?, ?)"
	insertResult, err := Db.ExecContext(ctx, query, user.Nom, user.Prenom, user.Email, user.Password, DatedeCreation, DatedeMaj)
	if err != nil {
		log.Fatalf("impossible insert user: %s", err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}
	log.Printf("inserted id: %d", id)

}

//func UpdateCar(u *User) {
//	u.UpdatedAt = time.Now()
//	stmt, err := Db.Prepare("UPDATE users SET nom=$1 , prenom = $2, email=$3, password=$4 ,CreatedAt=$6,UpdatedAt=$7, where id =$5 ")
//	//stmt, err := Db.Prepare("UPDATE users SET manufacturer=$1, design=$2, style=$3, doors=$4, updated_at=$5 WHERE id=$6;")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	_, err = stmt.Exec(u.Nom, u.Prenom, u.Email, u.Password, index, DatedeCreation, DatedeMaj)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//}
