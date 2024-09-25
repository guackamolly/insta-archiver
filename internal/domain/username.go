package domain

import (
	"errors"
	"regexp"
	"strings"
)

// A simple regex for matching valid Instagram usernames (can contain letters, numbers, dots (.) and underscores (_))
// It's not meant to be seen as a regex for creating a new username, just a regex to stop bad actors of doing naughty things
// on the server
var usernameRegex = regexp.MustCompile("[A-Za-z0-9._]+")

type PurifyUsername struct{}

func (u PurifyUsername) Invoke(username string) (string, error) {
	s := stripUsername(username)

	return s, validateUsername(s)
}

// Required to be called before [validateUsername], as users may submit an username starting with @
func stripUsername(username string) string {
	if strings.HasPrefix(username, "@") {
		return strings.TrimPrefix(username, "@")
	}

	return username
}

func validateUsername(username string) error {
	l := len(username)

	if l > 30 {
		return errors.New("username exceeds 30 characters")
	}

	if m := usernameRegex.FindStringSubmatch(username); len(m[0]) != l {
		return errors.New("username contains invalid characters")
	}

	return nil
}
