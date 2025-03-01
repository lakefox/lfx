package main

import (
	"embed"
	"io/fs"
	"lfx/auth"
	"lfx/contentpolicy"
	"lfx/db"
	"lfx/home"
	"lfx/ipban"
	"lfx/layout"
	"lfx/post"
	"lfx/profile"
	"lfx/spam"
	"lfx/submit"
	"lfx/utils"
	"lfx/vote"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//go:embed templates/*
var templateFiles embed.FS

//go:embed static/*
var staticFiles embed.FS

func checkBan(handler http.HandlerFunc, middleware func(http.Handler) http.Handler) http.Handler {
	return middleware(handler)
}

func main() {
	err := utils.ENV(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize BanManager with a timeout from the environment variable
	bantimeout, _ := strconv.Atoi(os.Getenv("BAN_TIMEOUT"))
	ipban.NewBanManager(time.Duration(bantimeout) * time.Minute)

	// Initialize other components
	db.Init(os.Getenv("DATABASE"))
	defer db.GetDB().Close()

	auth.Init(os.Getenv("JWT_SECRET"))
	layout.Init(os.Getenv("SITE_TITLE"), os.Getenv("THEME"), templateFiles)

	postsPerPage, _ := strconv.Atoi(os.Getenv("POSTS_PER_PAGE"))
	home.Init(postsPerPage)
	spam.Init()

	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Register routes with the ipban middleware
	http.Handle("/", checkBan(home.Handler, ipban.Middleware))
	http.Handle("/submit", checkBan(submit.Handler, ipban.Middleware))
	http.Handle("/profile", checkBan(profile.Handler, ipban.Middleware))
	http.Handle("/vote", checkBan(vote.Handler, ipban.Middleware))
	http.Handle("/register", checkBan(auth.Register, ipban.Middleware))
	http.Handle("/login", checkBan(auth.Login, ipban.Middleware))
	http.Handle("/logout", checkBan(auth.Logout, ipban.Middleware))
	http.Handle("/post", checkBan(post.Handler, ipban.Middleware))
	http.HandleFunc("/contentpolicy", contentpolicy.Handler)

	// Start the server
	log.Println("Server started at :80")
	http.ListenAndServe(":80", nil)
}
