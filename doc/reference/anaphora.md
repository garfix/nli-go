# Anaphora resolution

What is the algorithm for anaphora resulution? Simple question. Not so simple to answer. What can we say about it?

## The facts

Any noun phrase is a reference. It is a description of an entity and what that entity is, needs to be discovered. And since anaphora is so common for noun phrases, it makes sense to treat any noun phrase as anaphoric.

Anaphora has the following aspects

- references to discourse entities (anaphoric) vs references without discourse entities; these refer to the world (nonanaphoric)
- references to determinate entities ("the", names) (definite references) vs indeterminate entities ("a") (indefinite references)
- definite vs indefinite references ("the car", "a car") ("If Pedro owns some donkey, he beats it.")
- references to concepts ("whats the action radius of an electric car?") 
- plural vs singular references ("they", "she")
- pronouns ("it") vs descriptions ("the blue car")
- the type of entity is often specified ("the car")
- reflexive ("himself") vs nonreflexive ("him")
- within sentence references (intrasentential) vs references to other sentences (intersentential)
- back-references (anaphoric) vs forward references (cataphoric)
- direct references ("the car") vs indirect references (referring to an entity that can be deduced) ("My car broke down. The engine failed.")
- may require complex interpretation via world knowledge / selectional restrictions ("The city councilmen refused the demonstrators a permit because they [feared/advocated] violence.")
- it may refer to NP's but also to VP's and S's.
- one-anaphora. "How many blocks are in the box? Pick one out." ("one" refers to any block)
- bound variables. "Every student has received his grade"
- pleonastic/expletive use "It is raining". Is not a reference.

## Entities

What are the linguistically relevant features of the entities (anaphors) refered to? They have

- a gender: male, female
- animatedness: animate, inanimate
- number: singular, plural (entities that are known as a group)
- a type: car, person, thing, activity
- selectional restrictions: (the things they can do, "birds fly", "cars break down")
- a domain: in the discourse, in a conceptual structure, in the world

## References

What are the features of the references? They have

- a domain: in the discourse ("he"), conceptually related entities ("the engine of the car"), in the world ("My car")
- a number: singular, plural
- a phrase type: noun, verb
- animatedness: animate ("she"), inanimate ("it")
- a gender: male, female
- a specification: car, person, thing, on the left, that my mother liked
- a focus: single entity (definite), set of entities (indefinite)
- a lexical scope: reflexive (same clause), intersentential, intrasentential
- a direction: anaphora (backward), cataphora (forward)
- a phrase type: noun, verb

## Heuristics

Using inference to resolve anaphora quickly becomes very complicated (see references below for examples). It is interesting to know that this form of resolution leads to much more certain resolution, but in any case it would require the user of an NLI system to go to great trouble to handle the cases of anaphora resolution which she hoped would be tackled automatically. A second best solution is to go with the heuristics that provide good results in most cases. Which heuristics are known?

Hard constraints

- gender agreement ("she" can't match "the boy")
- number agreement ("they" can't match "a block")
- animatedness agreement ("it", can't match "the man")
- selectional restrictions: relations have argument type restrictions (a table can't be "picked up", a block can't be put "into a block")
- reflexivity in lexical scope: ("himself" can only match an entity in the same clause, and not a subclause)
- increasing focus: ("a block" can't refer to an entity before referred to as "the block") 
- centering: 

Soft constraints

- paralellism: a reference in subject position is more likely to refer to entity that is also in subject position; idem for object position

## Implementation

- gender: male, female, genderless: "he" P`male(E1)` ; "she" P`female(E1)` ; "it" P`neuter(E1)` 
- number: singular, plural: "it" S`count(E1, 1)` ; "they" S`greater_than(count(E1), 1)` OF `plural(E1)`
- animatedness: person, thing: <proper name>, "boy" P`person(E1)`, "object" P`thing(E1)`
- selectional restrictions: "pick up" -> S`pick_up(E1, E2)`, P`person(E1)`, P`object(E1)`
- reflexivity: "himself" P`reflexive(E1)` 

De grammar moet een onderscheid maken tussen constraints(sense, S) en productions (P). Productions zijn implicaties die aan de context worden toegevoegd op het moment dat de geparste zin geaccepteerd wordt.

## References

- https://en.wikipedia.org/wiki/Anaphora_(linguistics)

On using inference to determine anaphoric relations: 

Winograd schema challenge https://en.wikipedia.org/wiki/Winograd_Schema_Challenge
Coherence and-coreference - Jerry R. Hobbs (1979)
Coherence and Coreference Revisited - ANDREW KEHLER, LAURA KERTZ, HANNAH ROHDE AND JEFFREY L. ELMAN (2008)

Coreference and Bound variables

https://en.wikipedia.org/wiki/Coreference
