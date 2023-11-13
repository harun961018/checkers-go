package checkers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	BOARD_DIM = 8
	RED       = "red"
	BLACK     = "black"
	ROW_SEP   = "|"
)

type Player struct {
	Color string
}

type Piece struct {
	Player Player
	King   bool
}

var PieceStrings = map[Player]string{
	RED_PLAYER:   "r",
	BLACK_PLAYER: "b",
	NO_PLAYER:    "*",
}

var NO_PIECE = Piece{NO_PLAYER, false}

var StringPieces = map[string]Piece{
	"r": Piece{RED_PLAYER, false},
	"b": Piece{BLACK_PLAYER, false},
	"R": Piece{RED_PLAYER, true},
	"B": Piece{BLACK_PLAYER, true},
	"*": NO_PIECE,
}

type Pos struct {
	X int
	Y int
}

var NO_POS = Pos{-1, -1}

var BLACK_PLAYER = Player{BLACK}
var RED_PLAYER = Player{RED}
var NO_PLAYER = Player{
	Color: "NO_PLAYER",
}

var Players = map[string]Player{
	RED:   RED_PLAYER,
	BLACK: BLACK_PLAYER,
}

var Opponents = map[Player]Player{
	BLACK_PLAYER: RED_PLAYER,
	RED_PLAYER:   BLACK_PLAYER,
}

var Usable = map[Pos]bool{}
var Moves = map[Player]map[Pos]map[Pos]bool{}
var Jumps = map[Player]map[Pos]map[Pos]Pos{}
var KingMoves = map[Pos]map[Pos]bool{}
var KingJumps = map[Pos]map[Pos]Pos{}

func init() {
	// Initialize usable spaces
	for y := 0; y < BOARD_DIM; y++ {
		for x := (y + 1) % 2; x < BOARD_DIM; x += 2 {
			Usable[Pos{X: x, Y: y}] = true
		}
	}

	// Initialize deep maps
	for _, p := range Players {
		Moves[p] = map[Pos]map[Pos]bool{}
		Jumps[p] = map[Pos]map[Pos]Pos{}
	}

	// Compute possible moves, jumps and captures
	for pos := range Usable {
		KingMoves[pos] = map[Pos]bool{}
		KingJumps[pos] = map[Pos]Pos{}
		var directions = []int{1, -1}
		for i, player := range []Player{BLACK_PLAYER, RED_PLAYER} {
			Moves[player][pos] = map[Pos]bool{}
			Jumps[player][pos] = map[Pos]Pos{}
			movOff := 1
			jmpOff := 2
			for _, direction := range directions {
				mov := Pos{pos.X + (movOff * direction), pos.Y + (movOff * directions[i])}
				if Usable[mov] {
					Moves[player][pos][mov] = true
					KingMoves[pos][mov] = true
				}
				jmp := Pos{pos.X + (jmpOff * direction), pos.Y + (jmpOff * directions[i])}
				if Usable[jmp] {
					capturePos := Capture(pos, jmp)
					Jumps[player][pos][jmp] = capturePos
					KingJumps[pos][jmp] = capturePos
				}
			}
		}
	}
}

type Game struct {
	Pieces map[Pos]Piece
	Turn   Player
}

func (game *Game) String() string {
	var buf bytes.Buffer
	for y := 0; y < BOARD_DIM; y++ {
		for x := 0; x < BOARD_DIM; x++ {
			pos := Pos{x, y}
			if game.PieceAt(pos) {
				piece := game.Pieces[pos]
				val := PieceStrings[piece.Player]
				if piece.King {
					val = strings.ToUpper(val)
				}
				buf.WriteString(val)
			} else {
				buf.WriteString(PieceStrings[NO_PLAYER])
			}
		}
		if y < (BOARD_DIM - 1) {
			buf.WriteString(ROW_SEP)
		}
	}
	return buf.String()
}

func Parse(s string) (*Game, error) {
	if len(s) != BOARD_DIM*BOARD_DIM+(BOARD_DIM-1) {
		return nil, errors.New(fmt.Sprintf("invalid board string: %v", s))
	}
	pieces := make(map[Pos]Piece)
	result := &Game{pieces, BLACK_PLAYER}
	for y, row := range strings.Split(s, ROW_SEP) {
		for x, c := range strings.Split(row, "") {
			if x >= BOARD_DIM || y >= BOARD_DIM {
				return nil, errors.New(fmt.Sprintf("invalid board, piece out of bounds: %v, %v", x, y))
			}
			if piece, ok := ParsePiece(c); !ok {
				return nil, errors.New(fmt.Sprintf("invalid board, invalid piece at %v, %v", x, y))
			} else if piece != NO_PIECE {
				result.Pieces[Pos{x, y}] = piece
			}
		}
	}
	return result, nil
}
