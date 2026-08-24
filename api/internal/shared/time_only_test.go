package shared_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/SaidiBTW/appointment_booking_system_go/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestTimeOnly_String(t *testing.T) {
	parsedTime, err := time.Parse("15:04:05", "14:30:00")
	assert.NoError(t, err)

	to := shared.TimeOnly(parsedTime)
	assert.Equal(t, "14:30:00", to.String())
}

func TestTimeOnly_MarshalJSON(t *testing.T) {
	parsedTime, err := time.Parse("15:04:05", "09:15:00")
	assert.NoError(t, err)

	to := shared.TimeOnly(parsedTime)
	data, err := json.Marshal(to)

	assert.NoError(t, err)
	assert.Equal(t, `"09:15:00"`, string(data))
}

func TestTimeOnly_Scan(t *testing.T) {
	t.Run("Scan time.Time", func(t *testing.T) {
		now := time.Now()
		var to shared.TimeOnly
		err := to.Scan(now)

		assert.NoError(t, err)
		assert.Equal(t, time.Time(to), now)
	})

	t.Run("Scan []byte valid time string", func(t *testing.T) {
		var to shared.TimeOnly
		err := to.Scan([]byte("16:45:00"))

		assert.NoError(t, err)
		assert.Equal(t, "16:45:00", to.String())
	})

	t.Run("Scan []byte invalid time string", func(t *testing.T) {
		var to shared.TimeOnly
		err := to.Scan([]byte("invalid-time"))

		assert.Error(t, err)
	})

	t.Run("Scan unsupported type", func(t *testing.T) {
		var to shared.TimeOnly
		err := to.Scan(12345)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported type: int")
	})
}
