package Mesh

import "Photon/Math"

type Vertex struct {
	Position          Math.Vector3
	Normal            Math.Vector3
	TextureCoordinate Math.Vector2
	Sharp             bool
}

func NewVertex(position Math.Vector3, normal Math.Vector3, texCoord Math.Vector2) *Vertex {
	return &Vertex{
		Position:          position,
		Normal:            normal,
		TextureCoordinate: texCoord,
		Sharp:             false,
	}
}

func (v *Vertex) MatMul(m Math.Mat3) {
	v.Position = m.VecMul(v.Position)
	v.Normal = m.VecMul(v.Normal)
}

func (v *Vertex) Add(o Math.Vector3) {
	v.Position = v.Position.Add(o)
}

func (v *Vertex) Sub(o Math.Vector3) {
	v.Position = v.Position.Sub(o)
}
