package main

// every exported eleemnt must be documented

// Customer ...
type Customer struct{}

// ID ...
func (c Customer) ID() string { return "" }

// each comment should be a complete sentence that ends with punctuation

// highlight that the function intends to do; not how it does it

// document a variable or constant by indicating its purpose; documenting its contents shouldn't necessarily be public

// DefaultPermission is the default permission used by the storage engine.
const DefaultPermission = 0o644 // Need read and write access

// enough information should be provided to use the code without reading the implementation

// an exported element can be deprecated

// ComputePath ...
// Deprecated: ...
func ComputePath() {}

// also document the package

// Package math provides ...
//
// This package does not guarantee bit-identical results...
// package math



