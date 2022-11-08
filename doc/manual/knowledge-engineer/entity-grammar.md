# Entity grammar

## Introduction

The author of this framework created a special grammar because he needed one that was both simple to use and yet
expressive enough for a complex natural language interface. The grammar is not a traditional grammar which aims to keep
grammatical nonsensical utterances apart from well formed ones. Its sole purpose is to provide a highway to a semantic
representation of a sentence.

The grammar framework allows you to create both semantic grammars (whose categories are domain specific) and phrase
structure grammars (whose categories are syntactic categories). A combination of these types is not only possible but
encouraged.

## Difference with Montague grammar

The reason I created a new grammar was that I was not happy with the complexity of Montague grammar. It requires a good understanding of lambda calculus to use it, and I think this is a hurdle to its adoption to the mainstream of programmers. Also, I found that the calculus is not necessary. What's essential about handling semantics is not lambda calculus, but the principle of compositionality. Entity grammar provides a different way of implementing this principle. One in which the lambda calculus is largely implicit. 

A major difference between entity grammar and Montague grammar is that the latter involves the moving of senses, whereas the former just moves variables.

I created entity grammar because I wanted the application of techniques like lambda calculus, gap threading, and feature unification to be a problem for the framework, not for the developer.

## General form

The rules in this grammar have this form:

    { rule: s(P1) -> np(E1) vp(P1, E1),     sense: declare(P1) }

This is the main grammar rule for a sentence like

    Mary likes Jim

to be clear:

    Mary (np) likes Jim (vp)

The "rule" rewrites syntactic categories. Syntactic categories (s, np, vp) are familiar to anyone who has come into
contact with phrase structure grammars before. They stand for "sentence", "noun phrase" and "verb phrase".

The variables P1 and E1 stand for entities. Entities are special to this grammar. They play such a central role that it
is named after them. An entity can stand for anything: persons, objects, concepts, and even predications (things you can
say about something else). The names of the entity variables start with a letter and this letter can be chosen by you,
but it is good practise to make it represent what it stands for: P = predication, E = any entity. The number is meant to
distinguish between several variables of the same type. The scope of the variable is the rule. When another rule uses
the same variable, this does not mean they refer to the same thing.

The "sense" provides the meaning of the rule. It is a "semantic attachment" to the syntactic rule. The meaning of the
complete sentence is composed of these senses.

Each category in the rule is a syntactic or semantic category, and it is bound to argument variables. `noun(E1)` for
example is binds an entity to the variable `E1`.

The terminal rewrite rules have words in their right-hand side:

    pick(P1, E1, E2) -> 'pick'

There is no separate lexicon.

## Semantic composition

Each rule by itself just produces a single piece of semantics, or "sense". The "relationizer" combines these pieces to form the complete sense of the sentence. There are rules that govern this composition:

1. Append the sets of all the children, from child 1 to child N; then append the sense of the node
2. Exception to rule 1: force the location of the child sense with `$catIndex`  
3. Generate shared variables for the arguments
4. Turn regular expressions into string constants

## An imperative sentence

An imperative sentence instructs the system to do something, like play music:

    Play some jazz
    Play some folk music

The grammar for this sentence could be

     { rule: s(P1) -> 'play' 'some' music(E1),      sense: command_play(P1, E1) }
     { rule: music(E1) -> 'jazz',                   sense: type(E1, jazz) }
     { rule: music(E1) -> 'folk' 'music',           sense: type(E1, folk_music) }

`s` is the syntactic category for sentence. `'play'` and `'some'` are words. `music(E1)` (with variable) is a semantic
category. `'jazz'`, `'folk'` and `'music'` match words in the sentence.

The senses of the rules are application specific. This means that it's up to you to determine what works best for your
application.

## A declarative sentence

The rule given above what not sufficient to parse the complete sentence. Here is a simple grammar that will to the
trick:

~~~
Mary likes Jim
~~~

    { rule: s(P1) -> np(E1) vp(P1, E1),                 sense: declare(P1) }
    { rule: vp(P1, E1) -> iv(P1, E1) }
    { rule: vp(P1, E1) -> tv(P1, E1, E2) np(E2) }
    
    { rule: tv(P1, E1, E2) -> like(P1, E1, E2),         sense: like(P1, E1, E2) }
    { rule: like(P1, E1, E2) -> 'like' }
    { rule: like(P1, E1, E2) -> 'likes' }
    { rule: like(P1, E1, E2) -> 'liked' }
    
    { rule: np(E1) -> proper_noun(E1) }

The rules with an arrow ( -> ) are rewrite rules. When a sentence is parsed, the parser starts with the category `s` and
attempts to build the sentence from this `s` by rewriting it with the right hand categories `np` and `vp`. Then, it
tries to rewrite the `np` with `noun` and the `vp` with either `iv` or `tv` followed by `np`. `iv` is short for
"intransitive verb" (a verb without an object), while `tv` stands for "transitive verb" (which does have an object).

The result of the parse is the following parse tree

    s
        np
            proper_noun (Mary)
        vp
            tv
                like (likes)
            np
                proper_noun (Jim)

and the following sense

    declare(P1) like(P1, E1, E2)

## Open-ended categories

If you need numbers, email addresses or other forms of open-ended categories, you can specify them like this

    number(N1) -> /^[0-9]+$/ 
    email_address(E1) -> /^[a-z]+@[a-z]+\.[a-z]+$/          

The expression for email is very much incomplete. More importantly, you would need to change the tokenizer if an email
address is to be recognized as a single word by the parser.

Regular expressions can only be used as the single right hand position of a rewrite rule. The system will replace the variable of rewrite rule with a constant string that holds the word that was matched.
