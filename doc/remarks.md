# 2017 -10-01

My first DBPedia sentence is "Who married <person>". Just when I thought I had it finished I found out that people can have multiple spouses. Well, I _knew_ that of course :)

My "and" generator did not work with single np's. I fixed it temporarily.

I have to think about aggregation functions. My system works with a mix of breadth-first and depth first. I have not thought this through well enough.

## 2017-09-30

Finished n:m mapping for domain-kb. Yes! That was a lot of work! Had to rewrite the optimizer as well.

I use to have this long term task:

- Permanent goal: improve the grammar; extend it with new phrases, make it more precise. I think there's such a thing as an NLI-English grammar that exists of grammar rules commonly used when talking to a computer. It's a small subset of full English grammar, with an emphasis on questions.

I will abandon this goal, because a large grammar is actually detrimental to the performance of an NLI. It introduces ambiguity where it is not needed in individual domains. Plus it makes things much more complex. The downside is that grammars need to be copied from one domain to the next. But maybe I can create a basic set, or some sets, of grammar rules that are often needed.

I also introduced "solution routes". A solution route is one way to solve a problem. It consists of relation groups, together with a kb id and a cost of execution. A solution route is always ordered by least expensive first.

A consequence of these solution routes is an aggregate function in a solution route only applies to the bindings of _that route_. This may once become a problem, but I see no way around it.

## 2017-09-19

I have to change the ds2db mapping from 1-n to n-m. Which is quite a task! It introduced relation groups, that group a set of relations into relation groups. Each group server directly as input to a knowledge base.

It is theoretically possible that a set of relations could produce multiple sets of relation groups. I choose to ignore this possibility and select only the first row of relation groups. This is because I could not see how multiple sets of relation groups (which are in fact multiple queries) could work together with multiple bindings (aggregate) functions (i.e. number_of). Since they would need to work on the combined set of bindings of all sets. Which seems impossible, especially since other relations can (and will) depend on the outcome of the multiple bindings function.

## 2017-09-13

I have been writing some blogs about this on http://patrick-van-bergen.blogspot.nl/

I wrestled with the fact that the database could return no results. In which case the 'preparation' should not be applied, because it would have to many unbound variables.
This is very heavy on the database. I settled on a solution where a made the distinction between 'no-results' and 'some_results'. At least this will make the implementer
think about the issue.

## 2017-05-06

Added a cli command "nli" with two sub commands "answer" and "suggest" that produce JSON. Created a simple demo web app.

## 2017-04-01

First quantifier scoping tests passing!

## 2017-03-15

I have noticed one time too many that the "preferred reading" is simply the one where quantifiers decrease in scope from left to right. So I will now look for counter examples. And I need obvious examples, one you don't need to think or doubt about. If I can't find them, this whole exercise of creating laws and preferences for scoping is just too academic, and I will apply scoping in order of occurrence. Also, because if the order of scoping just changes after you had to think about it, in the end, no generic algorithm suffices, and hand coded quantifier rules will be needed.

> John visited every house on a street. (Quantifier Scoping in the SRI Core Language Engine)

 * EXISTS street(s) ALL house(h) [ in(h, s) visit(john, h) ]
 * Yes here the scope is inverted from the left-to-right one. But I think it is far-fetched. In normal use one would know the street. John visited every house on Baker Street.

I have found such a sentence:

> Name the mother of each child.

In this case the sense of applying wider -> narrower scope from left to right would be something like

    EXISTS m mother(m)
        ALL c child(c)
            child_of(c, m) name(q, m)

Which would ask the system to name the mother who has all children as a child. While 'every' has clearly wider scope, and this is a typical outcome of the quantifier scoping algorithm. So this is a good example to work with.

## 2017-03-14

Reading bits of Quantifiers in Language and Logic on Amazon:

Some Quantifier Phrases

    a few
    a little
    about five
    at least one
    at most six
    exactly three
    more than two
    fewer than four
    no more than five
    no fewer than five
    between six and twelve
    half of the
    at least a third of the
    at most two-thirds of the
    more than half of the
    Less than three-fifths of the

I thought of:

    two or three

Which reminds me to create conversions for the count-words (twelve => 12)

## 2017-03-13

Michael A Covington explains, in Natural Language Processing for Prolog Programmers, plainly what Quantifier Raising means:

some(X, dog(X), all(Y, cat(Y), chased(X, Y)))
=>
all(Y, cat(Y), some(X, dog(X), chased(X, Y)))

Q1(V1, R1, Q2(V2, R2, S2))
=>
Q2(V2, R2, Q1(V1, R1, S2))

So, is this a possible implementation for QS:

Collect all determiners.

* Go through all of their permutations.
* For each permutation of determiners:
    * Nest the determiners
    * Check their validity and calculate the score of this permutation.
* Pick the permutation with the highest score.
* Fill in the other relations at the outermost position where they still are scoped.

## 2017-03-12

I reviewed everything I wrote, and I reread some of the literature. And this is what I came up with, in the form of an example:

> Does every parent have 2 children?

1) raw parse

    rule: np(E1) -> determiner(D1) nbar(E1),                                   sense: determiner(D1, <1>, E1, <2>);

Where <2> means: move the relation set sense of the second consequent to this argument. The result of the parse is:

    isa(Q1, have)
    subject(Q1, S1)
    object(Q1, O1)
    determiner(D1, [ isa(D1, every) ], S1, [ isa(S1, parent) ])
    determiner(D2, [ isa(D2, 2) ], O1, [ isa(O1, child) ])

Then I'll introduce a generic step that converts to clumsy verb predicates to easier predicates. All occurrences of isa(Q1, PRED) subject() object() are turned into PRED(). The result:

    have(S1, O1)
    determiner(D1, [ isa(D1, every) ], S1, [ isa(S1, parent) ])
    determiner(D2, [ isa(D2, 2) ], O1, [ isa(O1, child) ])

And I'll introduce a step that helps remove vagueness ("have" is vague)

    [
        condition: isa(O1 child),
        old: have(S1, O1),
        new: have_child(S1, O1);
    ]

The step to remove vagueness is domain specific. I need to take this into account.

This allows you to replace single relations.

    have_child(S1, O1)
    determiner(D1, [ isa(D1, every) ], S1, [ isa(S1, parent) ])
    determiner(D2, [ isa(D2, 2) ], O1, [ isa(O1, child) ])


Then follows the old generic -> domain specific transformation, which we don't need here. Quantifier scoping gives:

    quant(D1, [ isa(D1, every) ], S1, [ isa(S1, parent) ], [
        quant(D2, [ isa(D2, 2) ], O1, [ isa(O1, child) ], [
            have_child(S1, O1)
        ])
    ])

Each quant is executed as follows:

* Range: Find all values for variable: S1
* Scope: Execute scope with S1 bound to one of its values (recursive process containing other quants)
* Quantification: execute the determiner. quant returns all bound variable sets on success, and clears them on failure

Question: how do I know that quants are always nested 1 by 1? It's not. Take for example: Roses are red and violets are blue. The scopes are not nested.

What does the quantifier scoping algorithm look like?

* Turn all determiners into quants

    have_child(S1, O1)
    quant(D1, [ isa(D1, every) ], S1, [ isa(S1, parent) ], [])
    quant(D2, [ isa(D2, 2) ], O1, [ isa(O1, child) ], [])

* Start binding all unbound relations (have_child). have_child has 2 variables: S1 and O1. Now look for quants that match and apply preferences.

* From CLE I can't really make out what happens next, I will have to try.

## 2017-03-11

The problem now is that the generic -> to domain specific conversion conflicts with the quantifier scoping. Which one goes first?

If quantifier scoping goes first I may loose transformation opportunities because the left-hand side relations are divided over different scopes.

If the domain specific conversion goes first, I may loose the "range". But since this is the best option and it is not clear to me, I will now try to show how it works out.

If I apply the conversion

    isa(Q1, have) subject(Q1, S1) object(Q1, O1) isa(O1, child) => have_child(S1, O1)

the range relation isa(O1, child) is lost. It that bad? Can't we do without a range? In theory, yes. The range, as I understand it, mainly serves to limit the possible values of the variable to an acceptable level. This is done for efficiency purposes. Let's see what this leads to

No range, just scope, with no conversion

    isa(Q1, have)
    quant(D1, S1, [ isa(D1, every) ], [ isa(S1, parent) subject(Q1, S1) ])
    quant(D2, O1, [ isa(D2, 2) ], [ isa(O1, child) object(Q1, O1) ])

No range, just scope, with conversion

    have_child(S1, O1) quantification(D1, S1) isa(D1, every) quantification(D2, S1) isa(D2, 2)

    quant(D1, S1, [ isa(D1, every) ], [ have_child(S1, O1) ])
    quant(D2, O1, [ isa(D2, 2) ], [])

Since O1 is unbound in the first quant, this leads to

    quant(D2, O1, [ isa(D2, 2) ], [
        quant(D1, S1, [ isa(D1, every) ], [ have_child(S1, O1) ])
    ])

I thought about this more. A range is required because it provides possible values for a variable.

---

A plan for what I need to do:

* change 'determiner()' to 'dp()' (syntactic rewrite)
* change 'determiner(E1, D1)' to 'determiner(E1, child-sense[E1], D1, child-sense[D1])'
* create a quantifier scoper that turns a relation set into a scoped relation set
* extend the answerer to make it answer scoped relation questions

---

Thinking about the problematic 'determiner(E1, child-sense[E1], D1, child-sense[D1])' it suddenly occurred to me that I could use an extra variable to separate the scope from the range:

rule: np(S) -> determiner(Q) nbar(R),                                     sense: quantification(Q, R, S);

This way, all relations with R form the range, all relations with Q form the quantifier, and all relations with S form the scope!

I like this solution (if it works?) because it allow me to keep the simple relation graph in tact, without resorting to nesting.

What would this do to the example sentence?

    isa(Q1, have)
    subject(Q1, E1)
    object(Q1, E2)
    quantification(D1, R1, E1)
    quantification(D2, R2, E2)
    isa(R1, parent)
    isa(R2, child)
    isa(D1, every)
    isa(D2, 2)

Let's presume a generic -> ds conversion before scoping:

    isa(Q1, have) subject(Q1, E1) object(Q1, E2) quantification(_, R1, E1) isa(R1, parent) quantification(_, R2, E2) isa(R2, child) -> have_child(R1, R2) ... etc

(omg! too long!)

    have_child(E1, E2)
    quantification(D1, R1, E1)
    quantification(D2, R2, E2)
    isa(R1, parent)
    isa(R2, child)
    isa(D1, every)
    isa(D2, 2)

When solving the question, quantifier scoping becomes active. 'every' takes preference and we find:

    scope(D1, R1, E1):
        execute: range -> scope -> quantifier
            range: foreach relation with R1: isa(R1, parent): R1 = 1, 2, 3, ...
            scope: all relations with E1: -
            quantification(D2, R2, E2) is placed in the scope of E1
                execute: range -> scope -> quantifier
                    range: R2 = 8, 9, 10, ...
                    scope: have_child(E1, E2)
                    quantifier: isa(D2, 2): succeeds for every R1 and R2 (or E1 and E2 ???)
            quantifier: isa(D1, every): succeeds for every R1

I really don't like this idea that the same entity now has two variables. This is asking for trouble. On the other hand, ... it may be more correct to do it this way.

Anyway, it is also possible to split the range from the scope from the fact that the range consists only of isa's and specification's.

## 2017-03-10

was

    rule: np(E1) -> determiner(D1) nbar(E1),                                   sense: determiner(E1, D1);

will be

    rule: np(E1) -> determiner(D1) nbar(E1),                                   sense: quantification(D1, E1, <D1>, <E1>);

Which means:

E1, which is a <E1> (the range) is quantified by <D1> (the sense of D1).

Interpreting "Has every parent 2 children?" yields:

    isa(Q1, have)
    subject(Q1, S1)
    object(Q1, O1)
    quantification(D1, S1, [ isa(D1, every) ], [ isa(S1, parent) ])
    quantification(D2, O1, [ isa(D2, 2) ], [ isa(O1, child) ])

This is heavily influenced by CLE, page 151.

The scope of S1 is formed by all relations with S1; likewise for O1.

In the quantifier scoping phase the scopes are created from the quantifications. Lets call the scopings 'quants' after CLE:

    isa(Q1, have)
    quant(D1, S1, [ isa(D1, every) ], [ isa(S1, parent) ], [ subject(Q1, S1) ])
    quant(D2, O1, [ isa(D2, 2) ], [ isa(O1, child) ], [ object(Q1, O1) ])

From what I've read, it appears that one and only scope is nested in another scope. These will need to be nested, but there are solutions for that.

## 2017-03-09

Working out my example sentence

    Does every parent have 2 children?

    rule: sInterrogative(S1) -> auxVerb(A) np(E1) vp(S1) questionMark(),       sense: question(S1, yesNoQuestion) subject(S1, E1);

    sInterrogative(S1)
        auxVerb(A) Does
        np(E1) every parent : scope( [ D1-sense ], E1, [ nbar-sense ], [ ] + object(S1, E1) )
            det(D1) every
            nbar(E1)
                noun(E1) parent
        vp(S1) have 2 children
            v(S1) have
            np(E2) : scope( [ D1-sense ], E1, [ nbar-sense ], [ ] + object(S1, E1) )
                det(D2) 2
                nbar(E2) children
        questionMark()

---

I tried to read "Representation and Inference" _again_, this time about Cooper and Keller storage. And I didn't understand.

But maybe I did find a way to extract the quantifier scoping data from the sense representation of a sentence.

Note that the 'determiner(E1, D1)' relation is central here. It connects the quantifier with the scope variable. And the range and scope can be deduced:

    rule: np(E1) -> determiner(D1) nbar(E1),                                   sense: determiner(E1, D1);

* quantifier: find all relations with D1 as an argument, except for 'determiner(E1, D1)', and all relations linked to it.
    These are all relations syntactically below the determiner node.
* variable: E1
* range: all relations 'isa(E1, X)' and 'specification(E1, X)', and all relations linked to it (except again for 'determiner(E1, D1)' and beyond)
    These are all relations syntactically below the nbar node.
* scope: all relations having E1 as an argument, except for the range relations

## 2017-03-08

The idea of a domain specific representation I got from TEAM. TEAM does not fully rewrite the generic representation, however. It just rewrites some domain specific parts. It calls this 'coercion' and it happens at several points in the process. Important is that the quantifier scoping is not influenced by these coercions.

Isn't it true that the possible readings of the determiners already surface at syntactic level? In that case, there's no need to juggle scopes after the parse. They are just possible interpretations of the same parse. It follows that it is probably a good idea to quantifier scoping at parse time already.

For the scoped expression

    quantifier variable range scope

the order of evaluation is

    range > scope > quantifier

for example:

    parent(x) good(x) > child(x, y) old(y) > number_of(x, 2)

The 'range' determines the possible values of x. That's why some range is required. Otherwise x could stand for all entities in the domain of discourse. Very unpractical.

---

It suddenly (finally) dawned on me that quantifier scoping NEEDS to be done by the parser, just because scoping information is only available on the syntactic level. It is possible that idioms and domain specific expressions change the scopes, but this is quite exceptional, I think, and it maybe hard but not impossible to handle. Having the parser handle quantifier scoping automatically makes things much easier for the user.

I am now attempting to make this work.

    rule: np(E1) -> determiner(D1) nbar(E1),                                   sense: determiner(E1, D1);
    rule: clause(S1) -> np(E1) vp(S1),                                         sense: object(S1, E1);

The quantification is not formed in the np-rule, as I expected. It is formed in the clause-rule, where all ingredients are available.

    rule: clause(S1) -> np(E1) vp(S1),                                         sense: object(S1, E1)                                scope: <E1, E1-determiner, np-sense, vp-sense>;

    quantifier: the determiner of E1 (find a determiner(E1, D1) relation)
    variable: E1
    range: the sense of np(E1)
    scope: the sense of vp(S1)

This all yields a sense like this:

    scope( [ quantifier-sense ], E1, [ np-sense ], [ vp-sense ] )

The sense of this clause is added to the vp-sense. In the grammar:

    rule: clause(S1) -> np(E1) vp(S1),                                         sense: scope( [ quantifier-sense ], E1, [ np-sense ], [ vp-sense ] object(S1, E1) );

Note: quantifiers are not just: ALL, SOME, NONE, but also, BETWEEN THREE AND FIVE, therefore quantifier sense is a compound relation set.

## 2017-03-07

Scoped expressions always look like this:

    quantifier variable range scope

for example

    some x man(x) [ immortal(x) ]

I have a problem: in my domain specific representation there's not always a range available. For example: when the query is

    [child(E2, E1) number_of(4, E2) every(E1) act(question, yesNo)]

the scoped representation would be

    every E1 ?
        number_of(4, E2) ?
            child(E2, E1)

But I'm missing ranges at the question marks. Both could be filled by something like person(X), but these are currently unavailable.

Note: range may be a conjunction, as in the case of "old man": man(x) && old(x)

## 2017-03-06

TEAM: With N (special) variables there are N! possible scope variations.

CLE: there are required conditions and preferences.

Natural Language Understanding, p. 109. The NLI system should not "work in the dark". Blindly trying many permutations.
In the end you still want a way to create an exception for a special case. Is it not possible to make scope rules "programmable" by the system user? Providing defaults for common cases?

## 2017-03-04

I read about scoping. There's a large chapter in the CLE book. But it's very hard to understand. I decided that I will just start programming. Creating tests. Then find out what problems I encounter, and then use the literature to help me.

The problem with Chat-80 is that it is as if determiners hold the answer to all questions. All types of questions are formatted in a determiner-based way. This holds true for yes/no questions (exists?) how many questions (number_of), but not for 'who' and 'which' questions. How are these handled by Chat-80? The identifier of the answer-entities serve an answers. So these are not the result of determiners.

I must turn my relation set representation into a scope-representation, where a scope is:

* a variable
* a relation set
* a set of sub-scopes

## 2017-03-03

First release was made. I got a few congratulations :) I had a list of things to do for release 2, but I am going to postpone them, because something more important came up: scoping.

I have read about scoping, but I never really got it. Yesterday I thought about a sentence like

> Does every parent have 2 children?

The stucture of this sentence is

    have(P)
        subject(P, S) parent(S) determiner(S, every)
        object (P, O) child(O) determiner(O, 2)

It understood that O has a smaller scope than S, and that O gets a different set of bindings for each element in S (!) This was a breakthrough for me.

I am late at understanding this, and I had much help understanding it, mainly from reading  about Chat-80, but also SHRDLU and Montague Semantics. Having a lot of background material helps enormously.

I want to place scoping _after_ generic -> domain specific interpretation, because this process may change the number and type of determiners. Determiners determine the type of scope.

This Chat-80 item also really struck me:

> every, all        =>          \+exists(X, R & \+S)

Which means, for example: every man is mortal: !exists(X, man(X) & !mortal(X))

Note that R is the main type of X (in my vocabulary: an isa(X, T) relation), and S are the remaining relations. Note that S is an relation set, so the relation exists() as a relation set as its third argument.

Before implementing this, I must reread everything about scoping to further understand the subject.

Possible sorting relations (chat-80):

* exists : the, some
* all ( !exists(x, !y) ) : all
* none: ( !exists )
* number_of: numeric

Isn't it just possible to automatically create a scope box for every variable that is subject to an aggregate function? Each variable needs its scope, but most variables are existentially scoped, and this is the default. If we take this path, there would be no need to create a syntactic representation for the nested relations. All of it could be handled by the answerer.

