package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type PPMPixel struct {
	red, green, blue int
}

type PPMImage struct {
	width, height int
	data          []PPMPixel
}

// Comment check in case of reading input image
func content_satrt_with_hash(content string) bool {
	if content[0] == '#' {
		return true
	}
	return false
}

//Read input image
func readPPM(file io.Reader) *PPMImage {
	image := PPMImage{}
	reader := bufio.NewReader(file)

	PPMType, _ := reader.ReadString('\n')

	for content_satrt_with_hash(PPMType) {
		PPMType, _ = reader.ReadString('\n')
	}
	if strings.Trim(PPMType, "\n\t") == "P6" {
		size, _ := reader.ReadString('\n')
		for content_satrt_with_hash(size) {
			size, _ = reader.ReadString('\n')
		}
		fmt.Sscanf(size, "%d %d", &image.width, &image.height)
		pixel, _ := reader.ReadString('\n')
		for content_satrt_with_hash(pixel) {
			pixel, _ = reader.ReadString('\n')
		}
		if strings.Trim(pixel, "\n\t") == "255" {
			//fmt.Println(pixel)
		}
	}

	size := image.width * image.height
	image.data = make([]PPMPixel, size)

	for i := 0; i < size; i++ {
		r, _ := reader.ReadByte()
		g, _ := reader.ReadByte()
		b, _ := reader.ReadByte()
		image.data[i] = PPMPixel{red: (int(r) * 4) / 256, green: (int(g) * 4) / 256, blue: (int(b) * 4) / 256}
	}
	return &image
}

//parallel function, distributed in cores
func parallel_task(image *PPMImage, x, j, k, l int, wg *sync.WaitGroup, h []float32) {
	defer wg.Done()
	count := 0
	for i := 0; i < (image.width * image.height); i++ {
		if image.data[i].red == j && image.data[i].green == k && image.data[i].blue == l {
			count++
		}
		h[x] = float32(count) / float32((image.width * image.height))
	}
}

func histogram(image *PPMImage) *[]float32 {
	h := make([]float32, 64)
	wg := sync.WaitGroup{}
	wg.Add(64)

	x := 0
	for j := 0; j <= 3; j++ { //red
		for k := 0; k <= 3; k++ { // green
			for l := 0; l <= 3; l++ { // blue
				go parallel_task(image, x, j, k, l, &wg, h) //parallel code calling
				x = x + 1
			}
		}
	}
	wg.Wait()
	return &h
}

func main() {
	image := readPPM(os.Stdin)
	h := histogram(image)
	for i := 0; i < 64; i++ {
		fmt.Printf("%0.3f ", (*h)[i])
	}
	fmt.Printf("\n")
}
