#version 450

// fix float calculations
const float epsilon = 0.0001;

// radius is always 1.0
// transform/scale will be done with matrix
const float radius = 1.0;

// -----------------

layout(location = 0) flat in vec2 inPos;
layout(location = 1) flat in vec2 inRadius;
//
//// 1.0  - 100% circle is visible
//// 0.1  - 10% of outer circle is visible
//// 0.01 - minumum value
//layout(location = 2) in float inThickness;
//
//// 1.0   - blur all circle
//// 0.005 - default value (minimum smooth)
//// 0.0   - without smooth
//layout(location = 3) in float inSmooth;

// -----------------

layout(location = 0) out vec4 outColor;

// -----------------

void main() {
//    vec2 viewport = vec2(960, 480); // todo: variable from UBO
//    vec2 uv = (gl_FragCoord.xy / viewport) * 2 -1;
//
//    float len = length(uv - inPos);

    outColor = vec4(inPos, 0.0, 1.0);
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
