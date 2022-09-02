package service

import (
	cfg "HelloFresh/config"
	"HelloFresh/pkg/model"
	"sync"
)

type processor struct {
	orderChan      chan model.Order
	recipeChan     chan string
	postalCodeChan chan model.Order
	config         *cfg.Config
}

func NewDistributorService(streamChannel chan model.Order, recipeChan chan string, postalCodeChan chan model.Order, config *cfg.Config) model.Runner {

	return &processor{
		orderChan:      streamChannel,
		config:         config,
		recipeChan:     recipeChan,
		postalCodeChan: postalCodeChan,
	}
}

func (p *processor) Run() error {
	var wg sync.WaitGroup
	for i := 0; i < p.config.MaxUniqueRecipeWorkersSize; i++ {
		wg.Add(1)
		go p.Worker(&wg)
	}
	wg.Wait()
	close(p.recipeChan)
	close(p.postalCodeChan)
	return nil
}

func (p *processor) Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range p.orderChan {
		p.recipeChan <- order.Recipe
		p.postalCodeChan <- order
		//p.recipeListChan <- order
	}

}
