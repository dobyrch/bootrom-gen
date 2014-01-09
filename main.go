package main

import (
	"image"
	"fmt"
	"log"
	"os"
	_ "image/png"
)

//TODO: Option for different output modes: byte array, string, or binary
//TODO: Include option for outputting the entire (patched) ROM
//TODO: Enable support for other (lossless) image formats
//TODO: Use more flexible condition for determining which pixels are black
func main() {
	reader, err := os.Open("logo.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()

	for y := bounds.Min.Y; y < 8; y += 4 {
		for x := bounds.Min.X; x < 48; x += 4 {
			fmt.Printf("0x%02X, ", encodeBlock(m, x, y))
			fmt.Printf("0x%02X, ", encodeBlock(m, x, y+2))
		}
	}

	reader, err = os.Open("notice.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err = image.Decode(reader)
        if err != nil {
                log.Fatal(err)
        }
        bounds = m.Bounds()

        for y := bounds.Min.Y; y < 8; y ++ {
		fmt.Printf("0x%02X, ", encodeLine(m, y))
        }

        fmt.Println()

}

func encodeBlock(m image.Image, x int, y int) byte {
	var block int = 0

	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			r, _, _, _ := m.At(x+i, y+j).RGBA()

			if r == 0 {
				block |= 1 << uint(7 - i - 4*j)
			}
		}
	}

	return byte(block)
}

func encodeLine(m image.Image, y int) byte {
	var line int = 0

	for x := 0; x < 8; x++ {
		r, _, _, _ := m.At(x, y).RGBA()

		if r == 0 {
			line |= 1 << uint(7 - x)
		}
	}

	return byte(line)
}
