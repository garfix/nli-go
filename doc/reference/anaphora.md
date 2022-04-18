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
- pick one from a group. "We found seven coins on the next dive. The oldes was dated 1823" (Natural Language Understanding, p348)
- one-anaphora. "How many blocks are in the box? Pick one out." ("one" refers to any block)
- bound variables. "Every student has received his grade"
- pleonastic/expletive use "It is raining". Is not a reference.
- sense and reference: "the morning star is the evening star"

## Referents

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

## Identity

If E2 references E1, then they are the same entity in some respect. If some information becomes available about E1, then the same should happen to E2.

Information about entities then exists on these layers:

- sense S
  - sentence: variables SS ("he", "himself")
  - dialog: dialog entities DG; multiple SS can reference the same DG ("a boy")
- reference R
  - database: tuples in DB; single DG can reference multiple DB
  - reality objects and concepts RR: multiple DB can reference the same RR

multiple S can reference the same R

## References

- https://en.wikipedia.org/wiki/Anaphora_(linguistics)

On using inference to determine anaphoric relations: 

Winograd schema challenge https://en.wikipedia.org/wiki/Winograd_Schema_Challenge
Coherence and-coreference - Jerry R. Hobbs (1979)
Coherence and Coreference Revisited - ANDREW KEHLER, LAURA KERTZ, HANNAH ROHDE AND JEFFREY L. ELMAN (2008)

Coreference and Bound variables

https://en.wikipedia.org/wiki/Coreference
