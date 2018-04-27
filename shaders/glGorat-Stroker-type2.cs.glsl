#version 430

#define CLOSETOZERO 0.00001
#define CLOSETOONE 0.99999
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
    ivec2 size = imageSize(ioutput);
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
    int yTo = min(int(ceil(to.y)), size.y);
    //
    for(int y = yFrom; y < yTo;y++){
        float deltaY = min(float(y + 1), to.y) - max(float(y), from.y);
        float xNext = xCurr + deltaY * deltaXY;
        if(y >= 0){
            float d = deltaY * dir;
            float x0 = min(xCurr, xNext), x1 = max(xCurr, xNext);
            float x0floor =floor(x0), x1ceil =ceil(x1);
            int x0i = int(x0floor), x1i = int(x1ceil);
            if (x1i <= x0i + 1){
                float xmf = 0.5 * (xCurr + xNext) - x0floor;
                imageAtomicAdd(ioutput, ivec2(iclamp(x0i, 0, size.x), y), int((d - d *xmf) * MAXUINT16));
                imageAtomicAdd(ioutput, ivec2(iclamp(x0i + 1, 0, size.x), y), int((d *xmf) * MAXUINT16));
            }else{
                float s = 1 / (x1 - x0);
                float x0f = x0 - x0floor;
                float oneMinusX0f = 1 - x0f;
                float a0 = 0.5 * s * oneMinusX0f * oneMinusX0f;
                float x1f = x1 - x1ceil + 1;
                float am = 0.5 *s * x1f* x1f;

                imageAtomicAdd(ioutput, ivec2(iclamp(x0i, 0, size.x), y), int((d * a0) * MAXUINT16));
                if ( x1i == x0i + 2){
                    imageAtomicAdd(ioutput, ivec2(iclamp(x0i + 1, 0, size.x), y), int((d * (1-a0-am)) * MAXUINT16));
                }else{
                    float a1 = s * (1.5 - x0f);
                    imageAtomicAdd(ioutput, ivec2(iclamp(x0i + 1, 0, size.x), y), int((d * (a1 - a0)) * MAXUINT16));
                    float dTimeS = d * s;
                    for (int xi = x0i; xi < x1i -1; xi++){
                        imageAtomicAdd(ioutput, ivec2(iclamp(xi, 0, size.x), y), int(dTimeS * MAXUINT16));
                    }
                    float a2 = a1 + s * float(x1i - x0i - 3);
                    imageAtomicAdd(ioutput, ivec2(iclamp(x1i - 1, 0, size.x), y), int((d * (1 - a2 - am)) * MAXUINT16));
                }
                imageAtomicAdd(ioutput, ivec2(iclamp(x1i, 0, size.x), y), int(d * am * MAXUINT16));
            }

        }
        xCurr = xNext;
    }
 }
int iclamp(int a, int min, int max){
    return min(max(a, min), max);
}