package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type String sql.NullString

func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
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
