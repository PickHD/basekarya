package constants

import (
	"database/sql/driver"
	"fmt"
)

type MaritalStatus string

const (
	MaritalStatusSingle  MaritalStatus = "TK"
	MaritalStatusMarried MaritalStatus = "K"
)

func (s MaritalStatus) Value() (driver.Value, error) {
	if s == "" {
		return nil, nil
	}
	return string(s), nil
}

func (s *MaritalStatus) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*s = MaritalStatus(v)
		return nil
	case string:
		*s = MaritalStatus(v)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into MaritalStatus", value)
	}
}
