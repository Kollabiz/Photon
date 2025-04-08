package Math

type Transform struct {
	position      Vector3
	rotation      Mat3
	scale         Vector3
	scaleMatrix   Mat3
	rotationEuler Vector3
}

func NewTransform(position Vector3, rotation Vector3, scale Vector3) *Transform {
	t := &Transform{}
	t.position = position
	t.rotationEuler = rotation
	t.rotation = Mat3Euler(rotation.X, rotation.Y, rotation.Z)
	t.scale = scale
	t.scaleMatrix = Mat3SeparateScale(scale.X, scale.Y, scale.Z)
	return t
}

func (transform *Transform) SetPosition(position Vector3) {
	transform.position = position
}

func (transform *Transform) SetRotation(rotation Vector3) {
	transform.rotationEuler = rotation
	transform.rotation = Mat3Euler(rotation.X, rotation.Y, rotation.Z)
}

func (transform *Transform) SetScale(scale Vector3) {
	transform.scale = scale
	transform.scaleMatrix = Mat3SeparateScale(scale.X, scale.Y, scale.Z)
}

func (transform *Transform) Move(offset Vector3) {
	transform.position = transform.position.Add(offset)
}

func (transform *Transform) Rotate(rotation Vector3) {
	transform.rotationEuler = transform.rotationEuler.Add(rotation)
	transform.rotation = Mat3Euler(
		DegToRad(transform.rotationEuler.X),
		DegToRad(transform.rotationEuler.Y),
		DegToRad(transform.rotationEuler.Z),
	)
}

func (transform *Transform) Resize(scale Vector3) {
	transform.scale = Vector3{
		X: transform.scale.X * scale.X,
		Y: transform.scale.Y * scale.Y,
		Z: transform.scale.Z * scale.Z,
	}
	transform.scaleMatrix = Mat3SeparateScale(transform.scale.X, transform.scale.Y, transform.scale.Z)
}

func (transform *Transform) GetRotation() Vector3 {
	return transform.rotationEuler
}

func (transform *Transform) GetPosition() Vector3 {
	return transform.position
}

func (transform *Transform) GetScale() Vector3 {
	return transform.scale
}

// Matrices

func (transform *Transform) GetScaleMatrix() Mat3 {
	return transform.scaleMatrix
}

func (transform *Transform) GetRotationMatrix() Mat3 {
	return transform.rotation
}

func (transform *Transform) Copy() *Transform {
	return &Transform{
		position:      transform.position,
		rotation:      transform.rotation,
		scale:         transform.scale,
		scaleMatrix:   transform.scaleMatrix,
		rotationEuler: transform.rotationEuler,
	}
}
