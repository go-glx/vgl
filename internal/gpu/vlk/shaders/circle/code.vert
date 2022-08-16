#version 450

layout(location = 0) in vec2 inPosition;
layout(location = 1) in vec2 inRadius;

layout(location = 0) flat out vec2 outPos;
layout(location = 1) flat out vec2 outRadius;

vec2 quad[6] = vec2[] (
    vec2(0.75,-0.75),
    vec2(0.75,0.75),
    vec2(-0.75,0.75),
    vec2(-0.75,0.75),
    vec2(-0.75,-0.75),
    vec2(0.75,-0.75)
);

void main() {
    gl_Position = vec4(quad[gl_VertexIndex], 0.0, 1.0);

    outPos = inPosition;
    outRadius = inRadius;
}
