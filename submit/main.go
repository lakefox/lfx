package submit

import (
	"lfx/auth"
	"lfx/db"
	"lfx/ipban"
	"lfx/layout"
	"lfx/spam"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	database := db.GetDB()

	username, err := auth.GetUserFromToken(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		url := r.FormValue("url")
		text := r.FormValue("text")

		if spam.ContainsBannedWord(title) || spam.ContainsBannedWord(text) {
			ipban.Ban(ipban.GetIP(r))
			http.Redirect(w, r, "/contentpolicy", http.StatusFound)
			return
		}

		_, err := database.Exec("INSERT INTO posts (title, url, text, username) VALUES (?, ?, ?, ?)", title, url, text, username)
		if err != nil {
			http.Error(w, "Failed to submit post", http.StatusInternalServerError)
			log.Println("Error inserting post:", err)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	layout.RenderPage(w, "submit.html", username, map[string]interface{}{
		"Username": username,
	})
}
