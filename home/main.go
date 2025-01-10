package home

import (
	"lfx/auth"
	"lfx/db"

	"lfx/layout"
	"lfx/utils"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	database := db.GetDB()

	rows, err := database.Query("SELECT id, title, url, score, timestamp, username FROM posts ORDER BY (score - strftime('%s', timestamp) / 10000) DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	username, _ := auth.GetUserFromToken(r)

	var posts []utils.Post
	for rows.Next() {
		var post utils.Post
		err := rows.Scan(&post.ID, &post.Title, &post.URL, &post.Score, &post.Timestamp, &post.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	layout.RenderPage(w, "home.html", username, posts)
}
