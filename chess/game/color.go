package game

type Color int

const (
  White Color = iota
  Black = iota
)

func (c Color) String() string {
  if c == White {
    return "white (lower case)"
  } else {
    return "black (upper case)"
  }
}
