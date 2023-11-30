package routes

import (
	"crud_golang_tuto/controllers"
	"net/http"
)

func Web() {
	// Définir les routes qui affichent les pages web
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/AddUser", controllers.AddUser)
	http.HandleFunc("/AddUserPost", controllers.AddUserPost)
	http.HandleFunc("/UpdateUser/", controllers.UpdateUser)
	http.HandleFunc("/UpdateUserPost/", controllers.UpdateUserPost)
	http.HandleFunc("/DeleteUser/", controllers.DeleteUser)
	http.HandleFunc("/DeleteUserPost/", controllers.DeleteUserPost)
	http.HandleFunc("/contact", controllers.Contact)
	http.HandleFunc("/ListeDesLivres", controllers.ListeDesLivres)
	http.HandleFunc("/AfficheLivresEtAuteurs", controllers.AfficheLivresEtAuteurs)
	http.HandleFunc("/AddLivre", controllers.AddLivre)
	http.HandleFunc("/AddLivrePost", controllers.AddLivrePost)
	http.HandleFunc("/UpdateLivre/", controllers.UpdateLivre)
	http.HandleFunc("/UpdateLivrePost/", controllers.UpdateLivrePost)
	http.HandleFunc("/DeleteLivre/", controllers.DeleteLivre)
	http.HandleFunc("/DeleteLivrePost/", controllers.DeleteLivrePost)
	http.HandleFunc("/AddEditeur", controllers.AddEditeur)
	http.HandleFunc("/AddEditeurPost", controllers.AddEditeurPost)
	http.HandleFunc("/ListeDesEditeurs", controllers.ListeDesEditeurs)
	http.HandleFunc("/UpdateEditeur/", controllers.UpdateEditeur)
	http.HandleFunc("/UpdateEditeurPost/", controllers.UpdateEditeurPost)
	http.HandleFunc("/ListeDesAuteurs", controllers.ListeDesAuteurs)
	http.HandleFunc("/AddAuteur/", controllers.AddAuteur)
	http.HandleFunc("/AddAuteurPost/", controllers.AddAuteurPost)
	http.HandleFunc("/UpdateAuteur/", controllers.UpdateAuteur)
	http.HandleFunc("/UpdateAuteurPost/", controllers.UpdateAuteurPost)
	http.HandleFunc("/AfficheAuteursduLivre/", controllers.AfficheAuteursduLivre)

}
