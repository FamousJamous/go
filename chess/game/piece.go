package game

import (
  "fmt"
  "unicode"
)

type Piece struct {
  // p, r, n, b, q, k
  name byte
  color Color
}

func (piece *Piece) GetName() byte {
  return piece.name
}

func (piece *Piece) GetPoints() int {
  if piece == nil {
    return 0
  }
  switch piece.name {
    case 'p': return 1
    case 'n': return 3
    case 'b': return 3
    case 'r': return 5
    case 'q': return 9
    case 'k': return 0
  }
  panic(fmt.Sprintf("Unexpected piece %c", piece.name))
}

func (piece *Piece) GetColor() Color {
  return piece.color
}

func (piece *Piece) String() string {
  if piece == nil {
    return " "
  }
  var char rune
  if piece.color == Black {
    char = unicode.ToUpper(rune(piece.name))
  } else {
    char = rune(piece.name)
  }
  return fmt.Sprintf("%c", char)
}
