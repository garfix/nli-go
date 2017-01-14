For the semantic language holds:
 
 * variables start with a capital (Predicate, Entity1)
 * atoms (0-arity predicates) are lowercase with underscores (snake-case)
 * predicates are lowercase with underscores (snake-case)
 * non-predicate atoms and numbers: "De Nachtwacht", 1.5
 * is het misschien nodig om predicates en constants te namespacen? Eigenlijk is de predicate al een namespace

## Senses

In the general relational representation I have these senses:

 * predication(P)           Based on a verb, this forms the main node of a sentence.
 * subject(P, S)            The semantic subject of a sentence (may differ from the syntactic subject)
 * object(P, O)             Idem for object.
 * indirectObject(P, I)     Idem for indirect object.
 * determiner(O, D)         Singles out which discourse entities are involved.
 * possession(E1, E2)       Based on "'s", it always denotes a possession relationship. Read: E1 is in possession by E2.
 * modality(P, M)           Based on modal auxiliary words like "can", "will" and "must", it denotes the modality of a predication.
 * relation(R, P, E)        Exposes the relation between two entities based on a preposition: the cat is on the mat (relation: on)
 * specification(E1, E2)    The specification of an entity (which is a set) is its intersection with another entity (set) 
 * conjunction(C, E1, E2)   A new entity (C) formed out of two other entities.
 * degree(E, D)             Based on a degree adverb, denotes the degree in which something is the case (i.e. very)
 * complement(P, C)

More of these means that it is easier to create specific transformations based on these relations.

So I try to be as specific as possible as much as the syntactic relation allows.

Note that the most important entity (the governor, is that what it's called?) is always the first argument.

## Lexicon
 
[
    {
        form: 'book'
        pos: noun
        sense: instance_of(this, book)
    }
    {
        form: 'read'
        pos: verb
        sense: predication(this, read)
    }
] 
 
## Transformation rules

[
    father(A, B) :- parent(A, B), male(A)
    mother(A, B), child(B, A) :- parent(A, B), female(A)
]

## Grammar

[
    {
        rule: s(P) :- np(E), vp(P)
        sense: subject(P, E)
    }
]

## Binding

{
    A: 'John',
    B: C 
}
