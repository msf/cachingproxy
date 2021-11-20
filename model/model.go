// Package model implements the application entities
package model

// Message represents a message for the application
type Message struct {
	ID      string
	Content string
}

// Valid checks and returns if the message is valid
func (m Message) Valid() (valid bool) {
	switch {
	case m.ID == "":
		return false
	case m.Content == "":
		return false
	case m.ID == m.Content:
		return false
	}

	return true
}
