---
title: "Code Design"
date: "2024-08-03"
tags: ["Programming"]
---


## SOLID


## Comments

## Documentation
- Document exported elements.


## Repositories


## Data Transfer Objects


## Encapsulation
- Avoid excessive use of getters and setters.


## Inheritance


## Interfaces
- Create an interface when it is needed, not when you forsee needing it, or at least prove the abstraction to be a valid one.
- Keep interfaces on the client side.
- Functions should accept interfaces and return concrete implementations.

## Coupling & Cohesion


## Polymorphism


## Control flow
- Avoid shadowing variables.
- Avoid deep nesting and keep the happy path aligned on the left; return as early as possible.

## Exceptions


## Dependencies


## Linters and formatters
- To improve code quality and consistency, use linters and formatters.

## Abstraction
- Abstractions should be discovered, not created.


## Generics

## Project structure
- Standardize a layout.
- Give packages meaningful and specific names.
- Avoid creating packages like common, util or shared.


## Primitive Types
- Integer overflows (https://en.wikipedia.org/wiki/Ariane_flight_V88).
- Comparing floats using a delta.


## Stack & Heap


## Dependency Injection


