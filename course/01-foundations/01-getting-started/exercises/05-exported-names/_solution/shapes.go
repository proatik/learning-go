package shapes

// Pi is a rough approximation of π used by CircleArea.
const Pi = 3.14159

// square returns x*x.
func square(x float64) float64 {
	return x * x
}

// SquareArea returns the area of a square with the given side.
func SquareArea(side float64) float64 {
	return square(side)
}

// RectangleArea returns the area of a width × height rectangle.
func RectangleArea(width, height float64) float64 {
	return width * height
}

// CircleArea returns the area of a circle with the given radius.
func CircleArea(radius float64) float64 {
	return Pi * square(radius)
}
