package vulkan

import "time"

// GPU timeout for render. After this, app will be crashed
const swapChainTimeout = time.Second * 10

// How many frames can be failed continuously before crash
const maxPresetFails = 100

// Capacity of each vertex buffer page
const vertexBufferSize = 2048 // 2048 Bytes
const indexBufferSize = 65536 // 65 MB

// always equal size of UBO ( 4 * mat4 )
// UBO always has 4 matrices:
//  1: projection
//  2: view
//  3: model
//  4: - reserved -
//
// !! UBO support up to 4*mat4 slots ( 256 bytes maximum ) !!
//
// each matrix is struct with ( 4 * Vec4 ):
//    [ 1 0 0 0 ]
//    [ 0 1 0 0 ]
//    [ 0 0 1 0 ]
//    [ 0 0 0 1 ]
//
// Sizes:
//  - SizeOfVec4 = 16 ( 4 bytes (float32) * 4 )
//  - SizeOfMat4 = 64 (16 * 4)
//  - SizeOfUBO = 256 (64 * 4)
const (
	uniformBufferSize = 256
)

const (
	shaderEntryPoint = "main"
	shaderTypeVert   = ".vert"
	shaderTypeFrag   = ".frag"
)
