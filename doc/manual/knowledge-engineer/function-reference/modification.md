# Modification functions

These functions allow you to modify a database by asserting and retracting relations. The `writeMap` of a fact base determines which predicates are writable.

Rules are written to a rulebase.  

## assert

Writes a relation to the database or adds a rule to the rule base.  

    go:assert(P)
    
* `P`: a relation (i.e. father(`luke`, `darth`))

The relation is offered to all fact bases; but only the ones that have the relation defined as head in their write map will write it. But only after the relation is converted to one or more rows of database tables.


    go:assert(R)

* `R`: a rule (i.e. fly(X) :- bird(X))

Note the following restrictions to the use of adding rules, that currently exist:

* A rule is only appended at the end of the first rule base that was added to the problem solver
* Since rule bases are in-memory, the rules are not saved when the system instance end.


See [common-sense-reasoning](common-sense-reasoning.md) for examples of default rules and exceptions.

## retract

Deletes a relation from the database  

    go:retract(P)
    
* `P`: a relation (i.e. father(`luke`, `darth`))

