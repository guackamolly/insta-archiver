package core

import "github.com/guackamolly/insta-archiver/internal/data/user"

// Main application DI container.
type Vault struct {
	UserRepository user.UserRepository
}
