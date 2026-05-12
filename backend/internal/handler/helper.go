package handler

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

func decodeJSONBody(bodyReader interface{ Read([]byte) (int, error) }, dst interface{}) error {
	decoder := json.NewDecoder(bodyReader)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func parseBirthday(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	layouts := []string{
		"2006-01-02",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return &t, nil
		}
	}

	return nil, errors.New("birthday format must be 2006-01-02 or RFC3339")
}
