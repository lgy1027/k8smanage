// @Author : liguoyu
// @Date: 2019/10/29 15:42
package errors

// Tipper is the interface that wraps the Tip method.
type Tipper interface {
	Tip() string
}

type withTipMessage struct {
	error
	msg string
}

func (e *withTipMessage) Tip() string {
	return e.msg
}

func (e *withTipMessage) Cause() error {
	return e.error
}

// WithTipMessage annotates err with a tip message.
// If err is nil, WithTipMessage returns nil.
func WithTipMessage(err error, msg string) error {
	// NOTE: 下面这句是必须的，否则：%!v(PANIC=runtime error: invalid memory address or nil pointer dereference)
	if err == nil {
		return nil
	}

	return &withTipMessage{
		error: err,
		msg:   msg,
	}
}
