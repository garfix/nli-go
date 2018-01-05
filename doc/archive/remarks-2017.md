# 2017-10-02

Another reason for not wanting a huge grammar, is that of scope. You don't want the grammar rules of one domain to affect another domain. It would only work if the perfect grammar was found.

There probably is no such thing. For all known languages.

# 2017-10-01

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

## 2017-02-27

I have solved all of the sentence puzzles of release 1 and I have fixed MySql access (which proved to be quite a hassle). Time to start writing documentation! (looking forward to that :)

## 2017-02-19

The !male(X) construction should be executed last, to prevent successive relations from introducing new values of X that match male(X) once more.

When I started to implement ! as a member of relation, I struck me as odd to do it this way. I _know_ the operator is not part of the relation, but it seemed easiest to do it this way. I will postpone this whole business to the second release, when I introduce grouped relations, not, and or.

I must not now introduce new stuff, and focus on the first release.

## 2017-02-18

I am thinking about introducing the not-operator. For example

    name(A, F, firstName) !name(A, I, insertion) name(A, L, lastName) join(N, ' ', F, L) => name(A, N);

How would that work? This makes me need to specifcy what a simple predicate-"operator" actually does. An operator like

    name(A, I, insertion)

Does this:

* It finds all predicates name() in the available sets, that match its variables A and I.
* If this is the first relation, A and I have no values yet, and a new binding is created for each relation that matches.
* For consequent relations, A and I do have values, and these limit the possible results. However, a new binding is created for each relation that matches the bound relation.

Now, what would this do?

    !name(A, I, insertion)

* If it is the first relation, and name() would match, !name() would return 0 bindings. The sequence ends.
* If it is the first relation, and name() does not match, !name() would return 0 bindings. But the sequence may continue.
* If it is a further relation, and name() matches, all these name() matches need to be removed from the set of bindings.
* If it is a further relation, and name() does not match, all earlier bindings continue to be used.

To make case 3 more clear:

    parent(X, Y) !male(X) => mother(X, Y)

parent() may yield [{X: 1, Y: 89} {X: 2, Y: 89}]

Now, !male(X) has male(X), and this may happen to succeed for X = 1; then this succeeding binding is removed.

## 2017-02-14

I want to keep database access as simple as possible. Only simple record retrieval. Laying complexity in database access leads to all sorts of complications. Think for instance of the inference rules that are applied. These should not become part of the sql query.

The problem with this approach is of course that we don't use the optimization techniques of the database to make the query faster. So the engine is not as fast as it could be. To this I objection I reply that the main use of an NLI engine is about questions that involve relatively little data. Questions like "Give me 5 bank transfers from Belgian customers in the last three years" are simple not the best use case for an NLI, and SQL will still be needed.

On the other hand, I will be adding optimization techniques to ensure that the nli queries are executed as efficient as possible (without rewriting the full query into an sql query).

## 2017-02-12

I answered the first release-1 question. Yay! But I took a shortcut. I still have a problem for both processing and generating proper nouns.
I will use this space to experiment.

    Sentence: 'Jaqueline de Boer'
    Generic: firstName = 'Jacqueline', middleName = 'de', lastName = 'Boer'
    Database: 'Jaqueline de Boer' 'Mark van Dongen'
    Answer: 'Mark' 'van' 'Dongen'

generic2ds

    fullName(A, N) :- name(A, F, firstName) + " " + name(A, M, middleName) + " " + name(A, L, lastName);

of

    fullName(A, N) :- name(A, F, firstName) name(A, M, middleName) name(A, L, lastName) serialize(F, M, S1) serialize(S1, L, N);
    fullName(A, N) :- name(A, F, firstName) name(A, M, middleName) name(A, L, lastName) concat(N, ' ', F, M, L);

Nieuw hier is het gebruik van systeem-predikaten in een transformation.

Dan moeten we misschien een db2ds introduceren:

    firstName(A, F) middleName(A, M) lastName(A, L) :- name(A, N) split(N, ' ', F, M, L)

## 2017-02-04

A sentence like this

    question(Q) isa(Q, have) subject(Q, S) name(S, 'Janice', fullName) object(Q, O) isa(O, child) specification(O, S) isa(S, many) specification(S, T) isa(T, how)

needs to be converted to a "program" and be executed. This is the essence of SHRDLU and this works. I want this to be done with as less human coding as possible. So how should we do it?

We have to convert the "how many" clause into a second order construct

    object(Q, O) specification(O, S) isa(S, many) specification(S, T) isa(T, how) -> number_of(O, N) focus(Q, N)

This forms

    question(Q) isa(Q, have) subject(Q, S) name(S, 'Janice', fullName) object(Q, O) isa(O, child) number_of(O, N)

Can we execute this? No we have to combine "have" with "child"

    isa(Q, have) subject(Q, S) object(Q, O) isa(O, child) -> child(S, O)

This gives us

    question(Q) child(S, O) name(S, 'Janice', fullName) number_of(O, N) focus(Q, N)

Can we execute this? Yes, after child() and name() are processed, there are 3 possible value for O left. Processing number_of() fills N with 3.

====

Can we do "largest"

Which is the largest block?

    question(Q) isa(Q, be) object(Q, O), determiner(O, D) isa(D, the) isa(O, block) specification(O, S) isa(S, largest)

How do we turn 'largest' into a program? (Note: this has to be domain-specific)

    isa(B1, block) specification(B1, Sp) isa(Sp, largest) -> block(B1) size(B1, S1) block(B2) size(B2, S2) greater(S2, S1, G) isFalse(G)

Does that work?

 * block(B1) : results for each block ID
 * size(B1, S1) : results for each block ID with its size
 * block(B2) : results for each block ID B1 cross joined with again each block ID B2, along with the sizes of B1
 * size(B2, S2) : the cross join of all blocks with all blocks, both containing size
 * greater(S2, S1, G) : goes through all results and keeps only the ones where S2 > S1, and sets G to (any entries left)

 No this doesn't work, but the version below does:

    isa(B1, block) specification(B1, Sp) isa(Sp, largest) -> block(B1) size(B1, S1) max(S1)

 * block(B1) : results for each block ID
 * size(B1, S1) : results for each block ID with its size
 * max(S1) : filter only the result with the highest S1

 Second order predicates like number_of() and max() act on result sets.

 ====

 When I started programming this I came across the problem that for some questions you have multiple answers. Can we handle these?

 Who were Mary's children?

    answer: name(C, N)


## 2017-02-02

There are several reasons why quantifier-constructs (exists, number_of) should not be added to the lexicon:

 * the word itself is not always enough to determine the quantifier ('how many': the combination of these words means 'number_of').
 * expressions can always give surface expressions another meaning than is apparent from the words. (every now and again, worth every penny)
 * some quantifiers cannot be deduced from the words alone and must be added later on ('not a lot', 'very little (people voted for Hillary)')

Trying some things:

How many children had Beatrice?

    solution: [
        condition: act(interrogation) focus(O) child(S, O)
        plan: number_of(child(S, O), N)
        answer: number_of_answers(N)
    ]

Was Mary a child of Charles?

    solution: [
        condition: act(interrogation) focus(O)  child(S, O)
        plan: ifExists(child(S, O), E) if(E, yes, no, A)
        answer: yesNoAnswer(A)
    ]

ds2generic

    yesNoAnswer(E) -> declaration(S1) specification(S1, Sp) isa(Sp, E)

The idea of a plan, though intriguing, is wrong. The question itself, rewritten in Domain Specific relations, is the plan. The reason is that the question contains many delicate details that are lost in a gross 'plan'. What is called a 'plan' is actually just a preparation for the answer.

Variables of the condition are populated by the matching variables of the question.

We find a new aspect of the domain specific representation: it is procedural. This makes the properties:

 * allow second order predicates
 * procedural: the representation is purposeful: it must contribute to the finding the answer

New question

    act(interrogation) focus(N) number_of(hasChild(P, C), N) name(P, 'Janice')

New solution

    solutions: [
        condition: act(interrogation) focus(N) number_of(hasChild(P, C), N),              // a question, about the number of children
        prep: gender(P, G),                                                              // look up the gender of the parent
        answer: gender(P, G) hasChild(P, C) number_of_answers(N);                           // "she, has children, number"
    ]

All relations of the solution are posed in the domain-specific language.

* Condition is matched against the input question. The first solution that matches is used.
* The variable set (S) used for the match is used for prep and answer.
* At that point the question itself is evaluated. Knowledge bases are used to look up answers.
* Then prep is evaluated and S is extended with its results.
* Finally the answer is formed by replacing the variables of answer with S. This answer is domain specific.

What needs to be done:

* second orderness in relations
* solutions
* processing solutions

## 2017-02-01

I am now looking at quantifiers and aggregations. Isn't it true that these are determined by determiners? May be, but I don't think you can link them at parse time. generic->domain specific would be fine. This means something like this:

generic:

    question(Q) isa(Q, have) subject(Q, S) name(S, 'Janice', fullName) object(Q, O) isa(O, child) specification(O, S) isa(S, many) specification(S, T) isa(T, how)

generic 2 domain specific:

    isa(O, child) specification(O, S) isa(S, many) specification(S, T) isa(T, how) -> act(interrogation) focus(N) number_of(O, isa(O, child), N)
    isa(Q, have) subject(Q, S) object(Q, O) isa(O, child) -> child(S, O)

By 'solution' I mean the matching of a question to an answer

    solution: [
        condition: act(interrogation) focus(N) number_of(O, isa(O, child), N) child(S, O)
        answer: declaration(S1) isa(S1, have) subject(S1, S) gender(S, female) object(S1, O) isa(O, child) determiner(O, Det) numeral(Det, N)
    ]

To answer a yes/no question I could use

    solution: [
        condition: act(interrogation) focus(N) exists(O, isa(O, child), E) child(S, O)
        answer: declaration(S1) specification(S1, Sp) isa(Sp, E)
    ]

If 'exists' yields a 'yes' or 'no' constant. Or if that's silly

    solution: [
        condition: act(interrogation) focus(O)  child(S, O)
        answer: declaration(S1) specification(S1, Sp) isa(Sp, E)
    ]

I could use the predicate 'focus(Entity)' to specify the activeness / passiveness of a sentences.

## 2017-01-31

surface:

    How many children had Janice?

generic:

    question(Q) isa(Q, have) subject(Q, S) name(S, 'Janice', fullName) object(Q, O) isa(O, child) specification(O, S) isa(S, many) specification(S, T) isa(T, how)

domain specific (no aggregation, no second order constructs):

    speechAct(question) questionType(howMany) child(A, B) fullName(A, "Janice")

conversion of domain specific to database:

    questionType(howMany) child(A, B) -> COUNT[ person(A, B) ]

database (variants, with and without aggregation):

    person(Id, "Janice", ParentId)      SELECT COUNT( ParentId ) FROM person
    person(Id, "Janice", ChildCount)    SELECT ChildCount FROM person

generic:

    declaration(S1) isa(S1, have) subject(S1, Subj) gender(Subj, female) object(S1, Obj) isa(Obj, child) determiner(Obj, Det) numeral(Det, 2)

surface:

    She had 2 children

Question: which types of aggregations do we need for NLI questions?

 * Is A married to B -> EXISTS
 * How many A -> COUNT
 * What is the total area -> SUM
 * Tallest child in the class -> MAX
 * Are some of the girls larger than all of the boys -> EXISTS

## 2017-01-29

And so it appears that even for the simplest of questions we need to resort to second order constructions. That I had wanted to postpone to release 2.

This is problem of aggregations, in database parlor. And the question is where in the chain first order forms are converted to second order ones. And back.

Let's have an example:

How many children had Janice?

The generic representation is:

    question(Q) isa(Q, have) subject(Q, S) name(S, 'Janice', fullName) object(Q, O) isa(O, child) specification(O, S) isa(S, many) specification(S, T) isa(T, how)

So the second order representation is not in the generic representation, and it should not be there either.

I don't really think it should be in the database representation either, because we want to keep the database layer as simple as possible as well. It's hard enough as it is. DB code should just be about retrieving simple records.

Let's imagine a domain specific representation for the question.

    speechAct(question) questionType(howMany) child(A, B)

next we need to find out what a proper response should look like

and how it would be turned into a generic representation.

    declaration(S1) isa(S1, have) subject(S1, Subj) gender(Subj, female) object(S1, Obj) isa(Obj, child) determiner(Obj, Det) numeral(Det, 2)

    (she had 2 children)

Note that the answer must be found by counting the number of child records. I mean: in _this_ case the answer is found by record counting. In another database, from the same domain, the answer could be stored directly (for example: person(id, name, numChildren)). This means that the aggregation must not be stored at the ds level. It should be stored at the database level.

## 2017-01-28

Insertions of Dutch persons: https://nl.wikipedia.org/wiki/Tussenvoegsel

1 woord: heel veel mogelijkheden (lijkt op lidwoord)
2 woorden: 1e woord: in, onder, op, over, uijt, uit, van, von, voor, vor (lijkt op voorzetsel)
3 woorden: de die le, de van der, uijt te de, uit te de, van de l, van de l', van van de, voor in 't, voor in t

I thought about solutions for multiple insertions, but currently I have none. The order of the insertions must be reconstructable from the semantic structure,
but I don't want to introduce several predicates for distinct insertion types. It gets too crowded that way.

Another question I must solve is how to represent questions. Questions are often of a meta level, second order predicate calculus. So we may think of

    act(question, who) who[A] married_to(A, B) :- question(Q) isa(Q, marry) subject(Q, A) object(Q, B)
    act(question, yesno) yesno[married_to(A, B)] :- question(Q) isa(Q, marry) subject(Q, A) object(Q, B)
    act(question, howmany) count[B] child(A, B) :- question(Q) isa(Q, marry) subject(Q, A) object(Q, B)

and how do I solve a second order problem?

## 2017-01-27

I added regular expressions as alternative for the word form. There are 2 sense variables now:

E            Will be replaced by the entity variable of current node (ex. E1)
Form         Will be replaced by the word-form in the sentence. Only to be used with regular expressions.

I replaced all occurrences of atom this with variable E.

The result of these changes:

		form: 'de',		    pos: insertion      sense: name(E, 'de', insertion);
		form: /^[A-Z]/,	    pos: lastName       sense: name(E, Form, lastName);
		form: /^[A-Z]/,	    pos: firstName      sense: name(E, Form, firstName);

## 2017-01-24

Pooh, I finally managed to port the Earley parser over to Go. Quite a bit of work, still.

I added terminal punctuation marks to the grammar (?.!). They need to be parsed, after all.

Now I am stuck with the following question: how to parse proper names?

In Echo I used the grammar to encode them:

    PN => propernoun1 insertion propernoun2,

I like this, because it does not require a database; the lexicon can determine that something is a proper name just because it starts with a capital.

Furthermore, it should be possible to parse sentences like "Is the name of your boss Charles?", even if "Charles" is not in the database.

My current solution:

    rule: properNoun(N1) -> fullName(N1);
    rule: properNoun(N1) -> firstName(N1) insertion(N1) lastName(N1);

    name(E2, 'Jacqueline', firstName) name(E2, 'de', insertion) name(E2, 'Boer', lastName)

fullName is used when there's only 1 name-word.
fullName, firstName, and lastName are recognized if they start with a capital letter
insertion must be part of the lexicon, i.e.

    form: 'de',		    pos: insertion      sense: name(this, 'de', insertion);

What about these possible syntaxes?

    form: '[A-Z]*',		    pos: lastName      sense: name(this, that, lastName);
    form: '<name>',		    pos: lastName      sense: name(this, name, lastName);

or allow full regexpses

    form: '/[A-Z]*/',	    pos: lastName      sense: name(this, name, lastName);

this would allow me to parse items like numbers and even e-mail addresses, given that the tokens created with the tokenizer would allow it.

LastName could be part of the lexicon.

## 2017-01-14

How to model "behind the door"; a PP?

Stanford says:

```
nsubj(looked-2, I-1)
root(ROOT-0, looked-2)
case(door-5, behind-3)
det(door-5, the-4)
nmod(looked-2, door-5)
```

Mainly: nmod(looked-2, door-5) case(door-5, behind-3); "door" modifies the verb "looked", "behind" modifies the noun "door".

The entity described as "behind the door" is a wedge-shaped place between the door and the wall. Currently I have no idea what the best way to model this is, but I think you shouldn't say that "behind" modifies "door", because I interpret "modifies" as "is a subset of", and behind the door is not a subset of door. I prefer to create a new entity that is formed from "door" and "behind"

```
PP(R1) -> Preposition(P1) NP(E1),         sense: relation(R1, P1, E1)

```

I name it "relation" (even though it is a wordt that already has too many meanings), because a PP is a

> Prepositions and postpositions, together called adpositions (or broadly, in English, simply prepositions), are a class of words that express spatial or temporal relations (in, under, towards, before) or mark various semantic roles (of, for).

(https://en.wikipedia.org/wiki/Preposition_and_postposition)

I am happy I have introduced the relations "declaration", "question" and "command", in stead of the more general sentence relation "predication". It is much more useful in transformations I think.

I had not heard of prepositional object, but today this way exactly what I needed.

https://en.wikipedia.org/wiki/Object_(grammar)

## 2017-01-13

I'm changing 'instance_of' into 'isa', just because it's shorter.

Lexicon: prime directive: senses take the name of their word form. Lexical inflexions are removed.

## 2017-01-12

I am working on a basic reusable grammar for English.

When modelling, you need to be careful as when to modify an existing entity variable and when to introduce a new entity.
In the clause "little rusty red book", the little rusty red book is an entity. It is formed in four steps.

instance_of(E1, book)

E1 is a book, or: E1 is a member of the set of books

instance_of(E1, red)

E1 is red, or E1 is a member of the set of red things
It is not necessary to say that E1 is a member of the red things that are also books (i.e. * instance_of(E2, E1), instance_of(E2, red)). Notably, when reasoning about red things, the red books must be found just as easily as the red vases.

But what to do about "rusty red"? I will use

```
modifier(E1, E2), instance_of(E2, red), modifier(E2, E3), instance_of(E3, rusty)
```

That is: both adjectives and adverbs are simply modifiers. Post-syntactic processes must determine the actual sense.

the stanford parser says

```
 (ROOT
   (S
     (NP (PRP I))
     (VP (VBP read)
       (NP (DT the) (JJ little) (JJ rusty) (JJ red) (NN book)))
     (. .)))

     nsubj(read-2, I-1)
     root(ROOT-0, read-2)
     det(book-7, the-3)
     amod(book-7, little-4)
     amod(book-7, rusty-5)
     amod(book-7, red-6)
     dobj(read-2, book-7)
```

I does not make a distinction between rusty and red. (I tried "bright" in stead of "rusty", same thing.)

Apparently, it is not clear that rusty is an adverb to red to the parser. This must be part of the semantic analysis.

What if I make the matter even more clear, and replace rusty with very?

"I read the little very red book."

```
(ROOT
  (S
    (NP (PRP I))
    (VP (VBD read)
      (NP (DT the) (JJ little)
        (ADJP (RB very) (JJ red))
        (NN book)))
    (. .)))

    nsubj(read-2, I-1)
    root(ROOT-0, read-2)
    det(book-7, the-3)
    amod(book-7, little-4)
    advmod(red-6, very-5)
    amod(book-7, red-6)
    dobj(read-2, book-7)
```

I cannot use these relations directly, though they are similar to what I need. Anyway, it says for "very red book":

```
    advmod(red-6, very-5)
    amod(book-7, red-6)
```

So I think I would make this into

```
instance_of(E1, E2, red)
modifier(E2, E3, rusty)
```

I appreciate the "root" relation of stanford's universal dependencies.

## 2017-01-07

I decided to work with releases. Each release has a goal functionality, and must be documented so as to be usable to others.

I cannot just use Erik T. Mueller's syntax rules (mueller-rewrites), because they have many constraints. I prefer to solve these constraints in the rules themselves (if that's possible). I keep them for inspiration.

I checked the grammar rules of The Structure of Modern English. It's quite amazing really. It is still the best book I know for rewrite rules. It says

>The version of the grammar presented here is not the most recent one, which has become highly theoretical and quite abstract, but takes those aspects of the various generative models which are most useful for empirical and pedagogical purposes.

This is very impressive. I think she refers to the Minimalist Program.

I reconsidered using a solely top-down or bottom-up parser. The top-down parsers can't handle left recursive grammars, and this is quite a heavy constraint. ThoughtTreasure uses a bottom-up parser, but I read in Speech and Language Processing that it can be quite inefficient. So I will recreate a Earley parser in Go. I love this :)

