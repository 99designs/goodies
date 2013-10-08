package goanna

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type ControllerFactoryFunc func() ControllerInterface

// ControllerHandler is a http.Handler for handling incoming requests
// and despatching to controllers
type ControllerHandler struct {
	factory    ControllerFactoryFunc
	methodName string
}

// ServeHTTP handles a http request
func (handler ControllerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.getResponse(r).Send(w)
}

// getResponse executes the specified controller's method using reflection
// and returns the response object
func (handler ControllerHandler) getResponse(r *http.Request) Response {
	controller := handler.factory()
	controller.SetRequest(r)
	controller.Init()
	rController := reflect.ValueOf(controller)
	method := rController.MethodByName(handler.methodName)

	// get args from gorilla mux
	var args []reflect.Value
	for _, val := range mux.Vars(r) {
		args = append(args, reflect.ValueOf(val))
	}

	// make sure number of args matches the controller method
	expected := len(args)
	actual := method.Type().NumIn()
	if expected != actual {
		panic(fmt.Sprintf("Method '%s' has %d args, expected %d", handler.methodName, actual, expected))
	}

	out := method.Call(args)
	if out[0].IsNil() {
		panic("Response from controller was nil")
	}

	resp := out[0].Interface().(Response)
	if resp == nil {
		panic("Response from controller was not Response interface")
	}

	controller.Session().WriteToResponse(resp)
	return resp
}

// isValid checks that the controller and method specifies
// will sucessfully execute if getResponse is called on it
func (handler ControllerHandler) isValid() bool {
	controller := handler.factory()
	rController := reflect.ValueOf(controller)
	method := rController.MethodByName(handler.methodName)
	if (method == reflect.Value{}) {
		panic("No such method: " + handler.methodName)
	}
	typeOfMethod := method.Type()

	var r *Response
	responseType := reflect.TypeOf(r).Elem()

	return (method.Kind() == reflect.Func) &&
		(typeOfMethod.NumMethod() == 0) &&
		(typeOfMethod.NumOut() == 1) &&
		typeOfMethod.Out(0) == responseType

}

// NewHandler creates a ControllerHandler from the factory and methodName
func NewHandler(factory ControllerFactoryFunc, methodName string) ControllerHandler {
	handler := ControllerHandler{factory: factory, methodName: methodName}
	if !handler.isValid() {
		panic("Invalid handler: " + methodName)
	}
	return handler
}
