package handlers

import (
	"StackService/internal/http_server"
	"StackService/internal/services"
	"io"
	"log/slog"
	"net/http"
)

type StackHandler struct {
	Stack  *services.Stack
	Logger *slog.Logger
}

func (s *StackHandler) PushToStackHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http_server.SendJSONError(w, http.StatusBadRequest, err)
		return
	}
	err = s.Stack.Push(body)
	if err != nil {
		http_server.SendJSONError(w, http.StatusBadRequest, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *StackHandler) PopFromStackHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	json, err := s.Stack.Pop()
	if err != nil {
		http_server.SendJSONError(w, http.StatusNotFound, err)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(json))
	if err != nil {
		http_server.SendJSONError(w, http.StatusNotFound, err)
	}
	return
}
