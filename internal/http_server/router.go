package http_server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type stackHandler interface {
	PushToStackHandler(w http.ResponseWriter, r *http.Request)
	PopFromStackHandler(w http.ResponseWriter, r *http.Request)
}

func NewStackServiceRouter(stack stackHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(
		"/",
		stack.PushToStackHandler,
	).Methods("POST")
	router.HandleFunc(
		"/",
		stack.PopFromStackHandler,
	).Methods("GET")
	return router
}
