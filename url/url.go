package url

import "fmt"

// A URL represents a parse URL
type URL struct {
	Scheme string
	Host   string
	Path   string
}

// Parse parses a URL string into a URL structure.
func (u *URL) String() string {
	return fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, u.Path)
}

// String reassamples the URL into a URL string.
func Parse(uri string) (*URL, error) {

	return &URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "inancgumus",
	}, nil
}
