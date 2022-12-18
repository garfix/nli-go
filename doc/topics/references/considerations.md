# Anaphora - considerations

If you're asking why we need all this: The reason we do all of this is to allow referring to unbound referents. This was not possible before.

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

## Reflective references

The word "himself" is a reflective reference, wheras "him" isn't. The difference between these words is that "himself" somehow has a very narrow scope, while "him" can span a great length. But what exactly determines this scope?

In "Syntax and Semantics" - Dalrymple (p.280) we find:

  "The antecedent of the English reflective himself must appear in the Minimal Complete Nucleus containing the pronoun".

Leaves us to determine this Nucleus. But since this definition presumes an f-structure, which we don't have, we must cast it in another form.

Currently I'm using this definition:

  "The Minimal Complete Nucleus can be described as the relation set that forms the body of a do/check predication"

Two variables that occur in the same nucleus are called co-arguments.

## References

- https://en.wikipedia.org/wiki/Anaphora_(linguistics)

On using inference to determine anaphoric relations: 

Winograd schema challenge https://en.wikipedia.org/wiki/Winograd_Schema_Challenge
Coherence and-coreference - Jerry R. Hobbs (1979)
Coherence and Coreference Revisited - ANDREW KEHLER, LAURA KERTZ, HANNAH ROHDE AND JEFFREY L. ELMAN (2008)

Coreference and Bound variables

https://en.wikipedia.org/wiki/Coreference
