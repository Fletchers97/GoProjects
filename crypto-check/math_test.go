package main

import "testing"

// Это функция, которую мы как будто тестируем
// В реальном проекте она была бы в другом файле
func CalculateDiff(currentPrice, averagePrice float64) float64 {
	if averagePrice == 0 {
		return 0
	}
	return ((currentPrice - averagePrice) / averagePrice) * 100
}

// А это сам тест
func TestCalculateDiff(t *testing.T) {
	// Входные данные
	current := 110.0
	average := 100.0
	expected := 10.0

	// Выполнение функции
	result := CalculateDiff(current, average)

	// Проверка результата
	if result != expected {
		t.Errorf("Ошибка! Ожидали %.2f, но получили %.2f", expected, result)
	}
}
