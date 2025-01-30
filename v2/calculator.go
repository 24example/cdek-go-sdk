package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CalculatorTariffListRequest struct {
	// Date Дата и время планируемой передачи заказа. По умолчанию - текущая
	Date string `json:"date,omitempty"`
	// Type Тип заказа: 1 - "интернет-магазин", 2 - "доставка". По умолчанию - 1
	Type string `json:"type,omitempty"`
	// Валюта, в которой необходимо произвести расчет. По умолчанию - валюта договора
	Currency int `json:"currency,omitempty"`
	// Lang Локализация офиса. По умолчанию "rus"
	Lang string `url:"lang,omitempty"`
	// FromLocation Адрес отправления
	FromLocation Location `json:"from_location,omitempty"`
	// ToLocation Адрес получения
	ToLocation Location `json:"to_location"`
	// Packages Список информации по местам (упаковкам)
	Packages []Package `json:"packages"`
}

type CalculatorTariffRequest struct {
	// Type Тип заказа: 1 - "интернет-магазин", 2 - "доставка". По умолчанию - 1
	Type int `json:"type,omitempty"`
	// 136 or 234
	TariffCode int `json:"tariff_code"`
	// FromLocation Адрес отправления
	FromLocation Location `json:"from_location,omitempty"`
	// ToLocation Адрес получения
	ToLocation Location `json:"to_location"`
	// Packages Список информации по местам (упаковкам)
	Packages []Package `json:"packages"`
}

type Tariff struct {
	TariffCode        int     `json:"tariff_code"`
	TariffName        string  `json:"tariff_name"`
	TariffDescription string  `json:"tariff_description"`
	DeliveryMode      int     `json:"delivery_mode"`
	DeliverySum       float64 `json:"delivery_sum"`
	PeriodMin         int     `json:"period_min"`
	PeriodMax         int     `json:"period_max"`
	CalendarMin       int     `json:"calendar_min"`
	CalendarMax       int     `json:"calendar_max"`
}

type CalculatorTrafiffListResponse struct {
	TariffCodes []Tariff `json:"tariff_codes"`
}

func (c *clientImpl) CalculatorTariffList(ctx context.Context, input *CalculatorTariffListRequest) (*CalculatorTrafiffListResponse, error) {
	if input.Lang == "" {
		input.Lang = "rus"
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.buildUri("/v2/calculator/tarifflist", nil),
		bytes.NewReader(payload),
	)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	return jsonReq[CalculatorTrafiffListResponse](req)
}

func (c *clientImpl) CalculatorTariff(ctx context.Context, input *CalculatorTariffRequest) (*Tariff, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.buildUri("/v2/calculator/tariff", nil),
		bytes.NewReader(payload),
	)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	return jsonReq[Tariff](req)
}
