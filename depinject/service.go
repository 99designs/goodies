package depinject

import (
	"reflect"
)

// ServiceContainer is a decorator for
// DependencyInjector that caches instances
// that the DI creates
type ServiceContainer struct {
	DependencyInjector
	instances map[reflect.Type]reflect.Value
}

func (sc *ServiceContainer) Create(avar interface{}) interface{} {
	return sc.CreateFromType(reflect.TypeOf(avar)).Interface()
}

// CreateFromType creates an instance of the given type
func (sc *ServiceContainer) CreateFromType(atype reflect.Type) reflect.Value {
	_, exists := sc.instances[atype]

	if !exists {
		sc.instances[atype] = sc.DependencyInjector.CreateFromType(atype)
	}

	return sc.instances[atype]
}

func NewServiceContainer() ServiceContainer {
	return ServiceContainer{
		DependencyInjector: NewDependencyInjector(),
		instances:          make(map[reflect.Type]reflect.Value),
	}
}
