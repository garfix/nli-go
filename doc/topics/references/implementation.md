# Anaphora - implementation

NLI-GO handles pronouns like "he", "she", "it", "they", but also expressions like "the block" and "the red one". These
expressions are called commonly called anaphora. Some refer to recent user input and some to a recent system response.

## The anaphora queue

To resolve an unbound variable of a sentence, it is matched to all variables in the anaphora queue, until a match is found.

The anaphora queue is built on demand, whenever it is needed. It is built from the entities stored in the clauses of the dialog context.

Each clause has a list of entities ordered as they appear in the sentence.

The anaphora queue is built from the last clause encountered to the first, and the entities within a clause in order. In general an entity that appears earlier in a sentence is more likely to be the referent.

It is also possible to tag entities as subject and object, and there functions will be used to determine the order of appearance. The subject is more important than the object, and this in turn is more important than other entities.

The anaphora queue is extended at the same of anaphora resolution itself. This is necessary because of intrasentential anaphora: a reference to an entity within the same sentence.

The anaphora queue is also extended when the system creates an answer sentence. At that time it creates a new clause, and its entities.

Each entity is added to the anaphora queue only once.

## The anaphora resolver

This resolver traverses the parse tree, and fully processes a node before proceeding to its children.

For each node it checks if it contains:

- a reference ("him", "the red block"): `go:reference(E1, person)`
- a "labeled" reference ("it"): `go:labeled_reference(E1, 'it', object)`
- one anaphora ("pick one"): `go:reference_slot(E1)`
- reflection ("himself") `go:reflective(E1)`

To find a referent, it checks

- does the sort in the reference predication match the one in the referent?
- if there is agreement in person, number, etc (using tag:category)
- if the reflection is correct (checking co-occurance)
- for definite references (with a non-empty definition): if the referent has an id, the definition should resolve with it (NB: we don't check for nonanaphoric references)

If the referent is a group, it tries to match each member of the group, the same way an individual is matched. References to the complete group are not yet supported.

In case if definite references, where an id is certainly available, we bind the reference variable to the id of the referent.
In all other cases, we replace the reference variable with the referent variable. This variable may or may not have a binding.

Ambiguity score:

- the older the clause, the lower the score (-10)
- a subject gets 5 points, an object 3 point
- first entity in the clause gets 1 point

## Features

When the sentence is parsed, the system does not only build the representation of the intention, it also produces "features" for each of the entities.

These features are used to constrain the options in anaphora resolution. They are:

- sort: `person`, `car`, `event`, ...
- gender: `male`, `female`, `neuter`
- number: `singular`, `plural`
- reflexivity: `true` ("himself"), `false` ("him")
- definiteness: `determinate` ("the") `indeterminate` ("a")
- resolved: `true`, `false` (a forward reference is unresolved for some time)

