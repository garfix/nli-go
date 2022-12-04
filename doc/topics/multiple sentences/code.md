# Code

## Root clause

In most cases the root clause coincides with the sentence. Some sentences, like the one below

    Find a block which is taller than the one you are holding and put it into the box

can more easily be processed as two different sentences:

    Find a block which is taller than the one you are holding.
    Put it into the box.

And in fact the result of processing is exactly the same. In these cases we split the sentence up in its root clauses.

Each root clause is processed separately by the system. The root clause coincides with the `solution` / `intent`.

## Syntax

Root clauses can be marked by adding the tag `go:root_clause()`.

    { rule: imperative_clause(C) -> imperative_clause(P1) and(_) imperative_clause(P2), sense: $imperative_clause1 $imperative_clause2,
        tag: go:root_clause(P1) go:root_clause(P2) }

