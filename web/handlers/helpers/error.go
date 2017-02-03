package helpers

type HttpError struct {
  Error string `json:"error"`
}

func NewHttpError(err error) HttpError {
  return HttpError{
    Error: err.Error(),
  }
}
