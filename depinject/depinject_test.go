package depinject

import (
	"fmt"
	"testing"
)

type Atype struct {
	bdep Btype
}

type Btype struct {
	myVar int
}

func NewAtype(myBtype Btype) Atype {
	return Atype{bdep: myBtype}
}

func NewBtype(mv int) Btype {
	return Btype{myVar: mv}
}

func TestInject(t *testing.T) {

	di := NewDependencyInjector()

	err1 := di.Register(NewAtype)
	if err1 != ConstructorArgsErr {
		t.Error("Expected an error")
	}

	err2 := di.Register(
		func() Btype { return NewBtype(5) },
	)
	if err2 != nil {
		t.Error("Didn't expect error")
	}
	err3 := di.Register(NewAtype)
	if err3 != nil {
		t.Error("Didn't expect error")
	}
	var instance Atype
	instance = di.Create(Atype{}).(Atype)

	fmt.Println(instance)
	if instance.bdep.myVar != 5 {
		t.Error("Expected 5")
	}
}

func TestService(t *testing.T) {

	di := NewServiceContainer()

	err1 := di.Register(NewAtype)
	if err1 != ConstructorArgsErr {
		t.Error("Expected an error")
	}

	err2 := di.Register(
		func() Btype { return NewBtype(5) },
	)
	if err2 != nil {
		t.Error("Didn't expect error")
	}
	err3 := di.Register(NewAtype)
	if err3 != nil {
		t.Error("Didn't expect error")
	}
	var instance Atype
	instance = di.Create(Atype{}).(Atype)

	fmt.Println(instance)
	if instance.bdep.myVar != 5 {
		t.Error("Expected 5")
	}
}
