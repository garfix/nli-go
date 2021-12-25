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

The dialog representation itself is a small database that consists of the entities that are used in the dialog.

Anaphora has the following aspects

- references to discourse entities (anaphoric) vs references without discourse entities; these refer to the world (nonanaphoric)
- references to determinate entities ("the", names) (definite references) vs indeterminate entities ("a") (indefinite references)
- plural vs singular references
- pronouns ("it") vs descriptions ("the blue car")  
- reflexive ("himself") vs nonreflexive ("him")
- within sentence references (intrasentential) vs references to other sentences (intersentential)
- back-references (anaphoric) vs forward references (cataphoric)
- direct references ("the car") vs indirect references (referring to an entity that can be deduced) ("related objects", NLU1, p346)
- it may refer to NP's but also to VP's and S's.
