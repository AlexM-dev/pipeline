package main

import (
	"fmt"
	pipeline2 "pipeline/pipeline"
	"pipeline/ring"
	"time"
)

const bufferTime = 15 * time.Second

const bufferSize int = 10

func main() {
	filterNeg := func(done <-chan int, input <-chan int) <-chan int {
		stream := make(chan int)
		go func() {
			defer close(stream)
			for {
				select {
				case <-done:
					return
				case i, ok := <-input:
					if !ok {
						return
					}
					if i >= 0 {
						select {
						case stream <- i:
						case <-done:
							return
						}
					}
				}
			}
		}()
		return stream
	}

	filterMulti3 := func(done <-chan int, input <-chan int) <-chan int {
		stream := make(chan int)
		go func() {
			defer close(stream)
			for {
				select {
				case <-done:
					return
				case i, ok := <-input:
					if !ok {
						return
					}
					if i%3 == 0 && i != 0 {
						select {
						case stream <- i:
						case <-done:
							return
						}
					}
				}
			}
		}()
		return stream
	}

	buffer := func(done <-chan int, c <-chan int) <-chan int {
		stream := make(chan int)
		buf := ring.NewRing(bufferSize)
		go func() {
			for {
				select {
				case data := <-c:
					buf.Write(data)
				case <-done:
					return
				}
			}
		}()
		go func() {
			for {
				select {
				case <-time.After(bufferTime):
					data := buf.ReadAll()
					if data != nil {
						for _, d := range data {
							select {
							case stream <- d:
							case <-done:
								return
							}
						}
					}
				case <-done:
					return
				}
			}
		}()
		return stream
	}

	user := func(done <-chan int, ch <-chan int) {
		for {
			select {
			case data, ok := <-ch:
				if !ok {
					return
				}
				fmt.Printf("Результат: %d\n", data)
			case <-done:
				return
			}
		}
	}

	pipe := pipeline2.NewPipeline()
	pipe.Init()
	pipe.AddStage(filterMulti3)
	pipe.AddStage(filterNeg)
	pipe.AddStage(buffer)
	c := pipe.Run()

	user(pipe.Done(), c)
}
