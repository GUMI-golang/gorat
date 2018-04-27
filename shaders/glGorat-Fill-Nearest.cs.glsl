#version 430

#define MAXUINT16 65535
#define SUPPORT 1

layout(r32i, binding = 0) uniform iimage2D mask;
layout(rgba32f, binding = 1) uniform image2D result;
layout(std430, binding = 2) buffer Bound {
    ivec4 bound;
};
layout(rgba32f, binding = 3) uniform image2D filler;
layout (local_size_x = 1, local_size_y = 1) in;

float fn(float x);

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
    ivec2 fillerpos = pos - bound.xy;
    ivec2 fillersize = imageSize(filler);
    vec2 delta = vec2(
        float(fillersize.x) / float(bound.z - bound.x),
        float(fillersize.y) / float(bound.w - bound.y)
    );
    vec2 scale = vec2(
        max(delta.x, 1),
        max(delta.y, 1)
    );
    vec2 rad = vec2(
        ceil(scale.x * SUPPORT),
        ceil(scale.y * SUPPORT)
    );
    vec2 value = vec2(fillerpos.x * delta.x, fillerpos.y  * delta.y);
    ivec2 rangeH = ivec2(int(value.x-rad.x+0.5), int(value.x+rad.x));
    ivec2 rangeV = ivec2(int(value.y-rad.y+0.5), int(value.y+rad.y));
    //
    vec4 color = vec4(0,0,0,0);
    float sum = 0;
    // Horizontal sum
    for(int x = rangeH.x; x < rangeH.y;x++){
        float normal = x - value.x;
        float res = fn(normal);
        color += imageLoad(filler, ivec2(x, int(value.y))) * res;
        sum += res;
    }
    // Vertical sum
    for(int y = rangeV.x; y < rangeV.y;y++){
        float normal = y - value.y;
        float res = fn(normal);
        color += imageLoad(filler, ivec2(int(value.x), y)) * res;
        sum += res;
    }
    color /= sum;
    float intense = float(imageLoad(mask, pos).x) / MAXUINT16;
    imageStore(result, pos, color * intense);
}

float fn(float x){
    x = max(x, 0);
    x = max(1 - x, 0);
    return x;
}