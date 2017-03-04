## 2017-03-03

First release was made. I got a few congratulations :) I had a list of things to do for release 2, but I am going to postpone them, because something more important came up: scoping.

I have read about scoping, but I never really got it. Yestersday I thought about a sentence like

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

