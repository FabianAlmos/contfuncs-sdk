package contfuncs

import "strconv"

type InvokeOptions struct {
	Version string
}

func NewInvokeOptions() *InvokeOptions {
	return &InvokeOptions{}
}

type InvokeOption func(*InvokeOptions)

func WithVersion(v uint64) InvokeOption {
	return func(o *InvokeOptions) {
		o.Version = strconv.FormatUint(v, 10)
	}
}
