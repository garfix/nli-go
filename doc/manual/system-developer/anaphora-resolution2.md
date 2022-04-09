# Anaphora resolution 

## Features 

When the sentence is parsed, the system does not only build the representation of the intention, it also produces "features" for each of the entities.

These features are used to constrain the options in anaphora resolution. They are:

- gender: male, female, genderless: "he" P`gender(E1, male)` ; "she" P`gender(E1, female)` ; "it" P`gender(E1, neuter)`
- number: singular, plural: "it" P`number(E1, plural)`
- animatedness: person, thing: <proper name>, "boy" P`sort(E1, person)`, "object" P`sort(E1, object)`
- selectional restrictions: "pick up" -> S`pick_up(E1, E2)`, P`sort(E1, person)`, P`sort(E1, object)`
- reflexivity: "himself" P`reflexive(E1, true)`
- determinate: "a" P`determinate(E1, false)` "the" P`determinate(E1, true)`
- domain: ?
- event: "why did you do that" "why" `why(P1)` "do that" `reference(E1, type(E1, event))`

These Features should be stored in the dialog context, but only when it is certain that this interpretation of the sentence is selected.



- anaphoric / nonanaphoric: the distinction is not made: all entities are treated as anaphoric
- concepts "whats the action radius of an electric car?" can be treated like objects: they must be represented in the database
- productions can be added to the _discourse entity_, not the sentence entity
- "My car broke down. The engine failed." - use frames
- "the morning star is the evening star" - will not support (see remark at 2022-04-02)

![entities](../../diagram/entities.png)

To do

