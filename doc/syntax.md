For the semantic language holds:
 
 * variables start with a capital (Predicate, Entity1)
 * atoms (0-arity predicates) are lowercase with underscores (snake-case)
 * predicates are lowercase with underscores (snake-case)
 * non-predicate atoms and numbers: "De Nachtwacht", 1.5
 * is het misschien nodig om predicates en constants te namespacen? Eigenlijk is de predicate al een namespace
 
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
