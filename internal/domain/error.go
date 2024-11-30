package domain

import (
	"errors"
	"net/http"
)

var (
	ErrorUnauthorized = errors.New("Unauthorized")
	a                 = http.StatusAccepted
)
