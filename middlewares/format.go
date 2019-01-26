package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
)

var formats = map[string]string{
	"html":             "text/html",
	"text/html":        "text/html",
	"json":             "application/json",
	"application/json": "application/json",
	"xml":              "application/xml",
	"application/xml":  "application/xml",
}

type ctxKey int

const responseFormat ctxKey = 0

//ResponseFormatFromCtx is function for getting response format (json/xml/html) from context
func ResponseFormatFromCtx(ctx context.Context) (string, bool) {
	format, ok := ctx.Value(responseFormat).(string)
	return format, ok
}

//SetResponseFormatToCtx is function for setting response format to context
func SetResponseFormatToCtx(ctx context.Context, format string) context.Context {
	return context.WithValue(ctx, responseFormat, format)
}

//FormatMiddleware is middleware for getting format for response from URL, Accept header or Query param "format"
func FormatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		format, _ := r.Context().Value(middleware.URLFormatCtxKey).(string)
		if format == "" {
			format = r.URL.Query().Get("format")
		}
		if format == "" {
			format = strings.Split(r.Header.Get("Accept"), ",")[0]
		}

		format, ok := formats[format]
		if !ok {
			format = "text/html"
		}

		ctx := SetResponseFormatToCtx(r.Context(), format)
		w.Header().Set("content-type", format+"; charset=utf-8")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
