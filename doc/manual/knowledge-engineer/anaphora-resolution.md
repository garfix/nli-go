# Anaphora resolution

NLI-GO keeps track of the entities that it last processed in a structure called "anaphora queue". This queue consists of
the id's of these entities.

Whenever a quant is processed, the system will first try to resolve the range of the quant with each of
these ids filled in. When one of these ids gives a match, this id will be used as the range of the quantification.

In order to allow pronouns like "he", "she" and "it" in the input, you need to model pronouns in a way that reflects
their function as a quantification:

    { rule: pronoun(E1) -> 'it',                                           sense: go:back_reference(E1, none) }
    
## Resolving abstract nouns    

To resolve an abstract noun like "one" (as in "Put a small one onto the green cube") to the concrete noun "block", specify "one" as a `sortal_back_reference()`

    { rule: noun(E1) -> 'one',                                           sense: go:sortal_back_reference(E1, none) }
    
When the `sortal_back_reference` is processed, it looks into the anaphora queue for the sort of the latest referent. When this sort is, say, `block`, the function will look for the `entity` relation set that belongs to the sort. Located in `sorts.yml` you will write   

    block:
      entity: block(Id)  

This tells us that the sort `block` is represented by the relation (set) `block()` in the domain. The variable `Id` will be used to identify the block.
