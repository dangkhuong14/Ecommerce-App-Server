package common

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func (id *UUID) Scan(src interface{}) error {
	var sqlID uuid.UUID

	if err := sqlID.Scan(src); err != nil {
		return err
	}

	*id = UUID(sqlID)

	return nil
}

func (id UUID) Value() (driver.Value, error) {
	return uuid.UUID(id).MarshalBinary()
}

// Method chuyển UUID từ byte slice thành chuỗi
func (id UUID) String() string {
	return uuid.UUID(id).String()
}

func (id UUID) MarshalJSON() ([]byte, error) {
    // Chuyển đổi UUID sang string bằng hàm String() của common.UUID
    s := id.String()
    return json.Marshal(s)
}