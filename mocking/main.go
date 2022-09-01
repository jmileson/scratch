package main

import (
	"os"

	"github.com/jmileson/scratch/mocking/step1"
	"github.com/jmileson/scratch/mocking/step2"
	"github.com/jmileson/scratch/mocking/step3"
)

func main() {
	log1 := step1.Logger{}
	log1.Info("hello world")

	log2 := step2.Logger{Out: os.Stdout}
	log2.Info("hello world, again")

	log3 := step3.Logger{
		Out: step3.AddFlush(os.Stdout),
	}
	log3.Info("hello world, again, again...")
}
