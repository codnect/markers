package visitor

type Tuple[T any] struct {
	items []T
}

func (t *Tuple[T]) Len() int {
	return len(t.items)
}

func (t *Tuple[T]) At(index int) (ret T) {
	if index >= 0 && index < len(t.items) {
		return t.items[index]
	}

	return
}
