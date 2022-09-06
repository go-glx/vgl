#version 450

// fix float calculations
const float epsilon = 0.0001;

// -----------------

layout(set=0, binding = 1) uniform UniformBufferObject {
    vec2 surfaceSize;
} ubo;

struct Circle {
    // >0.9999  - will be discard (because not visible)
    // 0.1      - small hole of 10% in center
    // 0.01     - micro hole of 1% in center
    float holeRadius;

    // 1.0   - blur all circle
    // 0.005 - default value (minimum smooth)
    // 0.0   - without smooth
    float smoothness;
};

layout(set=1, binding = 0) readonly buffer Props {
    Circle[] circles;
} props;

// -----------------

layout(location = 0) in vec4 fragColor;
layout(location = 1) flat in uint instanceID;
layout(location = 2) in vec2 UV;

layout(location = 0) out vec4 outColor;

// -----------------

void main() {
    Circle c = props.circles[instanceID];

    float len = length(UV);
    float thickness = 1 - c.holeRadius;

    // outer
    float circle = smoothstep(1, 1 - c.smoothness - epsilon, len);

    // inner
    circle *= smoothstep(1 - thickness - c.smoothness - epsilon, 1 - thickness, len);

    outColor = vec4(fragColor.rgb, fragColor.a * circle);
}

