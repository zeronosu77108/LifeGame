package main

import (
    // "fmt"
    "math/rand"
    "github.com/nsf/termbox-go"
    "time"
)

type board struct {
    board [][]bool
    w int
    h int
}

func newboard(w, h int) *board {
    b := new(board)
    b.w = w
    b.h = h
    b.board = make([][]bool, h)

    rand.Seed(time.Now().UnixNano())

    for i:=0; i<h; i++ {
        b.board[i] = make([]bool, w)
        for j:=0; j<w; j++ {
            if rand.Intn(7) == 0 {
                b.board[i][j] = true
            }
        }
    }
    return b
}

func  (b *board) render() {
    for i:=0; i<b.h; i++ {
        for j:=0; j<b.w; j++{
            if b.board[i][j] == true {
                termbox.SetCell(i,j,'■', termbox.ColorDefault, termbox.ColorDefault)
            } else {
                termbox.SetCell(i,j,'□', termbox.ColorDefault, termbox.ColorDefault)
            }
        }
    }
    termbox.Flush()
}

func (b *board) update() {
    board := make([][]bool, b.h)
    for i:=0; i<b.h; i++ {
        board[i] = make([]bool, b.w)
        for j:=0; j<b.w; j++ {
            board[i][j] = b.is_alive(i,j)
        }
    }
    b.board = board
}

func (b *board) count(i,j int) int {
    if i < 0 || b.h <= i {
        return 0
    } else if j < 0 || b.w <= j {
        return 0
    } else if b.board[i][j] {
        return 1
    } else {
        return 0
    }
}

func (b *board) is_alive(i,j int) bool {
    c := b.count(i-1,j-1) +
        b.count(i-1,j) +
        b.count(i-1,j+1) +
        b.count(i,j-1) +
        b.count(i,j+1) +
        b.count(i+1,j-1) +
        b.count(i+1,j) +
        b.count(i+1,j+1)

    if c == 2 {
        return b.board[i][j]
    } else if c ==3 {
        return true
    } else {
        return false
    }
}


func keyEventInter(kch chan termbox.Key) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			kch <- ev.Key
		default:
		}
	}
}

func timerInter(tch chan bool) {
	for {
		tch <- true
		time.Sleep(1500 * time.Millisecond)
		// time.Sleep(1500 * time.Millisecond)
	}
}



func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

    board := newboard(20,90)
    board.render()

    keyCh := make(chan termbox.Key)
    timerCh := make(chan bool)

    go keyEventInter(keyCh)
    go timerInter(timerCh)

    for {
        select {
        case key := <-keyCh:
            switch key{
            case termbox.KeyEsc, termbox.KeyCtrlC:
                return
            }
        case <-timerCh:
            board.update()
            board.render()
        }
    }
}
