/*
 Package monitor implements an instrumentation interface.

 Arbitrary functions can be profiled, and profile results are passed to the logger function (Monitor.Logger).
*/
package monitor

import (
	"reflect"
	"time"
)

type Monitor struct {
	Stopwatch
	Logger func(section string, duration time.Duration)
}

func NewMonitor() *Monitor {
	nullLogger := func(section string, duration time.Duration) {}
	monitor := Monitor{Stopwatch: NewStopwatch(), Logger: nullLogger}

	return &monitor
}

func Duration(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Now().Sub(start)
}

// stop the stopwatch
func (m *Monitor) StopAndLog(sectionName string) error {
	err := m.Stop(sectionName)
	if err == nil {
		m.Log(sectionName)
	}

	return err
}

func (m *Monitor) Log(section string) {
	duration, _ := m.Duration(section)
	m.Logger(section, duration)
}

// Function decorator to monitor the execution of an anonymous function
// returns an array of the returned Values
func (m *Monitor) MonitorFunc(section string, myfunc interface{}) []reflect.Value {
	return m.MonitorReflectedFunc(section, reflect.ValueOf(myfunc), []reflect.Value{})
}

func (m *Monitor) MonitorFuncWithArgs(section string, myfunc interface{}, args []reflect.Value) []reflect.Value {
	return m.MonitorReflectedFunc(section, reflect.ValueOf(myfunc), args)
}

// Function decorator
// to monitor a reflected function
func (m *Monitor) MonitorReflectedFunc(section string, reflectedFunc reflect.Value, args []reflect.Value) []reflect.Value {
	m.Start(section)
	defer m.StopAndLog(section)

	return reflectedFunc.Call(args)
}
