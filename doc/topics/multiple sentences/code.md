# Code

An old problem of a parser is that it can only parse one sentence at a time. Sometimes people enter two or more sentences on the same line, however. How can we solve this?

The sentence terminator period is problematic. Sentences contain periods in abbreviations as well.

There is however a simple solution: parse multiple sentences just like single sentences:

    { rule: S(P) -> S(P1) S(P2) }

This solution has a single problem: where do the subsentences start? In this example this is simple. But what if you want to treat sentences with top-level conjunctons in the same way?

    Find a block which is taller than the one you are holding and put it into the box.

    { rule: S(P) -> imperative(P1) 'and' imperative(P2) }

To solve this we can manually tag where the root clauses start:

    { rule: S(P) -> imperative(P1) 'and' imperative(P2),    tag: go:root_clause(P1) go:root_clause(P2) }

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

