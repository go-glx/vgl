#version 450

layout(location = 0) in vec2 inPosition;
layout(location = 1) in vec4 inColor;
layout(location = 2) in float inThickness;
layout(location = 3) in float inSmooth;

layout(location = 0) out vec2 outPos;
layout(location = 1) out vec4 outColor;
layout(location = 2) out float outThickness;
layout(location = 3) out float outSmooth;

void main() {
    gl_Position = vec4(inPosition, 0.0, 1.0);

    outPos = inPosition;
    outColor = inColor;
    outThickness = inThickness;
    outSmooth = inSmooth;
}
