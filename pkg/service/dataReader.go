package service

import (
	cfg "HelloFresh/config"
	"HelloFresh/pkg/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type dataReader struct {
	streamChannel chan model.Order
	config        *cfg.Config
}

func NewDataReaderService(streamChannel chan model.Order, config *cfg.Config) model.Runner {
	return &dataReader{
		streamChannel: streamChannel,
		config:        config,
	}
}

func (d *dataReader) Run() error {
	// Stop streaming channel as soon as nothing left to read in the file.
	defer close(d.streamChannel)

	// Open file to read.
	file, err := os.Open(d.config.DataFile)
	if err != nil {
		log.Fatal("Error in opening file to read")
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Read opening delimiter. `[` or `{`
	if _, err := decoder.Token(); err != nil {
		log.Fatal("decode opening delimiter")
		return err
	}

	// Read file content as long as there is something.
	i := 1
	for decoder.More() {
		var order model.Order
		if err := decoder.Decode(&order); err != nil {
			fmt.Errorf("deode line %d", i)
			return err
		}
		d.streamChannel <- order
		i++
	}

	// Read closing delimiter. `]` or `}`
	if _, err := decoder.Token(); err != nil {
		fmt.Errorf("decode closing delimiter: %w", err)
		return err
	}

	return nil
}
