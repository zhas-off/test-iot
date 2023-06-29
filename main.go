package main

import (
	"encoding/hex"
	"fmt"
	"log"
)

type Result struct {
	Temperature    float64
	Humidity       float64
	MagneticStatus string
}

func parseUint64(hexString string) (uint64, error) {
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return 0, err
	}

	n := uint64(0)
	mul := uint64(1)
	for i := 0; i < len(decoded); i++ {
		n = n + uint64(decoded[i])*mul
		mul *= 256
	}

	return n, nil
}

func parseMagneticStatus(hex string) (string, error) {
	switch hex {
	case "00":
		return "Close", nil
	case "01":
		return "Open", nil
	default:
		return "", fmt.Errorf("unexpected value")
	}
}

func Decode(hex string) (*Result, error) {
	res := &Result{}
	for i := 0; i < len(hex); {
		// read channel
		ch := hex[i : i+2]
		i += 2
		// type
		dataType := hex[i : i+2]
		i += 2
		// read data
		switch ch {
		case "03":
			if dataType == "67" { // Temperature
				n, err := parseUint64(hex[i : i+4])
				if err != nil {
					return nil, fmt.Errorf("cannot parse uint64 from channel 03 type 67")
				}
				res.Temperature = float64(n) / 10

				i += 4
			} else {
				return nil, fmt.Errorf("channel 03 has unknown type")
			}
		case "04":
			if dataType == "68" { // Humidity
				n, err := parseUint64(hex[i : i+2])
				if err != nil {
					return nil, fmt.Errorf("cannot parse uint64 from channel 05 type 68")
				}
				res.Humidity = float64(n) / 2

				i += 2
			} else {
				return nil, fmt.Errorf("channel 05 has unknown type")
			}
		case "06":
			if dataType == "00" { // MagneticStatus
				status, err := parseMagneticStatus(hex[i : i+2])
				if err != nil {
					return nil, fmt.Errorf("cannot parse status from channel 06 type 00")
				}
				res.MagneticStatus = status
				i += 2
			} else {
				return nil, fmt.Errorf("channel 06 has unknown type")
			}
		}
	}
	return res, nil
}

func main() {
	str := "0367F600046882060001"

	result, err := Decode(str)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("Temparature: %v C\nHumidity: %v%%\nMagneticStatus: %v", result.Temperature, result.Humidity, result.MagneticStatus)
}
