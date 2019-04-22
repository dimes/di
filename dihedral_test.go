package main

import (
	"testing"

	"github.com/dimes/dihedral/internal/example/bindings/digen"
	"github.com/dimes/dihedral/internal/example/dbstore"
	"github.com/stretchr/testify/assert"
)

func TestExampleInjection(t *testing.T) {
	component := digen.NewServiceComponent(&dbstore.DBProviderModule{
		Prefix: "Hello",
	})
	service := component.GetService()

	assert.NoError(t, service.SetValueInDBStore("World!"))
	assert.Equal(t, "Hello World!", service.GetValueFromDBStore())
}