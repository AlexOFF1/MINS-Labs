package strategy

// сделать динамический выбор

import "sort"

type AverageStrategy interface {
	Calculate(grades []int) float64
	Name() string
}

type ArithmeticMean struct{}

func (a *ArithmeticMean) Calculate(grades []int) float64 {
	if len(grades) == 0 {
		return 0
	}
	sum := 0
	for _, g := range grades {
		sum += g
	}
	return float64(sum) / float64(len(grades))
}

func (a *ArithmeticMean) Name() string {
	return "arithmetic"
}

type MedianStrategy struct{}

func (m *MedianStrategy) Calculate(grades []int) float64 {
	if len(grades) == 0 {
		return 0
	}
	sorted := make([]int, len(grades))
	copy(sorted, grades)
	sort.Ints(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return float64(sorted[mid])
	}
	return float64(sorted[mid-1]+sorted[mid]) / 2.0
}

func (a *MedianStrategy) Name() string {
	return "median"
}

type DropWorstStrategy struct{}

func (d *DropWorstStrategy) Calculate(grades []int) float64 {
	if len(grades) <= 1 {
		return 0
	}
	sorted := make([]int, len(grades))
	copy(sorted, grades)
	sort.Ints(sorted)
	trimmed := sorted[1:]
	sum := 0
	for _, g := range trimmed {
		sum += g
	}
	return float64(sum) / float64(len(trimmed))
}

func (a *DropWorstStrategy) Name() string {
	return "dropWorst"
}
