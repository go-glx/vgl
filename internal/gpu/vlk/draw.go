package vlk

import (
	"fmt"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/shader"
)

func (vlk *VLK) setSurface(ID surfaceID) {
	if ID > maxSurfaces {
		vlk.cont.logger.Error(fmt.Sprintf("failed set surface to %d, max surfaces is %d",
			ID,
			maxSurfaces,
		))
		return
	}

	if len(vlk.drawContext.surfaces) == 0 {
		// create new surface, if not already exist
		vlk.drawContext.surfaces = append(vlk.drawContext.surfaces, newDrawSurface(ID))
		return
	}

	current := vlk.drawContext.surfaces[len(vlk.drawContext.surfaces)-1]
	if current.surfaceID == ID {
		// switched to same surface
		return
	}

	// switched to new surface, break baking
	if len(vlk.drawContext.surfaces) >= maxSurfaces {
		vlk.cont.logger.Error(fmt.Sprintf("failed switch surface, max surfaces switch count %d is reached",
			maxSurfaces,
		))
		return
	}

	vlk.drawContext.surfaces = append(vlk.drawContext.surfaces, newDrawSurface(ID))
}

func (vlk *VLK) drawQueue(shader *shader.Shader, opts DrawOptions, instance shader.InstanceData) {
	// extend groups and bake params
	groupValid := vlk.autoBake(shader, opts)
	if !groupValid {
		return
	}

	// append new instance to current group
	surfInd := len(vlk.drawContext.surfaces) - 1
	grpInd := len(vlk.drawContext.surfaces[surfInd].groups) - 1
	vlk.drawContext.surfaces[surfInd].groups[grpInd].instances = append(
		vlk.drawContext.surfaces[surfInd].groups[grpInd].instances, instance,
	)
}

func (vlk *VLK) autoBake(shader *shader.Shader, opts DrawOptions) bool {
	brakeBaking := false

	// create default surface, if not set
	if len(vlk.drawContext.surfaces) == 0 {
		vlk.setSurface(surfaceIdMainWindow)
	}

	// get current surface
	currSurf := vlk.drawContext.surfaces[len(vlk.drawContext.surfaces)-1]

	// create new draw group, if not empty
	if len(currSurf.groups) == 0 {
		currSurf.groups = append(currSurf.groups, newDrawGroup(shader, opts.PolygonMode))
	}

	currGroup := currSurf.groups[len(currSurf.groups)-1]

	// brake: shader changed
	if currGroup.shader.Meta().ID() != shader.Meta().ID() {
		brakeBaking = true
	}

	// brake: polygon mode changed
	if currGroup.polygonMode != opts.PolygonMode {
		brakeBaking = true
	}

	// brake: blend mode changed (todo)

	if !brakeBaking {
		return true
	}

	// group changes, need create new group with different settings
	if len(currSurf.groups) >= maxGroups {
		vlk.cont.logger.Error(fmt.Sprintf("max drawing groups of %d is reached. Skip rendering..",
			maxGroups,
		))
		return false
	}

	currSurf.groups = append(currSurf.groups, newDrawGroup(shader, opts.PolygonMode))
	return true
}

// draw is actual drawing function, should be called every frame
// just before frame end.
func (vlk *VLK) draw() {
	vlk.drawExecution(vlk.drawContext)
}
