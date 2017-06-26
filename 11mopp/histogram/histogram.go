package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
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

func readPPM(file io.Reader) PPMImage {
	image := PPMImage{}
	reader := bufio.NewReader(file)

	PPMType, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	} else {
		if strings.Trim(PPMType, "\n\t") == "P6" {
			size, err := reader.ReadString('\n')
			//fmt.Println("Size: ", size)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Sscanf(size, "%d %d", &image.width, &image.height)
				//fmt.Println(image.width, image.height)
				if image.width == 0 || image.height == 0 {
					extra, _ := reader.ReadString('\n')
					fmt.Sscanf(extra, "%d %d", &image.width, &image.height)
					//fmt.Println("Inner Loop: ", image.width, image.height)
				}
				pixel, err := reader.ReadString('\n')
				//fmt.Println("Pixel:", pixel)
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
	image.data = make([]PPMPixel, size)

	for i := 0; i < size; i++ {
		r, _ := reader.ReadByte()
		g, _ := reader.ReadByte()
		b, _ := reader.ReadByte()
		image.data[i] = PPMPixel{red: int(r), green: int(g), blue: int(b)}
	}
	return image
}

//parallel function, distributed in cores
func split_task(image PPMImage, x, j, k, l int, wg *sync.WaitGroup, h []float32) {
	defer wg.Done()
	//fmt.Println("From Split task: ")
	count := 0
	for i := 0; i < (image.width * image.height); i++ {
		if image.data[i].red == j && image.data[i].green == k && image.data[i].blue == l {
			count++
		}
		h[x] = float32(count) / float32((image.width * image.height))
	}
}

func histogram(image PPMImage, h []float32) {

	for i := 0; i < (image.width * image.height); i++ {
		image.data[i].red = (image.data[i].red * 4) / 256
		image.data[i].blue = (image.data[i].blue * 4) / 256
		image.data[i].green = (image.data[i].green * 4) / 256
	}
	wg := sync.WaitGroup{}
	wg.Add(64)

	x := 0
	for j := 0; j <= 3; j++ {
		for k := 0; k <= 3; k++ {
			for l := 0; l <= 3; l++ {
				go split_task(image, x, j, k, l, &wg, h) //parallel code calling
				x = x + 1
			}
		}
	}
	wg.Wait()
}

func main() {
	image := readPPM(os.Stdin)
	h := make([]float32, 64)
	histogram(image, h)
	for i := 0; i < len(h); i++ {
		fmt.Printf("%0.3f ", h[i])
	}
	fmt.Printf("\n")
}
