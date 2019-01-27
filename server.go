package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"

	md "github.com/smertelny/ibanCheckerService/middlewares"
	"github.com/smertelny/ibanCheckerService/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/smertelny/ibanCheckerService/iban"

	"github.com/go-chi/chi"
)

// SuccessfulResponse is used in case it is all OK
type SuccessfulResponse struct {
	XMLName xml.Name `json:"-" xml:"response"`
	Result  string   `json:"result" xml:"result"`
}

func (r SuccessfulResponse) String() string {
	return r.Result
}

// ErrorContent is the content of error json object
type ErrorContent struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Msg     string   `json:"message" xml:"message"`
	Code    int      `json:"status_code"  xml:"status_code"`
}

// ErrorResponse is made for better json output
type ErrorResponse struct {
	XMLName xml.Name     `json:"-" xml:"errors"`
	Error   ErrorContent `json:"error" xml:"error"`
}

func (r ErrorResponse) String() string {
	return r.Error.Msg
}

func checkIban(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, HEAD")
	code := chi.URLParam(r, "code")

	err := iban.Check(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Render(w, r, ErrorResponse{Error: ErrorContent{Msg: err.Error(), Code: http.StatusBadRequest}})
	} else {
		w.WriteHeader(http.StatusOK)
		utils.Render(w, r, SuccessfulResponse{Result: fmt.Sprintf("IBAN %v is valid", code)})
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.Header().Set("Allow", "GET, HEAD")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello! Follow to <a href=\"/v1/iban/examplecode\">\"%v/v1/iban/YOUR_CODE_HERE\"</a> to check your iban", r.Host)
}

func custom404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	utils.Render(w, r, ErrorResponse{Error: ErrorContent{Msg: "Not found", Code: http.StatusNotFound}})
}

func main() {
	router := chi.NewRouter()
	router.NotFound(custom404)
	// Middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.GetHead)
	router.Use(middleware.Timeout(time.Second * 30))
	router.Use(middleware.RedirectSlashes)

	router.Route("/v1", func(router chi.Router) {
		router.Use(md.FormatMiddleware)
		router.Get("/", index)
		router.Get("/iban/{code}", checkIban)
	})

	log.Print("Server started")
	log.Fatal(http.ListenAndServe(":8000", router))
}
