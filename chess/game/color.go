package game

type Color int

const (
  White Color = iota
  Black = iota
)

func (color Color) String() string {
  if color == White {
    return "white (lower case)"
  } else {
    return "black (upper case)"
  }
}

func (color Color) Other() Color {
  if color == Black {
    return White
  }
  return Black
}
