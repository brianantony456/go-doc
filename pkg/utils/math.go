package utils

func Add(a, b int) int {
	return a + b
}

func Substract(a, b int) int {
	return a - b
}

type Calculator interface {
    Add(a, b int) int
}

type RealCalculator struct{}

func (r *RealCalculator) Add(a, b int) int {
    return a + b
}

// math_service.go
type MathService struct {
    calculator Calculator
}

func NewMathService(calculator Calculator) *MathService {
    return &MathService{calculator: calculator}
}

func (s *MathService) AddNumbers(a, b int) int {
    return s.calculator.Add(a, b)
}