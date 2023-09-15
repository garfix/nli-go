# Syntax

Mentalese has logical and imperative language constructs.

## Identifiers

* variable names are CamelCased: A, Verb, Entity1, OrderedBlocks
* predicates are snake_cased
* atoms are snake_cased
* string constants: use single quotes: 'De Nachtwacht'; also numbers: '25' '1.5'
* id: sort and identifier between backticks: `person:38911` `:http://dbpedia.org/page/Michael_Jackson_(actor)`

 ## Terms

 A term can be any of these

* variable (`HorLines`)
* atom (`plural`)
* string constant (`"Some text"`; also numbers: `"1.25"`)
* anonymous variable (`_`)
* regular expression (`r"a(b*)"`)
* relation list (`[apple(A) red(A)]`)
* id (&#96;`block:blue`&#96;)
* rule (`apple(A) if [fruit(A) round(A)]`)
* list of terms (`[2, 4, 6, 8]`)
* function (see below)

## Fact

A fact has no antecedents

    father(john, jack)

## Inference rule

An inference rule is a procedure that generates 0 or more facts.

An inference results in multiple variable bindings. When a variable is bound, it cannot be bound again, instead it should match the previous binding. Unbound variables should be bound when done. An inference fails when one of its conditions fails.

* There may be multiple inference rules with the same name.
* A condition is matched against one of the inference rules in code, but also against the data in the databases.

An inference rule has a fact, and zero or more conditions.

    father(A, B) :- [ parent(A, B) male(A) ]

A condition can also be an assignment or boolean expression. An assignment doesn't change the number of bindings. When the boolean expression returns false, the binding is dropped.

    too_old(A) :- [ birth(A, Birth) A := age(Birth) A > 40 ]

## Function

A function is a procedure to calculate a value.

It's defined by its name, its parameters, the return variable, and a body.

* A function execution results in one binding.
* All variables in a function are mutable.
* A function can't call a procedure directly.
* There can only be 1 function with a given name.
* The body contains only statements: conditionals, assignments, and loops

Example:

    hypothenuse(Width, Height) => Hypo {
        WidthSquared := Width * Width
        HeightSquared := Height * Height
        Hypo := go:sqrt(WidthSquared + HeightSquared)
    }

The result can be assigned to a variable. A returned list can be destructured:

    [HorLines, VerLines] := createLines(Inputs)

## Loops

Using a relation list in a function:

    for [ father(A, B) father(B, C) ] {

    }

Loop through array

    for E in List {

    }

    for Index, E in List {

    }

## Conditional

If-then

    if Exp {

    }

## List processing

Lists are initialized on demand.

Get

    X := List[I]
    X := List[Y][X]

If an element does not exist, it is created

Set

    List[I] := 3

If an element does not exist, it is created

Extend

    List append 5
    List prepend 5

    List = ListA + ListB

Split (with start and end-index + 1)

    L2 := L1[3:5]
    L2 := L1[:5]
    L2 := L1[5:]

Find index of element

    I := List find E

## Relation list

Relations lists can be concatened

    [relation() relation()] + VarB + [ relation() relation()]

When binding the varables of this relation set, `{{ VarB }}` will be expanded to the value of variable `VarB`, which must be a relation set.

## Comments

 In any file (except the json files) comments may be placed on any position, like this:

    /* much ado about nothing! */

## Syntax

The syntax of the language, in Extended Backus-Naur Form:

    variable-name-list = variable-name [",", variable-name-list]
    literal = predicate, "(", [ variable-name-list ], ")"

    expression = literal
    expression = expression, "+", expression
    expression = expression, "-", expression
    expression = expression, "/", expression
    expression = expression, "*", expression
    expression = "(", expression, ")"
    expression = boolean-expression

    boolean-expression = not boolean-expression
    boolean-expression = "(", boolean-expression, ")"
    boolean-expression = boolean-expression, "and", boolean-expression
    boolean-expression = boolean-expression, "or", expression
    boolean-expression = boolean-expression, "xor", boolean-expression
    boolean-expression = expression, ">", expression
    boolean-expression = expression, ">=", expression
    boolean-expression = expression, "<", expression
    boolean-expression = expression, "<=", expression
    boolean-expression = expression, "==", expression
    boolean-expression = expression, "!=", expression

    assignment = variable-name, { ",", variable-name } ":=", expression
    assignment = literal                  // discard any results
    for = "for", goal-list function-body
    if = "if" expression function-body

    statement = for | if | assignment

    goal = literal | assignment | boolean-expression
    goal-list = "[", { goal }, "]"

    function = function-header, function-body, function-footer
    function-header = literal
    function-body = "{", { statement }, "}"
    function-footer = "(", [ variable-name-list ], ")"

    implication = literal, "if", goal-list

## Struct

For later. Some ideas: structs must be defined

struct Person {
    Age,
    Name
}

P1 := Person {
    age: 41,
    name: "John"
}

P1[age] = 42
P1["age"] = 42

Alternatively, structs are immutable, and any change produces a new struct.
