For the semantic language holds:
 
 * variables start with a capital, followed by zero or more lowercase or underscore characters: A, Verb, Entity1, Noun_entity
 * predicates start with a lower case character, followed by zero or more upper, lower or underscore characters
 * atoms, like predicates
 * string constants: use single quotes: 'De Nachtwacht'
 * numbers: 25 1.5
 * id: anything between backticks: `38911` `http://dbpedia.org/page/Michael_Jackson_(actor)`

 ## Comments

 In any file (except the json files) comments may be placed on any position, like this:

    == much ado about nothing! ==

## Syntactic Relations

Syntactic relations are formed in the relationizer phase, when the semantic attachments from the lexicon and grammar are combined.

I try to use the Universal Dependencies (used by the Stanford Parser), but there are some exceptions, needed for semantic processing.

http://universaldependencies.org/u/dep/index.html
http://nlp.stanford.edu:8080/parser/index.jsp

In the general relational representation I have these senses:

 * root(S)                              The root of a sentence
 * subject(P, S)               nsubj    The syntactic subject of a sentence
 * object(P, O)                obj      The syntactic object
 * ind_object(P, I)            iobj     The syntactic indirect object
 * aux(S, A)                   aux      Auxiliary verb relation
 * copula(S, C)                cop      Copula relation
 * name(E, N, P)                        The name of an entity. E is an entity, N is a name string constant (i.e. 'Charles') and P is its position (1, 2, 3)
 * quantification(E, R, Q)     -        A quantification. E is the scoped relations. R is the range. Q is the quantifier variable.
 * determiner(E1, D1)          det      The relation between a noun phrase and its determiner
 * case(E, C)                  case     The relation between a prepositional object and a preposition
 * mod(E1, E2)                 nmod     Modifies the meaning of a noun phrase with an attribute
 * mod(E1, E2)                 adjmod   Modifies the meaning of a noun phrase with an adjective
 * mod(E1, E2)                 nummod   Modifies the meaning of a noun with a number
 * mod(E1, E2)                 advmod   Modifies the meaning of a verb phrase with an adverb
 * mod(E1, E2)                 obl      Modifies the meaning of a verb phrase with a proposition
 * conjunction(C, E1, E2)      conj     A new entity (C) formed out of two other entities.

 * degree(E, D)                         Based on a degree adverb, denotes the degree in which something is the case (i.e. very)
 * declaration(P)                       Based on a verb, this forms the main node of a sentence.
 * command(I)                           Based on the sentence structure, the verb must be interpretered as a command (Go!)
 * question(Q)                          Based on a verb, this sentence forms a question.
 * possession(E1, E2)          case     Based on "'s", it always denotes a possession relationship. Read: E1 is in possession by E2.
 * ERROR prepositional_object(P, PO)          Like object, but linked via a preposition. "to the teacher"
 * modality(P, M)                       Based on modal auxiliary words like "can", "will" and "must", it denotes the modality of a predication.
 * complement(P, C)

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
