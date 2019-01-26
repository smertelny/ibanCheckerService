package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	mw "github.com/smertelny/ibanCheckerService/middlewares"
)

//Render is a function for rendering response from format context variable
func Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	format, ok := mw.ResponseFormatFromCtx(r.Context())
	if !ok {
		format = "html"
	}

	switch format {
	case "application/json":
		json.NewEncoder(w).Encode(data)

	case "application/xml":
		xml.NewEncoder(w).Encode(data)

	case "text/html":
		fmt.Fprint(w, data)
	}
}
