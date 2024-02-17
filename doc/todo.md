# Todo

I'm creating an if2 that only works within functions. When all procedures that use if have been rewritten, the if can go and if2 be renamed to if.

## go as scripting language

start: ProcessRunner::RunProcessLevel(process, level)
    process has relations
input: a quantized semantic relation tree
vind een manier om een standaard relatie door Go te laten uitvoeren
it's possible to call a method by its name in Go

Test-driven!

First up:

- create an if/then that takes an expression as a condition
- go:list_get(VerLines, XIndex, rv) => xIndex := VerLines[XIndex]
- turn all "functions" into functions

Later:
- remove all functionality for multiple binding functions

- abc(x, y) {} // x and y are input or output variables (either bound or unbound)
- abc(x, y) => z {} // x and y are (bound) input variables

- disallow loops in rules
- disallow rules as function results ( a := abs(X) ) // abs should be a function, not an implication


## Development client

- send errors and productions as separate messages

## Command line application

- repair the command-line tools and document them

## Conditional sentences

Like "If I pick up a blue block, would there be a blue block on the table?" Can temporal questions like "When I picked up the block, were they any blocks on the table?" be treated similarly?

- Add temporary facts? Add rules that fire in this condition?
- Remove or hide the existing facts (blue block is on the table) so they don't interfere with the conditional facts?

## Integrate a Speech recognizer and generator

Open source
Accessible from Go

See also:
https://deepspeech.readthedocs.io/en/r0.9/?badge=latest

## Skills / games

Start / stop games, skills
Multiple skills in a dialog
Multiple games at the same time: start playing a song, do something else at the same time

See also:
MyCroft https://en.wikipedia.org/wiki/Mycroft_(software)   https://github.com/MycroftAI/mycroft-skills

## Stop - resume

Stop/pause, continue/resume : new topic

## Initiate dialog

Initiate conversation: new topic

## Documentation

* write a good tutorial

- create proper function diagrams for all processes
  - input
  - output
  - dependencies
  - side effects (which data stores are involved?)

Describe documentation files per linguistic feature, with the following sections

- topic description (anaphora, multiple sentences, conjunctions, etc)
- examples that should be covered
- possible approaches
- nli-go's approach

This is how to write a book about it.

List the biggest constraints; essential for an architecture:

- groups of entities
- or
- ...? questions


## Scripts and Frames

The restaurant script. Certain phrases invoke a script. (How?)

    "We went into a restaurant"

This phrase causes a number of discourse entities and relationships to be created at once, in the dialog context.

If the next sentence is "The waiter showed us our seat.", "the waiter" refers to an discourse entity in from the script.

## The programming language "mentalese"

Language: the difference between relations that are executed, and declarative relations must be syntactically clear

The combination of multiple results per relation and returning a single relation is absolutely counter-intuitive and leads to hard to detect and fix errors.

Make it consistent, complete, robust, etc. Have it conform existing paradigms.

- maybe remove `result` from `responses` in the intent; it is not used now
- I must implement all entities with atoms. Currently they are variables, but it means that variables are used as values, and this is clumsy. Then there must be a mapping from these atoms to database ids.
- typed arguments
- operators > = [H|T]
- n-dimensional arrays as local variables
- extend a module with another module
* quant_foreach: add as second parameter the variable to which the ids must be bound

## mentalese 2 Go

NB: it will always be necessary to add more language features. Can't I just use Go for the producedural language? Yes, this means that the application needs to be recompiled each time the code changes, and that the application gets linked to one language.

What brings this up? The reason is that when you need to call and process an API call in response to an action, you will need many functions built-in into Go, but not available in the mentalese language.

Perhaps make a separate Go module for each library that a skill produces. Then guide the process of building a custom server for a user. Teach the user that's it's simple to build a new server.

## mentalese: variables in loop-functions

When this list_foreach is done, the binding set has bindings for all variations of F C and V, while these should have been temporary
Because there are too many bindings, much extra calculations are done.

    go:list_foreach(List, E1,
        form(E1, F)
        color(E1, C)
        volume(E1, V)

I solved this for now by changing only list_foreach, and returning it with the binding is started out with. This is too limiting, so it needs to be properly solved, and for all other body-functions as well.

What's the problem?

When the loop fills one of the out-parameters of the function (ColSpan), and returns, this variable is now lost.

find_span(Width, VerLines, ColIndex, ColSpan) :-
    go:list_get(VerLines, ColIndex, X1)

    go:list_foreach(VerLines, Index, Line,
        [W := [Line - X1]]
        [W >= Width]
        [ColSpan := [Index - ColIndex]]
        break
    )
;

I solved this for now using a mutable variable.

## Code

Tags: maak duidelijk onderscheid tussen this-sentence-only tags en dialog-wide tags

* binding set -> results / binding list
* relation set -> relation list
* better validation for built-in functions; especially multi-binding ones
* use functional programming: make all data immutable; use copy-on-write everywhere; stop making deep copies

- turn the type `id` into `entity`. It's id and type are required. Is it possible that the id is optional?
- when listed values are used in the next sentence, they should be converted to multiple bindings

Adding info to the database is problematic:

tag: dom:tell(`some_event`, `:friend`, `:shrdlu`, P1)
or
go:uuid(P1, event) go:uuid(P2, event) go:assert(dom:tell(P2, `:friend`, `:shrdlu`, P1))

## Parsing source files

- check the correctness of all relations after parsing

## Interactive application

- interactive: arrow up/down for history

## The built-in mentalese application

- use parse tree as slot

## Sense -> use

To clarify that the meaning of a word is its use, we can use the word "use" in stead of "sense".

## Check -> find

For a Prolog user it would be more natural to use "find" for a function that produces multiple bindings.

## Combine the sentence and morphology grammars

I think it should be possible to create a single grammar (constructicon) that is able to parse both sentences and morphemes. Space would then be an explicit token in the grammar, and the tokenization phase must be skipped.

## Anaphora

I have not given any attention yet to "bound variable anaphora" https://en.wikipedia.org/wiki/Bound_variable_pronoun

- implement reflective_reference()

## Collect solution types

For each linguistics feature, there is a problem: how to put it into the process? Experience learns that there are procedures for this. And I'd like to collect them for later use:

- insert a new step (preferably operate on parse tree, if needed on semantic structure)
- insert several steps, for different aspects of the same feature (syntactic, semantic)
- create a tag
- create a relation

## Producing data at parse time

produce information in the parsing process. If A is parsed, information B is implied and can thus be added to the dialog context knowledge base.

## Performance

- the internal factbase is inefficient; for every new and removed fact, all facts are matched
- would be nice to have `is_first()` and `is_last()`: a check if the current binding is the first / last of the active bindings; such a function takes both a single binding and all bindings as input

## Database

* database mappings: allow a rule to be used only for given sorts; for performance
* SparqlFactBase: todo predicates does not contain database relations (just ontology relations), so this needs to be solve some other way

## Intents

Intent detection: provide access to the main clause variable: condition(P1)

## Language games / skills

A system should be able to handle multiple games in a single dialog.

## Intents

Multiple conditions should be allowed in a single intent.

- How can i buy X?
- help me buy X?

## Blocks demo

- put a stack in the largest open space
- provide a description of a steeple (are there any steeples now? Yes, the one based on the green block)
- A blob is a green block (doesn't work at the moment)
- When the demo is done. Do it in German as well, as proof of multilinguality.

The animation also reveals another problem: when the system builds a stack, it first decides on a location, then builds it. When building the first block, it may need to place the objects on top of it in some location. And it chooses the exact location where the stack should be. Later, the rest of the stack is still placed there. A solution could be to exclude this intended location from free space.

## Rules

Test if this works or make it work. Create a stack of current relations to be solved, and check if the stack already contains the bound relation.

    married_to(A, B) :- married_to(B, A);

* change rewrite rules from categories with variables to relations (see also Generator)

## Relations

Find a way to ensure completeness of information about all relations used in a system. An interpretation should not even be attempted if not all conversions have a chance to succeed.

* convert number words into numbers

## Planning

- replanning: once a plan is being executed, it may need to be discarded and replanned, due to new circumstances
- stop / continue commands

## Quantifiers

- create special cases for quantifiers, so that they don't look so complicated
- Make "more than" "less than" work

## Quantifier Scoping

Should be reintroduced. Syntactic or semantic? Find a good test-case.

- A range itself can contain quantified nouns (the oldest child in every family). The algorithm is not up to it. (See CLE)

## Dialog context

* the entire function DialogContext::ReplaceVariable is bad; I should not change structs that are in the history list already

## Interesting stuff

The catena concept may be interesting (for ellipsis a.o.):

    https://en.wikipedia.org/wiki/Catena_(linguistics)
