package delivery

type Error struct {
	Message string `json:"message"`
}

func GetErrorMessage(err error) *Error {
	return &Error{Message: err.Error()}
}
