package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type PPMPixel struct {
	red, green, blue int
}

type PPMImage struct {
	width, height int
	data          []PPMPixel
}

func readPPM(rd io.Reader) PPMImage {
	image := PPMImage{}
	reader := bufio.NewReader(rd)

	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	} else {
		if strings.Trim(line, "\n\t") == "P6" {
			size, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Sscanf(size, "%d %d", &image.width, &image.height)
				//fmt.Println("Width: ", image.width)
				//fmt.Println("Height: ", image.height)
				pixel, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				} else {
					if strings.Trim(pixel, "\n\t") == "255" {
						//fmt.Println(pixel)
					}
				}
			}
		}
	}

	size := image.width * image.height
	//fmt.Println("size: ", size)

	image.data = make([]PPMPixel, size)

	for i := 0; i < size; i++ {
		r, _ := reader.ReadByte()
		g, _ := reader.ReadByte()
		b, _ := reader.ReadByte()
		image.data[i] = PPMPixel{red: int(r), green: int(g), blue: int(b)}
	}

	/*for i := 0; i < size; i++ {
		fmt.Println(image.data[i])
	}*/
	return image
}

func Histogram(image PPMImage, h []float32) {
	cols := image.width
	rows := image.height
	n := rows * cols
	//fmt.Println("value of N: ", n)
	for i := 0; i < n; i++ {
		image.data[i].red = (image.data[i].red * 4) / 256
		image.data[i].blue = (image.data[i].blue * 4) / 256
		image.data[i].green = (image.data[i].green * 4) / 256
		//fmt.Println(image.data[i])
	}

	count := 0
	x := 0
	for j := 0; j <= 3; j++ {
		for k := 0; k <= 3; k++ {
			for l := 0; l <= 3; l++ {
				for i := 0; i < n; i++ {
					if image.data[i].red == j && image.data[i].green == k && image.data[i].blue == l {
						count++
					}
				}
				h[x] = float32(count) / float32(n)
				count = 0
				x = x + 1
			}
		}
	}
	//return h
}

func main() {
	image := readPPM(os.Stdin)
	h := make([]float32, 64)
	Histogram(image, h)
	for i := 0; i < len(h); i++ {
		fmt.Printf("%0.3f ", h[i])
	}
	fmt.Printf("\n")
}
