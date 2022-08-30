#version 450

layout(set=0, binding = 0) uniform UniformBufferObject {
    mat4 view;
    mat4 proj;
} ubo;

layout(location = 0) in vec2 inPosition;
layout(location = 1) in vec4 inColor;

layout(location = 0) out vec4 outColor;
layout(location = 1) out flat uint outInstanceID;

void main() {
    gl_Position = ubo.view * ubo.proj * vec4(inPosition, 0.0, 1.0);
    outColor = inColor;
    outInstanceID = gl_InstanceIndex;
}
