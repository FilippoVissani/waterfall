package main

import "sync"

type InputCell[T any] interface {
	Value() T
	Subscribe(ch chan T)
	Update(value T)
	Close()
}

type ComputeCell[T any] interface {
	Value() T
	Subscribe(ch chan T)
	Close()
}

type Cell[T any] struct {
	value       T
	subscribers []chan<- T
}

func MakeInputCell[T any](value T) InputCell[T] {
	return &Cell[T]{
		value:       value,
		subscribers: make([]chan<- T, 0),
	}
}

func MakeComputeCell2[T any](input1 *Cell[T], input2 *Cell[T], f func(a T, b T) T) ComputeCell[T] {
	return MakeComputeCell([]*Cell[T]{input1, input2}, func(inputs []T) T {
		return f(inputs[0], inputs[1])
	})
}

func MakeComputeCell[T any](inputs []*Cell[T], f func(inputs []T) T) ComputeCell[T] {
	channels := make([]chan T, len(inputs))
	values := make([]T, len(inputs))
	mu := sync.Mutex{}
	for i, input := range inputs {
		values[i] = input.Value()
		channels[i] = make(chan T)
		input.Subscribe(channels[i])
	}
	cell := &Cell[T]{
		value:       f(values),
		subscribers: make([]chan<- T, 0),
	}
	for i, input := range inputs {
		go func(input *Cell[T], position int) {
			for {
				newValue := <-channels[position]
				mu.Lock()
				values[position] = newValue
				mu.Unlock()
				cell.Update(f(values))
			}
		}(input, i)
	}
	return cell
}

func (c *Cell[T]) Value() T {
	return c.value
}

func (c *Cell[T]) Subscribe(ch chan T) {
	c.subscribers = append(c.subscribers, ch)
}

func (c *Cell[T]) Update(value T) {
	c.value = value
	for _, subscriber := range c.subscribers {
		subscriber <- value
	}
}

func (c *Cell[T]) Close() {
	for _, subscriber := range c.subscribers {
		close(subscriber)
	}
}
