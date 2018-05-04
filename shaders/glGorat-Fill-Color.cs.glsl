#version 430

#define MAXUINT16 65535

layout(r32i, binding = 0) uniform iimage2D mask;
layout(rgba32f, binding = 1) uniform image2D result;
layout (std430, binding = 2) buffer Bound {
    ivec4 resultBound;
};
layout (std430, binding = 3) buffer Color {
    vec4 color;
};

layout (local_size_x = 1, local_size_y = 1) in;

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);

    float intense = float(imageLoad(mask, pos).x) / MAXUINT16;
    ivec2 resultsize = imageSize(result);
    ivec2 resultpos = ivec2(pos + resultBound.xy);
    resultpos.y = resultsize.y - resultpos.y - 1;
    //
    vec4 prev = imageLoad(result, resultpos);
    vec4 need = color * intense;
    prev *= 1 - need.w;
    imageStore(result, resultpos, prev + need);
}
