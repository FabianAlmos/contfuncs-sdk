package handler

import (
	"context"
	"contfunc-sdk/fn_http"
)

type Handler[In any, Out any] func(ctx context.Context, req fn_http.Request[In]) (fn_http.Response[Out], error)
