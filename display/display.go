package display

import (
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

const (
	WIDTH      = 160
	HEIGHT     = 144
	PIXEL_SIZE = 4
)

type Display struct {
	data     []byte
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

	d.data = make([]byte, HEIGHT*WIDTH*PIXEL_SIZE)
}

func (d *Display) Disconnect(wg *sync.WaitGroup) {
	defer d.renderer.Destroy()
	defer d.window.Destroy()
	wg.Done()
}

func (d *Display) Refresh() {
	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			r := PIXEL_SIZE*(i*WIDTH+ j)
			g := PIXEL_SIZE*(i*WIDTH + j) + 1
			b := PIXEL_SIZE*(i*WIDTH + j) + 2
			a := PIXEL_SIZE*(i*WIDTH + j) + 3
			d.data[r] = byte((i+j)*int(d.cycle))
			d.data[g] = byte((i+j+2)*2*int(d.cycle))
			d.data[b] = byte((i+j+5)*3*int(d.cycle))
			d.data[a] = byte(255)
		}
	}
	d.texture.Update(nil, d.data, WIDTH*4)
	d.renderer.Copy(d.texture, nil, nil)
	d.renderer.Present()
	d.cycle++
}

func (d *Display) Run(wg *sync.WaitGroup) {
	d.renderer.Present()
	return
}
