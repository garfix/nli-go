# Syntax

## Identifiers

 * variables start with a capital, followed by zero or more lowercase or underscore characters: A, Verb, Entity1, Noun_entity
 * predicates start with a lower case character, followed by zero or more upper, lower or underscore characters
 * atoms, like predicates
 * string constants: use single quotes: 'De Nachtwacht'
 * numbers: 25 1.5
 * id: entity type and identifier between backticks: `person:38911` `:http://dbpedia.org/page/Michael_Jackson_(actor)`
 * relation: a predicate, followed by an argument list of terms: kick(john, jake)

 A term can be any of these

 * atom
 * number
 * variable
 * anonymous variable
 * string constant
 * regular expression
 * id
 * relation
 * relation set

## Comments

 In any file (except the json files) comments may be placed on any position, like this:

    /* much ado about nothing! */

## Grammar

    [
        { rule: s(P1) -> np(E1) vp(P, E1) }
        { rule: tv(P1, E1, E2) -> like(P1, E1, E2),         sense: like(P1, E1, E2) }
        { rule: like(P1, E1, E2) -> 'likes' }
        { rule: noun(E1) -> 'cat',                          sense: cat(E1) }
        { rule: number(E1) -> /^[0-9]+$/ }
    ]

## Inference

    [
        father(A, B) :- parent(A, B) male(A);
    ]

## Solutions

    [
        {
            condition: relationSet,
            transformations: [
                transformation
            ],
            result: variable | anonymous variable,
            responses: [
                {
                    condition: relationSet, 
                    transformation: transformation,
                    preparation: relationSet,
                    answer: relationSet
                }
                {
                    condition: relationSet, 
                    transformation: transformation,
                    preparation: relationSet,
                    answer: relationSet
                }
            ]
        } 
        {
            condition: relationSet,
            responses: [
                {
                    condition: relationSet, 
                    answer: relationSet
                }
                {
                    answer: relationSet
                }
            ]        
        }
    ]

## Binding

{
    A: 'John',
    B: C 
}
