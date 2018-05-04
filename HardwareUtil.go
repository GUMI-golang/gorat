package gorat

func Mixing(to, a, b HardwareResult) {
	c := call()
	defer back()
	//
	prog := c.UtilLoadMixing()
	defer c.UtilUnloadMixing(prog)
	w, h := to.Size()
	prog.Use()
	prog.BindResult(0, to)
	prog.BindResult(1, a)
	prog.BindResult(2, b)
	prog.Compute(w, h, 1)
}
