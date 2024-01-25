package enum

import (
	"encoding/json"
	"errors"
	"strings"
)

type Algorithm uint

const (
	RSA Algorithm = iota + 1
	ECC
)

func (algorithm *Algorithm) UnmarshalJSON(b []byte) error {
	var algorithmString string
	if err := json.Unmarshal(b, &algorithmString); err != nil {
		return err
	}
	switch strings.ToUpper(algorithmString) {
	case "ECC":
		*algorithm = ECC
	case "RSA":
		*algorithm = RSA
	default:
		return errors.New("Algorithm is not supported!")
	}

	return nil
}
