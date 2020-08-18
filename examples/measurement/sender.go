package measurement

import (
	"errors"
	"github.com/tarent/gomulocity"
	"github.com/tarent/gomulocity/measurement"
	"log"
	"os"
	"sync"
	"time"
)

type Sender struct {
	In      chan measurement.NewMeasurement
	ErrChan chan error
	Stop    chan os.Signal
	Timer   time.Duration
	Client  gomulocity.Gomulocity
	wg      *sync.WaitGroup
	Source  measurement.Source
}

func NewMeasurementSender(client gomulocity.Gomulocity, timer time.Duration, wg *sync.WaitGroup, source measurement.Source) *Sender {
	return &Sender{
		In:      make(chan measurement.NewMeasurement),
		ErrChan: make(chan error),
		Stop:    make(chan os.Signal),
		Timer:   timer,
		Client:  client,
		wg:      wg,
		Source:  source,
	}
}

func (s Sender) Fill() {
	for {
		s.In <- setSource(s.Source, Example1NewMeasurements)
		time.Sleep(s.Timer)
	}
}

func (s Sender) Send() {
	for {
		select {
		case measurementEntry := <-s.In:
			m, err := s.Client.MeasurementApi.Create(&measurementEntry)
			if err != nil {
				s.ErrChan <- err
			} else {
				log.Println(m)
			}
		case <-time.After(20 * time.Second):
			s.ErrChan <- errors.New("sender ran into timeout... Send() forced to stop")
			s.wg.Done()
			break
		case <-s.Stop:
			s.wg.Done()
			break
		}
	}
}

func (s Sender) Errors() {
	for {
		select {
		case err := <-s.ErrChan:
			log.Println(err)
		}
	}
}

func setSource(source measurement.Source, newMeasurement measurement.NewMeasurement) measurement.NewMeasurement {
	newMeasurement.Source = source
	return newMeasurement
}
