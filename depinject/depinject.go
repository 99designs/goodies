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

type DependencyInjector interface {
	Register(interface{}) error
	MustRegister(interface{})
	Create(interface{}) interface{}
	CreateFromType(reflect.Type) reflect.Value
}

// DependencyInjector is a simple IOC container
// that allows registering constructor functions
// and creating types
type dependencyInjector struct {
	caching   bool
	instances map[reflect.Type]reflect.Value
	registry  map[reflect.Type]interface{}
}

// Register registers a constructor function
// with the DI container
func (di *dependencyInjector) Register(constructorFunc interface{}) error {
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
func (di *dependencyInjector) MustRegister(constructorFunc interface{}) {
	err := di.Register(constructorFunc)
	if err != nil {
		panic(err)
	}
}

// Create creates an instance of the type of the given parameter
func (di *dependencyInjector) Create(avar interface{}) interface{} {
	varType := reflect.TypeOf(avar)

	if di.caching {
		return di.cachedCreateFromType(varType).Interface()
	} else {
		return di.CreateFromType(varType).Interface()
	}
}

// CreateFromType creates an instance of the given type
func (di *dependencyInjector) createFromType(atype reflect.Type) reflect.Value {
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

// cachedCreateFromType creates an instance of the given type
func (di *dependencyInjector) cachedCreateFromType(atype reflect.Type) reflect.Value {
	_, exists := di.instances[atype]

	if !exists {
		di.instances[atype] = di.createFromType(atype)
	}

	return di.instances[atype]
}

// CreateFromType creates an instance of the given type
func (di *dependencyInjector) CreateFromType(atype reflect.Type) reflect.Value {
	if di.caching {
		return di.cachedCreateFromType(atype)
	} else {
		return di.createFromType(atype)
	}
}

// NewDependencyInjector returns a new DependencyInjector
func NewDependencyInjector() DependencyInjector {
	return &dependencyInjector{
		registry:  make(map[reflect.Type]interface{}),
		instances: make(map[reflect.Type]reflect.Value),
	}
}

// ServiceContainer is a type of DependencyInjector that caches instances
type ServiceContainer struct {
	DependencyInjector
}

// NewServiceContainer returns a new ServiceContainer
func NewServiceContainer() ServiceContainer {
	return ServiceContainer{
		DependencyInjector: &dependencyInjector{
			registry:  make(map[reflect.Type]interface{}),
			instances: make(map[reflect.Type]reflect.Value),
			caching:   true,
		},
	}
}
