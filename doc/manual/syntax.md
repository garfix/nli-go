For the semantic language holds:
 
 * variables start with a capital, followed by zero or more lowercase or underscore characters: A, Verb, Entity1, Noun_entity
 * predicates start with a lower case character, followed by zero or more upper, lower or underscore characters
 * atoms, like predicates
 * string constants: use single quotes: 'De Nachtwacht'
 * numbers: 25 1.5

## Senses

In the general relational representation I have these senses:

 * declaration(P)               Based on a verb, this forms the main node of a sentence.
 * command(I)                   Based on the sentence structure, the verb must be interpretered as a command (Go!)
 * question(Q)                  Based on a verb, this sentence forms a question.
 * subject(P, S)                The semantic subject of a sentence (may differ from the syntactic subject)
 * object(P, O)                 Idem for object.
 * indirect_object(P, I)         Idem for indirect object.
 * prepositional_object(P, PO)   Like object, but linked via a preposition. "to the teacher"
 * quantification(O, [], D, []) Creates a non-scoped quantifier. O is the quantified variable. The first [] will hold the _range_, the second [] will hold the _quantifier_.
 * possession(E1, E2)           Based on "'s", it always denotes a possession relationship. Read: E1 is in possession by E2.
 * modality(P, M)               Based on modal auxiliary words like "can", "will" and "must", it denotes the modality of a predication.
 * relation(R, P, E)            Exposes the relation between two entities based on a preposition: the cat is on the mat (relation: on)
 * specification(E1, E2)        The specification of an entity (which is a set) is its intersection with another entity (set)
 * modifier(V, P)               Modifies the meaning of a verb because it has a particle (look at is different from look into)
 * conjunction(C, E1, E2)       A new entity (C) formed out of two other entities.
 * degree(E, D)                 Based on a degree adverb, denotes the degree in which something is the case (i.e. very)
 * complement(P, C)
 * name(E, N, T)                Here E is an entity, N is a name string constant (i.e. 'Charles') and T is its type (fullName, firstName, lastName, insertion)

More of these means that it is easier to create specific transformations based on these relations.

So I try to be as specific as possible as much as the syntactic relation allows.

Note that the most important entity (the governor, is that what it's called?) is always the first argument.

## Lexicon
 
[
    { form: 'book',           pos: noun,              sense: isa(E, book) }
    { form: 'read',           pos: verb,              sense: isa(E, read) }
    { form: /^[A-Z]/,         pos: firstName,         sense: name(E, Form, first_name) }
]

Lexicon definitions may use either a string constant or an expression for the form and use these variables in the sense:

E            Will be replaced by the entity variable of current node (ex. E1)
Form         Will be replaced by the word-form in the sentence. Only to be used with regular expressions. Form is a string, except when part-of-speech is 'number', then it is a number.

## Transformation rules

[
    parent(A, B) male(A) => father(A, B);
    parent(A, B) female(A) => mother(A, B) child(B, A);
]

It is possible to add a condition that applies to all relation of the question

[
    IF male(A) THEN parent(A, B) => father(A, B);
]

this transforms 'parent' into 'father', if the relation 'male' is present. But it does not affect 'male'. It is not removed.

## Grammar

[
    { rule: s(P) -> np(E) vp(P),     sense: subject(P, E) }
]

## Inference

[
    father(A, B) :- parent(A, B) male(A);
]

## Solutions

[
    {
        condition: relationSet,
        no_results: {
            answer: relationSet
        },
        some_results: {
            preparation: relationSet,
            answer: relationSet
        }
    } {
        condition: relationSet,
        no_results: {
            answer: relationSet
        },
        some_results: {
            preparation: relationSet,
            answer: relationSet
        }
    }
]

## Binding

{
    A: 'John',
    B: C 
}
