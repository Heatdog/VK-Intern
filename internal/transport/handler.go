package transport

import "net/http"

type Handler interface {
	Register(router *http.ServeMux)
}
