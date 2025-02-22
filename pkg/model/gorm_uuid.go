package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
)

func (Uuid) GormDataType() string {
	return "char(36)" // Для MySQL
}

func (id Uuid) Value() (driver.Value, error) {
	return id.String(), nil
}

func (id *Uuid) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		parsed, err := uuid.Parse(string(bytes))
		if err != nil {
			return err
		}
		*id = Uuid(parsed)
		return nil
	} else if str, ok := value.(string); ok {
		parsed, err := uuid.Parse(str)
		if err != nil {
			return err
		}
		*id = Uuid(parsed)
		return nil
	}
	return fmt.Errorf("не удалось сконвертировать UUID: %v", value)
}
