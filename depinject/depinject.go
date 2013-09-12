// Package depinject implements a simple
// Dependency Injector
package depinject

import (
	"errors"
	"log"
	"reflect"
)

var (
	ConstructorErr     = errors.New("Constructor must be a function and output 1 parameter")
	ConstructorArgsErr = errors.New("Can't resolve function arguments")
)

// DependencyInjector is a simple IOC container
// that allows registering constructor functions
// and creating them
type DependencyInjector struct {
	registry map[reflect.Type]interface{}
}

// Register registers a constructor function
// with the DI container
func (di *DependencyInjector) Register(constructorFunc interface{}) error {
	constructorType := reflect.TypeOf(constructorFunc)

	if (constructorType.Kind() != reflect.Func) || (constructorType.NumOut() != 1) {
		return ConstructorErr
	}
	outType := constructorType.Out(0)

	// make sure we can resolve the constuctor arguments
	for i := 0; i < constructorType.NumIn(); i++ {
		inType := constructorType.In(i)
		_, ok := di.registry[inType]
		if !ok {
			log.Printf("Can't find a %s for a %s\n", inType, outType)
			return ConstructorArgsErr
		}
	}

	di.registry[outType] = constructorFunc

	return nil
}

// Create creates an instance of the type of the given parameter
func (di *DependencyInjector) Create(avar interface{}) interface{} {
	return di.CreateFromType(reflect.TypeOf(avar)).Interface()
}

// CreateFromType creates an instance of the given type
func (di *DependencyInjector) CreateFromType(atype reflect.Type) reflect.Value {
	constructor, exists := di.registry[atype]
	if !exists {
		log.Panicf("Can't find a mapping to create a %s", atype)
	}

	constructorType := reflect.TypeOf(constructor)
	constructorArgs := []reflect.Value{}

	for i := 0; i < constructorType.NumIn(); i++ {
		t := constructorType.In(i)
		v := di.CreateFromType(t)
		constructorArgs = append(constructorArgs, v)
	}

	newObj := reflect.ValueOf(constructor).Call(constructorArgs)

	return newObj[0]
}

func NewDependencyInjector() DependencyInjector {
	return DependencyInjector{
		registry: make(map[reflect.Type]interface{}),
	}
}
