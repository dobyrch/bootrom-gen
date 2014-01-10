package main

import (
	"errors"
	"fmt"
	"image"
	"os"
	"encoding/binary"
	_ "image/png"
)

//TODO: Option for different output modes: byte array, string, or binary
//TODO: Enable support for other (lossless) image formats
//TODO: Use more flexible condition for determining which pixels are black
func main() {
	imageData := make([]byte, 0, 56)
	var err error

	imageData, err = generateLogo(imageData)
	if (err != nil) {
		fmt.Printf("%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	imageData, err = generateNotice(imageData)
	if (err != nil) {
		fmt.Printf("%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	copy(bootrom[0x22:], []byte {0xa8, 0x00})
	copy(bootrom[0xA8:], imageData)

	// Ordinarily the boot ROM compares its internal copy of the Nintendo
	// logo with a copy of the logo stored in the cartridge; if the two do
	// not match, the boot procedure hangs permanently.  In order to allow
	// for a custom logo, the machine code must be modified to prevent to
	// prevent the boot procedure from hanging.
	bootrom[0xEA] = 0x01
	bootrom[0xFB] = 0x01

	binary.Write(os.Stdout, binary.LittleEndian, bootrom)
}

func generateLogo(imageData []byte) (data []byte, err error) {
	data = imageData

	reader, err := os.Open("logo.png")
	if err != nil {
		return
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return
	}

	bounds := m.Bounds()
	if bounds.Min.X != 0 || bounds.Max.X != 48 ||
	   bounds.Min.Y != 0 || bounds.Max.Y != 8 {
		err = errors.New("logo.png must be 48x8 pixels")
		return
	}

	for y := 0; y < 8; y += 4 {
		for x := 0; x < 48; x += 4 {
			data = append(data, encodeBlock(&m, x, y))
			data = append(data, encodeBlock(&m, x, y+2))
		}
	}

	return
}

func encodeBlock(m *image.Image, x, y int) byte {
	var block int = 0

	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			r, _, _, _ := (*m).At(x+i, y+j).RGBA()

			if r == 0 {
				block |= 1 << uint(7 - i - 4*j)
			}
		}
	}

	return byte(block)
}

func generateNotice(imageData []byte) (data []byte, err error) {
	data = imageData

	reader, err := os.Open("notice.png")
	if err != nil {
		return
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
        if err != nil {
                return
        }

        bounds := m.Bounds()
	if bounds.Min.X != 0 || bounds.Max.X != 8 ||
	   bounds.Min.Y != 0 || bounds.Max.Y != 8 {
		err = errors.New("notice.png must be 8x8 pixels")
		return
	}

        for y := 0; y < 8; y ++ {
		//TODO: Get registration symbol replacement working
		//(Zeroed out for now)
		data = append(data, 0x00)
        }

	return
}

func encodeLine(m *image.Image, y int) byte {
	var line int = 0

	for x := 0; x < 8; x++ {
		r, _, _, _ := (*m).At(x, y).RGBA()

		if r == 0 {
			line |= 1 << uint(7 - x)
		}
	}

	return byte(line)
}

var bootrom []byte = []byte{
	0x31, 0xFE, 0xFF, 0xAF, 0x21, 0xFF, 0x9F, 0x32, 0xCB, 0x7C, 0x20, 0xFB, 0x21, 0x26, 0xFF, 0x0E,
	0x11, 0x3E, 0x80, 0x32, 0xE2, 0x0C, 0x3E, 0xF3, 0xE2, 0x32, 0x3E, 0x77, 0x77, 0x3E, 0xFC, 0xE0,
	0x47, 0x11, 0x04, 0x01, 0x21, 0x10, 0x80, 0x1A, 0xCD, 0x95, 0x00, 0xCD, 0x96, 0x00, 0x13, 0x7B,
	0xFE, 0x34, 0x20, 0xF3, 0x11, 0xD8, 0x00, 0x06, 0x08, 0x1A, 0x13, 0x22, 0x23, 0x05, 0x20, 0xF9,
	0x3E, 0x19, 0xEA, 0x10, 0x99, 0x21, 0x2F, 0x99, 0x0E, 0x0C, 0x3D, 0x28, 0x08, 0x32, 0x0D, 0x20,
	0xF9, 0x2E, 0x0F, 0x18, 0xF3, 0x67, 0x3E, 0x64, 0x57, 0xE0, 0x42, 0x3E, 0x91, 0xE0, 0x40, 0x04,
	0x1E, 0x02, 0x0E, 0x0C, 0xF0, 0x44, 0xFE, 0x90, 0x20, 0xFA, 0x0D, 0x20, 0xF7, 0x1D, 0x20, 0xF2,
	0x0E, 0x13, 0x24, 0x7C, 0x1E, 0x83, 0xFE, 0x62, 0x28, 0x06, 0x1E, 0xC1, 0xFE, 0x64, 0x20, 0x06,
	0x7B, 0xE2, 0x0C, 0x3E, 0x87, 0xE2, 0xF0, 0x42, 0x90, 0xE0, 0x42, 0x15, 0x20, 0xD2, 0x05, 0x20,
	0x4F, 0x16, 0x20, 0x18, 0xCB, 0x4F, 0x06, 0x04, 0xC5, 0xCB, 0x11, 0x17, 0xC1, 0xCB, 0x11, 0x17,
	0x05, 0x20, 0xF5, 0x22, 0x23, 0x22, 0x23, 0xC9, 0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
	0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
	0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
	0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E, 0x3C, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x3C,
	0x21, 0x04, 0x01, 0x11, 0xA8, 0x00, 0x1A, 0x13, 0xBE, 0x20, 0xFE, 0x23, 0x7D, 0xFE, 0x34, 0x20,
	0xF5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xFB, 0x86, 0x20, 0xFE, 0x3E, 0x01, 0xE0, 0x50,
}
