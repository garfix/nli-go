# Considerations

Each entity is represented by a variable in the dialog. The variable may hold multiple values.

## Multiple bindings or single binding with list?

A variable may have multiple bindings, and this is the best way to deal with groups while processing a sentence.

After the intent is executed, and an answer needs to be generated, it's easier to deal with a group as a list of ids. This way we only need to generate a sentence based on a single binding.

## Problem

When the sentence is complete, the group is added to the dialog bindings as a list. The next sentence then needs to deal with this list, but it doesn't. This is currently an unresolved issue.

## Representing a group of entities

A sentence entity variable regularly refers not just to a single entity, but to multiple entities.

    "John, Jack and Joe played football."

This may not just be a simple collection, like in the sentence above, but it may have a more complex internal structure as well.

    "John and Jack, or Jack's brother, played football."

When a question is answered, each of the entities has its own binding.

When the group of entities is to be used in a single binding, each entity must have its own variable.

"Jack's brother" may not be bound to an actual ID. In that case, it is only represented by its dialog variable.

Lesson: use entity groups in your examples as soon as you can. Don't start by assuming that an entity is always a single person.

## Groups and elements

A group has elements, and it is useful to store the relation between the group as-an-entity, and its elements.

    stack_of(S, E)

Note that each AND and OR also forms a group, and that its operands are different entities.

    and(A, X, Y)
    or(O, X, Y)
