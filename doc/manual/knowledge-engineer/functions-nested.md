# Nested functions

## Quantification

The relations `find()`, `do()` and `quant()` apply here. Check [quantification](quantification.md) for more information.

## Negation

It is possible to use "not" in a simple case.

Here's an example from the blocks world: "How many blocks are not in the box?"

"not" is modelled in the grammar:

    { rule: how_many_clause(E1) -> np(E1) copula() not() pp(E1),           sense: not(sem(4)) }

not() is a "nested structure" that wraps a relation set.

This set is specified in the example as "sem(4)". This means: the combined senses of all syntactic structures that were
linked to the fourth consequent (which is "pp(E1)").

A not() predicate can only be evaluated correctly when it is evaluated as part of a quant scope.

## Or

`or(C1, A, B)` processes `A` first. If it yields results, `or` stops. Otherwise it processes `B`.
