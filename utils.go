package main

func ToArrayCoords(x, y, width, bytesPerPixel int) int {
	return (y*width + x) * bytesPerPixel
}

func ToArrayCoordsYUp(x, y, width, height, bytesPerPixel int) int {
	screenY := height - 1 - y
	return (screenY*width + x) * bytesPerPixel
}
