package Mesh

import "Photon/Math"

type Triangle struct {
	Edge1 *Edge
	Edge2 *Edge
	Edge3 *Edge
	// TODO: Add materials
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

func (triangle *Triangle) SecondVertPosition() Math.Vector3 {
	return triangle.Edge2.Vertex1.Position
}

func (triangle *Triangle) ThirdVertPosition() Math.Vector3 {
	return triangle.Edge3.Vertex1.Position
}
