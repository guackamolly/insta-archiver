package http

import (
	"encoding/json"
	"io"
)

type HttpBody struct {
	read io.ReadCloser

	value *[]byte
	err   error
}

func (b *HttpBody) Value() ([]byte, error) {
	if b.err != nil {
		return nil, b.err
	}

	if b.value == nil {
		defer b.read.Close()
		r, err := io.ReadAll(b.read)

		b.err = err
		b.value = &r
	}

	return *b.value, b.err
}

func Typed[T interface{}](b *HttpBody) (T, error) {
	var res T
	r, err := b.Value()

	if err != nil {
		return res, err
	}

	err = json.Unmarshal(r, &res)

	return res, err
}
