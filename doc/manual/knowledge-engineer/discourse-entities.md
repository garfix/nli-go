# Discourse entities

Discourse entities is the solution for having indefinite descriptions (i.e. "a block", "someone") and referring to them ("I saw a big car. It had a painting of a dragon on the side.")

The entity is represented by a variable, that may or may not be bound to a variable.

An indefinite description is represented by a variable without a value. 

## Dialogizer

The dialogizer class turns a parse tree into a tree that has variable names that are unique in the discourse / dialog.

I.e. 

    S => Sentence$12
    E1 => E$813

This turns each of the entities introduced by the parser into a unique discourse entity.

## Ellipsizer

The ellipsizer completes the parse tree with bits of parse tree from other trees. It keeps the variables (discourse entities) from the original trees, so that the variable-to-id bindings remain active in the new parse tree. 

This may mean that the variables in the rest of the new parse tree may be changed to the ellipsized variables as well. 

## Back reference

The back reference operation checks if the variable has been bound before, and if so, take this existing binding.

A back reference is an discourse entity that refers to another discourse entity. However, they are different variables, because there is, in general, no way of knowing beforehand what entities are the same. 

## References to variables in stead of ids?

A back reference, or any anaphoric reference for that matter, currently copies the id of its referent. In a future extension, the back reference may hold the _variable_ of its referent. This would be useful for cataphoric (forward) references: once the id of the entity becomes known, earlier references will inherit this binding automatically.


