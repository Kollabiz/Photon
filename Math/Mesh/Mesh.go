package Mesh

import (
	"Photon/Math"
	"Photon/Math/BoundingVolumes"
)

type Mesh struct {
	Transform   *Math.Transform
	Triangles   []Triangle
	BoundingBox *BoundingVolumes.BVHNode
	MeshName    string
	// Cache
	cachedMesh []Triangle
}
