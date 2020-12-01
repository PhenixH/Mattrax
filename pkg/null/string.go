package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type String sql.NullString

func (ns String) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *String) UnmarshalJSON(data []byte) error {
	var str *string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	ns.Valid = str != nil
	if ns.Valid {
		ns.String = *str
	}

	return nil
}

func (ns *String) Scan(value interface{}) error {
	if value == nil {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	ns.String = value.(string)
	return nil
}

func (ns String) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}
