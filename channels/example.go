package main

import "fmt"

type Service struct {
	resultChan chan error
	closed     chan int
}

func (s *Service) StartProducer() {
	go producer(s.resultChan)
}

func (s *Service) StartConsumer() {
	go consumer(s.resultChan, s.closed)
}

func (s *Service) Done() {
	cnt := <-s.closed
	fmt.Printf("handled %d errors\n", cnt)
}

func NewService() Service {
	return Service{
		resultChan: make(chan error),
		closed:     make(chan int, 1),
	}
}

func producer(ch chan<- error) {
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			fmt.Println("publishing")
			ch <- fmt.Errorf("error on %d", i)
		}
	}

	close(ch)
}

func consumer(ch <-chan error, closed chan<- int) {
	var cnt int
	for range ch {
		fmt.Println("consuming")
		cnt += 1
	}

	closed <- cnt
	close(closed)
}

func main() {
	s := NewService()

	s.StartProducer()
	s.StartConsumer()

	s.Done()

}
