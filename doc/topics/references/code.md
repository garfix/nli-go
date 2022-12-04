# Anaphora - code

## Pronouns

To handle pronouns like "he", "she" and "it" in the input, you need to tag them as a reference:

    { rule: pronoun(E1) -> 'it',                                           sense: object(E1), tag: go:reference(E1) }
    { rule: pronoun(E1) -> 'she',                                          sense: person(E1) female(E1), tag: go:reference(E1) }
    
## Definite references

The determiners "the" can mean that an NP is a reference. "I found a block and a pyramid. Give me the block." "the" refers to a single entity that may be in the dialog, but it doesn't have to be. If the reference can be resolved, it will be. If it cannot, the system will ask: "I don't understand which one you mean".

    { rule: np(E1) -> 'the' nbar(E1),                                      sense: go:quant(none, E1, $nbar), tag: go:reference(E1) }

## One anaphora    

To resolve words like "one" (as in "Put a small one onto the green cube") to the concrete noun "block", specify "one" as a `sortal_reference()`

    { rule: noun(E1) -> 'one',                                           tag: go:sortal_reference(E1) }
    
When the `sortal_reference` is processed, the system looks into the anaphora queue for the sort of the latest referent. When this sort is, say, `block`, the system will add the relation for the sort `block` (`block(E1)`) to the sense of the sentence.

