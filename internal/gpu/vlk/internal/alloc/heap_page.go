package alloc

// Buffer Types:
// - index
// - vertex
// - uniform
//
// Layout:
// - Coherent         (must be in host memory)
// - LocalWritable    (must be in device memory, can be rewritten (has mapped memory and coherent space))
// - LocalImmutable   (must be in device memory, write-only, cannot be changed (less space used))
//
// Flags
// - OneTime  (can be overridden with generationID > allocID)
//
// | heap coherent                                          |
//
//             Real memory chunk (some vk buffer)
// | ------------------------------------------------------ |
//

// Api example:
//
// alloc(Uniform, LocalImmutable, size, OneTime & Flag) Allocation
// free(Allocation)
// write(Allocation, []byte(data))

// todo: page is collection of areas
// Responsibility:
// - create new areas when needed
// - physically map vulkan buffers to logical areas
// - contain list of flags and requirements
type h3Page struct {
}
