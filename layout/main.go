package layout

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"
)

func add(a, b int) int { return a + b }
func sub(a, b int) int { return a - b }
func formatTimestamp(t time.Time) string {
	return t.Format("Jan 2, 2006 at 3:04pm")
}

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"safeHTML": func(s string) template.HTML {
		return template.HTML(s)
	},
	"add":             add,
	"sub":             sub,
	"formatTimestamp": formatTimestamp,
}).ParseFiles(
	"templates/layout.html",
	"templates/home.html",
	"templates/profile.html",
	"templates/submit.html",
	"templates/login.html",
	"templates/register.html",
	"templates/post.html",
))

var title string

func Init(t string) {
	title = t
}

func RenderPage(w http.ResponseWriter, page string, username string, content any) {
	var contentBuffer bytes.Buffer
	err := templates.ExecuteTemplate(&contentBuffer, page, content)
	if err != nil {
		http.Error(w, "Failed to render content template", http.StatusInternalServerError)
		log.Println("Error rendering content template:", err)
		return
	}
	// Render "layout.html" with the rendered content
	err = templates.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Content":  contentBuffer.String(),
		"Username": username,
		"Title":    title,
	})
	if err != nil {
		http.Error(w, "Failed to render layout template", http.StatusInternalServerError)
		log.Println("Error rendering layout template:", err)
	}
}
