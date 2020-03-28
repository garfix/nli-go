# Todo

* phase out the lexicon
* remove the lexicon from the generation phase
* remove cost function and everything related
* document 'explicit references'
* remove solution routes altogether; rebuild the solver
* map relations to database => not many-to-many, but 1-to-many
* NewRelationizer: pointer, geen kopie maken
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

## Aggregation

My system works with a mix of breadth-first and depth first. I have not thought this through well enough.

Want to have aggregation functions that work reliably.

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

## Optimizer

Check out the optimizer. It can probably be improved beyond the level of using table sizes to calculate cost.

write(A, B) number_of(B, N) => number_of() should come later based on a dependency. This is not worked out at all.

## Domains

The domain tests are not a goal in themselves, but only help to make up test cases.

## Long term goals

- Permanent goal: improve the grammar; extend it with new phrases, make it more precise. I think there's such a thing as an NLI-English grammar that exists of grammar rules commonly used when talking to a computer. It's a small subset of full English grammar, with an emphasis on questions.
