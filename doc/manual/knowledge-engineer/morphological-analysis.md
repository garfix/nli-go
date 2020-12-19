# Morphological analysis

It is possible to attach semantics to the individual morphemes of a word.

When the sentence is parsed and a word is not found in the lexicon, and the sought for category is not a proper noun, the parser transfers control to the morphological analyser to try and find the meaning.

The morphological analyser takes as input

- the word itself
- the syntactic category of the word
- the variables connected to the category

and receives as output

- semantic relations, to be integrated in the rest of the sentence
- or it may fail to find a sense

 ## The process
 
The morphological analyser performs three steps:

- segmentation of the word into segments using orthographical extraction rules
- parsing of these segments
- combining the senses to a semantic representation

## Segmentation

Examples of word segmentations:

    bigger -> big er
    littlest -> little est
    cities -> city s
    sleeps -> sleep s
    unbelievable -> un believe able
    thieves -> thief s

Place segmentation rules in a grammar, in a file named, for example `morpho.segment`, and name the file in the `index.yml`.

    morphology:
        segmentation: morpho.segment

Example segmentation rules

    /* character classes */
    vowel: ['a', 'e', 'i', 'o', 'u', 'y']
    consonant: ['b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'z']

    /* rewrite rules */
    relation: '*' -> noun: '*'
    noun: '*s' -> noun: '*', suffix: 's'
    /* example: biggest -> big est */
    super: '*{consonant1}{consonant1}est' -> adj: '*{consonant1}', suffix: 'est'
    comp: '*{consonant1}{consonant1}er' -> adj: '*{consonant1}', suffix: 'er'
    /* example: little -> little est */ 
    super: '*est' -> adj: '*e', suffix: 'est'
 
    /* terminals */
    adj: 'big'
    suffix: 'er'
    suffix: 'est'
 
 A word goes through the rules one by until one matches, and, when processed, results in one or more segments. Each resulting segments is one of the terminals.
 
 ## Parsing
 
 Once the segments are created, they form a string of morphemes that can now be parsed into a word, just like words are parsed into a sentence. It is done by the same (Earley) parser. For this, you need to create a grammar file for the morphological rewrite rules:
 
    morphology:
         segmentation: morpho.segment
         parsing: morpho.segment
         
Here's an example

    { rule: super(E1) -> super_base(E1) 'est' }
    
    { rule: super_base(E1) -> 'little',                                       sense: dom:height(E1, H1) go:order(H1, asc) }
         
This morphology grammar parses "little est" back into a single entry again, which accumulating sense.

Note that the `super(E1)` matches the category in the normal grammar:

    { rule: superlative(E1) -> super(E1),                                  sense: $super go:first() }

Every time the system handles a word that is not a proper noun, it invokes the morphological analyzer to check if it can create a sense for it. 
