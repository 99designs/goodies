package depinject_test

import (
	"fmt"
	"github.com/99designs/goodies/depinject"
)

func ExampleDependencyInjector_Create() {
	type Atype struct {
		MyVar int
	}
	type Btype struct {
		a Atype
	}

	atypeBuilder := func(myvar int) Atype {
		return Atype{MyVar: myvar}
	}

	di := depinject.NewDependencyInjector()

	di.MustRegister(func() Atype {
		return atypeBuilder(5)
	})

	di.MustRegister(func(myAtype Atype) Btype {
		return Btype{a: myAtype}
	})

	newBtypeInstance := di.Create(Btype{}).(Btype)

	fmt.Println(newBtypeInstance.a.MyVar) // returns 5
}
