package gen

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/dimes/dihedral/resolver"
	"github.com/dimes/dihedral/structs"
	"github.com/pkg/errors"
)

// GeneratedModuleProvider is a single generated provider method on the component
// from a module source
type GeneratedModuleProvider struct {
	resolvedType *resolver.ModuleResolvedType
	assignments  []Assignment
	dependencies []*injectionTarget
}

// NewGeneratedProvider generates a provider function for the given resolved type
// The generated function has the form:
//
// func (generatedComponent *GeneratedComponent) provides_Name() *SomeType {
//     return someModule.providerFunc(
//	       component.provides_ProvidedType(),
//         InjectableFactory(component),
//     )
// }
func NewGeneratedProvider(
	resolvedType *resolver.ModuleResolvedType,
	providers map[string]resolver.ResolvedType,
	bindings map[string]*structs.Struct,
) (*GeneratedModuleProvider, error) {
	assignments := make([]Assignment, 0)
	dependencies := make([]*injectionTarget, 0)
	signature := resolvedType.Method.Type().(*types.Signature)
	for i := 0; i < signature.Params().Len(); i++ {
		param := signature.Params().At(i)
		assignment, err := AssignmentForFieldType(param.Type(), providers, bindings)
		if err != nil {
			return nil, errors.Wrapf(err, "Error generating binding for %+v", resolvedType)
		}

		assignments = append(assignments, assignment)
		dependencies = append(dependencies, newInjectionTarget(param.Type()))
	}

	return &GeneratedModuleProvider{
		resolvedType: resolvedType,
		assignments:  assignments,
		dependencies: dependencies,
	}, nil
}

// ToSource returns the source code for this provider.
func (g *GeneratedModuleProvider) ToSource(componentPackage string) string {
	moduleVariableName := SanitizeName(g.resolvedType.Module.Name)
	returnType := "target_pkg." + g.resolvedType.Name.Obj().Name()
	if g.resolvedType.IsPointer {
		returnType = "*" + returnType
	}

	var builder strings.Builder
	builder.WriteString("// Code generated by go generate; DO NOT EDIT.\n")
	builder.WriteString("package " + componentPackage + "\n")
	builder.WriteString("import target_pkg \"" + g.resolvedType.Name.Obj().Pkg().Path() + "\"\n")
	builder.WriteString(
		"func (" +
			componentName + " *" + componentType + ",\n" +
			") " + ProviderName(g.resolvedType.Name) + "() (" + returnType + ", error) {\n")

	for i, assignment := range g.assignments {
		varName := fmt.Sprintf("param%d", i)
		builder.WriteString("\t" + varName + ", err := " + assignment.GetSourceAssignment() + "\n")
		builder.WriteString("\tif err != nil {\n")
		builder.WriteString("\t\treturn nil, err\n")
		builder.WriteString("\t}\n")
	}

	builder.WriteString(
		"\treturn " + componentName + "." + moduleVariableName + "." + g.resolvedType.Method.Name() + "(\n")

	for i := range g.assignments {
		varName := fmt.Sprintf("param%d", i)
		builder.WriteString("\t\t" + varName + ",\n")
	}
	builder.WriteString("\t)")

	if g.resolvedType.HasError {
		builder.WriteString("\n")
	} else {
		builder.WriteString(", nil\n")
	}

	builder.WriteString("}\n")
	return builder.String()
}
