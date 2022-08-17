package main

import (
	"fmt"
	"log"
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
							log.Printf("Write in stream")
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
							log.Printf("Write in stream")
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
		log.Printf("Create buffer")
		go func() {
			for {
				select {
				case data := <-c:
					buf.Write(data)
					log.Printf("Write in buffer")
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
					log.Printf("Read buffer")
					if data != nil {
						for _, d := range data {
							select {
							case stream <- d:
								log.Printf("Write in stream")
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
				log.Printf("Get result3")
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
