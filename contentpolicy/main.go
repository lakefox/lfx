package contentpolicy

import (
	"lfx/layout"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	layout.RenderPage(w, "contentpolicy.html", "", nil)
}
