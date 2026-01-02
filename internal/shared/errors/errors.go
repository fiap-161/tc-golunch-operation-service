package errors

type ErrorDTO struct {
	Message      string `json:"message"`
	MessageError string `json:"message_error"`
}

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}

type UnauthorizedError struct {
	Msg string
}

func (e *UnauthorizedError) Error() string {
	return e.Msg
}

type InternalError struct {
	Msg string
}

func (e *InternalError) Error() string {
	return e.Msg
}

type NotFoundError struct {
	Msg string
}

func (e *NotFoundError) Error() string {
	return e.Msg
}
