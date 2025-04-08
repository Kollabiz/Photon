package Structs

import (
	"Photon/Math"
	"Photon/Utils"
	"fmt"
	"strconv"
)

// TODO: Remove me
func Debug_TraverseBVHTree(tree *BVHNode) {
	// Basic inorder traversal
	nodes_to_traverse := []BVHNode{*tree}
	max_depth := 10
	i := 0
	Utils.LogWarning("DEBUG_BVH_TRAVERSAL (Scene.go:13 -> Debug_TraverseBVHTree())")
	for len(nodes_to_traverse) > 0 {
		if i > max_depth {
			return
		}
		node := nodes_to_traverse[0]
		nodes_to_traverse = nodes_to_traverse[1:]
		if !node.IsALeaf() {
			i++
			fmt.Print("Node(" + strconv.Itoa(node.NodeID) + "): ")
			if node.Child1.Mesh != nil {
				fmt.Print(" Child1: " + node.Child1.Mesh.MeshName + "(" + strconv.Itoa(node.Child1.NodeID) + ")")
			} else {
				fmt.Print(" Child1: NO_MESH" + "(" + strconv.Itoa(node.Child1.NodeID) + ")")
			}

			if node.Child2.Mesh != nil {
				fmt.Print(" Child2: " + node.Child2.Mesh.MeshName + "(" + strconv.Itoa(node.Child2.NodeID) + ")")
			} else {
				fmt.Print(" Child2: NO_MESH" + "(" + strconv.Itoa(node.Child2.NodeID) + ")")
			}
			fmt.Println()
			nodes_to_traverse = append(nodes_to_traverse, *node.Child1, *node.Child2)
		} else {
			fmt.Println("Node: Leaf (" + node.Mesh.MeshName + ")" + " (" + strconv.Itoa(node.NodeID) + ")")
		}
	}
	Utils.LogSuccess("Done DEBUG_BVH_TRAVERSAL")
}

type Scene struct {
	objects       []Mesh
	lightSources  []LightSource
	camera        *Camera
	sceneSettings *SceneSettings
	baseNode      *BVHNode
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

func (scene *Scene) AddObject(object *Mesh) {
	if scene.containsSameObject(object.MeshName) {
		Utils.LogError("trying to add duplicate mesh")
		panic("there is already an object with the same name")
	}
	scene.objects = append(scene.objects, *object)
}

func (scene *Scene) AddObjectOrCopy(object *Mesh) {
	if scene.containsSameObject(object.MeshName) {
		Utils.LogWarning("mesh with the same name found, the mesh was copied")
		scene.objects = append(scene.objects, *object.Copy())
		return
	}
	scene.objects = append(scene.objects, *object)
}

func (scene *Scene) AddObjectOrLinkedCopy(object *Mesh) {
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

func (scene *Scene) GetObject(name string) *Mesh {
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

func (scene *Scene) RebuildBVH() {
	Utils.Log("rebuilding scene BVH...")
	var nodeLayer []BVHNode
	for j := 0; j < len(scene.objects); j++ {
		nodeLayer = append(nodeLayer, *BVHFromMesh(&scene.objects[j], scene.sceneSettings.KNearestPointRatio))
	}
	// Iterate through nodes until there is but one left
	for len(nodeLayer) > 1 {
		// Take the first node. We don't really care about which node that is
		node := &nodeLayer[0]
		// Then remove it from the array
		nodeLayer = nodeLayer[1:]
		// And find the closest node
		var closestNode BVHNode
		var closestIndex int = -1
		for i := 0; i < len(nodeLayer); i++ {
			if closestIndex == -1 {
				closestNode = nodeLayer[i]
				closestIndex = i
				continue
			}
			if closestNode.AABB.MiddlePoint().Sub(node.AABB.MiddlePoint()).LenSq() > nodeLayer[i].AABB.MiddlePoint().Sub(node.AABB.MiddlePoint()).LenSq() {
				closestNode = nodeLayer[i]
				closestIndex = i
			}
		}
		if closestIndex != -1 {
			nodeLayer = append(nodeLayer[:closestIndex], nodeLayer[closestIndex+1:]...)
		}
		jNode := JoinedNode(node, &closestNode)
		nodeLayer = append(nodeLayer, jNode)
	}
	// When len(nodeLayer) reaches zero, we are done and can save the resulting root node to the scene object
	Utils.LogSuccess("done rebuilding BVH!")
	scene.baseNode = &nodeLayer[0]
}

func (scene *Scene) GetSceneSettings() SceneSettings {
	return *scene.sceneSettings
}

func (scene *Scene) GetLightSources() []LightSource {
	return scene.lightSources
}
