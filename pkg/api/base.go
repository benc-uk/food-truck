package api

// This is intended to be wrapped & extended by the application API struct
type Base struct {
	Healthy bool   `json:"healthy"` // Whether the server is healthy
	Version string `json:"version"` // Version of the API
	Name    string `json:"name"`    // Name of this API	or web service
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
