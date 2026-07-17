package steam

import "errors"

var (
	ErrNotFound     = errors.New("Steam is not installed")
	ErrGameNotFound = errors.New("Among Us was not found in any Steam library")
)
