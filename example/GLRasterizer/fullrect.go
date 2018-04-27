package GLRasterizer

var FULLRECTUV = [][5]float32{
	{-1., -1., 0, 0, 1},
	{+1., -1., 0, 1, 1},
	{+1., +1., 0, 1, 0},
	{-1., +1., 0, 0, 0},
	{-1., -1., 0, 0, 1},
	{+1., +1., 0, 1, 0},
}
var FULLRECTUVSIZE = len(FULLRECTUV) * 5 * 4