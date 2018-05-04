#version 430

layout(rgba32f, binding = 0) uniform image2D to;
layout(rgba32f, binding = 1) uniform image2D a;
layout(rgba32f, binding = 2) uniform image2D b;

layout (local_size_x = 1, local_size_y = 1) in;

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
    vec4 ca = imageLoad(a, pos);
    vec4 cb = imageLoad(a, pos);
    imageStore(to, pos, mix(ca, cb, cb.w));
}
