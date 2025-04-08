package PhotonMapping

import (
	"Photon/Structs"
	"Photon/Utils"
	"strconv"
	"sync"
	"sync/atomic"
)

type PhotonThreadHandler struct {
	maxThreads int
	busy       bool
	wg         sync.WaitGroup
	phCount    atomic.Int64
}

func (handler *PhotonThreadHandler) AllocThreads(scene *Structs.Scene, pointCloud *CameraPointCloud, env *Environment) {
	Utils.Log("Allocating threads for async photon mapping...")
	if handler.busy {
		Utils.LogError("Trying to allocate threads while busy!")
		return
	}

	if len(pointCloud.NonCameraPoints) == 0 {
		Utils.LogWarning("No intersection points. Aborting")
		return
	}

	handler.busy = true
	lsCount := len(scene.GetLightSources())
	for i := 0; i < handler.maxThreads; i++ {
		lIdx := 0
		if lsCount != 0 {
			lIdx = i % lsCount
		}
		go AsyncPhotonCast(scene, env, pointCloud, lIdx, handler)
		handler.wg.Add(1)
		Utils.LogSuccess("Allocated thread #" + strconv.Itoa(i))
	}
}

func (handler *PhotonThreadHandler) UnsafeFinish() {
	Utils.LogWarning("Finishing async photon mapping without thread exit checks")
	handler.busy = false
}

func (handler *PhotonThreadHandler) Finish() {
	Utils.Log("Finishing async photon mapping")
	handler.busy = false
	Utils.Log("Waiting for the threads to exit")
	handler.wg.Wait()
	Utils.LogSuccess("All threads exited")
}

func (handler *PhotonThreadHandler) IsFinished() bool {
	return !handler.busy
}
