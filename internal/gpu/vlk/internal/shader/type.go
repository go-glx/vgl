package shader

import (
	"fmt"

	"github.com/vulkan-go/vulkan"
)

type Type string

const (
	TypeVertexBit                 Type = "VertexBit"
	TypeTessellationControlBit    Type = "TessellationControlBit"
	TypeTessellationEvaluationBit Type = "TessellationEvaluationBit"
	TypeGeometryBit               Type = "GeometryBit"
	TypeFragmentBit               Type = "FragmentBit"
	TypeComputeBit                Type = "ComputeBit"
	TypeAllGraphics               Type = "AllGraphics"
	TypeAll                       Type = "All"
	TypeRaygenBitNvx              Type = "RaygenBitNvx"
	TypeAnyHitBitNvx              Type = "AnyHitBitNvx"
	TypeClosestHitBitNvx          Type = "ClosestHitBitNvx"
	TypeMissBitNvx                Type = "MissBitNvx"
	TypeIntersectionBitNvx        Type = "IntersectionBitNvx"
	TypeCallableBitNvx            Type = "CallableBitNvx"
	TypeTaskBitNv                 Type = "TaskBitNv"
	TypeMeshBitNv                 Type = "MeshBitNv"
)

var typesMap = map[Type]vulkan.ShaderStageFlagBits{
	TypeVertexBit:                 vulkan.ShaderStageVertexBit,
	TypeTessellationControlBit:    vulkan.ShaderStageTessellationControlBit,
	TypeTessellationEvaluationBit: vulkan.ShaderStageTessellationEvaluationBit,
	TypeGeometryBit:               vulkan.ShaderStageGeometryBit,
	TypeFragmentBit:               vulkan.ShaderStageFragmentBit,
	TypeComputeBit:                vulkan.ShaderStageComputeBit,
	TypeAllGraphics:               vulkan.ShaderStageAllGraphics,
	TypeAll:                       vulkan.ShaderStageAll,
	TypeRaygenBitNvx:              vulkan.ShaderStageRaygenBitNvx,
	TypeAnyHitBitNvx:              vulkan.ShaderStageAnyHitBitNvx,
	TypeClosestHitBitNvx:          vulkan.ShaderStageClosestHitBitNvx,
	TypeMissBitNvx:                vulkan.ShaderStageMissBitNvx,
	TypeIntersectionBitNvx:        vulkan.ShaderStageIntersectionBitNvx,
	TypeCallableBitNvx:            vulkan.ShaderStageCallableBitNvx,
	TypeTaskBitNv:                 vulkan.ShaderStageTaskBitNv,
	TypeMeshBitNv:                 vulkan.ShaderStageMeshBitNv,
}

func (t Type) VulkanShaderStage() vulkan.ShaderStageFlagBits {
	if stage, exist := typesMap[t]; exist {
		return stage
	}

	panic(fmt.Errorf("unknown shader type: %s", t))
}
