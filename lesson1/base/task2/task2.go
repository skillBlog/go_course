package task2

import (
	"errors"
	"fmt"
	"math/rand"
)

var (
	ErrInvalidAmount       = errors.New("некорректная сумма платежа")
	ErrProviderUnavailable = errors.New("провайдер недоступен")
)

type PaymentProcessor interface {
	ProcessPayment(amount float64) error
}

type SberbankPaymentProcessor struct {
	APIKey string
}

func (s *SberbankPaymentProcessor) ProcessPayment(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if rand.Float64() < 0.5 {
		return ErrProviderUnavailable
	}

	return nil
}

type TbankPaymentProcessor struct {
	APIKey string
}

func (t *TbankPaymentProcessor) ProcessPayment(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if rand.Float64() < 0.5 {
		return ErrProviderUnavailable
	}

	return nil
}

type AlfabankPaymentProcessor struct {
	APIKey string
}

func (a *AlfabankPaymentProcessor) ProcessPayment(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if rand.Float64() < 0.5 {
		return ErrProviderUnavailable
	}

	return nil
}

func NewAlfabankPaymentProcessor(APIKey string) *AlfabankPaymentProcessor {
	return &AlfabankPaymentProcessor{APIKey: APIKey}
}

func NewTbankPaymentProcessor(APIKey string) *TbankPaymentProcessor {
	return &TbankPaymentProcessor{APIKey: APIKey}
}

func NewSberbankPaymentProcessor(APIKey string) *SberbankPaymentProcessor {
	return &SberbankPaymentProcessor{APIKey: APIKey}
}

func main() {
	processors := []PaymentProcessor{
		&SberbankPaymentProcessor{APIKey: "SBER-123"},
		&TbankPaymentProcessor{APIKey: "TBANK-456"},
		&AlfabankPaymentProcessor{APIKey: "ALFA-789"},
	}

	amount := 1000.0

	for _, p := range processors {
		err := p.ProcessPayment(amount)

		switch err {
		case nil:
			fmt.Println("Платеж успешно обработан")
		case ErrInvalidAmount:
			fmt.Println("Ошибка: некорректная сумма")
		case ErrProviderUnavailable:
			fmt.Println("Ошибка: провайдер недоступен")
		}
	}
}
