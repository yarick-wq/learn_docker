package handlers

import (
	"errors"
	"strings"
)

type CreateEmployeeRequest struct {
	FullName string `json:"full_name"`
	Position string `json:"position"`
}

func (c CreateEmployeeRequest) Validate() error {
	full_name := strings.TrimSpace(c.FullName)
	posit := strings.TrimSpace(c.Position)
	switch {
	case full_name == "" && posit == "":
		return errors.New("fullname and position are empty")
	case full_name == "":
		return errors.New("fullname is empty")
	case posit == "":
		return errors.New("position is empty")
	default:
		return nil
	}
}
