package models

import (
	"fmt"
	"strings"
)

type Quote struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Text   string `json:"quote"`
}

func (q *Quote) Validate() error {
	if strings.TrimSpace(q.Author) == "" {
		return fmt.Errorf("author is required")
	}
	if strings.TrimSpace(q.Text) == "" {
        return fmt.Errorf("quote is required")
    }
    return nil
}
