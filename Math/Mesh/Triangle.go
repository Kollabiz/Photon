package Mesh

import (
	"Photon/Math"
	"Photon/Structs"
)

type Triangle struct {
	V1             Vertex
	V2             Vertex
	V3             Vertex
	Material       *Structs.Material
	TriangleNormal Math.Vector3
}

func (triangle *Triangle) FirstVertPosition() Math.Vector3 {
	return triangle.V1.Position
}

func (triangle *Triangle) SecondVertPosition() Math.Vector3 {
	return triangle.V2.Position
}

func (triangle *Triangle) ThirdVertPosition() Math.Vector3 {
	return triangle.V3.Position
}

func (triangle *Triangle) Edge12() Math.Vector3 {
	return triangle.V2.Position.Sub(triangle.V1.Position)
}

func (triangle *Triangle) Edge23() Math.Vector3 {
	return triangle.V3.Position.Sub(triangle.V2.Position)
}

func (triangle *Triangle) Edge13() Math.Vector3 {
	return triangle.V3.Position.Sub(triangle.V1.Position)
}

func (triangle *Triangle) Middle() Math.Vector3 {
	return triangle.FirstVertPosition().Add(triangle.SecondVertPosition()).Add(triangle.ThirdVertPosition()).FDiv(3)
}

func (triangle *Triangle) ApplyTransform(t *Math.Transform, m Math.Vector3) {
	triangle.V1.Sub(m)
	triangle.V1.MatMul(t.GetScaleMatrix())
	triangle.V1.MatMul(t.GetRotationMatrix())
	triangle.V1.Add(t.GetPosition())
	triangle.V2.Sub(m)
	triangle.V2.MatMul(t.GetScaleMatrix())
	triangle.V2.MatMul(t.GetRotationMatrix())
	triangle.V2.Add(t.GetPosition())
	triangle.V3.Sub(m)
	triangle.V3.MatMul(t.GetScaleMatrix())
	triangle.V3.MatMul(t.GetRotationMatrix())
	triangle.V3.Add(t.GetPosition())
}

func (triangle *Triangle) RecalcNormal() {
	triangle.TriangleNormal = triangle.Edge12().Normalized().Cross(triangle.Edge23().Normalized()).Normalized()
}
