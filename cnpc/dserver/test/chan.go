package main

import (
    "time"
    "fmt"
)

func Process(ch chan int, id int) {
    //Do some work...
    time.Sleep(5*time.Second)

    ch <- id
}

func main() {

    channels := make([]chan int, 10)

    for i:= 0; i < 10; i++ {
        channels[i] = make(chan int)
        go Process(channels[i], i)
    }

    for _, ch := range channels {
	    ex := <-ch
        fmt.Println("Routine ", ex, " quit!")
    }
}
