package main

import (
	"errors"
	"io/ioutil"
)

type Game struct {
	headers map[string]string
	moves   string
}

type PgnReader struct {
	games         []Game
	state         pgnState
	headers       map[string]string
	moves         string
	currentHeader string
	currentValue  string
}

func NewPgnReader() *PgnReader {

	return &PgnReader{
		games:         make([]Game, 0),
		state:         LINE_START,
		headers:       make(map[string]string),
		moves:         "",
		currentHeader: "",
		currentValue:  "",
	}
}

type pgnState int

const (
	LINE_START pgnState = iota
	IN_HEADER_NAME
	OUT_HEADER_NAME
	IN_HEADER_VALUE
	OUT_HEADER_VALUE
	LINE_END
	IN_MOVES
)

func (r *PgnReader) Read(p []byte) (int, error) {
	for i, c := range p {
		switch r.state {
		case LINE_START:
			switch c {
			case '[':
				r.state = IN_HEADER_NAME
			case '\n':
				if r.moves == "" {
					r.state = LINE_START
				} else {
					game := Game{
						headers: r.headers,
						moves:   r.moves}
					r.games = append(r.games, game)
					r.Reset()
					return i, nil
				}
			default:
				r.state = IN_MOVES
				r.moves += string(c)
			}
		case IN_HEADER_NAME:
			if c == ' ' {
				r.state = OUT_HEADER_NAME
			} else {
				r.currentHeader += string(c)
			}
		case OUT_HEADER_NAME:
			if c != '"' {
				return 0, errors.New("Bad pgn form: " + string(c))
			}
			r.state = IN_HEADER_VALUE
		case IN_HEADER_VALUE:
			if c == '"' {
				r.headers[r.currentHeader] = r.currentValue
				r.currentHeader = ""
				r.currentValue = ""
				r.state = OUT_HEADER_VALUE
			} else {
				r.currentValue += string(c)
			}
		case OUT_HEADER_VALUE:
			if c != ']' {
				return 0, errors.New("Bad pgn form: " + string(c))
			}
			r.state = LINE_END
		case LINE_END:
			if c != '\n' {
				return 0, errors.New(string(c))
			}
			r.state = LINE_START
		case IN_MOVES:
			if c == '\n' {
				r.state = LINE_START
				r.moves += " "
			} else {
				r.moves += string(c)
			}
		}
	}
	game := Game{
		headers: r.headers,
		moves:   r.moves}
	r.games = append(r.games, game)
	r.Reset()
	return len(p), nil
}

func (r *PgnReader) Reset() {
	r.state = LINE_START
	r.headers = make(map[string]string)
	r.moves = ""
	r.currentHeader = ""
	r.currentValue = ""
}

func (r *PgnReader) ReadFile(name string) ([]Game, error) {
	content, err := ioutil.ReadFile(name)

	if err != nil {
		return make([]Game, 0), err
	}

	r.Reset()
	offset := 0

	for {
		if offset == len(content) {
			return r.games, nil
		}
		read, err := r.Read(content[offset:])
		if err != nil {
			return make([]Game, 0), err
		}
		offset += read
	}
}
