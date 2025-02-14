package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONMap map[string]string

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = JSONMap{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("не удалось конвертировать в JSON")
	}

	return json.Unmarshal(bytes, j)
}
