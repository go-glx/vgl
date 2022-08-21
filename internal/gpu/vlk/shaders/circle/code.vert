#version 450

layout(set=0, binding = 0) uniform UniformBufferObject {
    mat4 view;
    mat4 proj;
} ubo;

layout(location = 0) in vec2 inPosition;

void main() {
    gl_Position = ubo.view * ubo.proj * vec4(inPosition, 0.0, 1.0);
}
