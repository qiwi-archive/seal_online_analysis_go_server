package safeerror
import "errors"

type ISafeError interface {
	error
	Code() string
}

type SafeError struct {
	error
	code string
}

func NewByCode(code string)(*SafeError) {
	return &SafeError{errors.New(code), code}
}

func New(err error, code string)(*SafeError) {
	return &SafeError{err, code}
}

func (self *SafeError)Code()(string) {
	return self.code
}