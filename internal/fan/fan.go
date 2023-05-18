package fan

import (
	"fmt"

	"gobot.io/x/gobot/platforms/raspi"
)

const address = 0x1A

func SetSpeed(speed int) error {
	if speed < 0 || speed > 100 {
		return fmt.Errorf("speed (%d) is out of range (0-100)", speed)
	}

	adapter := raspi.NewAdaptor()
	defer adapter.Finalize()

	bus := adapter.GetDefaultBus()
	conn, err := adapter.GetConnection(address, bus)
	if err != nil {
		return fmt.Errorf("failed to connect to fan: %w", err)
	}
	defer conn.Close()

	err = conn.WriteByte(byte(speed))
	if err != nil {
		return fmt.Errorf("failed to write speed (%d) to fan: %w", speed, err)
	}

	return nil
}
