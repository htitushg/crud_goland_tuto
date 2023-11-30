package main

import (
	"crud_golang_tuto/config"
	"crud_golang_tuto/controllers"
	"crud_golang_tuto/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	Port = ":8090"
)

func main() {
	path, err2 := os.Getwd()
	if err2 != nil {
		log.Fatal(err2)
	}
	// On relie le fichier css
	fs := http.FileServer(http.Dir(path + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	templateCache, err := controllers.CreateTemplateCache()
	if err != nil {
		panic(err)
	}
	var appConfig config.Config
	appConfig.TemplateCache = templateCache
	appConfig.Port = Port
	controllers.CreateTemplates(&appConfig)

	// Définir les urls de notre application web
	routes.Web()
	fmt.Printf("http://localhost%v , Cliquez sur le lien pour lancer le navigateur", appConfig.Port)
	http.ListenAndServe(Port, nil)
}
