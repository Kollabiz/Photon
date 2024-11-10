package Math

import "math"

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func ZeroVector3() Vector3 {
	return Vector3{0, 0, 0}
}

func InfiniteVector3() Vector3 {
	return Vector3{math.Inf(1), math.Inf(1), math.Inf(1)}
}

func NegativeInfiniteVector3() Vector3 {
	return Vector3{math.Inf(-1), math.Inf(-1), math.Inf(-1)}
}

// Simple operations

func (v Vector3) FMul(m float64) Vector3 {
	return Vector3{v.X * m, v.Y * m, v.Z * m}
}

func (v Vector3) FDiv(m float64) Vector3 {
	return Vector3{v.X / m, v.Y / m, v.Z / m}
}

func (v Vector3) IMul(m int) Vector3 {
	return Vector3{v.X * float64(m), v.Y * float64(m), v.Z * float64(m)}
}

func (v Vector3) Mul(o Vector3) Vector3 {
	return Vector3{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}

func (v Vector3) IDiv(m int) Vector3 {
	return Vector3{v.X / float64(m), v.Y / float64(m), v.Z / float64(m)}
}

func (v Vector3) Abs() Vector3 {
	return Vector3{
		math.Abs(v.X),
		math.Abs(v.Y),
		math.Abs(v.Z),
	}
}

// Vector operations

func (v Vector3) LenSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector3) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) Normalized() Vector3 {
	l := v.Len()
	return Vector3{
		X: v.X / l,
		Y: v.Y / l,
		Z: v.Z / l,
	}
}

func (v Vector3) Dot(o Vector3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vector3) Cross(o Vector3) Vector3 {
	return Vector3{
		X: v.Y*o.Z - v.Z*o.Y,
		Y: v.Z*o.X - v.X*o.Z,
		Z: v.X*o.Y - v.Y*o.X,
	}
}

func (v Vector3) Inverse() Vector3 {
	return Vector3{-v.X, -v.Y, -v.Z}
}

func (v Vector3) Add(o Vector3) Vector3 {
	return Vector3{
		X: v.X + o.X,
		Y: v.Y + o.Y,
		Z: v.Z + o.Z,
	}
}

func (v Vector3) Sub(o Vector3) Vector3 {
	return Vector3{
		X: v.X - o.X,
		Y: v.Y - o.Y,
		Z: v.Z - o.Z,
	}
}

func InterpolateVector3(f, s Vector3, t float64) Vector3 {
	return s.FMul(t).Add(f.FMul(1 - t))
}

// Comparing vectors

func (v Vector3) Equal(o Vector3) bool {
	if v.X == o.X && v.Y == o.Y && v.Z == o.Z {
		return true
	}
	return false
}
