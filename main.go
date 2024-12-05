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
	defer input1.Close()
	defer input2.Close()
	defer computeCell.Close()
	computeCh := make(chan int)
	computeCell.Subscribe(computeCh)
	fmt.Println(computeCell.Value())
	go input1.Update(7) // Sum should now be 7 + 10 = 17
	time.Sleep(100 * time.Millisecond)
	fmt.Println(computeCell.Value())
	go input2.Update(3) // Sum should now be 7 + 3 = 10
	time.Sleep(100 * time.Millisecond)
	fmt.Println(computeCell.Value())
	go input1.Update(20) // Sum should now be 20 + 3 = 23
	time.Sleep(100 * time.Millisecond)
	fmt.Println(computeCell.Value())
}
