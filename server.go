package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/smertelny/ibanCheckerService/iban"
)

type SuccessfulResponse struct {
	Result string `json:"result"`
}

type ErrorContent struct {
	Msg  string `json:"message"`
	Code int    `json:"status_code"`
}

type ErrorResponse struct {
	Error ErrorContent `json:"error"`
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientAddres string
		if r.Header.Get("X-Forwarded-For") != "" {
			clientAddres = r.Header.Get("X-Forwarded-For")
		} else {
			clientAddres = r.RemoteAddr
		}
		log.Printf("Sender: %v; Recieved: %v;", clientAddres, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-type", "text/html")
	fmt.Fprintf(w, "Oops... We can't find such page. Sorry =(")
	return
}

func checkIban(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(urlParts) > 2 {
		pageNotFound(w, r)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Allow", "GET")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{ErrorContent{"Method not allowed", http.StatusMethodNotAllowed}})
		return
	}

	if code := urlParts[len(urlParts)-1]; code != "iban" {
		err := iban.Check(code)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{ErrorContent{err.Error(), http.StatusBadRequest}})
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SuccessfulResponse{"IBAN code is valid"})
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	if r.URL.Path != "/" {
		pageNotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Hello! Follow to <a href=\"/iban/examplecode/\">\"%v/iban/YOUR_CODE_HERE\"</a> to check your iban", r.Host)
}

func main() {
	mux := http.NewServeMux()
	handlerChecker := http.HandlerFunc(checkIban)
	mux.Handle("/", logMiddleware(http.HandlerFunc(index)))
	mux.Handle("/iban/", logMiddleware(handlerChecker))
	log.Print("Server started")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
