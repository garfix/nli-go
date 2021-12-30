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
 * json: any structure that is json serializable 
 
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

A grammar rule looks like this:

    { rule: r, sense: s, ellipsis: e, tag: t }

Only `rule` is required.

`rule` forms the syntactic rewrite rule, `sense` forms the meaning, `ellipsis` is a path to the missing phrase, `tag` adds syntactic information.

Here are some examples:

    { rule: s(P1) -> np(E1) vp(P, E1) }
    { rule: tv(P1, E1, E2) -> like(P1, E1, E2),         sense: like(P1, E1, E2) }
    { rule: like(P1, E1, E2) -> 'likes' }
    { rule: noun(E1) -> 'cat',                          sense: cat(E1) }
    { rule: number(E1) -> ~^[0-9]+$~ }

    { rule: ...,                                        ellipsis: [root] [prev] .. catname(V) }

    { rule: noun(E1) -> 'cats',                         tag: number(E1, plural) }
    { rule: noun(E1) -> 'eats',                         tag: number(E1, singular) person(E1, 2) }
    { rule: S(P) -> np(E1) vp(P),                       tag: function(E1, subject) }
    { rule: S(P) -> s(P1) 'and' s(P2),                  tag: root_clause(P1) root_clause(P2) }

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

## Category path

Ellipsis uses category paths to specify the route from one node to another. `/` is the path separator. The nodes may be

- `..` up one node
- `..<cat>` up until a node has a `<cat>` category
- `<cat>` all `<cat>` categories directly below
- `/<cat>` all `<cat>` categories directly and indirectly below
- `-` previous sibling node
- `-<cat>` previous sibling node with a `<cat>` category
- `+` next sibling node with a `<cat>` category
- `+<cat>` next sibling node with a `<cat>` category
- `+-` any sibling node with a `<cat>` category
- `+-<cat>` any sibling node with a `<cat>` category
- `[prev_sentence]` to the root node of the sentence previous in the dialog context

## Binding

{
    A: 'John',
    B: C 
}
