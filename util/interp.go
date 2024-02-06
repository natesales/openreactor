package util

// Point represents a 2D coordinate pair (x, y)
type Point struct {
	X, Y float64
}

// linearInterpolation performs a linear interpolation between two points
func linearInterpolation(p1, p2 Point, x float64) float64 {
	if p1.X == p2.X {
		// Avoid division by zero
		return p1.Y
	}
	return p1.Y + (x-p1.X)*(p2.Y-p1.Y)/(p2.X-p1.X)
}

// Interpolate interpolates between a list of points at a given x value
func Interpolate(x float64, points []Point) float64 {
	for i := 1; i < len(points); i++ {
		if x <= points[i].X {
			return linearInterpolation(points[i-1], points[i], x)
		}
	}
	// If x is beyond the last point, return the y value of the last point
	return points[len(points)-1].Y
}
