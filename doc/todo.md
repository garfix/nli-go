# Todo

## Next release

* Create a next-word autosuggest
* Create an autosuggest app
* Find or create a javascript autosuggest line editor and create an example web app
* Build an application from a configuration file

## omschrijving (?)

* grouped relations in matcher and solver: (), and, or, not
* aggregations, handle literature cases
* handle common questions
* declaratives and imperatives that update the database
* Blocks World examples


* Names with and without insertion
* name(A, F, firstName) !name(A, I, insertion) name(A, L, lastName) join(N, ' ', F, L) => name(A, N);
* married(A, B) :- married(B, A)

## Syntax

- Perhaps replace the syntax of functions like numberOf(N, X) to
    numberOf(X: N)
    join('', firstName, lastName: name)
    join('', firstName, lastName -> name)
    name = join('', firstName, lastName)
- disable the use of underscores in predicates? don't want to encourage both snake and camelcasing
- no: disable capitals! underscores are much better readable!
- should you be allowed to mix predicates of several sets? Is this confusing or a necessity to keep things manageable?
- Must be able to write whword in place of whword(); but wait, maybe we need multiple variables as well?
- is het misschien nodig om predicates en constants te namespacen? Eigenlijk is de predicate al een namespace
- find a solution for multiple (2, 3) insertions
- check if transformations are complete: non-mapping predicates may not be ignored, but should give an error

## Developing

* Introduce a "show intermediate representations" mode that shows all temp results

## Aggregations

- The aggregate base has currently more power than it needs. It can all bindings completely.
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

## Solutions

If you remove 

    focus(E1) => focus(E1);

all hell breaks loose, even though it should be removed.

In this solution:

		condition: act(question, who) child(A, B) focus(A),
		preparation: name(A, N),
		answer: name(A, N);

If there are no children, or if the DB mapping is not defined, preparation is still executed (and needs to be so, for 'exists' clauses), and yields ALL names

## Domains

The domain tests are not a goal in themselves, but only help to make up test cases.

## Long term goals

- Permanent goal: improve the grammar; extend it with new phrases, make it more precise. I think there's such a thing as an NLI-English grammar that exists of grammar rules commonly used when talking to a computer. It's a small subset of full English grammar, with an emphasis on questions.
