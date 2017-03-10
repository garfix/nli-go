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

    [child(E2, E1) numberOf(4, E2) every(E1) act(question, yesNo)]

the scoped representation would be

    every E1 ?
        numberOf(4, E2) ?
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
* numberOf : numeric

Isn't it just possible to automatically create a scope box for every variable that is subject to an aggregate function? Each variable needs its scope, but most variables are existentially scoped, and this is the default. If we take this path, there would be no need to create a syntactic representation for the nested relations. All of it could be handled by the answerer.

