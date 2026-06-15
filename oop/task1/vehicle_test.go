package main

import (
	"errors"
	"strings"
	"testing"
)

func TestStartEngine(t *testing.T) {
	tests := []struct {
		name       string
		vehicle    Vehicle
		wantErr    error
		wantEngine bool
	}{
		{
			name:       "Car: успешный запуск",
			vehicle:    &Car{Brand: "Toyota"},
			wantErr:    nil,
			wantEngine: true,
		},
		{
			name: "Car: двигатель уже работает",
			vehicle: &Car{Brand: "Toyota", engineOn: true},
			wantErr:    ErrEngineAlreadyRunning,
			wantEngine: true,
		},
		{
			name:       "Truck: успешный запуск через композицию",
			vehicle:    &Truck{Car: Car{Brand: "Volvo"}, cargoCapacity: 20},
			wantErr:    nil,
			wantEngine: true,
		},
		{
			name: "Truck: двигатель уже работает",
			vehicle: &Truck{
				Car:           Car{Brand: "Volvo", engineOn: true},
				cargoCapacity: 20,
			},
			wantErr:    ErrEngineAlreadyRunning,
			wantEngine: true,
		},
		{
			name:       "ElectricCar: успешный запуск при заряде > 5%",
			vehicle:    &ElectricCar{Car: Car{Brand: "Tesla"}, batteryLevel: 80},
			wantErr:    nil,
			wantEngine: true,
		},
		{
			name:       "ElectricCar: запуск при заряде ровно 6%",
			vehicle:    &ElectricCar{Car: Car{Brand: "Tesla"}, batteryLevel: 6},
			wantErr:    nil,
			wantEngine: true,
		},
		{
			name:       "ElectricCar: низкий заряд 5%",
			vehicle:    &ElectricCar{Car: Car{Brand: "Tesla"}, batteryLevel: 5},
			wantErr:    ErrLowBattery,
			wantEngine: false,
		},
		{
			name:       "ElectricCar: низкий заряд 0%",
			vehicle:    &ElectricCar{Car: Car{Brand: "Tesla"}, batteryLevel: 0},
			wantErr:    ErrLowBattery,
			wantEngine: false,
		},
		{
			name: "ElectricCar: двигатель уже работает",
			vehicle: &ElectricCar{
				Car:          Car{Brand: "Tesla", engineOn: true},
				batteryLevel: 80,
			},
			wantErr:    ErrEngineAlreadyRunning,
			wantEngine: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vehicle.StartEngine()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("StartEngine() error = %v, want %v", err, tt.wantErr)
			}

			status := engineStatus(tt.vehicle)
			if status != tt.wantEngine {
				t.Fatalf("engine status = %v, want %v", status, tt.wantEngine)
			}
		})
	}
}

func TestStopEngine(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() Vehicle
		wantErr    error
		wantEngine bool
	}{
		{
			name: "Car: успешная остановка",
			setup: func() Vehicle {
				c := &Car{Brand: "Toyota", engineOn: true}
				return c
			},
			wantErr:    nil,
			wantEngine: false,
		},
		{
			name:       "Car: двигатель уже выключен",
			setup:      func() Vehicle { return &Car{Brand: "Toyota"} },
			wantErr:    ErrEngineOff,
			wantEngine: false,
		},
		{
			name: "Truck: успешная остановка",
			setup: func() Vehicle {
				return &Truck{Car: Car{Brand: "Volvo", engineOn: true}, cargoCapacity: 15}
			},
			wantErr:    nil,
			wantEngine: false,
		},
		{
			name: "ElectricCar: успешная остановка",
			setup: func() Vehicle {
				return &ElectricCar{Car: Car{Brand: "Tesla", engineOn: true}, batteryLevel: 50}
			},
			wantErr:    nil,
			wantEngine: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.setup()
			err := v.StopEngine()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("StopEngine() error = %v, want %v", err, tt.wantErr)
			}

			if status := engineStatus(v); status != tt.wantEngine {
				t.Fatalf("engine status = %v, want %v", status, tt.wantEngine)
			}
		})
	}
}

func TestVehiclePolymorphism(t *testing.T) {
	tests := []struct {
		name         string
		vehicle      Vehicle
		wantBrand    string
		wantInfoPart string
	}{
		{
			name:         "Car через интерфейс Vehicle",
			vehicle:      &Car{Brand: "BMW"},
			wantBrand:    "BMW",
			wantInfoPart: "Состояние двигателя: true",
		},
		{
			name:         "Truck через интерфейс Vehicle",
			vehicle:      &Truck{Car: Car{Brand: "MAN"}, cargoCapacity: 12.5},
			wantBrand:    "MAN",
			wantInfoPart: "Грузоподъемность: 12 тонн",
		},
		{
			name:         "ElectricCar через интерфейс Vehicle",
			vehicle:      &ElectricCar{Car: Car{Brand: "Nissan"}, batteryLevel: 42},
			wantBrand:    "Nissan",
			wantInfoPart: "Уровень заряда батареи: 42%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.vehicle.StartEngine(); err != nil {
				t.Fatalf("StartEngine() unexpected error: %v", err)
			}

			info := tt.vehicle.GetInfo()
			if !strings.Contains(info, tt.wantBrand) {
				t.Fatalf("GetInfo() = %q, want substring %q", info, tt.wantBrand)
			}
			if !strings.Contains(info, tt.wantInfoPart) {
				t.Fatalf("GetInfo() = %q, want substring %q", info, tt.wantInfoPart)
			}

			if err := tt.vehicle.StopEngine(); err != nil {
				t.Fatalf("StopEngine() unexpected error: %v", err)
			}

			if strings.Contains(tt.vehicle.GetInfo(), "true") {
				t.Fatalf("GetInfo() after stop = %q, engine should be off", tt.vehicle.GetInfo())
			}
		})
	}
}

func TestHonk(t *testing.T) {
	tests := []struct {
		name     string
		honker   interface{ Honk() string }
		wantHonk string
	}{
		{
			name:     "Car",
			honker:   &Car{Brand: "Toyota"},
			wantHonk: "Beep beep!",
		},
		{
			name:     "Truck переопределяет Honk",
			honker:   &Truck{Car: Car{Brand: "Volvo"}},
			wantHonk: "Honk Honk!",
		},
		{
			name:     "ElectricCar наследует Honk от Car",
			honker:   &ElectricCar{Car: Car{Brand: "Tesla"}, batteryLevel: 100},
			wantHonk: "Beep beep!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.honker.Honk(); got != tt.wantHonk {
				t.Fatalf("Honk() = %q, want %q", got, tt.wantHonk)
			}
		})
	}
}

func TestGetters(t *testing.T) {
	tests := []struct {
		name  string
		setup func() any
		check func(t *testing.T, v any)
	}{
		{
			name:  "Car GetEngineStatus",
			setup: func() any { return &Car{Brand: "Toyota"} },
			check: func(t *testing.T, v any) {
				c := v.(*Car)
				if c.GetEngineStatus() {
					t.Fatal("GetEngineStatus() = true, want false")
				}
				if err := c.StartEngine(); err != nil {
					t.Fatalf("StartEngine() error: %v", err)
				}
				if !c.GetEngineStatus() {
					t.Fatal("GetEngineStatus() = false, want true")
				}
			},
		},
		{
			name: "Truck GetCargoCapacity",
			setup: func() any {
				return &Truck{Car: Car{Brand: "Kamaz"}, cargoCapacity: 25.5}
			},
			check: func(t *testing.T, v any) {
				truck := v.(*Truck)
				if got := truck.GetCargoCapacity(); got != 25.5 {
					t.Fatalf("GetCargoCapacity() = %v, want 25.5", got)
				}
			},
		},
		{
			name: "ElectricCar GetBatteryLevel",
			setup: func() any {
				return &ElectricCar{Car: Car{Brand: "Tesla"}, batteryLevel: 73}
			},
			check: func(t *testing.T, v any) {
				ev := v.(*ElectricCar)
				if got := ev.GetBatteryLevel(); got != 73 {
					t.Fatalf("GetBatteryLevel() = %d, want 73", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, tt.setup())
		})
	}
}

func engineStatus(v Vehicle) bool {
	switch x := v.(type) {
	case *Car:
		return x.GetEngineStatus()
	case *Truck:
		return x.GetEngineStatus()
	case *ElectricCar:
		return x.GetEngineStatus()
	default:
		return false
	}
}
