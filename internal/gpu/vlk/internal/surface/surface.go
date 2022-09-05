package surface

import (
	"fmt"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/internal/gpu/vlk/internal/instance"
	"github.com/go-glx/vgl/shared/vlkext"
)

type Surface struct {
	logger vlkext.Logger
	inst   *instance.Instance

	ref vulkan.Surface
}

func NewSurface(logger vlkext.Logger, inst *instance.Instance, wm vlkext.WindowManager) *Surface {
	surface, err := wm.CreateSurface(inst.Ref())
	if err != nil {
		panic(fmt.Errorf("failed create vulkan surface: %w", err))
	}

	return &Surface{
		logger: logger,
		inst:   inst,
		ref:    surface,
	}
}

func (s *Surface) Free() {
	vulkan.DestroySurface(s.inst.Ref(), s.ref, nil)
	s.logger.Debug("freed: surface")
}

func (s *Surface) Ref() vulkan.Surface {
	return s.ref
}
