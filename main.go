package main

import (
	"fmt"
	"time"
)

func main() {
	input1 := MakeInputCell(5)
	input2 := MakeInputCell(10)
	computeCell := MakeComputeCell2(input1.(*Cell[int]), input2.(*Cell[int]), func(a int, b int) int {
		return a + b
	})
	computeCh := make(chan int)
	computeCell.Subscribe(computeCh)
	input1.Update(7)
	input2.Update(3)
	input1.Update(20)
	time.Sleep(1 * time.Second)
	fmt.Println(computeCell.Value()) // 23
}
