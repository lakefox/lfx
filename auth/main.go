package auth

import (
	"lfx/db"
	"lfx/layout"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

// Init sets the JWT secret
func Init(secret string) {
	if secret == "" {
		log.Fatal("JWT_SECRET is not set in the environment")
	}
	jwtSecret = []byte(secret)
}

func GetUserFromToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}

	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", jwt.ErrSignatureInvalid
		}
		return username, nil
	}

	return "", jwt.ErrSignatureInvalid
}

func Login(w http.ResponseWriter, r *http.Request) {
	database := db.GetDB()

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Fetch the user from the database
		var hashedPassword string
		err := database.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Check the password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			log.Println("Error generating token:", err)
			return
		}

		// Set token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	layout.RenderPage(w, "login.html", "", nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func Register(w http.ResponseWriter, r *http.Request) {
	database := db.GetDB()

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			log.Println("Error hashing password:", err)
			return
		}

		// Insert user into the database
		_, err = database.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(w, "Failed to create account", http.StatusInternalServerError)
			log.Println("Error inserting user:", err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	layout.RenderPage(w, "register.html", "", nil)
}
