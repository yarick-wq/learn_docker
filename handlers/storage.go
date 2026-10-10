package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	storage EmployeeStorage
}

func NewHandler(stor EmployeeStorage) *Handler {
	return &Handler{
		storage: stor,
	}
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var emp CreateEmployeeRequest
	if err := decoder.Decode(&emp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON body"))
		return
	}
	if err := emp.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	employee := Employee{
		FullName: strings.TrimSpace(emp.FullName),
		Position: strings.TrimSpace(emp.Position),
	}
	created, err := h.storage.Create(r.Context(), employee)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Problem in server"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(created); err != nil {
		fmt.Println("err", err)
	}
}

func (h *Handler) CollectEmployee(w http.ResponseWriter, r *http.Request) {
	list, err := h.storage.Collect(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Problem in server"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(list); err != nil {
		fmt.Println("err", err)
		return
	}
}

func (h *Handler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idStr := query.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id must be a number"))
		return
	}
	err = h.storage.Delete(r.Context(), id)
	if errors.Is(err, ErrIdNotFound) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("if not found"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("problem in server"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
