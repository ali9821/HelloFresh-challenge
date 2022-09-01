package service

import (
	cfg "HelloFresh/config"
	"HelloFresh/pkg/model"
	"sync"
)

type processor struct {
	streamChannel             chan model.Order
	uniqueRecipeChan          chan model.Order
	mostPostCodeDeliveredChan chan model.Order
	specificPostCodeChan      chan model.Order
	recipeListChan            chan model.Order
	config                    *cfg.Config
}

func NewDistributorService(streamChannel chan model.Order, uniqueRecipeChan chan model.Order, mostPostCodeDeliveredChan chan model.Order, specificPostCodeChan chan model.Order, recipeListChan chan model.Order, config *cfg.Config) model.Runner {

	return &processor{
		streamChannel:             streamChannel,
		config:                    config,
		uniqueRecipeChan:          uniqueRecipeChan,
		mostPostCodeDeliveredChan: mostPostCodeDeliveredChan,
		specificPostCodeChan:      specificPostCodeChan,
		recipeListChan:            recipeListChan,
	}
}

func (p *processor) Run() error {
	var wg sync.WaitGroup
	for i := 0; i < p.config.MaxUniqueRecipeWorkersSize; i++ {
		wg.Add(1)
		go p.Worker(&wg)
	}
	wg.Wait()
	return nil
}

func (p *processor) Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range p.streamChannel {
		p.uniqueRecipeChan <- data
		p.mostPostCodeDeliveredChan <- data
		p.specificPostCodeChan <- data
		p.recipeListChan <- data
	}
}
