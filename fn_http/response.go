package fn_http

type Response[T any] struct {
	Data       T     `json:"data,omitempty"`
	StatusCode int   `json:"status,omitempty"`
	Err        error `json:"error,omitempty"`
}

func Success[T any](statusCode int, data T) Response[T] {
	return Response[T]{StatusCode: statusCode, Data: data}
}

func Failure[T any](statusCode int, err error) Response[T] {
	return Response[T]{StatusCode: statusCode, Err: err}
}
