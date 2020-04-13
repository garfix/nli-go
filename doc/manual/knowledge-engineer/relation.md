# Relation

The relation construct is very versatile. It has a remarkable set of uses. And it is such a simple form: it has a name
and some arguments. 

    relation_name(argument1, argument2)

The arguments can be variables, values or even sets of relations.

    relation_name(argument1, [relation_name() relation_name()])  
    
The variables in the relation can be bound by both a binding

    { A: 11, B: `b6`, C: 'Byron' }

## Fact

A relation can represent a factual relation between two objects. Like in this case where is represents the fact that
`r11` is the father of `r23`.

    father_of(`r11`, `r23`)

## Link

Two relations can be linked together by a common argument. The argument is usually a variable:

    grandparent(A, C) :- parent(A, B) parent(B, C);

Here `B` is the linking argument.

## Test

A relation can be used to test if a certain condition holds:

     greater_than(X1, X2)
     
If the test does not yield any results, processing stops.

## Search

A relation can be used as a template to select all relations that fulfil its argument values

    parent_of(`b1`, X)

## Goal

A single relation can express a goal

    pick_up(A)

The goal can be reached by processing subgoals

    pick_up(A) :- position(A, P) move_hand(P) grab(P) up(P, Q) move_hand(Q); 

## Function

A function has a number of inputs (arguments), and an output (the return value). It can be represented as a relation, by
treating both arguments and return value as arguments of the relation

    greater_than(X1, X2)
    
Some functions can have any number of arguments

    concat(Result, X1, X2, ...)
    
## Dependencies

An argument that consists of a relation set creates a dependency on these other relations

    c(argument1, [a() b()])
    
In this case a() and b() must be processed before c() can be processed, and this creates a dependency relation between
`c()` and `[a() b()]`.

## Assertion of a New Fact

The `assert` relation has a special built-in function. It causes the relation set of its single argument to be added to
the database.

    assert([on_top_of(A1, B1)])

## Quantification

This relation expresses the sentence "all children have a father":

    quant(
        Q5, [all(Q5)],
        R5, [child(R5)], 
        [have_father(R5, P1)]
    )

The `quant` relation has 5 parameters. These form the `range`, the `quantifier` and the `scope` of the quantification.

When processing the quant, first all values for the range `R5` are collected. Then for each value, the scope
`[have_father(R5, P1)]` is processed. Finally the results are compared to the quantifier `Q5`. In the case of "all" this
means that the number of R5 values is compared to the number of scope values. It should match.

## Negation

It is possible to use "not" in a simple case.

Here's an example from the blocks world: "How many blocks are not in the box?"

"not" is modelled in the grammar:

    { rule: how_many_clause(E1) -> np(E1) copula() not() pp(E1),           sense: not(sem(4)) }

not() is a "nested structure" that wraps a relation set.

This set is specified in the example as "sem(4)". This means: the combined senses of all syntactic structures that were
linked to the fourth consequent (which is "pp(E1)").

A not() predicate can only be evaluated correctly when it is evaluated as part of a quant scope.
