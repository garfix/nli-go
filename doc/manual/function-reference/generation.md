# Generation

In the generation phase, a very specific function is available:

## already_generated

Checks if the entity of the variable has already been generated (put into words).

Can be used to create an anaphoric reference in a generated sentence. 

    go:already_generated(Variable)

* `Variable`: a variable holding an id

Example:

    { rule: qp(E1) -> 'that',                                                   condition: go:already_generated(E1) }

