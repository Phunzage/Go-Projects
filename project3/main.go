package main

import (
	rb "github.com/go-vgo/robotgo"
)

func main() {
	bit := rb.CaptureScreen(115, 130, 70, 70)
	defer rb.FreeBitmap(bit)

	img := rb.ToImage(bit)
	rb.Save(img, "test.png")
}
