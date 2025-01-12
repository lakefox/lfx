package profile

import (
	"lfx/auth"
	"lfx/db"
	"lfx/layout"
	"lfx/utils"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	database := db.GetDB()
	username := r.URL.Query().Get("user")
	if username == "" {
		http.NotFound(w, r)
		return
	}

	// Fetch posts by the user
	postRows, err := database.Query("SELECT id, title, url, score, timestamp FROM posts WHERE username = ? ORDER BY timestamp DESC", username)
	if err != nil {
		http.Error(w, "Failed to fetch user posts", http.StatusInternalServerError)
		log.Println("Error fetching user posts:", err)
		return
	}
	defer postRows.Close()

	var posts []utils.Post
	for postRows.Next() {
		var post utils.Post
		err := postRows.Scan(&post.ID, &post.Title, &post.URL, &post.Score, &post.Timestamp)
		if err != nil {
			http.Error(w, "Failed to scan database results", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return
		}
		posts = append(posts, post)
	}

	// Fetch comments by the user
	commentRows, err := database.Query("SELECT post_id, content, timestamp FROM comments WHERE username = ? ORDER BY timestamp DESC", username)
	if err != nil {
		http.Error(w, "Failed to fetch user comments", http.StatusInternalServerError)
		log.Println("Error fetching user comments:", err)
		return
	}
	defer commentRows.Close()

	var comments []utils.Comment
	for commentRows.Next() {
		var comment utils.Comment
		err := commentRows.Scan(&comment.PostID, &comment.Content, &comment.Timestamp)
		if err != nil {
			http.Error(w, "Failed to scan comments", http.StatusInternalServerError)
			log.Println("Error scanning comment row:", err)
			return
		}
		comments = append(comments, comment)
	}

	user, _ := auth.GetUserFromToken(r)

	layout.RenderPage(w, "profile.html", user, map[string]interface{}{
		"Username": username,
		"Posts":    posts,
		"Comments": comments,
	})
}
