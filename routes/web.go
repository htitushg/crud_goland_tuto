package routes

import (
	"crud_golang_tuto/controllers"
	"net/http"
)

func Web() {
	// Définir les routes qui affichent les pages web
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/AddUser", controllers.AddUser)                // Formulaire AddUser
	http.HandleFunc("/AddUserPost", controllers.AddUserPost)        // Traitement du formulaire AddUser
	http.HandleFunc("/UpdateUser/", controllers.UpdateUser)         //Editer UpdateUser
	http.HandleFunc("/UpdateUserPost/", controllers.UpdateUserPost) // Traitement du formulaire UpdateUser
	http.HandleFunc("/DeleteUser/", controllers.DeleteUser)         // Effacer DeleteUser
	http.HandleFunc("/DeleteUserPost/", controllers.DeleteUserPost) // Traitement du formulaire Delete
}
