package depinject

import (
	"fmt"
	"testing"
)

func TestInject(t *testing.T) {

	// set up types and constructors
	type Btype struct {
		myVar int
	}
	type Atype struct {
		bdep Btype
	}
	NewAtype := func(myBtype Btype) Atype {
		return Atype{bdep: myBtype}
	}
	NewBtype := func(mv int) Btype {
		return Btype{myVar: mv}
	}

	di := NewDependencyInjector()

	err1 := di.Register(NewAtype)
	if err1 == nil {
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

	// set up types and constructors
	type Btype struct {
		myVar int
	}
	type Atype struct {
		bdep Btype
	}
	NewAtype := func(myBtype Btype) Atype {
		return Atype{bdep: myBtype}
	}
	NewBtype := func(mv int) Btype {
		return Btype{myVar: mv}
	}

	di := NewServiceContainer()

	err1 := di.Register(NewAtype)
	if err1 == nil {
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

func TestServiceDoesntCreateDepsIfAlreadyExist(t *testing.T) {
	// set up types and constructors
	type Btype struct {
		myVar int
	}
	type Atype struct {
		b *Btype
	}
	type Ctype struct {
		b *Btype
	}
	NewA := func(myBtype *Btype) Atype { return Atype{b: myBtype} }
	NewC := func(myBtype *Btype) Ctype { return Ctype{b: myBtype} }
	NewB := func() *Btype { return &Btype{} }

	di := NewServiceContainer()

	err := di.Register(NewB)
	if err != nil {
		t.Error("Didn't expect error")
	}

	err = di.Register(NewA)
	if err != nil {
		t.Error("Didn't expect error")
	}

	err = di.Register(NewC)
	if err != nil {
		t.Error("Didn't expect error")
	}

	a1 := di.Create(Atype{}).(Atype)
	c1 := di.Create(Ctype{}).(Ctype)

	if a1.b != c1.b {
		t.Error("Expected the same instance")
	}
}
