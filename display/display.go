package display

import (
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

const (
	WIDTH  = 160
	HEIGHT = 144
)

type Display struct {
	data     [HEIGHT][WIDTH]int
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	cycle    uint64
}

func (d *Display) Init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("yesSGMB", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		256, 256, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	d.window = window

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	d.renderer = renderer

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, WIDTH, HEIGHT)
	if err != nil {
		panic(err)
	}
	d.texture = texture
}

func (d *Display) Disconnect(wg *sync.WaitGroup) {
	defer d.renderer.Destroy()
	defer d.window.Destroy()
	wg.Done()
}

func (d *Display) Refresh() {
	var pixels []byte
	for i := 0; i < WIDTH*HEIGHT; i++ {
		r := byte(i * int(d.cycle))
		g := byte(i * 2 * int(d.cycle))
		b := byte(i * 3 * int(d.cycle))
		pixels = append(pixels, r, g, b, 255)
	}
	d.texture.Update(nil, pixels, WIDTH*4)
	d.renderer.Copy(d.texture, nil, nil)
	d.renderer.Present()
	d.cycle++
}

func (d *Display) Run(wg *sync.WaitGroup) {
	d.renderer.Present()
	return
}
