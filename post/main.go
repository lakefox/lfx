package post

import (
	"lfx/db"
	"lfx/ipban"
	"lfx/spam"

	"lfx/auth"
	"lfx/layout"
	"lfx/utils"
	"log"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	database := db.GetDB()

	// Parse the post ID
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Fetch the post details
	var post utils.Post
	err = database.QueryRow("SELECT id, title, url, text, score, timestamp, username FROM posts WHERE id = ?", postID).
		Scan(&post.ID, &post.Title, &post.URL, &post.Text, &post.Score, &post.Timestamp, &post.Username)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Fetch comments for the post
	rows, err := database.Query("SELECT username, content, timestamp FROM comments WHERE post_id = ? ORDER BY timestamp ASC", postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		log.Println("Error fetching comments:", err)
		return
	}
	defer rows.Close()

	var comments []utils.Comment
	for rows.Next() {
		var comment utils.Comment
		err := rows.Scan(&comment.Username, &comment.Content, &comment.Timestamp)
		if err != nil {
			http.Error(w, "Failed to scan comments", http.StatusInternalServerError)
			log.Println("Error scanning comment row:", err)
			return
		}
		comments = append(comments, comment)
	}

	// Handle adding a new comment
	if r.Method == http.MethodPost {
		username, err := auth.GetUserFromToken(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		content := r.FormValue("content")

		if spam.ContainsBannedWord(content) {
			ipban.Ban(ipban.GetIP(r))
			http.Redirect(w, r, "/contentpolicy", http.StatusFound)
			return
		}

		_, err = database.Exec("INSERT INTO comments (post_id, username, content) VALUES (?, ?, ?)", postID, username, content)
		if err != nil {
			http.Error(w, "Failed to add comment", http.StatusInternalServerError)
			log.Println("Error inserting comment:", err)
			return
		}

		http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusFound)
		return
	}
	username, _ := auth.GetUserFromToken(r)

	// Render the post page
	layout.RenderPage(w, "post.html", username, map[string]interface{}{
		"Post":     post,
		"Comments": comments,
	})
}
