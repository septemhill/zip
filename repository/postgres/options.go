package postgres

type options struct {
	txHolder bool
}

type Option interface {
	apply(*options)
}

type txHolder bool

func (o txHolder) apply(opts *options) {
	opts.txHolder = bool(o)
}

func WithTxHolder(b bool) Option {
	return txHolder(b)
}
