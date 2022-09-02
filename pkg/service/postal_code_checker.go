package service

import (
	cfg "HelloFresh/config"
	"HelloFresh/pkg/model"
	"HelloFresh/tools"
	"fmt"
	"sync"
)

type PostalCodeChecker struct {
	config                            *cfg.Config
	postalCodeChan                    chan model.Order
	uniquePostalCodes                 sync.Map
	specificPostalCodeDeliveriesCount int
	rendererChannel                   model.RendererChannel
}

func NewPostalCodeChecker(config *cfg.Config, postalCodeChan chan model.Order, rendererChannel *model.RendererChannel) model.Runner {
	return &PostalCodeChecker{
		config:                            config,
		postalCodeChan:                    postalCodeChan,
		uniquePostalCodes:                 sync.Map{},
		specificPostalCodeDeliveriesCount: 0,
		rendererChannel:                   *rendererChannel,
	}
}

func (p *PostalCodeChecker) Run() error {
	var wg sync.WaitGroup
	for i := 0; i < p.config.MaxMostPostCodeDeliveredWorkersSize; i++ {
		wg.Add(1)
		go p.worker(&wg)
	}
	wg.Wait()
	p.sendToRenderer()
	return nil
}

func (p *PostalCodeChecker) worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range p.postalCodeChan {
		if order.Postcode == "10120" {
			if p.isInValidDateRange(order.Delivery) {
				p.specificPostalCodeDeliveriesCount++
			}
		}
		numberOfOccurrence, ok := p.uniquePostalCodes.Load(order.Postcode)
		if !ok {
			p.uniquePostalCodes.Store(order.Postcode, 1)
		} else {
			p.uniquePostalCodes.Store(order.Postcode, numberOfOccurrence.(int)+1)
		}
	}
}

func (p *PostalCodeChecker) sendToRenderer() {
	busiestPostalCode, busiestPostalCodeOccurrence := tools.FindMaxInSyncMap(&p.uniquePostalCodes)
	fmt.Println(busiestPostalCode, busiestPostalCodeOccurrence)
	p.rendererChannel.BusiestPostcode <- model.BusiestPostcode{
		Postcode:      busiestPostalCode,
		DeliveryCount: busiestPostalCodeOccurrence,
	}
	close(p.rendererChannel.BusiestPostcode)
	p.rendererChannel.CountPerPostcodeAndTime <- model.CountPerPostcodeAndTime{
		Postcode:      "10120",
		From:          "10AM",
		To:            "3PM",
		DeliveryCount: p.specificPostalCodeDeliveriesCount,
	}
	close(p.rendererChannel.CountPerPostcodeAndTime)

}

func (p *PostalCodeChecker) isInValidDateRange(delivery string) bool {
	deliveryHours := tools.ExtractNumbersFromString(delivery)
	if len(deliveryHours) != 2 {
		return false
	}
	startTime := tools.ConvertStringToInt(deliveryHours[0])
	endTime := tools.ConvertStringToInt(deliveryHours[1])
	if startTime >= 10 && endTime <= 3 {
		return true
	}
	return false
}
