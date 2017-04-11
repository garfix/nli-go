# Release remarks

## 1.1 Quantifier Scoping

* ok: change 'determiner()' to 'dp()' (syntactic rewrite)
* ok: introduce relation set as an argument type
* ok: change 'determiner(E1, D1)' to 'quantification(E1, [], D1, [])'
* ok: introduce a new step that subsumes determiner's relations
* nok: introduce a generic step that converts to clumsy verb predicates to easier predicates. All occurrences of isa(Q1, PRED) subject() object() are turned into PRED().
* nok: introduce a step that helps remove vagueness ("have" is vague)
* ok: create a quantifier scoper that turns a relation set into a scoped relation set
* ok: extend the answerer to make it answer scoped relation questions
* ok apply quantifier scoper to integration tests
* extend README
