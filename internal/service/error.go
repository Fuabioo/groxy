package service

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/charmbracelet/log"
)

var ErrPlaceholder = errors.New("error by chance")

type (
	ErrorChance float32
	Every       struct {
		Nth     int `yml:"nth"`
		Amount  int `yml:"amount"`
		counter int `yml:"-"`
	}
)

func (e *ErrorChance) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var chance float32
	if err := unmarshal(&chance); err != nil {
		return err
	}
	if chance < 0 || chance > 1 {
		return fmt.Errorf("invalid chance value: %f", chance)
	}
	*e = ErrorChance(chance)
	return nil
}

func (e *ErrorChance) Take() error {
	if e == nil || *e == 0 {
		return nil
	}
	if rand.Float32() < float32(*e) {
		percentage := float32(*e) * 100
		percentageString := fmt.Sprintf("%.1f%%", percentage)
		log.Info("simulated error by chance",
			"percentage", percentageString,
			"err", ErrPlaceholder,
		)
		return fmt.Errorf("simulated error by chance (%s): %w",
			percentageString,
			ErrPlaceholder,
		)
	}
	return nil
}

func (e *ErrorChance) String() string {
	if e == nil {
		return "0.0%"
	}
	return fmt.Sprintf("%.1f%%", float32(*e)*100)
}

func (e *Every) Take() error {
	if e == nil {
		return nil
	}

	e.counter++

	if e.counter == 0 {
		return nil
	}

	if e.counter >= e.Nth {
		log.Info("simulated error every n requests",
			"every", e.Nth,
			"amount", e.Amount,
			"counter", e.counter,
			"err", ErrPlaceholder,
		)
		if e.counter == e.Nth+e.Amount-1 {
			e.counter = 0
		}
		return fmt.Errorf("simulated error every %d requests %d times, current %d: %w",
			e.Nth,
			e.Amount,
			e.counter,
			ErrPlaceholder,
		)
	}

	return nil
}

func (e *Every) String() string {
	if e == nil {
		return "0"
	}
	return fmt.Sprintf("%d times on the %dnth", e.Amount, e.Nth)
}
