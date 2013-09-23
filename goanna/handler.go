package goanna

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
)

type ControllerFactoryFunc func() ControllerInterface

type ControllerHandler struct {
	factory    ControllerFactoryFunc
	methodName string
}

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
	expected := len(args)
	actual := method.Type().NumIn()
	if expected != actual {
		log.Panic(fmt.Sprintf("Method '%s' has %d args, expected %d", handler.methodName, actual, expected))
	}

	out := method.Call(args)
	if out[0].IsNil() {
		return NewErrorResponse("Response from controller was nil", 500)
	}
	resp := out[0].Interface().(Response)
	if resp == nil {
		return NewErrorResponse("Response from controller was not Response interface", 500)
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

func NewHandler(factory ControllerFactoryFunc, methodName string) ControllerHandler {
	handler := ControllerHandler{factory: factory, methodName: methodName}
	if !handler.isValid() {
		log.Panic("Invalid handler: " + methodName)
	}
	return handler
}
