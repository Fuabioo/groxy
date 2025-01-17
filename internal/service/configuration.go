package service

import (
	"strings"
	"sync"
)

type (
	EndpointConfiguration struct {
		sync.RWMutex `yaml:"-"`
		Response     *Response           `yml:"response"`
		Delay        *Delay              `yml:"delay"`
		Error        *ErrorConfiguration `yml:"error"`
	}

	ErrorConfiguration struct {
		Chance   *ErrorChance `yml:"chance"`
		Every    *Every       `yml:"every"`
		Response *Response    `yml:"response"`
	}

	Response struct {
		ContentType string `yml:"contenttype"`
		StatusCode  int    `yml:"status"`
		Body        string `yml:"body"`
	}
)

func (e *ErrorConfiguration) GetChance() *ErrorChance {
	if e == nil {
		return nil
	}
	return e.Chance
}

func (e *ErrorConfiguration) GetEvery() *Every {
	if e == nil {
		return nil
	}
	return e.Every
}

func (e *ErrorConfiguration) String() string {
	if e == nil {
		return "<nil>"
	}
	var builder strings.Builder
	if e.Chance != nil {
		builder.WriteString("Chance: ")
		builder.WriteString(e.Chance.String())
		builder.WriteString(", ")
	}
	if e.Every != nil {
		builder.WriteString("Every: ")
		builder.WriteString(e.Every.String())
		builder.WriteString(", ")
	}
	if e.Response != nil {
		builder.WriteString("Response: <not nil>")
	} else {
		builder.WriteString("Response: <nil>")
	}
	return e.Chance.String()
}
