package Structs

import (
	"Photon/Math"
	"Photon/Math/BoundingVolumes"
	"Photon/Math/Mesh"
	"Photon/Utils"
	"fmt"
)

type Scene struct {
	objects       []Mesh.Mesh
	lightSources  []LightSource
	camera        *Camera
	sceneSettings *SceneSettings
	baseNode      *BoundingVolumes.BVHNode
	renderBuffer  *RenderBuffer
}

func NewScene(resolutionX, resolutionY int, FOV float64, settings *SceneSettings) *Scene {
	Utils.Log("creating a Scene...")
	scene := &Scene{}
	Utils.Log("creating a scene Camera...")
	scene.camera = NewCamera(
		Math.Vector3{0, 0, 0},
		Math.Vector3{0, 0, 0},
		Math.Vector2{
			float64(resolutionX),
			float64(resolutionY),
		},
		FOV,
	)
	Utils.LogSuccess("scene Camera created successfully!")
	scene.sceneSettings = settings
	scene.renderBuffer = MakeRenderBuffer(resolutionX, resolutionY)
	Utils.LogSuccess("Scene created successfully!")
	return scene
}

func (scene *Scene) containsSameObject(name string) bool {
	for i := 0; i < len(scene.objects); i++ {
		if scene.objects[i].MeshName == name {
			return true
		}
	}
	return false
}

func (scene *Scene) containsSameLight(lightID int) bool {
	for i := 0; i < len(scene.lightSources); i++ {
		if scene.lightSources[i].GetID() == lightID {
			return true
		}
	}
	return false
}

// Adding objects

func (scene *Scene) AddObject(object *Mesh.Mesh) {
	if scene.containsSameObject(object.MeshName) {
		Utils.LogError("trying to add duplicate mesh")
		panic("there is already an object with the same name")
	}
	scene.objects = append(scene.objects, *object)
}

func (scene *Scene) AddObjectOrCopy(object *Mesh.Mesh) {
	if scene.containsSameObject(object.MeshName) {
		Utils.LogWarning("mesh with the same name found, the mesh was copied")
		scene.objects = append(scene.objects, *object.Copy())
		return
	}
	scene.objects = append(scene.objects, *object)
}

func (scene *Scene) AddObjectOrLinkedCopy(object *Mesh.Mesh) {
	if scene.containsSameObject(object.MeshName) {
		Utils.LogWarning("mesh with the same name found, the mesh was linked")
		scene.objects = append(scene.objects, *object.LinkedCopy())
		return
	}
	scene.objects = append(scene.objects, *object)
}

// Adding Lights

func (scene *Scene) AddLightSource(light LightSource) {
	if scene.containsSameLight(light.GetID()) {
		Utils.LogError("trying to add duplicate light")
		panic("there is already a light source with same ID")
	}
	scene.lightSources = append(scene.lightSources, light)
}

// Getters

func (scene *Scene) GetCamera() *Camera {
	if scene.camera == nil {
		Utils.LogWarning("no camera found, GetCamera() returned <nil>. Make sure to instantiate a camera before rendering!")
	}
	return scene.camera
}

func (scene *Scene) GetObject(name string) *Mesh.Mesh {
	for i := 0; i < len(scene.objects); i++ {
		if scene.objects[i].MeshName == name {
			return &scene.objects[i]
		}
	}
	Utils.LogWarning(fmt.Sprintf("no object with name \"%s\" found", name))
	return nil
}

func (scene *Scene) GetLight(id int) LightSource {
	for i := 0; i < len(scene.lightSources); i++ {
		if scene.lightSources[i].GetID() == id {
			return scene.lightSources[i]
		}
	}
	Utils.LogWarning(fmt.Sprintf("no light with ID \"%d\" found", id))
	return nil
}

func (scene *Scene) rebuildBVH() {
	Utils.Log("rebuilding scene BVH...")
	nodeLayer := make([]BoundingVolumes.BVHNode, len(scene.objects))
	for i := 0; i < len(scene.objects); i++ {
		nodeLayer = append(nodeLayer, *BoundingVolumes.BVHFromMesh(&scene.objects[i], scene.sceneSettings.KNearestPointRatio))
	}
	// Iterate through nodes until there is but one left
	for len(nodeLayer) > 1 {
		// Take the first node. We don't really care about which node that is
		node1 := &nodeLayer[0]
		// Then remove it from the array
		nodeLayer = nodeLayer[1:]
		// And find the closest node
		var closest = &nodeLayer[0]
		var closestDistance = node1.AABB.MiddlePoint().Sub(closest.AABB.MiddlePoint()).LenSq() // LenSq is good enough
		closestIndex := 0                                                                      // Needed solely for the sake of removing that node from the list later on
		for j := 1; j < len(nodeLayer); j++ {
			node := &nodeLayer[j]
			d := node.AABB.MiddlePoint().Sub(node1.AABB.MiddlePoint()).LenSq()
			if d < closestDistance {
				closest = node
				closestDistance = d
				closestIndex = j
			}
		}
		// Remove the closest node from the array
		nodeLayer = append(nodeLayer[:closestIndex], nodeLayer[closestIndex+1:]...)
		// And join selected nodes together
		nodeLayer = append(nodeLayer, BoundingVolumes.JoinedNode(node1, closest))
	}
	// When len(nodeLayer) reaches zero, we are done and can save the resulting root node to the scene object
	Utils.LogSuccess("done rebuilding BVH!")
	scene.baseNode = &nodeLayer[0]
}
