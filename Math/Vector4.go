package Math

type Vector4 struct {
	X float64
	Y float64
	Z float64
	W float64
}

func ZeroVector4() Vector4 {
	return Vector4{0, 0, 0, 0}
}

func (vec Vector4) Add3(other Vector3) Vector4 {
	return Vector4{
		X: vec.X + other.X,
		Y: vec.Y + other.Y,
		Z: vec.Z + other.Z,
		W: vec.W + 1,
	}
}

func (vec Vector4) ToVec3() Vector3 {
	return Vector3{
		X: vec.X / vec.W,
		Y: vec.Y / vec.W,
		Z: vec.Z / vec.W,
	}
}
