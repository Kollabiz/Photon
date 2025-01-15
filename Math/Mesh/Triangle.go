package Mesh

import (
	"Photon/Math"
	"Photon/Structs"
)

type Triangle struct {
	V1Pos          Math.Vector3
	V2Pos          Math.Vector3
	V3Pos          Math.Vector3
	V1Normal       Math.Vector3
	V2Normal       Math.Vector3
	V3Normal       Math.Vector3
	V1Tex          Math.Vector2
	V2Tex          Math.Vector2
	V3Tex          Math.Vector2
	Smooth         bool
	Material       *Structs.Material
	TriangleNormal Math.Vector3
}

func (triangle *Triangle) Edge12() Math.Vector3 {
	return triangle.V2Pos.Sub(triangle.V1Pos)
}

func (triangle *Triangle) Edge23() Math.Vector3 {
	return triangle.V3Pos.Sub(triangle.V2Pos)
}

func (triangle *Triangle) Edge13() Math.Vector3 {
	return triangle.V3Pos.Sub(triangle.V1Pos)
}

func (triangle *Triangle) Middle() Math.Vector3 {
	return triangle.V1Pos.Add(triangle.V2Pos).Add(triangle.V3Pos).FDiv(3)
}

func (triangle *Triangle) ApplyTransform(t *Math.Transform, m Math.Vector3) {
	triangle.V1Pos = t.GetRotationMatrix().VecMul(t.GetScaleMatrix().VecMul(triangle.V1Pos.Sub(m))).Add(t.GetPosition())
	triangle.V2Pos = t.GetRotationMatrix().VecMul(t.GetScaleMatrix().VecMul(triangle.V2Pos.Sub(m))).Add(t.GetPosition())
	triangle.V3Pos = t.GetRotationMatrix().VecMul(t.GetScaleMatrix().VecMul(triangle.V3Pos.Sub(m))).Add(t.GetPosition())
}

func (triangle *Triangle) RecalcNormal() {
	triangle.TriangleNormal = triangle.Edge12().Normalized().Cross(triangle.Edge23().Normalized()).Normalized()
}

func (triangle *Triangle) InterpolateTexcoords(uv Math.Vector2) Math.Vector2 {
	x, y, z := uv.U, uv.V, 1-uv.U-uv.V
	return triangle.V2Tex.FMul(x).Add(triangle.V3Tex.FMul(y)).Add(triangle.V3Tex.FMul(z))
}

func (triangle *Triangle) InterpolateNormals(uv Math.Vector2) Math.Vector3 {
	if !triangle.Smooth {
		return triangle.TriangleNormal
	}
	x, y, z := uv.U, uv.V, 1-uv.U-uv.V
	return triangle.V1Normal.FMul(x).Add(triangle.V2Normal.FMul(y)).Add(triangle.V3Normal.FMul(z))
}
