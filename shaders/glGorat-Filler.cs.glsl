#version 430

#define THERSHOLD 2
#define MAXUINT16 65535



layout(r32i, binding = 0) uniform coherent iimage2D ioutput;
layout (std430, binding = 1) buffer Bound {
    ivec4 bound;
};
layout (local_size_x = 1, local_size_y = 1) in;
void main() {
	// Draw Rect
	int width = imageSize(ioutput).x;
	int y = int(gl_GlobalInvocationID.x);
	int acc = 0;
	int a;
	for(int x = 0; x < width;x ++){
	    acc += imageLoad(ioutput, ivec2(x, y)).x;
	    a = acc;
	    if (a < 0){
	        a = -a;
	    }
	    if (a > MAXUINT16){
	        a = MAXUINT16;
	    }
	    if (a > THERSHOLD){
	        atomicMin(bound.x, x);
	        atomicMin(bound.y, y);
	        atomicMax(bound.z, x);
	        atomicMax(bound.w, y);
	    }else{
	        a = 0;
	    }
	    imageStore(ioutput, ivec2(x, y), ivec4(a, 0,0,0));
	}
 }