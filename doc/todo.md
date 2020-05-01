# Todo

* actually evaluate the quantifier / generalized quantification
* seq => and
* quant => for
* for(range(), quant()) ?
* append or prepend child senses? ; see ## Semantic composition in entity-grammar.md
* namespaces for relations
* do not allow zero valued predicates in the grammar
* document 'explicit references'
* SparqlFactBase: todo predicates does not contain database relations (just ontology relations), so this needs to be
    solve some other way
* entity type (multiple) inheritance

## Misc

* Provide common user queries
* Separate interfaces (api) from implementations (model)
* grouped relations in matcher and solver: (), and, or, not
* Blocks World examples
* Solutions: do not just create cases for 0 or any results, but allow for arbitrary conditions (a relation set), for the response

## Rules

Test if this works or make it work:

    married_to(A, B) :- married_to(B, A);

## Syntax

- Perhaps replace the syntax of functions like number_of(N, X) to
    number_of(X: N)
    join('', firstName, lastName: name)
    join('', firstName, lastName -> name)
    name = join('', firstName, lastName)
- Namespace predicates?
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

## Domains

The domain tests are not a goal in themselves, but only help to make up test cases.
