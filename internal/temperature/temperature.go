package temperature

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const path = "/sys/class/thermal/thermal_zone0/temp"

func CurrentTemperature() (int, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %w", path, err)
	}

	temperatureString := strings.TrimSuffix(string(content), "\n")
	temperature, err := strconv.Atoi(temperatureString)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", content, err)
	}

	return temperature / 1000, nil
}
