package cmd

import (
	cfg "HelloFresh/config"
	"HelloFresh/pkg/model"
	"HelloFresh/pkg/service"
)

type Factory struct {
	PipelineStages []model.Runner
}

func NewFactory() *Factory {
	var streamChannel = make(chan model.Order, 1000)
	var uniqueRecipeChan = make(chan model.Order, 1000)
	var mostPostCodeDeliveredChan = make(chan model.Order, 1000)
	var specificPostCodeChan = make(chan model.Order, 1000)
	var recipeListChan = make(chan model.Order, 1000)

	config := cfg.NewConfig()

	dataReaderService := service.NewDataReaderService(streamChannel, config)

	distributorService := service.NewDistributorService(streamChannel, uniqueRecipeChan, mostPostCodeDeliveredChan, specificPostCodeChan, recipeListChan, config)

	return &Factory{
		PipelineStages: []model.Runner{dataReaderService, distributorService},
	}
}
