package game

type Set struct {
  elements map[string]interface{}
  toKey func (interface {}) string
}

func MakeSet(toKey func (interface {}) string) *Set {
  return &Set{make(map[string]interface{}), toKey}
}

func (set *Set) Insert(element interface{}) {
  set.elements[set.toKey(element)] = element
}

func (set *Set) AllElements() []interface{} {
  elements := make([]interface{}, 0, set.Size())
  for _, element := range set.elements {
    elements = append(elements, element)
  }
  return elements
}

func (set *Set) IsEmpty() bool {
  return set.Size() == 0
}

func (set *Set) Size() int {
  return len(set.elements)
}

func (set *Set) Erase(element interface{}) {
  delete(set.elements, set.toKey(element))
}

func (set *Set) Contains(element interface{}) bool {
  _, ok := set.elements[set.toKey(element)]
  return ok
}
