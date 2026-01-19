package github

import (
	"github.com/google/go-github/v58/github"
)

// NewClient tworzy i zwraca nowego klienta API GitHub.
// Obecnie jest to klient nieuwierzytelniony, co wystarcza do odczytu publicznych repozytoriów.
func NewClient() *github.Client {
	return github.NewClient(nil)
}
