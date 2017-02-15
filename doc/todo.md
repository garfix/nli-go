# TODO

## Releases

### 1: simple full-circle nli

* language: english
* question types: yes/no, who
- solutions
* second order predicates, aggregations
* proper nouns
* real database access (MySQL)
* a few simple questions
* simple natural language responses
* working example

### 2: omschrijving

# Known issues

## Syntax

- De syntax voor 'domain specific 2 database conversion' is dezelfde als inference, dat klopt niet :- should be -> (?) Can this not be a transformation?
- disable the use of underscores in predicates? don't want to encourage both snake and camelcasing
- should you be allowed to mix predicates of several sets? Is this confusing or a necessity to keep things manageable?
- Must be able to write whword in place of whword(); but wait, maybe we need multiple variables as well?
- is het misschien nodig om predicates en constants te namespacen? Eigenlijk is de predicate al een namespace
- find a solution for multiple (2, 3) insertions
- check if transformations are complete: non-mapping predicates may not be ignored, but should give an error

## Aggregations

- The aggregate base has currently more power than it needs. It can all bindings completely.
- Add min, max

## Domains

The domain tests are not a goal in themselves, but only help to make up test cases.

## Long term goals

- Permanent goal: improve the grammar; extend it with new phrases, make it more precise. I think there's such a thing as an NLI-English grammar that exists of grammar rules commonly used when talking to a computer. It's a small subset of full English grammar, with an emphasis on questions.
