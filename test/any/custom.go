package any

import "net/http"

func (HttpHandler[Z, K, V, M]) CustomMethod(ctx Z, value V, req http.Request) http.Response {
	return http.Response{}
}
