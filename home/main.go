package home

import (
	"lfx/auth"
	"lfx/db"
	"lfx/layout"
	"lfx/utils"
	"net/http"
	"strconv"
)

var postsPerPage int

func Init(posts int) {
	postsPerPage = posts
}

func Handler(w http.ResponseWriter, r *http.Request) {
	database := db.GetDB()

	// Get the "page" query parameter and calculate the offset
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if the parameter is invalid or missing
	}
	offset := (page - 1) * postsPerPage

	// Fetch paginated posts
	query := `
		SELECT id, title, url, score, timestamp, username
		FROM posts
		ORDER BY (score - strftime('%s', timestamp) / 10000) DESC
		LIMIT ? OFFSET ?`
	rows, err := database.Query(query, postsPerPage, offset)
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

	// Determine if there are more pages
	var totalPosts int
	err = database.QueryRow("SELECT COUNT(*) FROM posts").Scan(&totalPosts)
	if err != nil {
		http.Error(w, "Failed to count posts", http.StatusInternalServerError)
		return
	}
	hasNextPage := page*postsPerPage < totalPosts

	// Add pagination data
	data := map[string]interface{}{
		"Posts":       posts,
		"CurrentPage": page,
		"HasNextPage": hasNextPage,
		"HasPrevPage": page > 1,
	}

	layout.RenderPage(w, "home.html", username, data)
}
