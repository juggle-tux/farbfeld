package imagefile

import "image/color"

var (
	imageData = []byte("imagefile" +
		"\x00\x00\x00\x03" +
		"\x00\x00\x00\x03" +
		"\xff\x00\x00\xff\x00\xff\x00\xff\x00\x00\xff\xff" +
		"\x00\x00\xff\xff\x80\x80\x80\x80\x00\xff\x00\xff" +
		"\x00\xff\x00\xff\x00\x00\xff\xff\xff\x00\x00\xff")
	red         = color.RGBA{255, 0, 0, 255}
	green       = color.RGBA{0, 255, 0, 255}
	blue        = color.RGBA{0, 0, 255, 255}
	grayTr      = color.RGBA{128, 128, 128, 128}
	imagePixels = [3][3]color.RGBA{
		{red, green, blue},
		{blue, grayTr, green},
		{green, blue, red},
	}
	red64         = color.RGBA64{65535, 0, 0, 65535}
	green64       = color.RGBA64{0, 65535, 0, 65535}
	blue64        = color.RGBA64{0, 0, 65535, 65535}
	grayTr64      = color.RGBA64{32768, 32768, 32768, 32768}
	imagePixels64 = [3][3]color.RGBA64{
		{red64, green64, blue64},
		{blue64, grayTr64, green64},
		{green64, blue64, red64},
	}
)
