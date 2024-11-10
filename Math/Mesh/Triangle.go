package Mesh

import (
	"Photon/Math"
	"Photon/Structs"
)

type Triangle struct {
	Edge1          *Edge
	Edge2          *Edge
	Edge3          *Edge
	Material       *Structs.Material
	TriangleNormal Math.Vector3
}

func (triangle *Triangle) VertByID(id int) *Vertex {
	if id == 0 {
		return triangle.Edge1.Vertex1
	}
	if id == 1 {
		return triangle.Edge2.Vertex1
	}
	if id == 2 {
		return triangle.Edge3.Vertex1
	}
	return nil
}

func (triangle *Triangle) FirstVertPosition() Math.Vector3 {
	return triangle.Edge1.Vertex1.Position
}

func (triangle *Triangle) FirstVert() *Vertex {
	return triangle.Edge1.Vertex1
}

func (triangle *Triangle) SecondVertPosition() Math.Vector3 {
	return triangle.Edge2.Vertex1.Position
}

func (triangle *Triangle) SecondVert() *Vertex {
	return triangle.Edge2.Vertex1
}

func (triangle *Triangle) ThirdVertPosition() Math.Vector3 {
	return triangle.Edge3.Vertex1.Position
}

func (triangle *Triangle) ThirdVert() *Vertex {
	return triangle.Edge3.Vertex1
}

func (triangle *Triangle) Middle() Math.Vector3 {
	return triangle.FirstVertPosition().Add(triangle.SecondVertPosition()).Add(triangle.ThirdVertPosition()).FDiv(3)
}

func (triangle *Triangle) ApplyTransform(t *Math.Transform, m Math.Vector3) {
	v1 := triangle.FirstVert()
	v2 := triangle.SecondVert()
	v3 := triangle.ThirdVert()
	v1.Sub(m)
	v1.MatMul(t.GetScaleMatrix())
	v1.MatMul(t.GetRotationMatrix())
	v1.Add(t.GetPosition())
	v2.Sub(m)
	v2.MatMul(t.GetScaleMatrix())
	v2.MatMul(t.GetRotationMatrix())
	v2.Add(t.GetPosition())
	v3.Sub(m)
	v3.MatMul(t.GetScaleMatrix())
	v3.MatMul(t.GetRotationMatrix())
	v3.Add(t.GetPosition())
}

func (triangle *Triangle) RecalcNormal() {
	triangle.TriangleNormal = triangle.Edge1.Vector().Normalized().Cross(triangle.Edge2.Vector().Normalized()).Normalized()
}
