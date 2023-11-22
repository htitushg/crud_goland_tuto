package routes

import (
	"crud_golang_tuto/controllers"
	"net/http"
)

func Web() {
	// Définir les routes qui affichent les pages web
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/AddUser", controllers.AddUser)                     // Formulaire AddUser
	http.HandleFunc("/AddUserPost", controllers.AddUserPost)             // Traitement du formulaire AddUser
	http.HandleFunc("/UpdateUser/", controllers.UpdateUser)              //Editer UpdateUser
	http.HandleFunc("/UpdateUserPost/", controllers.UpdateUserPost)      // Traitement du formulaire UpdateUser
	http.HandleFunc("/DeleteUser/", controllers.DeleteUser)              // Effacer DeleteUser
	http.HandleFunc("/DeleteUserPost/", controllers.DeleteUserPost)      // Traitement du formulaire Delete
	http.HandleFunc("/contact", controllers.Contact)                     // Formulaire de contact
	http.HandleFunc("/ListeDesLivres", controllers.ListeDesLivres)       // Formulaire ListeLivres
	http.HandleFunc("/AddLivre", controllers.AddLivre)                   // Formulaire AddLivre
	http.HandleFunc("/AddLivrePost", controllers.AddLivrePost)           // Traitement du formulaire AddLivre
	http.HandleFunc("/UpdateLivre/", controllers.UpdateLivre)            // Formulaire UpdateLivre
	http.HandleFunc("/UpdateLivrePost/", controllers.UpdateLivrePost)    // Traitement du formulaire UpdateLivre
	http.HandleFunc("/DeleteLivre/", controllers.DeleteLivre)            // Formulaire UpdateLivre
	http.HandleFunc("/DeleteLivrePost/", controllers.DeleteLivrePost)    // Traitement du formulaire UpdateLivre
	http.HandleFunc("/AddEditeur", controllers.AddEditeur)               // Formulaire AddUser
	http.HandleFunc("/AddEditeurPost", controllers.AddEditeurPost)       // Traitement du formulaire AddUser
	http.HandleFunc("/ListeDesEditeurs", controllers.ListeDesEditeurs)   // Traitement du formulaire AddUser
	http.HandleFunc("/UpdateEditeur", controllers.UpdateEditeur)         // Formulaire AddUser
	http.HandleFunc("/UpdateEditeurPost", controllers.UpdateEditeurPost) // Traitement du formulaire AddUser

}
