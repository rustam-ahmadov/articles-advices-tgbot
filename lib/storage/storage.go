package storage

import (
	"crypto/sha1"
	"fmt"
)

type Storage interface {
	Save(p *Page) error
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
	PickRandom(userName string) (*Page, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	h.Write([]byte(p.URL))
	h.Write([]byte(p.UserName))

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
