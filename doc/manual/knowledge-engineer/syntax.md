# Syntax

## Identifiers

 * variables are CamelCased: A, Verb, Entity1, OrderedBlocks
 * predicates are snake_cased
 * atoms are snake_cased
 * string constants: use single quotes: 'De Nachtwacht'; also numbers: '25' '1.5'
 * id: sort and identifier between backticks: `person:38911` `:http://dbpedia.org/page/Michael_Jackson_(actor)`
 
 ## Terms
 
 A term can be any of these

 * atom
 * variable
 * anonymous variable
 * string constant
 * regular expression
 * id
 * relation set
 * a rule
 * a reference to a rule
 * list of terms
 
Integer numbers are stored as a string, but recognized as integer by the function IsNumber(). 
 
 ## Relation set 
 
A relation is a predicate followed by zero or more arguments, separated by commas and wrapped in parentheses.

    predicate(arg1, arg2, ...)
    
A relation set is a list of relations. A relation set may contain relation tags like this:

    relation() relation() {{ VarB }} relation() relation()
    
When binding the varables of this relation set, `{{ VarB }}` will be expanded to the value of variable `VarB`, which must be a relation set.         

## Comments

 In any file (except the json files) comments may be placed on any position, like this:

    /* much ado about nothing! */

## Grammar

    { rule: s(P1) -> np(E1) vp(P, E1) }
    { rule: tv(P1, E1, E2) -> like(P1, E1, E2),         sense: like(P1, E1, E2) }
    { rule: like(P1, E1, E2) -> 'likes' }
    { rule: noun(E1) -> 'cat',                          sense: cat(E1) }
    { rule: number(E1) -> /^[0-9]+$/ }

## Inference rules

    father(A, B) :- parent(A, B) male(A);

## Solutions

    {
        condition: relationSet,
        transformations: 
            transformation
        ,
        result: variable | anonymous variable,
        responses: 
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
    } 
    {
        condition: relationSet,
        responses: 
            {
                condition: relationSet, 
                answer: relationSet
            }
            {
                answer: relationSet
            }
    }

## Binding

{
    A: 'John',
    B: C 
}
