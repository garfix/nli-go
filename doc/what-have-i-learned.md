# What have I learned?

Since this is a research project, its important result is knowledge. What have I learned from it?

## Representing a group of entities

A sentence entity variable regularly refers not just to a single entity, but to multiple entities. 

    "John, Jack and Joe played football."

This may not just be a simple collection, like in the sentence above, but it may have a more complex internal structure as well.

    "John and Jack, or Jack's brother, played football."

When a question is answered, each of the entities has its own binding.

When the group of entities is to be used in a single binding, each entity must have its own variable.

"Jack's brother" may not be bound to an actual ID. In that case, it is only represented by its dialog variable. 

Lesson: use entity groups in your examples as soon as you can. Don't start by assuming that an entity is always a single person.

## Discourse entities

A discourse entity may be

- unbound: a boy (E1: null)
- bound by one database entry, or entries in multiple databases (E1: 13, 211)
- refer to another entity ("it", E1: E2)

## Parsing multiple sentences

An old problem of a parser is that it can only parse one sentence at a time. Sometimes people enter two or more sentences on the same line, however. How can we solve this?

The sentence terminator period is problematic. Sentences contain periods in abbreviations as well.

There is however a simple solution: parse multiple sentences just like single sentences:

    { rule: S(P) -> S(P1) S(P2) }

This solution has a single problem: where do the subsentences start? In this example this is simple. But what if you want to treat sentences with top-level conjunctons in the same way?

    Find a block which is taller than the one you are holding and put it into the box.

    { rule: S(P) -> imperative(P1) 'and' imperative(P2) }

To solve this we can manually tag where the root clauses start:

    { rule: S(P) -> imperative(P1) 'and' imperative(P2),    tag: go:root_clause(P1) go:root_clause(P2) }

