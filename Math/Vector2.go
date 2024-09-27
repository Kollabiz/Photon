package Math

import "math"

type Vector2 struct {
	U float64
	V float64
}

func ZeroVector2() Vector2 {
	return Vector2{0, 0}
}

// Basic functions

func (v Vector2) FMul(m float64) Vector2 {
	return Vector2{v.U * m, v.V * m}
}

func (v Vector2) FDiv(m float64) Vector2 {
	return Vector2{v.U / m, v.V / m}
}

func (v Vector2) Add(o Vector2) Vector2 {
	return Vector2{v.U + o.U, v.V + o.V}
}

func (v Vector2) Sub(o Vector2) Vector2 {
	return Vector2{v.U - o.U, v.V - o.V}
}

// Vector operations

func (v Vector2) LenSq() float64 {
	return v.U*v.U + v.V*v.V
}

func (v Vector2) Len() float64 {
	return math.Sqrt(v.U*v.U + v.V*v.V)
}
