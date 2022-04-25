package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width     = 500
	height    = 200
	num_iters = 2000
)

func get_next_color_index(current_color_index int, colors [3]color.RGBA) int {
	var out int = 0
	if current_color_index == len(colors)-1 {
		out = 0
	} else {
		out = current_color_index + 1
	}
	return out
}

func get_moore_neighborhood(idx_h int, idx_w int, width int, height int) [][2]int {
	var neighbours_indices [][2]int
	var checkList []int
	moore_neighborhood := [8][2]int{
		{idx_w - 1, idx_h - 1},
		{idx_w - 1, idx_h},
		{idx_w - 1, idx_h + 1},
		{idx_w, idx_h - 1},
		{idx_w, idx_h + 1},
		{idx_w + 1, idx_h - 1},
		{idx_w + 1, idx_h},
		{idx_w + 1, idx_h + 1},
	}

	for i := 0; i < len(moore_neighborhood); i++ {
		checkList = append(checkList, 1)
	}
	if idx_w == width-1 {
		checkList[7] = 0
		checkList[6] = 0
		checkList[5] = 0
	}
	if idx_w == 0 {
		checkList[0] = 0
		checkList[1] = 0
		checkList[2] = 0
	}
	if idx_h == height-1 {
		checkList[7] = 0
		checkList[4] = 0
		checkList[2] = 0
	}
	if idx_h == 0 {
		checkList[5] = 0
		checkList[3] = 0
		checkList[0] = 0
	}
	for i := 0; i < len(moore_neighborhood); i++ {
		if checkList[i] == 1 {
			neighbours_indices = append(neighbours_indices, moore_neighborhood[i])
		}
	}
	return neighbours_indices
}

func take_a_step(img_array [width][height]int, colors [3]color.RGBA) [width][height]int {
	var new_img_array [width][height]int
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			new_color_in_cycle_index := get_next_color_index(img_array[x][y], colors)
			neighbours_indices := get_moore_neighborhood(y, x, width, height)
			var num_neighbours_in_next_cycle int = 0
			for i := 0; i < len(neighbours_indices); i++ {
				if img_array[neighbours_indices[i][0]][neighbours_indices[i][1]] == new_color_in_cycle_index {
					num_neighbours_in_next_cycle += 1
				}
			}
			if num_neighbours_in_next_cycle > 2 {
				new_img_array[x][y] = new_color_in_cycle_index
			} else {
				new_img_array[x][y] = img_array[x][y]
			}
		}
	}
	return new_img_array
}

func main() {

	var err error
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}

	window, err := sdl.CreateWindow(
		"cyclic cellular automaton", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	red := color.RGBA{255, 0, 0, 0}
	green := color.RGBA{0, 255, 0, 0}
	blue := color.RGBA{0, 0, 255, 0}
	colors := [3]color.RGBA{red, green, blue}

	rand.Seed(time.Now().UnixNano())

	var img_array [width][height]int
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img_array[x][y] = rand.Intn(len(colors))
		}
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			surface.Set(x, y, colors[img_array[x][y]])
		}
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		for iter := 0; iter < num_iters; iter++ {
			img_array = take_a_step(img_array, colors)
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {
					surface.Set(x, y, colors[img_array[x][y]])
				}
			}
			window.UpdateSurface()
			fmt.Printf("step " + strconv.Itoa(iter) + "/" + strconv.Itoa(num_iters) + "\n")
		}
		running = false
	}

}
