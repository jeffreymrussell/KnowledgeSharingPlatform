package internal

import "errors"

// ErrNotFound is used when a resource is not found in the repository.
var ErrNotFound = errors.New("resource not found")

// Add other custom errors here as needed.
