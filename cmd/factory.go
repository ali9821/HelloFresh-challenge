package cmd

import (
	cfg "HelloFresh/config"
	"HelloFresh/pkg/model"
	"HelloFresh/pkg/service"
)

type Factory struct {
	PipelineStages []model.Runner
	DoneChan       chan bool
}

func NewFactory() *Factory {
	config := cfg.NewConfig()

	var orderChannel = make(chan model.Order, 10000)
	var recipeChan = make(chan string, 1000)
	var postalCodeChan = make(chan model.Order, 1000)
	var rendererChannel = &model.RendererChannel{
		UniqueRecipeCount:       make(chan int, 1),
		CountPerRecipe:          make(chan model.CountPerRecipe),
		BusiestPostcode:         make(chan model.BusiestPostcode),
		CountPerPostcodeAndTime: make(chan model.CountPerPostcodeAndTime),
		MatchByName:             make(chan string),
	}
	done := make(chan bool)

	dataReaderService := service.NewDataReaderService(orderChannel, config)

	distributorService := service.NewDistributorService(orderChannel, recipeChan, postalCodeChan, config)

	recipeCheckerService := service.NewRecipeChecker(config, recipeChan, rendererChannel)
	postalCodeCheckerService := service.NewPostalCodeChecker(config, postalCodeChan, rendererChannel)
	rendererService := service.NewRenderer(config, rendererChannel, done)
	return &Factory{
		PipelineStages: []model.Runner{
			dataReaderService,
			distributorService,
			recipeCheckerService,
			postalCodeCheckerService,
			rendererService,
		},
		DoneChan: done,
	}
}
