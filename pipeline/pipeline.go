package pipeline

import (
	"bufio"
	"fmt"
	"log"
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
	log.Printf("Create pipeline")
	return &pipeline{}
}

func (p *pipeline) Init() {
	ch := make(chan int)
	done := make(chan int)
	log.Printf("Init")
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		log.Printf("Create scanner")
		var data string
		for {
			scanner.Scan()
			log.Printf("Scan")
			data = scanner.Text()
			if strings.EqualFold(data, "exit") {
				fmt.Println("Работа завершена")
				log.Printf("Exit")
				close(done)
				return
			}
			i, err := strconv.Atoi(data)
			if err != nil {
				fmt.Println("Только целые числа")
				log.Printf("Error: only int")
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
	log.Printf("Add new stage")
	p.stages = append(p.stages, s)
}

func (p *pipeline) Run() <-chan int {
	var c <-chan int = p.ch
	log.Printf("Run")
	for i := range p.stages {
		c = p.stages[i](p.done, c)
	}
	return c
}

func (p *pipeline) End() {
	log.Printf("End")
	close(p.done)
}
