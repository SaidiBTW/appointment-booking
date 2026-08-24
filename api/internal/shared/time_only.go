package shared

import (
	"fmt"
	"time"
)

type TimeOnly time.Time

func (t TimeOnly) String() string {
	return time.Time(t).Format("15:04:05")
}

func (t TimeOnly) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).Format("15:04:05")
	return fmt.Appendf(nil, "\"%s\"", formatted), nil
}

func (t *TimeOnly) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*t = TimeOnly(v)
	case []byte:
		parsed, err := time.Parse("15:04:05", string(v))
		if err != nil {
			return err
		}
		*t = TimeOnly(parsed)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}
