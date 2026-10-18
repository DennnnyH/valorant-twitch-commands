package upstream

import "fmt"

// Error is a non-success response from Riot or Henrik.
type Error struct {
	Source     string
	Status     int
	RetryAfter string
	Detail     string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: status %d: %s", e.Source, e.Status, e.Detail)
}
