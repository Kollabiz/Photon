package Mesh

import "Photon/Math"

type Vertex struct {
	Position          Math.Vector3
	Normal            Math.Vector3
	TextureCoordinate Math.Vector2
}

func NewVertex(position Math.Vector3, normal Math.Vector3, texCoord Math.Vector2) *Vertex {
	return &Vertex{
		Position:          position,
		Normal:            normal,
		TextureCoordinate: texCoord,
	}
}
