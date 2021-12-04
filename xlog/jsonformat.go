// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

import (
	"encoding/json"
	"time"
)

type JSONFormat struct {
	FormatType TextFormatType
	timeFormat string // should no change for JSON
}

func (j *JSONFormat) Format(e *Entry) ([]byte, error) {
	j.timeFormat = time.RFC3339Nano

	r, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	r = append(r, '\n')
	return r, nil
}
