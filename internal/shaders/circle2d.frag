#version 450

// fix float calculations
const float epsilon = 0.0001;

// -----------------

layout(set=0, binding = 1) uniform UniformBufferObject {
    vec2 surfaceSize;
} ubo;

layout(set=1, binding = 0) uniform Props {
    // center of circle
    vec2 center;

    // 1.0  - 100% circle is visible
    // 0.1  - 10% of outer circle is visible
    // 0.01 - minumum value
    float thickness;

    // 1.0   - blur all circle
    // 0.005 - default value (minimum smooth)
    // 0.0   - without smooth
    float smoothness;
} props;

// -----------------

layout(location = 0) out vec4 outColor;

// -----------------

void main() {
    vec2 viewport = vec2(ubo.surfaceSize.x, ubo.surfaceSize.y);
    vec2 uv = (gl_FragCoord.xy / viewport) * 2 -1;

    float len = length(uv - props.center);

    outColor = vec4(vec3(len), 1);
}

//void main() {
//    // circle
//    float len = length(inLocalPosition.xy);
//    float circle = len;
//
//    // outer
//    circle = smoothstep(radius, radius-inSmooth-epsilon, len);
//
//    // inner
//    circle *= smoothstep(1.0-inThickness-inSmooth-epsilon, 1.0-inThickness, len);
//
//    outColor = vec4(inColor.rgb, circle * inColor.a);
//}
