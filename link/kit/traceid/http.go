package traceid

import "net/http"

// Middleware adds a new trace ID to [http.Request.Context]
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := FromContext(r.Context()); !ok {
			r = r.WithContext(WithContext(r.Context(), New()))
		}
		next.ServeHTTP(w, r)
	})
}
