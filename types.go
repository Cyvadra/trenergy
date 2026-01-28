package trenergy

// APIResponse represents a generic API response wrapper.
type APIResponse[T any] struct {
	Status bool   `json:"status"`
	Data   T      `json:"data"`
	Links  *Links `json:"links,omitempty"`
	Meta   *Meta  `json:"meta,omitempty"`
}

// Links represents pagination links.
type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}

// Meta represents pagination metadata.
type Meta struct {
	CurrentPage int        `json:"current_page"`
	From        int        `json:"from"`
	LastPage    int        `json:"last_page"`
	Links       []MetaLink `json:"links"`
	Path        string     `json:"path"`
	PerPage     int        `json:"per_page"`
	To          int        `json:"to"`
	Total       int        `json:"total"`
}

// MetaLink represents a link in the metadata.
type MetaLink struct {
	URL    *string `json:"url"`
	Label  string  `json:"label"`
	Active bool    `json:"active"`
}

// ErrorResponse represents an error response from the API (if structured differently, though usually errors are just non-200 with maybe a message).
// Based on samples, success is status: true. We'll need to handle status: false or HTTP errors.
