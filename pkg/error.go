package pkg

import "fmt"

type NilConnError struct {
	Service Service
}

func (e NilConnError) Error() string {
	return e.Service.String() + " has nil conn"
}

type TimeoutError struct {
	err     error
	timeout int
	service string
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("%s spended out of %ds, %s", e.service, e.timeout, e.err.Error())
}

func (e TimeoutError) Unwrap() error { return e.err }
