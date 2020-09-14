// @Author : liguoyu
// @Date: 2019/10/29 15:42
package errors

// Errnoer is the interface that wraps the Tip method.
type Errnoer interface {
	Errno() int
}

type withErrno struct {
	error
	errno int
}

func (e *withErrno) Errno() int {
	return e.errno
}

func (e *withErrno) Cause() error {
	return e.error
}

// WithErrno annotates err with a errno.
// If err is nil, WithErrno returns nil.
func WithErrno(err error, errno int) error {
	// NOTE: 下面这句是必须的，否则：%!v(PANIC=runtime error: invalid memory address or nil pointer dereference)
	if err == nil {
		return nil
	}

	return &withErrno{
		error: err,
		errno: errno,
	}
}
