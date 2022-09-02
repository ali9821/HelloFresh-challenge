package main

import (
	"HelloFresh/cmd"
	"HelloFresh/pkg/model"
	"log"
)

func main() {
	log.Println("task started")
	factory := cmd.NewFactory()

	for _, runner := range factory.PipelineStages {
		go func(runner model.Runner) {
			err := runner.Run()
			if err != nil {
				log.Fatal(err)
			}
		}(runner)
	}

}
