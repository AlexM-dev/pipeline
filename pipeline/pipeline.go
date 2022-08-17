package pipeline

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type stage func(<-chan int, <-chan int) <-chan int

type pipeline struct {
	stages []stage
	done   chan int
	ch     chan int
}

func NewPipeline() *pipeline {
	return &pipeline{}
}

func (p *pipeline) Init() {
	ch := make(chan int)
	done := make(chan int)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		var data string
		for {
			scanner.Scan()
			data = scanner.Text()
			if strings.EqualFold(data, "exit") {
				fmt.Println("Работа завершена")
				close(done)
				return
			}
			i, err := strconv.Atoi(data)
			if err != nil {
				fmt.Println("Только целые числа")
				continue
			}
			ch <- i
		}
	}()
	p.ch = ch
	p.done = done
}

func (p *pipeline) Done() <-chan int {
	return p.done
}

func (p *pipeline) AddStage(s stage) {
	p.stages = append(p.stages, s)
}

func (p *pipeline) Run() <-chan int {
	var c <-chan int = p.ch
	for i := range p.stages {
		c = p.stages[i](p.done, c)
	}
	return c
}

func (p *pipeline) End() {
	close(p.done)
}
