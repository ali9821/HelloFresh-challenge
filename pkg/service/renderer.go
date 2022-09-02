package service

import (
	cfg "HelloFresh/config"
	"HelloFresh/pkg/model"
	"encoding/json"
	"fmt"
	"sync"
)

type Renderer struct {
	config          *cfg.Config
	rendererChannel *model.RendererChannel
	doneChannel     chan bool
}

func NewRenderer(config *cfg.Config, rendererChannel *model.RendererChannel, done chan bool) model.Runner {
	return &Renderer{
		config:          config,
		rendererChannel: rendererChannel,
		doneChannel:     done,
	}
}

func (r *Renderer) Run() error {
	var wg sync.WaitGroup
	var response = model.Response{}

	wg.Add(1)
	go r.HandleRecipeResponses(&response, &wg)
	wg.Add(1)
	go r.HandlePostalCodeResponses(&response, &wg)
	wg.Wait()

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonResponse))
	r.doneChannel <- true
	return nil
}

func (r *Renderer) HandleRecipeResponses(response *model.Response, wg *sync.WaitGroup) {
	defer wg.Done()

	for recipe := range r.rendererChannel.MatchByName {
		response.MatchByName = append(response.MatchByName, recipe)
	}

	for recipe := range r.rendererChannel.CountPerRecipe {
		response.CountPerRecipe = append(response.CountPerRecipe, recipe)
	}
	uniqueRecipeCount := <-r.rendererChannel.UniqueRecipeCount
	response.UniqueRecipeCount = uniqueRecipeCount
}

func (r *Renderer) HandlePostalCodeResponses(response *model.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	//fill busiest postal code
	busiestPostalCode := <-r.rendererChannel.BusiestPostcode

	response.BusiestPostcode = busiestPostalCode
	//fill count per postcode and time
	specificPostalCodeDeliveriesCount := <-r.rendererChannel.CountPerPostcodeAndTime
	response.CountPerPostcodeAndTime = specificPostalCodeDeliveriesCount
}
