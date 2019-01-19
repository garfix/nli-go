# Quants

When you work with this system you will come across quants. What are they? They are used for aggregation relations like 'all', 'at least two', and none.

A quant looks like this

    quant(
        R5,
        [isa(R5, child)], Q5,
        [isa(Q5, how) isa(Q5, many)],
        [subject(S5, R5)]
    )

These are the parts:

    quant(
        RangeVariable,
        RangeRelations,
        QuantifierVariable,
        QuantifierRelations,
        ScopedRelations
    )

A quant has a range, a quantifier, and scoped relations.

The _Range_ is a domain of objects to which this quant applies. It can be men, houses, orders, or three legged dogs. The _RangeVariable_ is the variable used for the range. The _RangeRelations_ limit the domain.

The _Quantifier_ is a specifier like 'all', 'some', 3, or 'at most two'. The _QuantifierVariable_ is used for the _QuantifierRelations_. These relations describe the quantifier. Using a set of relations allows this quantifier to be more complex than just 'all' or 'none', and be like 'two to five'.

The _ScopedRelations_ are the body part of the quant.

## Evaluation

When the quant is executed, or evaluated, the system goes through 3 steps.

First, the range relations are processed. This results in a set of bindings for the range variable. For instance, when the range consists of men, and the range variable is R, R may have as bindings john, jack and jill.

Then the scoped relations are processed. This is done in a loop for each of the bindings in the range. This results in a set of variable bindings.

Finally, these scoped relations bindings are evaluated with respect to the quantifier. When the quantifier is 'at least two' and the number of variable bindings of the scoped variables is one or less, all bindings for the quant are dropped.
