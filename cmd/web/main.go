package main

import (
	"crud_golang_tuto/config"
	"crud_golang_tuto/controllers"
	"crud_golang_tuto/routes"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	Port = ":8090"
)

func main() {
	// On relie le fichier css
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	templateCache, err := controllers.CreateTemplateCache()
	if err != nil {
		panic(err)
	}
	var appConfig config.Config
	appConfig.TemplateCache = templateCache
	appConfig.Port = ":8070"
	controllers.CreateTemplates(&appConfig)

	// Définir les urls de notre application web
	routes.Web()
	fmt.Printf("http://localhost%v , Cliquez sur le lien pour lancer le navigateur", Port)
	http.ListenAndServe(Port, nil)
}
