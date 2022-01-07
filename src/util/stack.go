package util

type Stack []string

func NewStack() *Stack {
  s := Stack([]string{})
  return &s
}

// Creates a new stack
func NewStackWith(first string) *Stack {
  arr := []string{first}
  stack := Stack(arr)
  return &stack
}

// Top returns the element at the top of the stack
func (s *Stack) Top() string {
  if len(*s) == 0 {
    return ""
  }
  return (*s)[len(*s)-1]
}

// Push pushes an element to the stack
func (s *Stack) Push(str string) {
  *s = append(*s, str)
}

// Pop pops an element from the stack
func (s *Stack) Pop() string {
  if len(*s) == 0 {
    return ""
  }
  ret := (*s)[len(*s)-1]
  *s = (*s)[:len(*s)-1]
  return ret
}

