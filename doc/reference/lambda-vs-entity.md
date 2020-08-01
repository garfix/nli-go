# Comparing schemes of semantic composition

In this text I want to compare several types of semantic composition. I do this using the sentence

    a man kisses a woman
    
I want to show that these different schemes produce different semantic representations.    
    
## Montague grammar    

[Montague grammar](https://en.wikipedia.org/wiki/Montague_grammar) uses lambda calculus to combine the senses of phrases. Each phrase has a semantic attachment in the form of a lambda function. Composition is done by applying functions recursively.

Here's an example grammar  

    S -> NP VP              (NP VP)
    VP -> TV NP             λx. (NP (λy. (TV y x)))
    TV -> "kisses"          λx. λy. kiss(x, y)
    NP -> DET CN            (DET CN)
    DET -> "a"              λP. λQ. ∃x((P x) ∧ (Q x))
    CN -> "man"             λx. man(x)                       
    CN -> "woman"           λx. woman(x)
  
Some notes about this representation:

* parents are linked to children using function application
* the sense of a determiner (i.e. "a") contains elements of both the noun phrase and the verb phrase it is linked to
* the sense is a function that produces quantified relations
* quantifiers are not generalized. "a" is represented by "∃" which does not mean "a", but rather "at least a"     

Montague grammar has been very popular in computational linguistics, but it is, I think, generally believed to be quite verbose, and difficult to use. 

## CLE grammar

The Core Language Engine, which I see regularly mentioned as the state of the art, explicitly criticises Montague Grammar for being complicated. CLE grammar itself uses feature unification for several purposes, both syntax and semantic. This makes their expressions quite verbose in themselves. I will try to create a small grammar that keeps only the semantic parts. Parts left out are shown as '...'. The rules are written in Prolog, but I will give them another form in order to be able to compare them better with the others. 

    S -> Np Vp                          ( Vp, s: [...]) -> [(Np, np: [...]) (Vp, vp, [subjVal=Np, ...])]
    Np -> Det NBar                      ( qterm(<...>, V, NBar), np: [...] ) -> [ (Det, det:[...]), (NBar, nbar: [...]) ]
    "kisses"                            v: [arglist=[(B, np: [...]), (C, np: [...]]], [ kiss1, qterm(<...>, E, [event, E], A, B, C)]
    "the"                               det: [quantform=ref, reftype=def, ...], the

Some notes about this representation:

* parents are linked to children using feature unification
* the syntax excels in dealing with large types of features: agreement of number, gap threading, reference type, etc  
* subject and object are represented as a feature, which is an explicit representation of these roles
* the syntax is more verbose than Montague Grammar, because every feature that is passed needs to be named explicitly
* the sense of a determiner does not have a reference to the noun phrase or the verb phrase it is linked to

This syntax looks very complicated to me. It will probably take an expert many months to learn it, and it will continue to be a cognitive burden even then.

from: The Core Language Engine, chapter 5: Semantic rules for English

## Entity grammar

The grammar proposed by NLI-GO treats the sentence like this:

    s(P1) -> np(E1) vp(E1)              quant_check($np, $vp)
    vp(E1) -> tv(E1, E2) np(E2)         quant_check($np, $tv)
    tv(E1, E2) -> "kisses"              kiss(E1, E2)
    np(E1) -> det(_) cn(E1)             quant($det, E1, $cn)
    det(_) -> "a"                       quantifier(Result, Range, equals(Result, 1))
    cn(E1) -> "man"                     man(E1)                       
    cn(E1) -> "woman"                   woman(E1)

* parents are linked to children by concatenation (which is implicit) and variables (like `$np`) (explicit)
* the quantification is created at the np level, just like in CLE, and unlike Montague

This syntax is simpler than that of CLE, but it lacks features like agreement.
