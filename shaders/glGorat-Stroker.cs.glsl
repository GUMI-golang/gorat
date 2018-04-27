#version 430

#define CLOSETOZERO 0.00001
#define CLOSETOONE 0.99999
#define MAXUINT16 65535.



layout(r32i, binding = 0) uniform coherent iimage2D ioutput;
layout (std430, binding = 1) buffer Points {
    vec2 points[];
};
layout (local_size_x = 1, local_size_y = 1) in;


void main() {
    float tempf;
    ivec2 size = ivec2(imageSize(ioutput).xy);
    ivec2 pos;
    vec2 from, to, tempv2;
    from = points[gl_GlobalInvocationID.x];
    to =  points[(gl_GlobalInvocationID.x + 1)];
    if (from.x < 0 || to.x < 0){
        return;
    }
    float d = 1;
//    float pix;
//    imageStore(ioutput, ivec2(1, 1), vec4(1, 0, 0, 0));
    if (from.y > to.y){
        d = -1;
        tempv2 = from;
        from = to;
        to = tempv2;
    }
    if((to.y - from.y) > CLOSETOZERO){
        int yFrom = int(from.y);
        int yTo = int(to.y + CLOSETOONE);
        float x = from.x;
        float dx = (to.x - from.x) / (to.y - from.y);
        for(int y = yFrom; y < yTo; y++){
            float xNext = x + dx;
            float x0 = x;
            float x1 = xNext;
            if (x0 > x1){
                tempf = x1;
                x1 = x0;
                x0 = tempf;
            }
            if (x0 + 1 >= x1){
                float xmf = float(0.5*(x+xNext)) - x;
                pos = ivec2(clamp(int(x0)+0, 0, size.x), y);
                imageAtomicAdd(ioutput, pos, int((d - d *xmf) * MAXUINT16));
                pos = ivec2(clamp(int(x0)+1, 0, size.x), y);
                imageAtomicAdd(ioutput, pos, int((d *xmf) * MAXUINT16));
            }else{
                float s = 1 / float(x1 - x0);
			    float x0under = x - float(floor(x0));
			    float oneMinusx0under = 1 - x0under;
			    float a0 = float(0.5 * s * oneMinusx0under * oneMinusx0under);
			    float x1under = x1 - float(floor(x1));
			    float am = float(0.5 * s * x1under * x1under);
			    //
			    pos = ivec2(clamp(int(x0)+1, 0, size.x), y);
			    imageAtomicAdd(ioutput, pos, int((d * a0) * MAXUINT16));
//			    pix = imageLoad(ioutput, pos).x;
//                imageStore(ioutput, pos, vec4(pix + (d * a0), 0,0,0));
                if (int(x1) == int(x0)+2) {
                    pos = ivec2(clamp(int(x0)+1, 0, size.x), y);
                    imageAtomicAdd(ioutput, pos,int( (d * (1 - a0 - am)) * MAXUINT16));
//                    pix = imageLoad(ioutput, pos).x;
//                    imageStore(ioutput, pos, vec4(pix + d * (1 - a0 - am), 0,0,0));
                } else {
                    float a1 = float(s * (1.5 - x0under));
                    pos = ivec2(clamp(int(x0)+1, 0, size.x), y);
                    imageAtomicAdd(ioutput, pos, int((d * (a1 - a0)) * MAXUINT16));
//                    pix = imageLoad(ioutput, pos).x;
//                    imageStore(ioutput, pos, vec4(pix + d * (a1 - a0), 0,0,0));

                    float dTimesS = float(d * s);
                    for(int xi = int(x0) + 2; xi < int(x1); xi++ ){
                        pos = ivec2(clamp(xi, 0, size.x), y);
                        imageAtomicAdd(ioutput, pos, int((dTimesS) * MAXUINT16));
//                        pix = imageLoad(ioutput, pos).x;
//                        imageStore(ioutput, pos, vec4(pix + dTimesS, 0,0,0));
                    }
                    float a2 = a1 + float(s*float(int(x1)-int(x0)-3));
                    pos = ivec2(clamp(int(x1), 0, size.x), y);
                    imageAtomicAdd(ioutput, pos, int((d * (1 - a2 - am)) * MAXUINT16));
//                    pix = imageLoad(ioutput, pos).x;
//                    imageStore(ioutput, pos, vec4(pix + d * (1 - a2 - am), 0,0,0));
                }
                pos = ivec2(clamp(int(x1), 0, size.x), y);
                imageAtomicAdd(ioutput, pos, int((d * (d * am)) * MAXUINT16));
//                pix = imageLoad(ioutput, pos).x;
//                imageStore(ioutput, pos, vec4(pix + d * (d * am), 0,0,0));
            }
            //
            x = x + dx;
        }
    }
 }