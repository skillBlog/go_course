package task3

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrUnsupported = errors.New("обновление недоступно")
)

func versionToFloat(version string) (float64, error) {
	return strconv.ParseFloat(version, 64)
}

type Device interface {
	UpdateOS(version string) error
	GetInfo() string
}

type Smartphone struct {
	OSVersion string
	Model     string
}

func (s *Smartphone) UpdateOS(version string) error {
	currentVersion, err := versionToFloat(s.OSVersion)
	if err != nil {
		return ErrUnsupported
	}

	if currentVersion >= 12.0 {
		return ErrUnsupported
	}

	s.OSVersion = version
	return nil
}

func (s *Smartphone) GetInfo() string {
	return fmt.Sprintf("Модель: %s, ОС: %s", s.Model, s.OSVersion)
}

type Laptop struct {
	OSVersion string
	Model     string
}

func (l *Laptop) UpdateOS(version string) error {
	if !strings.HasPrefix(version, "Windows") {
		return ErrUnsupported
	}
	l.OSVersion = version
	return nil
}

func (l *Laptop) GetInfo() string {
	return fmt.Sprintf("Модель: %s, ОС: %s", l.Model, l.OSVersion)
}

type Smartwatch struct {
	OSVersion string
	Model     string
}

func (s *Smartwatch) UpdateOS(version string) error {
	if len(version) < 5 {
		return ErrUnsupported
	}
	s.OSVersion = version
	return nil
}

func (s *Smartwatch) GetInfo() string {
	return fmt.Sprintf("Модель: %s, ОС: %s", s.Model, s.OSVersion)
}

func main() {
	devices := []struct {
		device  Device
		version string
	}{
		{
			device: &Smartphone{
				OSVersion: "11.0",
				Model:     "iPhone 11",
			},
			version: "12.0",
		},
		{
			device: &Laptop{
				OSVersion: "Windows 10",
				Model:     "Lenovo ThinkPad",
			},
			version: "Windows 11",
		},
		{
			device: &Smartwatch{
				OSVersion: "1.0",
				Model:     "Apple Watch",
			},
			version: "watchOS 10",
		},
	}

	for _, d := range devices {
		fmt.Println("До обновления:")
		fmt.Println(d.device.GetInfo())

		if err := d.device.UpdateOS(d.version); err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			fmt.Println("Обновление успешно")
		}

		fmt.Println("После обновления:")
		fmt.Println(d.device.GetInfo())
		fmt.Println("--------------------")
	}
}
