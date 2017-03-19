# Todo

## This release

todo

* ok: change 'determiner()' to 'dp()' (syntactic rewrite)
* ok: introduce relation set as an argument type
* ok: change 'determiner(E1, D1)' to 'determiner(E1, [], D1, [])'
* introduce a new step that subsumes determiner's relations
* convert number words into numbers
* introduce a generic step that converts to clumsy verb predicates to easier predicates. All occurrences of isa(Q1, PRED) subject() object() are turned into PRED().
* introduce a step that helps remove vagueness ("have" is vague)
* create a quantifier scoper that turns a relation set into a scoped relation set
* extend the answerer to make it answer scoped relation questions

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

## Solutions

In this solution:

		condition: act(question, who) child(A, B) focus(A),
		preparation: name(A, N),
		answer: name(A, N);

If there are no children, or if the DB mapping is not defined, preparation is still executed (and needs to be so, for 'exists' clauses), and yields ALL names

## Domains

The domain tests are not a goal in themselves, but only help to make up test cases.

## Long term goals

- Permanent goal: improve the grammar; extend it with new phrases, make it more precise. I think there's such a thing as an NLI-English grammar that exists of grammar rules commonly used when talking to a computer. It's a small subset of full English grammar, with an emphasis on questions.
