package main

import (
	"errors"
	"fmt"
)

var (
	ErrEngineAlreadyRunning = errors.New("двигатель уже работает")
	ErrEngineOff            = errors.New("двигатель не запущен")
	ErrLowBattery           = errors.New("низкий заряд батареи")
)

type Vehicle interface {
	StartEngine() error
	StopEngine() error
	GetInfo() string
}

type Car struct {
	Brand    string
	engineOn bool
}

func (c *Car) StartEngine() error {
	if c.engineOn {
		return ErrEngineAlreadyRunning
	}
	c.engineOn = true
	return nil
}

func (c *Car) StopEngine() error {
	if !c.engineOn {
		return ErrEngineOff
	}
	c.engineOn = false
	return nil
}

func (c *Car) GetInfo() string {
	return fmt.Sprintf("Марка: %s, Состояние двигателя: %t", c.Brand, c.engineOn)
}

func (c *Car) GetEngineStatus() bool {
	return c.engineOn
}

func (c *Car) Honk() string {
	return "Beep beep!"
}

type Truck struct {
	Car
	cargoCapacity float64
}

func (t *Truck) GetInfo() string {
	return fmt.Sprintf(
		"Марка: %s, Состояние двигателя: %t, Грузоподъемность: %.0f тонн",
		t.Brand, t.engineOn, t.cargoCapacity,
	)
}

func (t *Truck) Honk() string {
	return "Honk Honk!"
}

func (t *Truck) GetCargoCapacity() float64 {
	return t.cargoCapacity
}

type ElectricCar struct {
	Car
	batteryLevel int
}

func (e *ElectricCar) StartEngine() error {
	if e.engineOn {
		return ErrEngineAlreadyRunning
	}
	if e.batteryLevel <= 5 {
		return ErrLowBattery
	}
	e.engineOn = true
	return nil
}

func (e *ElectricCar) GetInfo() string {
	return fmt.Sprintf(
		"Марка: %s, Состояние двигателя: %t, Уровень заряда батареи: %d%%",
		e.Brand, e.engineOn, e.batteryLevel,
	)
}

func (e *ElectricCar) GetBatteryLevel() int {
	return e.batteryLevel
}

func main() {
	fmt.Println("Hello, World!")
}
