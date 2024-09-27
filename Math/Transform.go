package Math

type Transform struct {
	Position      Vector3
	Rotation      Mat3
	Scale         Vector3
	scaleMatrix   Mat3
	rotationEuler Vector3
}

func NewTransform(position Vector3, rotation Vector3, scale Vector3) *Transform {
	t := &Transform{}
	t.Position = position
	t.rotationEuler = rotation
	t.Rotation = Mat3Euler(rotation.X, rotation.Y, rotation.Z)
	t.Scale = scale
	t.scaleMatrix = Mat3SeparateScale(scale.X, scale.Y, scale.Z)
	return t
}
