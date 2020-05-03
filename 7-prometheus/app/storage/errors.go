package storage

import (
	"errors"
	"fmt"
)

var (
	// Sentinel errors
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")

	ErrBadInput = errors.New("bad input")

	// User
	ErrUserNotFound      = fmt.Errorf("user %w", ErrNotFound)
	ErrUserAlreadyExists = fmt.Errorf("user %w", ErrAlreadyExists)

	// User settings
	ErrUserSettingsNotFound = fmt.Errorf("user settings %w", ErrNotFound)

	// User contacts
	ErrUserContactNotFound           = fmt.Errorf("user contact %w", ErrNotFound)
	ErrUserContactAlreadyExists      = fmt.Errorf("user contact %w", ErrAlreadyExists)
	ErrUserContactBelongsToOtherUser = errors.New("user contact belongs to other user")
)
