# Anaphora resolution (AR) 

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
