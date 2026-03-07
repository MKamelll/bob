package main

import (
	"os"
	"slices"

	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	buffer [][]rune
	row    int
	col    int
}

func NewEditor() *Editor {
	return &Editor{buffer: [][]rune{{}}, row: 0, col: 0}
}

func (e *Editor) AddRune(r rune) *Editor {
	e.buffer[e.row] = slices.Insert(e.buffer[e.row], e.col, r)
	e.col++
	return e
}

func (e *Editor) RemoveRune() *Editor {
	e.buffer[e.row] = append(e.buffer[e.row][:e.col-1], e.buffer[e.row][e.col:]...)
	e.col = max(0, e.col-1)
	return e
}

func (e *Editor) HandleBackspace() *Editor {
	if e.col >= 1 {
		e.RemoveRune()
	} else if e.row >= 1 && e.col == 0 {
		after := e.buffer[e.row]
		e.col = len(e.buffer[e.row-1])
		e.buffer[e.row-1] = append(e.buffer[e.row-1], after...)
		e.buffer = append(e.buffer[:e.row], e.buffer[e.row+1:]...)
		e.row--
	}
	return e
}

func (e *Editor) HandleEnter() *Editor {
	before := e.buffer[e.row][:e.col]
	after := e.buffer[e.row][e.col:]
	e.buffer[e.row] = before
	e.buffer = slices.Insert(e.buffer, e.row+1, []rune{})
	e.row++
	e.buffer[e.row] = after
	e.col = 0
	return e
}

func (e *Editor) HandleUp() *Editor {
	e.row = max(0, e.row-1)
	e.col = min(e.col, len(e.buffer[e.row]))
	return e
}

func (e *Editor) HandleDown() *Editor {
	e.row = min(e.row+1, len(e.buffer))
	return e
}

func (e *Editor) HandleKey(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyUp:
		e.HandleUp()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		e.HandleBackspace()
	case tcell.KeyEnter:
		e.HandleEnter()
	case tcell.KeyRune:
		e.AddRune(ev.Rune())
	}
}

func (e *Editor) Draw(s tcell.Screen) {

	s.Clear()

	defaultStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defaultStyle)

	for y, line := range e.buffer {
		for x, char := range line {
			s.SetContent(x, y, char, nil, defaultStyle)
		}
	}

	s.ShowCursor(e.col, e.row)
	s.Show()
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := s.Init(); err != nil {
		panic(err)
	}

	e := NewEditor()

	for {
		e.Draw(s)
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				os.Exit(0)
			}
			e.HandleKey(ev)
		}
	}
}
