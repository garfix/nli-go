# Bases

There are several types of knowledge bases that can be used in the inference process.

## Fact base

A fact base matches an input relation to one or more relations in a database. Its main function is

    Bind(goal mentalese.Relation) ([]mentalese.RelationSet, []mentalese.Binding)



## Rule base

A rule base treats an input relation as a goal, and returns sets of subgoals that need to succeed in order for the goal to succeed.

Its main function is

    Bind(goal mentalese.Relation) ([]mentalese.RelationSet, []mentalese.Binding)

## Multiple bindings base / Aggregate base

The main function

    Bind(goal mentalese.Relation, bindings []mentalese.Binding) ([]mentalese.Binding, bool)

An example of an aggregate function is

    numberOf(X, N)

In the function Bind, two things happen:

* the number of different value of X in bindings is calculated
* all bindings are extended with a value for N

Bind returns false if ''goal'' is not one of the aggregate functions of the base.
