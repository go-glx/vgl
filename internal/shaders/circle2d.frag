#version 460

// fix float calculations
const float epsilon = 0.0001;

// -----------------

layout(set=0, binding = 1) uniform UniformBufferObject {
    vec2 surfaceSize;
} ubo;

layout(set=1, binding = 0) readonly buffer Props {
    // center of circle
    vec2 center;

    // radius of circle
    float radius;

    // 1.0  - 100% circle is visible
    // 0.1  - 10% of outer circle is visible
    // 0.01 - minumum value
    float thickness;

    // 1.0   - blur all circle
    // 0.005 - default value (minimum smooth)
    // 0.0   - without smooth
    float smoothness;
} c;

// -----------------

layout(location = 0) in vec4 fragColor;

layout(location = 0) out vec4 outColor;

// -----------------

void main() {
    // todo: refactor this mess

    vec2 viewport = vec2(ubo.surfaceSize.x, ubo.surfaceSize.y);
    float aspectRatio = ubo.surfaceSize.x / ubo.surfaceSize.y;
    vec2 uv = (gl_FragCoord.xy / viewport) * 2 -1;
    vec2 line = uv - c.center;
    line.y /= aspectRatio;

    float len = length(line) / 2;
    float thickness = 1 - (c.thickness * c.radius);
    float smoothness = c.smoothness * c.radius;

    // outer
    float circle = smoothstep(c.radius, c.radius - smoothness - epsilon, len);

    // inner
    circle *= smoothstep(1.0 - thickness - smoothness - epsilon, 1.0 - thickness, len);

    outColor = vec4(fragColor.rgb, fragColor.a * circle);
}
