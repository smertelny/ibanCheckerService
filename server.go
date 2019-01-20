package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/smertelny/ibanCheckerService/iban"

	"github.com/go-chi/chi"
)

// SuccessfulResponse is used in case it is all OK
type SuccessfulResponse struct {
	Result string `json:"result" xml:"result"`
}

// ErrorContent is the content of error json object
type ErrorContent struct {
	Msg  string `json:"message" xml:"message"`
	Code int    `json:"status_code"  xml:"status_code"`
}

// ErrorResponse is made for better json output
type ErrorResponse struct {
	Error ErrorContent `json:"error" xml:"error"`
}

var formats = []string{
	"html",
	"text/html",
	"json",
	"application/json",
	"xml",
	"application/xml",
}

func checkIban(w http.ResponseWriter, r *http.Request) {
	format, _ := r.Context().Value(middleware.URLFormatCtxKey).(string)
	if format == "" {
		format = r.URL.Query().Get("format")
	}
	if format == "" {
		format = strings.Split(r.Header.Get("Accept"), ",")[0]
	}
	log.Print(format)

	var flag bool
	for _, v := range formats {
		if format == v {
			flag = true
			break
		}
	}
	if !flag {
		format = "json"
	}

	switch format {
	case "json", "application/json":
		w.Header().Set("Content-type", "application/json")
	case "html", "text/html":
		w.Header().Set("Content-type", "text/html")
	case "xml", "application/xml":
		w.Header().Set("Content-type", "application/xml")
	}
	code := chi.URLParam(r, "code")

	err := iban.Check(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		switch format {
		case "json", "application/json":
			json.NewEncoder(w).Encode(ErrorResponse{ErrorContent{err.Error(), http.StatusBadRequest}})
		case "xml", "application/xml":
			xml.NewEncoder(w).Encode(ErrorResponse{ErrorContent{err.Error(), http.StatusBadRequest}})
		case "html", "text/html":
			fmt.Fprintf(w, err.Error())
		}
	} else {
		w.WriteHeader(http.StatusOK)
		switch format {
		case "json", "application/json":
			json.NewEncoder(w).Encode(SuccessfulResponse{"IBAN code is valid"})
		case "xml", "application/xml":
			xml.NewEncoder(w).Encode(SuccessfulResponse{"IBAN code is valid"})
		case "html", "text/html":
			fmt.Fprintf(w, "IBAN code is valid")
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello! Follow to <a href=\"/iban/examplecode/\">\"%v/iban/YOUR_CODE_HERE\"</a> to check your iban", r.Host)
}

func main() {
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(time.Second * 30))
	// router.Use(middleware.AllowContentType()

	router.Get("/", index)
	router.Get("/iban/{code}", checkIban)

	log.Print("Server started")
	log.Fatal(http.ListenAndServe(":8000", router))
}
