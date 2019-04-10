package main

import (
	"fmt"

	"github.com/dimes/dihedral/example/bindings/digen"
	"github.com/dimes/dihedral/example/dbstore"
)

func main() {
	component := digen.NewServiceComponent(&dbstore.DBProviderModule{
		Prefix: "Hello",
	})
	service := component.GetService()

	if err := service.SetValueInDBStore("World!"); err != nil {
		panic(err)
	}

	fmt.Println(service.GetValueFromDBStore())
}
