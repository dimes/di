// Code generated by go generate; DO NOT EDIT.
package digen

import target_pkg "github.com/dimes/dihedral/internal/example/dbstore"

func factory_github_com_dimes_dihedral_internal_example_dbstore_MemoryDBStore(generatedComponent *GeneratedComponent) *target_pkg.MemoryDBStore {
	target := &target_pkg.MemoryDBStore{}
	target.Prefix = generatedComponent.provides_github_com_dimes_dihedral_internal_example_dbstore_Prefix()
	return target
}
