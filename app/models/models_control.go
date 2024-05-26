package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

// Value implements the driver.Valuer interface for decimal.Decimal.
func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements the sql.Scanner interface for decimal.Decimal.
func (d *Decimal) Scan(value interface{}) error {
	if v, ok := value.([]byte); ok {
		dec, err := decimal.NewFromString(string(v))
		if err != nil {
			return err
		}
		d.Decimal = dec
		return nil
	} else if v, ok := value.(string); ok {
		dec, err := decimal.NewFromString(v)
		if err != nil {
			return err
		}
		d.Decimal = dec
		return nil
	} else if v, ok := value.(float64); ok {
		d.Decimal = decimal.NewFromFloat(v)
		return nil
	} else if v, ok := value.(int64); ok {
		d.Decimal = decimal.NewFromInt(v)
		return nil
	} else {
		return fmt.Errorf("cannot scan type %T into decimal.Decimal", value)
	}
}