package main

import "github.com/veandco/go-sdl2/sdl"
import "github.com/veandco/go-sdl2/ttf"
import "image/color"
import "log"
import "fmt"

var test_text = "foo bar baz quux こんにちは"

type Text struct {
	font *ttf.Font
	text []string 
	row int 
	col int
}

func (t *Text) render(onto *sdl.Surface) {
	for r, line := range t.text {
		var with_cursor string 
		if r == t.row { 
			with_cursor = line[:t.col] + "|" + line[t.col:]
		} else {
			with_cursor = line
		}
		if len(with_cursor) == 0 {
			continue
		}
		text, err := t.font.RenderUTF8Solid(with_cursor, sdl.Color(color.RGBA{0, 255, 0, 0}))
		if err != nil {
			panic(err)
		}
		rect := sdl.Rect{0, int32(r*20), 800, 800}
		text.Blit(nil, onto, &rect)
	}
}

func (t *Text) move_cursor(deltax, deltay int) {
	t.col += deltax
	t.row += deltay
	if t.row <= 0 {
		t.row = 0
	} else if t.row >= len(t.text) {
		t.row = len(t.text) - 1
	}
	if t.col <= 0 {
		t.col = 0
	} else if t.col > len(t.text[t.row]) {
		t.col = len(t.text[t.row])
	}
}

func (t *Text) write(r rune) {
	t.text[t.row] = fmt.Sprintf("%s%c%s", t.text[t.row][:t.col], r, t.text[t.row][t.col:])
	t.col++
}

func (t *Text) delete() {
	if t.col == 0 { // combine with previous row
		t.row--
		t.col = len(t.text[t.row])
		t.text[t.row] += t.text[t.row+1]
		t.text = append(t.text[:t.row+1], t.text[t.row+2:]...)
	} else if t.col > 0 {
		t.text[t.row] = t.text[t.row][:t.col-1] + t.text[t.row][t.col:]
		t.col--
	}
}

func (t *Text) newline() {
	end := make([]string, len(t.text)-t.row+1)
	end[0] = t.text[t.row][t.col:]
	i := copy(end[1:], t.text[t.row+1:])
	log.Println("...", i)
	end = end[:1+i]
	new_text := append(t.text[:t.row], t.text[t.row][:t.col])
	t.text = append(new_text, end...)
	t.col = 0
	t.row++
	for i, line := range t.text {
		log.Println(i, line)
	}
}

func get_rune(key sdl.Keysym) rune {
	keycode := key.Sym
	shift := 0 != ((key.Mod & sdl.KMOD_LSHIFT) | (key.Mod & sdl.KMOD_RSHIFT))
	if (keycode >= 'a') && (keycode <= 'z') && shift {
		keycode = keycode - 'a' + 'A'
	}
	if shift {
		switch keycode {
		case '1': keycode = '!'
		case '2': keycode = '@'
		case '3': keycode = '#'
		case '4': keycode = '$'
		case '5': keycode = '%'
		case '6': keycode = '^'
		case '7': keycode = '&'
		case '8': keycode = '*'
		case '9': keycode = '('
		case '0': keycode = ')'
		case '-': keycode = '_'
		case '=': keycode = '+'
		case '[': keycode = '{'
		case ']': keycode = '}'
		case '\\': keycode = '|'
		case ';': keycode = ':'
		case '\'': keycode = '"'
		case ',': keycode = '<'
		case '.': keycode = '>'
		case '/': keycode = '?'
		case '`': keycode = '~'
		}
	}
	return rune(keycode)
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	defer ttf.Quit()

	window, err := sdl.CreateWindow("ooo", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	rect := sdl.Rect{0, 0, 800, 800}
	surface.FillRect(&rect, 0xffff0000)
	window.UpdateSurface()

	var text Text
	text.font, err = ttf.OpenFont("lsu.ttf", 18)
	if err != nil {
		panic(err)
	}
	text.row = 0
	text.col = 0
	text.text = []string{test_text}
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				keyevent := event.(*sdl.KeyboardEvent)
				if keyevent.GetType() == sdl.KEYDOWN {
					sym := keyevent.Keysym.Sym
					if sym == sdl.K_LEFT {
						text.move_cursor(-1, 0)
					} else if sym == sdl.K_RIGHT {
						text.move_cursor(1, 0)
					} else if sym == sdl.K_UP {
						text.move_cursor(0, -1)
					} else if sym == sdl.K_DOWN {
						text.move_cursor(0, 1)
					} else if sym == sdl.K_BACKSPACE {
						text.delete()
					} else if sym == sdl.K_RETURN {
						text.newline()
					} else if sym == sdl.K_TAB {
						text.write(' ')
						text.write(' ')
					} else if sym <= 127 {
						r := get_rune(keyevent.Keysym)
						text.write(r)
					} else {
						continue
					}
				}
			}
			surface.FillRect(&rect, 0xffff0000)
			text.render(surface)
			window.UpdateSurface()
		}
	}
}