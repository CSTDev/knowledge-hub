package types

import "fmt"

type Field struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Order int    `json:"order"`
}

type FieldNotFoundError struct {
	ID      string
	Message string
}

func (fnf FieldNotFoundError) Error() string {
	return fmt.Sprintf("%s : %s", fnf.Message, fnf.ID)
}
