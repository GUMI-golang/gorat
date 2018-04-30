#version 430

#define CLOSETOZERO 0.00001
#define MAXUINT16 65535



layout(r32i, binding = 0) uniform coherent iimage2D ioutput;
layout (std430, binding = 1) buffer Points {
    vec2 points[];
};
layout (std430, binding = 2) buffer Bound {
    ivec4 bound;
};
layout (local_size_x = 1, local_size_y = 1) in;

int iclamp(int a, int min, int max);

void main() {
    vec2 from, to, tempv2;
    from = points[gl_GlobalInvocationID.x];
    to =  points[(gl_GlobalInvocationID.x + 1)];
    if (isnan(from.x )|| isnan(to.x)){
        return;
    }
    float dir = 1;
//    float pix;
//    imageStore(ioutput, ivec2(1, 1), vec4(1, 0, 0, 0));
    if (from.y > to.y){
        dir = -1;
        tempv2 = from;
        from = to;
        to = tempv2;
    }
    if(to.y - from.y < CLOSETOZERO) {
        return;
    }
    tempv2 = (to - from);
    float deltaXY = tempv2.x / tempv2.y;
    float xCurr = from.x;
    int yFrom = int(floor(from.y));
    int yTo = int(floor(to.y));
    //
    for(int y = yFrom; y < yTo;y++){
        float deltaY = min(float(y + 1), to.y) - max(float(y), from.y);
        float xNext = xCurr + deltaY * deltaXY;
        float x0 = min(xCurr, xNext), x1 = max(xCurr, xNext);
        float xDiff = 1 / (ceil(x1) - floor(x0) + 1);
        for(int x = int(floor(x0)); x <= int(ceil(x1));x++){
            imageAtomicAdd(ioutput, ivec2(iclamp(x, bound.x, bound.z), iclamp(y, bound.y, bound.w)), int(xDiff*dir * MAXUINT16));
        }
        xCurr = xNext;
    }
 }
int iclamp(int a, int imin, int imax){
    return min(max(a, imin), imax);
}