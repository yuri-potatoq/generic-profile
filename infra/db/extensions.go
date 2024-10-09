package db

import (
	"bytes"
	"errors"
	"io"
)

func ScanEnum[T ~string](t *T, tName string) func(any) error {
	scanErr := errors.New("Incompatible type for " + tName)
	return func(src any) error {
		var source []byte
		switch src.(type) {
		case string:

			source = []byte(src.(string))
		case []byte:
			source = src.([]byte)
		default:
			return scanErr
		}
		b, err := io.ReadAll(bytes.NewReader(source))
		if err != nil {
			return errors.Join(scanErr, err)
		}
		*t = T(b)
		return nil
	}
}
