package Mesh

import "Photon/Math"

type Edge struct {
	Vertex1 *Vertex
	Vertex2 *Vertex
	Sharp   bool
}

func NewEdge(v1 *Vertex, v2 *Vertex, sharp bool) *Edge {
	return &Edge{
		Vertex1: v1,
		Vertex2: v2,
		Sharp:   sharp,
	}
}

func (edge Edge) Vector() Math.Vector3 {
	return edge.Vertex2.Position.Sub(edge.Vertex1.Position)
}
