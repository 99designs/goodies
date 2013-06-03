package monitor

import (
	"errors"
	"time"
)

type Section struct {
	start *time.Time
	stop  *time.Time
}

type Stopwatch struct {
	sections map[string]*Section
}

func NewStopwatch() Stopwatch {
	t := Stopwatch{sections: make(map[string]*Section)}

	return t
}

func now() *time.Time {
	now := time.Now()
	return &now
}

func (t *Stopwatch) Start(sectionName string) error {
	_, ok := t.sections[sectionName]
	if ok {
		return errors.New("Section " + sectionName + " already started")
	}

	section := Section{start: now()}
	t.sections[sectionName] = &section

	return nil
}

func (t *Stopwatch) Stop(sectionName string) error {
	section, ok := t.sections[sectionName]
	if !ok {
		return errors.New("Section " + sectionName + " not yet started")
	} else {
		section.stop = now()
	}

	return nil
}

func (t *Stopwatch) Duration(sectionName string) (time.Duration, error) {
	section, ok := t.sections[sectionName]

	var err error
	var duration time.Duration

	if !ok {
		err = errors.New("Section " + sectionName + " not yet started")
	} else if section.start == nil || section.stop == nil {
		err = errors.New("Section " + sectionName + " not yet stopped")
	} else {
		duration = section.stop.Sub(*section.start)
	}

	return duration, err
}
