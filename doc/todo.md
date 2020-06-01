# Todo

* should boolean functions have P1 as argument? different or for read/write?
* check if the nested functions are called correctly
* namespaces for relations
* do not allow zero valued predicates in the grammar
* SparqlFactBase: todo predicates does not contain database relations (just ontology relations), so this needs to be solve some other way
* entity type (multiple) inheritance

## Misc

* Separate interfaces (api) from implementations (model)
* Blocks World examples

## Rules

Test if this works or make it work:

    married_to(A, B) :- married_to(B, A);
    
* Allow the dynamically added rules to be saved (in the session).
* Specify which predicates a rule base allows to be added.    

## Syntax

- Perhaps replace the syntax of functions like number_of(N, X) to
    number_of(X: N)
    join('', firstName, lastName: name)
    join('', firstName, lastName -> name)
    name = join('', firstName, lastName)
- should you be allowed to mix predicates of several sets? Is this confusing or a necessity to keep things manageable?
- Must be able to write whword in place of whword(); but wait, maybe we need multiple variables as well?

## Aggregations

- Add min, max

## Relations

Find a way to ensure completeness of information about all relations used in a system. An interpretation should not even be attempted if not all conversions have a chance to succeed.

* convert number words into numbers

# Multiple languages

- Introduce a second language
- Constants like "all", are they universal, or english?

# Quantifier Scoping

- Make "more than" "less than" work
- A range itself can contain quantified nouns (the oldest child in every family). The algorithm is not up to it. (See CLE)
