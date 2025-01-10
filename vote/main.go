package vote

import (
	"lfx/auth"
	"lfx/db"
	"log"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetUserFromToken(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Use the database instance from the database package
	database := db.GetDB()
	_, err = database.Exec("UPDATE posts SET score = score + 1 WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to update vote", http.StatusInternalServerError)
		log.Println("Error updating vote:", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
