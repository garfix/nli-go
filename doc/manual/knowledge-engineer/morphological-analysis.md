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

    vowel: ['a', 'e', 'i']
    consonant: ['b', 'c', 'd']

    relation: '*' -> noun: '*'
    noun: '*s' -> noun: '*', suffix: 's'
    super: '*{consonant1}{consonant1}est' -> adj: '*{consonant1}', suffix: 'est'
    comp: '*{consonant1}{consonant1}er' -> adj: '*{consonant1}', suffix: 'er'
 
 