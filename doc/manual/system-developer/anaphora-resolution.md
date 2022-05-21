# Anaphora resolution (AR)

NLI-GO handles pronouns like "he", "she", "it", "they", but also expressions like "the block" and "the red one". These
expressions are called commonly called anaphora. Some refer to recent user input and some to a recent system response.
It also handles references to entities mentioned earlier within the same sentence (intra-sentential).

## The anaphora queue

NLI-GO uses a structure called the anaphora queue to store the most recent references. The queue is located in the
dialog context.

The queue is simply a queue (first in first out) of entity variables. The variables are grouped by the clause they belong to.

More precise, it is a queue of entity reference groups. Multiple entities can be added as a single group. The result of
a query for example, that exists of 3 persons. So that you can later refer to these as "they".

## Solution

The variable where the query results are stored must be specified in the field "result" of a solution:

    result: E1

This is how the system knows which entities to refer to later.

## Adding entities

Any entity that is named in the input is added to the queue. This allows a user to type the name just once and
henceforth refer to him/her with a pronoun.

Any entity that is part of the result set of a question is added to the queue. This allows a user to refer to a previous
response. Only the entities designated by the variable named in the "result" field from the solution are added.

If the entity had been part of the queue before, it would be removed from the queue before being added again at the
front. Each entity can be in the queue only once.

## Features

When the sentence is parsed, the system does not only build the representation of the intention, it also produces "features" for each of the entities.

These features are used to constrain the options in anaphora resolution. They are:

- sort: `person`, `car`, `event`, ...
- gender: `male`, `female`, `neuter`
- number: `singular`, `plural`
- reflexivity: `true` ("himself"), `false` ("him")
- determinacy: `determinate` ("the") `indeterminate` ("a")
- resolved: `true`, `false` (a forward reference is unresolved for some time)

These Features should be stored in the dialog context, but only when it is certain that this interpretation of the sentence is selected.

## Replacing variables

The implementation of anaphora resolution we take here involves replacing the variable of the reference with the variable of the referent.

This implementation effectively adds the constraint that the reference variable equals the referent variable [B = A] in Discourse Representation Theory. Only working with equalities is very hard, and reducing the number of variables is logically equivalent.

Patrick: if you're asking why we need all this: The reason we do all of this is to allow referring to unbound referents. This was not possible before.

## Anahpora resolution step

There needs to be a separate AR step. The algorithm is like this:

- go through the relational structure of the sentence, quant by quant, from bottom to top
- handle `back_reference`-like relations as special cases
- handle all quants E1 (like "the box")
    - go through all entities E2 in the queue
    - check if the features of E1 match those of E2
    - check if E2 matches the scope of E1
    - if E2 is a group, check its members individually
    - check for unresolved references
    - when in doubt, use parallellism (a reference in subject position is more likely to refer to entity that is also in subject position; idem for object position)
    - a local reference (same clause) is more likely than a remote reference

- anaphoric / nonanaphoric: the distinction is not made: all entities are treated as anaphoric
- concepts "whats the action radius of an electric car?" can be treated like objects: they must be represented in the database
- forward references: "He picked up a block. Jack."
- "My car broke down. The engine failed." - use frames (this will be a future extension)
- "the morning star is the evening star" - will not support (see remark at 2022-04-02)

![entities](../../diagram/entities2.png)

## Splitting up a referent group

In this interaction:

    What does the box contain?
    The blue pyramid and the blue block

the two objects are bound to the same variable.

But when only one of them is referenced

    What is the pyramid supported by?

the variable in the new sentence should not be replaced by the one holding the 2 objects.

This is an interesting issue. The two objects should be referencable together ("move them out of the box"). But in our case it seems as if a new referent is created out of the existing referent set:

    E18: [blue-block1, blue-pyramid]
    =>
    E23: blue-pyramid

Solution: if a reference refers to a single entity from a referent group, keep the reference variable unchanged, but bind it single referent's value.

