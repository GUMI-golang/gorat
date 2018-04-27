#version 430

#define MAXUINT16 65535

layout(r32i, binding = 0) uniform iimage2D mask;
layout(rgba32f, binding = 1) uniform image2D result;
layout(std430, binding = 2) buffer Bound {
    ivec4 bound;
};
layout(rgba32f, binding = 3) uniform image2D filler;
layout (local_size_x = 1, local_size_y = 1) in;

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
    ivec2 fillerpos = pos - bound.xy;
    ivec2 fillersize = imageSize(filler);
    fillerpos.x %= fillersize.x;
    fillerpos.y %= fillersize.y;
    float intense = float(imageLoad(mask, pos).x) / MAXUINT16;
    imageStore(result, pos, imageLoad(filler, fillerpos) * intense);
}
