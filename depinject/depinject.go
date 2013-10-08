// Package depinject provides a simple Dependency Injector.
//
// Dependency injection is a useful pattern that allows you to standardise and
// centralise the way types are constructed.
//
// This DI container allows you to Register constructor closures for the
// types that you need, and Create types on demand.
//
// With well defined constructors that use dependency injection for all required
// dependencies, a full dependency tree can be built up. This means that dependencies
// will cascade from your initial Create. You should aim to keep calls to Create to a
// minimum and allow dependencies to cascade.
package depinject

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	constructorErr = errors.New("Constructor must be a function and output 1 parameter")
)

// DependencyInjector is a simple IOC container
// that allows registering constructor functions
// and creating types
type DependencyInjector struct {
	registry map[reflect.Type]interface{}
}

// Register registers a constructor function
// with the DI container
func (di *DependencyInjector) Register(constructorFunc interface{}) error {
	constructorType := reflect.TypeOf(constructorFunc)

	if (constructorType.Kind() != reflect.Func) || (constructorType.NumOut() != 1) {
		return constructorErr
	}
	outType := constructorType.Out(0)

	// make sure we can resolve the constuctor arguments
	for i := 0; i < constructorType.NumIn(); i++ {
		inType := constructorType.In(i)
		_, ok := di.registry[inType]
		if !ok {
			return fmt.Errorf("Can't resolve function arguments - can't find a %s for a %s\n", inType, outType)
		}
	}

	di.registry[outType] = constructorFunc

	return nil
}

// MustRegister is a helper that calls Register and panics if it returns an error
func (di *DependencyInjector) MustRegister(constructorFunc interface{}) {
	err := di.Register(constructorFunc)
	if err != nil {
		panic(err)
	}
}

// Create creates an instance of the type of the given parameter
func (di *DependencyInjector) Create(avar interface{}) interface{} {
	return di.CreateFromType(reflect.TypeOf(avar)).Interface()
}

// CreateFromType creates an instance of the given type
func (di *DependencyInjector) CreateFromType(atype reflect.Type) reflect.Value {
	constructor, exists := di.registry[atype]
	if !exists {
		panic(fmt.Sprintf("Can't find a mapping to create a %s", atype))
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

// NewDependencyInjector returns a new DependencyInjector
func NewDependencyInjector() DependencyInjector {
	return DependencyInjector{
		registry: make(map[reflect.Type]interface{}),
	}
}
