package Mesh

import (
	"Photon/Math"
	"Photon/Utils"
)

type Mesh struct {
	Transform   *Math.Transform
	Triangles   []Triangle
	MeshName    string
	middle      Math.Vector3
}

func (mesh *Mesh) applyTransform() {
	for i := 0; i < len(mesh.Triangles); i++ {
		t := mesh.Triangles[i]
		t.ApplyTransform(mesh.Transform, mesh.middle)
	}
}

// Transform stuff

func (mesh *Mesh) Move(offset Math.Vector3) {
	mesh.Transform.Move(offset)
	mesh.applyTransform()
	mesh.middle.Add(offset)
}

func (mesh *Mesh) Scale(scale Math.Vector3) {
	mesh.Transform.Resize(scale)
	mesh.applyTransform()
}

func (mesh *Mesh) Rotate(rotation Math.Vector3) {
	mesh.Transform.Rotate(rotation)
	mesh.applyTransform()
}

func (mesh *Mesh) LinkedCopy() *Mesh {
	return &Mesh{
		Transform:   mesh.Transform,
		Triangles:   mesh.Triangles,
		MeshName:    Utils.IncrementName(mesh.MeshName),
		middle:      mesh.middle,
	}
}

func (mesh *Mesh) Copy() *Mesh {
	var trianglesCopy = make([]Triangle, len(mesh.Triangles))
	copy(mesh.Triangles, trianglesCopy)
	var transformCopy = mesh.Transform.Copy()
	return &Mesh{
		Transform:   transformCopy,
		Triangles:   trianglesCopy,
		MeshName:    Utils.IncrementName(mesh.MeshName),
		middle:      mesh.middle,
	}
}
