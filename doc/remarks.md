# 3030-05-02

Generalized quantifiers are described in Core Language Engine (p.17) As an example they use the sentence:

    "All representatives voted"
    
"all representatives" is called the `restriction set` and their cardinality (number of distinct entities) is `N`. "voted" is not named, but the _intersection_ of these sets is called the `intersection set`. Its cardinality is called `M`. Now the quantifiers can be described as follows:     

    every:              N^M^[eq, N, M]
    more than two:      N^M^[geq, M, 2]

In NLI-GO this can now be written as

    every:              quantifier(Result, Range, equals(Result, Range))
    more than two:      quantifier(Result, Range, greater_than(Result, 2))

My "range" corresponds to their "restriction set" and my "result" to their "intersection set". 

I documented all nested functions.

# 2020-05-01

I decided it was a good time to implement generalized quantifiers. This had been "possible" for a long time, but I never actually got to it. You can now use any relation set to describe a quantifier. That is, as long as it takes two arguments: the number of entity ids in the range set (Range_count), and the number of (distinct) entity ids in the result set (Result_count). 

A generalized quantifier allows you to express not only "a", "some", and "all", but also "more than two" and "two" or "three".

It was also necessary to change the references for this reason, and so I turned `reference()` into `back_reference()` and I added `definite_reference()`, a construct that expresses "the red block" (a block that was either mentioned before, or is present in the scene.

I renamed `sequence` to `and` because it resembles the `or` function.

# 2020-04-24

I rewrote the functions library and added documentation. The use of the arguments is now more strict and I added error messages for wrong use.

One of those messages was that `greater_than(A, B)` requires two integer numbers. But when I ran that against the DBPedia tests I got an error because the dbp:populationCensus field of http://dbpedia.org/page/Republika_Srpska is not a simple integer, but is in scientific notation: `1.146520224E11` (!) So I left this type check out for now.

I noted that `quant()` is an iterator(!). It can be passed as an object (or predication) to other functions, and this brings interesting new possibilities!

Now I want to make the following possible, for the problem "stack up N blocks":

    equals(Q, quant(Q1, [], R1, []))
    
Which means: make sure that Q identifies with the quant on the right, and _bind_ the variables to values within the quant. This is a form of destructuring. It is also available in Prolog, and I found that there are two equals operators: `=` and `==`.   

https://www.swi-prolog.org/pldoc/man?section=compare

https://stackoverflow.com/questions/8219555/what-is-the-difference-between-and-in-prolog

`=` means "unify", while `==` means check identity. 

It is a good idea to make this difference already. So I will use `equals` for `==` and `unify` for `=`. So then this becomes:

    equals(A, B)
    unify(Q, quant(Q1, [], R1, []))

Note that unify is different from assignment, because assignment allows one to assign a new value to a variable that already has one. I am not sure that I want to enable that. Not now at least.

# 2020-04-18

For "Do all parents have 2 children", why use

    find(
        [
            quant(Q5, [all(Q5)], E5, [parent(E5, _)]) 
            quant(_, [sem(3)], E2, [sem(4)])
        ], 
        [
            have_child(E1, E2)
        ]
        
and why not 

    for([quant(Q5, [all(Q5)], E5, [parent(E5, _)])], [
        for([quant(_, [sem(3)], E2, [sem(4)])], [
            have_child(E1, E2)])])
    
i.e. why not a single relation for each quant? because the scope of variables of the the innermost `for` does not break out to the top level.
          

# 2020-04-17

So there are multiple forms of quantification. And I think I have found a way to represent them. I created a page `shrdlu-theorems.md` that describes the important PLANNER functions. SHRDLU's `THFIND` is particularly interesting here.

First off, the `np` will always produce a `quant`, but this quant will only have a quantifier and a scope:

    { rule: np(E1) -> qp(Q1) nbar(E1),                                     sense: quant(Q1, sem(1), E1, sem(2)) }
    
Note that the `sem(parent)` has gone. Such a quant will then serve as an argument in a second order relation. 

Next, it is necessary to make a distinction between `too little`, `enough`, and `too many` matches. The quantifier determines this. The quantification relations are iterators. They try every permutation of the ranges of elements. As long as too little matches are found, keep looking. The first match that makes the result enough will do. When the system finds too many results, it fails. 

Now, this is the relation for commands; it says: go through all permutations of the ranges and for each of them do `do_pick_up()` until `enough` matches are found. `do` only fails if `not enough` permutations could be found. 

    'pick' 'up' np(E1) -> do(sem(3), do_pick_up(e1))
    
For example: "pick up 2 or 3 boxes" will iterate over all boxes, tries to pick up each of them, but stop as soon as 2 is reached. "pick up a red block" will iterate over all red boxes, try to pick up each of them and stops if one succeeds.  
    
This is the relation for questions; it says: go through all permutations of the ranges and for of them do `support()` until `too many` matches are found. It fails both if `too little` or `too many` matches were found. 

    np(E1) support() np(E2) -> find(sem(1) sem(3), support(E1, E2))
    
For example: "does every parent have 2 children" will go through all parents, and for each parent, go through all children; it stops when 3 children were found for a parent; then fails immediately. Also fails if some parents had less then 2 children.    
    
As for declarations, this is a whole different ball-game. I will deal with that later on.

# 2020-04-15

There is a user! Vladyslav Nesterovskyi has taken an interest in the library. He asks the right questions and I hope I can help him out. The library is now conceptually rather stable, but it lacks many features and it easiliy cracks with the wrong input. I bid him to be patient and offered to assist him if I can (and the requests are in line with the intent of the library). 

===

Still thinking about quantification (handling words like "all", "some", "a", "2") for commands. It is interesting what Winograd has to say about this (possibly Winograd is the _only_ one who has something to say about this). NLU (p141)

> The system can accept commands in the form of IMPERATIVE sentences. These are handled somewhat different from questions. If they contain only definite objects, they can be treated in the way mentioned above for questions with no focus. The command "Pick up the red ball" is translated into the relationshop (#PICKUP :B7).
>
> However if we say "Pick up _a_ red ball", the situation is different. We could first use THFIND to find a red ball then put this object in a simple goal statement as we did with "_the_ red ball". This, however, might be a bad idea. In choosing a red ball arbitrarily, we may choose one which is out of reach or which is supporting a tower. The robot might fail or be forced to do a lot of work which it could have avoided with a little thought.
>
> Instead, we send the theorem which works on the goal a description rather than an object name, and let the theorem choose the specific object to be used, according to the criteria which best suit it. Remember that each OSS (Object Semantic Structure) has a name like "NG45". Before a clause is related to its objects, these are the symbols used in the relationship.
>
> When we analyze "Pick up a red ball", it will actually procedure (#PICKUP N45), where NG45 names an OSS describing "a red ball". We use this directly as a goal statement, calling a special theorem named TC-FINDCHOOSE, which uses the descriptions of the object, along with a set of "desirable properties" associated with objects used for trying to achieve that goal. #PICKUP may specify that it would prefer picking upsomething which doesn't support anything, or something near the hand's current location. Each theorem can ask for whatever it wants. Of course it may be impossible to find an object which fits all of the requirements, and the theorem has to be satisfied with what it can get. TC-FINDCHOOSE tries to meet the full specifications first, but if it can't find an object (or enough objects in the case of plural), it gradually removes the restrictions in the order they were listed in the theorem. It must always keep the full requirements of the description input in English in order to carry out the specified command. The robot simply tries to choose those objects which fit the command but are also the easiest to use.

# 2020-04-11

I am now reading "Semantic Syntax" which is written by my former Language Philosophy teacher Pieter Seuren. Semantic
syntax is part of the field of generative semantics https://en.wikipedia.org/wiki/Generative_semantics

I am too new here to say anything useful about it, but he mentions one thing that is interesting for the present
purpose. He proposes to use predicates for quantifiers. The sentence

    Nobody here speaks two languages

can be represented like this

    not(
        some(
            person(x)
            two(
                language(y)
                present(
                    speak(X, y)
                )
            )
        )
    )

While this is an interesting take, the problem is that quantifiers can be complex (at least two, two or three), and then
this scheme breaks down.

Nevertheless generative semantics is the area of linguistics that is most connected to semantics; and it is interesting
to read more about it.

===

I removed all transformations. They were only used in the solutions, and there I replaced them with rules.

Transformations are n-to-m replacements; they are _too_ powerful, they can break stuff easily and are hard to debug.
Since I introduced nested structures they have been trouble.

Now I am writing documentation, it is time for them to go.

# 2020-04-06

It is the verb that must be able to modify the quantifications of its np's.

Presently this is not possible. I want to change

    { rule: np(E1) -> qp(Q1) nbar(E1),                                     sense: quant(Q1, sem(1), E1, sem(2), sem(parent)) }

to

    { rule: np(E1) -> qp(Q1) nbar(E1),                                     sense: quant(Q1, sem(1), E1, sem(2)) }

and define verbs as such

    { rule: tv_infinitive(P1, E1, E2) -> 'support',                        sense: support(P1, sem(E1), sem(E2)) }

and so the semantics will be

    support(
        P1,
        quant(Q1, quantification, R1, range),
        quant(Q2, quantification, R2, range)
    )

a nice side effect is that I will loose the unintuitive sem(parent)

This will mean that when `support` is processed, it will have to go through all its quants and create nested loops for
them:

    foreach Rx in range R1
        foreach Ry in range R2
            support(P1, Rx, Ry)
        do quantification for R2    
    do quantification for R1

As a consequence, `support` can add a modifier to _its_ quants:

    do_quantification(command)
    do_quantification(declare)
    do_quantification(interrogate)

# 2020-04-05

I wrote about calculating the cost of executing relations. This technique actually has a very limited advantage in
certain circumstances. But it also has a disadvantage. Relations are no longer executed in the order that the programmer
intended. And the programmer can actually use this order for performance reasons, and this is quite common an natural.

The use of the cost function lost most of its worth already when I started nesting queries (quants). That reduced the
scope of the ordering to the relations in a very restricted area: that of a range, for example. Turns out that the
number of relations in such a scope is very small, and that it is very important that the programmer can decide which
order the relations have.

So I will remove the cost function.

===

Relation groups and solution routes were introduced in order to combine relations, which was needed for complex mappings
of generic semantics to domain specific semantics. A concept that I let go. Also it was used in combination with the
cost function. I just removed that as well. Now these solution routes are just a complex construct that serves no
purpose.

The only thing I need to resolve is this: can it be necessary to map 2 semantic relations to 1 database relation? Like
this:

    firstname(N1) lastname(N2) -> database_fullname(N)

or is it always possible to create a single semantic relation for this situation? I think so, because the programmer can
both influence the semantic language and the database language. Now the example with the names given here is a bad one,
which was the original reason, but has become obsolete. I have no better example.

I removed the optimizer, relation groups and solution routes.

===

Being busy, I also changed the many-to-many mapping of semantics to database to 1-to-many, as planned for a long time.

===

Another thing: the difference in meaning of number-quantifiers in questions and commands:

If I ask you:

    How many candy bars did you have?

You will mentally select all the bars you had, count them and reply: 2.

But when I say:

    You can have 2 candy bars.

You select 2 bars and eat them. You can't select them all, eat them and then find out you ate the whole cookie jar.

This difference is important for the QuantSolver. In questions, the quantifier operates after the scope was executed. In
commands, the quantifier operates during the processing of the scope. It continues trying the next entity in the range
until n scope executions have succeeded.

I could introduce a quant modifier.

# 2020-03-31

    Stack up two pyramids.

When I would keep my routine this sentence would be interpreted as:

    Find all pyramids, stack each of them up (?) and check if you have stacked up 2

This will actually succeed, since there are only two pyramids.

But what the sentence actually means is more like:

    Create a new stack of 2 (distinct) pyramids.

or

    do_create_stack(2, pyramids)

and this can be subgoaled as something like:

    place 1 ObjectType on the :table
    while N != 0 {
        select a topmost ObjectType that is not the topmost object on our stack and put it on top of the stack
    }

# 2020-03-29

Removed lexicons altogether; both from the parse and the generation side.

===

Question 9:

    Can the table pick up blocks?
    No

Winograd notes: The semantic rules make this a nonsensical question which therefore must be false.

This is a strange reaction, because the question, which starts with "can", is clearly about capabilities. These
capabilities could be explicitly modelled.

The following questions also use capabilities:

- 10. can a pyramid be supported by a block? --- yes
- 11. can a pyramid support a pyramid --- i don't know
- 12. stack up two pyramids --- i can't

The few lines Winograd adds to these questions suggest that there is no explicit capability model, no meta-knowledge
model, present. Questions 10 and 11 use induction based on the situation in the scene to come to a conclusion. Question
12, however, still requires some sort of knowledge about what each object can or cannot hold. Winograd offers _no_
explanation for 12: "the robot tries and fails".

The information for 12 can be built into the S-Selection constraints of `stack_up`: `stack_up(block, block)`. This will
cause the attempt to fail.

So I am going to implement "can" as "does the scene have an example of". This is quite easy. I just did number 9.

Fun fact: the system did actually find an instance of "pick_up(`:table`, X)", namely the memory of picking up the big
red block. Since I did not store _who_ picked up that block, but only that it was picked up, it matched. I made the
memory of picking up more explicit: "shrdlu picked up the big red block".

So 9 and 10 were no problem. 11 is a problem because of the subtle difference between "semantic rules" disallowing an
action and not finding an actual example. Have to work on this.

Number 12 is tough because of the expression "two pyramids". This actually means "a pyramid A and a pyramid B which is
not A". But "two pyramids" does not always mean that, of course. Just in this sentence.

# 2020-03-28

I made it possible to use strings in the grammar:

    vp(P1) -> 'put' np(E1) 'into' np(E2)

I can now do without the lexicon. Reasons for abandoning the lexicon were:

- I could place 'put' into the lexicon, but not 'put in', which is a verb as well; 'put in' needed to be in the grammar
- The 'part of speech' that was important in a lexicon was often used for semantic grammar purposes, and rightly so, so
    the concept turned into that was more naturally performed by the grammar

===

I solved question 8. I added the fact in the grammar that

    passive form uses the past participle verb form (supported)

It is different from the past tense form for irregular verbs.

# 2020-03-25

Question 8:

    Is it supported?
    Yes, by the table

Looks pretty easy. We'll see.

# 2020-03-24

Replacing "the one" with "the block" solves the problem completely! This is because "the block" creates a normal range
that results in a bounded range variable _before_ the dive into the relative clause.

But I can't do that of course. What I can do is create the sense "object()" for "one". But that still leaves me the
problem that object will (eventually) find all objects in the (restricted) universe. That is not funny!

    object(E1) :- block(E1); 

I could also use "reference()" as a sense for "one" but than the block must still be in the anaphora queue, and that is
not always acceptable.

===

It is not funny in the real world, but you can grin about it in the blocks world. It works. When the system handles
"one", it reads `object()` and it matches this to all blocks and pyramids in the microworld. There are only 8 of them so
going through them all is acceptable practise here. The process stops at the second block, so we're lucky. As I said
this would be too slow in a real database, but remember this:

    Simply questions should be answered fast; it should be possible (but not necessarily fast) to ask complicated questions.

And this question is undeniably complicated.

Anyway, my system would be almost equally slow had the question been about "block" in stead of "one"; the difference is
just two pyramids. That's why I decided to accept this (simple!) solution.

# 2020-03-22

To solve yesterday's problem, lacking a better solution, I will introduce "primed variables". These are variables whose
variable was primed (pre-activated) with a null-value (the anonymous variable):

    E1 { E1: _ }

When this variable is bound to anything else, the binding will _change_ to the new value:

    E1 { E1: `:id9755` }

This way, I can introduce a variable early and it will be active during the whole process of problem solving.

In this case, it will make sure that range variable will hold its binding.

===

I tried this, but it doesn't work very well. Now a matching with a primed variable always produces results, so I would
have to make changes in a large number of places.

# 2020-03-21

Next problem: in the semantic structure there is a range

    E19, [quant(
        _, [all(_)], 
        P9, [i(P9)], 
        [
            tell(P8, [
                quant(
                    _, [all(_)], 
                    E20, [you(E20)], 
                    [
                        pick_up(P9, E19)
                    ]
                )
            ])
        ]
    )]

And this range does not deliver the right E19's. This is because E19 only occurs first at a deep level (`pick_up(P9,
E19)`), and it doesn't make it all the way up to the root of the range. I think I should solve this by adding E19 to the
binding from the start, giving it an "empty binding", or something. But I don't like it.

# 2020-03-19

My problem is solved in Prolog by using the meta predicate `call`. So I can solve it this way:

    tell(P1, P2) :- call(P2);

I just implemented this.

Now there is the thing with "to pick up". Up until now I interpreted this part as `pick_up()`, but now there's a
problem. `pick_up()` is a command that is executed immediately. I don't want that in this case. In this case "pick up"
must be interpreted declaratively. And the syntax helps here, in the form of the infinitive.

https://en.wikipedia.org/wiki/Infinitive#Uses_of_the_infinitive

So I must make another distinction between a declarative relation and a imperative relation.

# 2020-03-18

Is at least one of them narrower than the one which I told you to pick up ?

    [
        question(S11) 
        quant(
            Q13, [some(Q13)], 
            E18, [reference(E18)], 
            [
                select(E18) 
                quant(
                    Q14, [the(Q14)], 
                    E19, [quant(
                        _, [all(_)], 
                        P9, [i(P9)], 
                        [
                            tell(P8, [
                                quant(
                                    _, [all(_)], 
                                    E20, [you(E20)], 
                                    [
                                        pick_up(P9, E19)
                                    ]
                                )
                            ])
                        ]
                    )], 
                    [
                        narrower(E18, E19)
                    ]
                )
            ]
        )
    ]

Note the impressive "range" of this question. Where it is usually just "men" or something. Here it is "the one which I
told you to pick up"

Current problem: when executing

    tell(me,
        pick_up(P5, you, E3))

I have the rule:

    tell(P1, P2) :- ;

I now must find a way to have the second argument of tell() executed.

# 2020-03-10

I started using explicit reference predicates. I need to document this.

# 2020-03-07

It's now possible to provide multiple arguments to the grammar rewrite rules. I wrote a test "filler stack test" that
tests if this works. It passes! :D

# 2020-03-06

When the rewrite rules get multiple arguments, these arguments have specific functions:

The proper_noun must have only one argument. The name will be always be linked to the first argument.

    proper_noun(N1)

NP predicates must have one argument: the entity

    np(E1)

VP predicates have a predication argument (that represents the predication itself) and zero or more entity arguments. We
keep them apart by giving them different letters

    vp(P1, E1, E2)

VP dependency passing dependencies have both the VP arguments, and also the dependency arguments. Let's give these
long-term dependencies variable names like L1, L2. These dependencies can be thought of as a stack. L2 was added later,
and should be removed first.

    dep_vp(P1, E1, E2, L1, L2) 

Dependencies must be passed _only_ via `dep_vp` relations, not via `vbar`, `verb` or other relations. Also: dependencies
are not passed via `np` relations. At least not until I have a use case for it, in which case I will create `dep_np`.

# 2020-02-29

The last few days I have been going through every book I have on the subject of extraposition. Trying to find the
optimal syntax to effectively parse several tough questions. I was willing to use some sort of stack scheme if that
meant that I could parse them well. However, the stack technique has its own problems. The problem is mainly that you
don't know in what order the entities can be found on the stack. And also, the stack does not allow you to pop entities
multiples times. And sometimes that is what you need. An entity may be used several times in the same sentence.
Interestingly, my own technique happens to doing be rather well in handling such sentences. Therefore I tend to leave
the stack path again and continue with what I had, using the extension of allowing multiple arguments (dependencies) in
the rule's antecedent.

Let me show you some of the tough questions I found in my books and a quick scan of how I plan to handle them.

    > ... the one I told you to pick up (SHRDLU)
    
    np(E1) -> np(E1) rel(E1)
        np(E1) -> the one
        rel(E1) -> np(E2) vp_dep(P2, E1)
            np(E2) -> I
            vp_dep(P2, E1) -> iv(P2, E2, P3) S2(P3, E1)       // pass E1 as dependency              
                iv(P2, E2, P3) -> told                        // tell(P2, E2, P3)
                s_dep(P3, E1) -> np(E3) inf(P3, E3, E1)
                    np(E3) -> you
                    inf(P3, E3, E1) -> to pick up             // pick_up(P3, E3, E1)

The construct vp_dep(P2, E1), along with S_dep(P3, E1) is still tough, but it is necessary to pass E1 down to its user.
The stack structure would not make it any easier, so this is as simple as I know how to make it.

Notice that I am not using concrete verbs in the abstract phrases (vp does not say "tell", just the more abstract "iv").
This should cut down the number of rules.

vp_dep is just an arbitrary name for a strange vp. Maybe it has a proper name. I hope so.

iv = intransitive verb (subject only); tv = transitive verb (subject + object)

    > which babies were the toys easiest to take from (CLE)

    s(P1) -> wh_np(E1) aux_be() s_dep(P1, E1)
        wh_np(E1) -> which np(E1)
            np(E1) -> baby(E1)                                  // baby(E1)
        s_dep(P1, E1) -> np(E2) vp_dep(P2, E1, E2)
            np(E2) -> the toys                                  // toy(E2)
            vp_dep(P2, E1, E2) -> adv(P2) inf(P2, E2, E1)
                adv(P2) -> easiest                              // easiest(P2)
                inf(P2, E2, E1) -> to take from                 // take_from(P2, E2, E1)    

Again, a tough question, with a double gap. On the one hand its a bit awkward to pass dependencies like this. On the
other hand, there are simple to track down, and the syntactic burden is manageable.

`vp_dep(P2, E2, E1)` simply means: this is a verb phrase that passes two extra entities / dependencies, without
consuming them.

To structure the order of the dependencies, let's say that the oldest dependency comes first. So if E1 is passed first,
it comes left. So it's a little bit like a stack, but not really.

    > the man who we saw cried
    
    s(P1) -> np(E1) iv(P1, E1)
        np(E1) -> np(E1) rel(E1)
            rel(E1) -> who s(P2, E1)
                s(P2, E1) -> np(E2) tv(P2, E2, E1)
                    np(E2) -> we
                    tv(P2, E2, E1) -> saw                       // see(P2, E2, E1)
        iv(P1, E1) -> cried                                     // cry(P1, E1) 

This is not a difficult sentence. It just has a relative clause.

    > what apple was assumed to be eaten by me?
    
    s(P1) -> wh_np(E1) aux_be() vp_dep(P1, E1)
        wh_np(E1) -> what np(E1)
            np(E1) -> apple                                     // apple(E1)
        vp_dep(P1, E1) -> iv(P2, P3) partp(P3, E1)
            iv(P2, P3) -> assumed to be                         // assume(P2, P3)
            partp(P3, E1) -> partic(P3, E2, E1) by np(E2)
                partic(P3, E2, E1) -> eaten                     // eat(P3, E2, E1)
                np(E2) -> me     

As you can see in the `assume(P2, P3)`, a predication can be an entity to be passed around as a dependency.

    > which woman wanted john but chose mary?
    
    s(P1) -> wh_np(E1) vp(P1, E1)
        wh_np(E1) -> which np(E1)
            np(E1) -> woman                                     // woman(E1)
        vp(P1, E1) -> vp(P1, E1) but() vp(P2, E1)    
            vp(P1, E1) -> tv(P1, E1, E2) np(E2)
                tv(P1, E1, E2) -> want                          // want(P1, E1, E2)
            vp(P2, E1) -> tv(P2, E1, E3) np(E3)    
                tv(P2, E1, E3) -> chose                         // chose(P2, E1, E3)

The woman (a gap filler) serves as subject in two predications. This is problematic for a stack system, but simple for
my system because semantic entities are its main (only) "features".

    > john sold or gave the book to mary
    
    s(P1) -> np(E1) vp(P1, E1)
        np(E1) -> john
        vp(P1, E1) -> tv(P1, E1, E2) np(E2)
            tv(P1, E1, E2) -> tv(P1, E1, E2) or() tv(P2, E1, E2)    // or(P1, P2)
                tv(P1, E1, E2) -> sold                              // sell(P1, E1, E2)
                tv(P2, E1, E2) -> gave                              // give(P2, E1, E2)     

Again, problematic for a stack system, because the filler is wanted in two places.

    > terry read every book that bertrand wrote (Prolog and natural language analysis)
    
    s(P1) -> np(E1) vp(P1, E1)
        np(E1) -> terry
    vp(P1, E1) -> tv(P1, E1, E2) np(E2)    
        tv(P1, E1, E2) -> read                              // read(P1, E1, E2)
        np(E2) -> np(E2) rel(E2)
            np(E2) -> every book                            // book(E2)
            rel(E2) -> that s(P2, E2)
                s(P2, E2) -> np(E3) tv(P2, E3, E2)
                    np(E3) -> bertrand
                    tv(P2, E3, E2) -> wrote                 // write(P2, E3, E2)

I think I'll have a beer now :P

# 2020-02-25

The filler stack that I worked out only works for left extra-position. Also the syntax I used doesn't work.

    { rule: s(P1) -> which() np(E1) vp(P) [E1], 			sense: which(E1) }

The `[E1]` must say after which consequent the variable is pushed. Like this for example

    { rule: s(P1) -> which() np(E1)^ vp(P), 			sense: which(E1) }

Does CLE have a right extra-position example? I didn't find one, but Wikipedia has

https://en.wikipedia.org/wiki/Discontinuity_(linguistics)#Topicalization

    It surprised us that it rained.     // surprise(E1, rain(Ev1))

# 2020-02-24

The cost function for the heaviness of using a relation to query a database (see Warren) becomes less important.

Why? It becomes less important when all NP's are quantified. If that is the case then few relations remain anyway, and
they are already highly constrained.

I already suspected this and now I read in a CHAT-80 paper that they don't quantify for SOME (since it isn't strictly
necessary to do so). It makes no difference for the results, but it does make a difference for the efficiency.

---

I read in Winograd's "Language as a cognitive process" (p. 367) about his divisions of syntactic systems. One that hit
me was:

    "Systematic non-syntactic structure": In many systems, there is no need to produce a complete syntactic sentence. Each constituent can be analyzed for its 'content' as it is parsed, and that content used to build up an overall semantic analysis. In systems like MARGIE, the structure is based on a semantic theory and each syntactic constituent simply contributes its part to the evolving semantic structure. In data base systems, the structure being built up may be a request in some query language associated with the data base. In general, systematic non-syntactic systems are organized to produce on overall structure that is determined by the syntactic pieces but is not organized along syntactic lines.  

Not that this is not semantic grammar. Semantic grammar has domain concepts as constituents, and cannot be reused for
other domains. The type of systems Winograd names don't care about syntactic formalism. It is subordinate to the easy
retrieval of the semantic content. I can relate to that very much. He also says about MARGIE somewhere that it was only
able to parse shallow syntactic structures, though. But I hope this is not essential to these kinds of systems. Because
I need "deep".

# 2020-02-23

CLE uses a stack to allow multiple gap-fillers to pass up and down the tree. My system, in the form I worked out
yesterday, uses antecedents to store these fillers. So I might think of them as a stack. This would allow arbitrarily
complex sentences, but they might take up a lot of syntactic variants eventually.

My system is simpler for simple sentences, and most sentences are simple, and allows more complex ones. This, I think,
is how it should be. The CLE system is complex from the start and does not get more complicated in complex ones.

But I could make the stack idea more visible as in

    vp-passing-fillers(P1, [E3, E1]) -> ask(P2) vp-passing-fillers(P2, [E1])      // ask(P1, E3, P2)

I think I must make an extra syntactic construction like this to help the maintainer of a system see what's what.

---

I am beginning to understand the syntactically complicated system of CLE that combines feature unification with
stack-based gap-threading. The concept is not that hard. It's the syntax that makes it hard. And the syntax is hard
because it tries to create a stack in a feature structure, which pushes the unification technique beyond the limits for
which it was originally designed. Feature unification is simply not meant to do that.

Feature unification is a good technique for number and tense agreement. But it stops there. Adding semantics is
stretching it. Adding gap-threading is step too far.

The concept is simply that specific phrases use the filler-stack to push or pop constituents. It is not necessary to
build all kinds of rules for different types of stack. The stack is not present in the rules. Just the top-elements of
the stack are shown. So this rule below:

    vp(P1) [E3] -> ask(P2) vp(P2)      // ask(P1, E3, P2)

now means: this rule only applies when the filler-stack holds _at least_ 2 elements. It uses the topmost element, called
E3.

And this rule

    s(P1) -> what() np(E1) vp(P1) [E1]                                        // what-question(E1)

means that E1 is pushed to the stack. With this addition the system would just be a syntactic variant of the CLE system
with the same power. It should be able to handle arbitrarily complex sentences. Let's see what this does to the sample
sentence:

    "Which babies were the toys easiest to take _ from _ ?"
    
    (Which babies) were (the toys) easiest (to take _ from _) ?
    
    s(P1) -> which() np(E1) vp(P) [E1]                                      // what-question(E1) -> pushes E1
        np(E1) -> baby(E1)                                                  // baby(E1)
        vp(P1) -> were() NP(E2) advp(P1) vp(P1) [E2]                        // -> pushed E2
            np(E2) -> toy(E2)                                               // toy(E2)
            advp(P1) -> adverb(P1)
                adverb(P1) -> easiest(P1)                                   // easiest(P1)
            vp(P1) [E1, E2] -> to() take() from()                           // take(P1, E2, E1) -> first pops E2, then E1  

Nice :)

At the end of the parse, the filler-stack must be empty for a sentence to be complete.

---

I could use this notation to denote the order in which the stack grows:

    <E1, E2]

This may or may not make it clearer that E2 is the topmost element on the stack. But it looks pretty ugly.

---

I could also say that the arguments I already used before can be part of the stack

    s(P1) -> which() np(E1) vp(E1, P)                                      // what-question(E1) -> pops P1, pushes E1 and then P 

# 2020-02-22

I "discovered" this brilliant Wikipedia article:

    https://en.wikipedia.org/wiki/Discontinuity_(linguistics)

It tells you everything about extraposition, wh-fronting and the way these are handled by grammars. For the first time I
get a little bit the idea that I understand something about this topic.

---

I want to make a point about syntactic grammar. I use this sentence

    "What book did your teacher ask you to read"
    
    S
        what book (m)
        did
        your teacher
        ask
            you to read (t)
            
    s
        your teacher
        ask
        you
        to read
            what book                 

Here (t) is the trace of a missing object. The object is "fronted" to the wh-clause at the start.

My point is that I always see "you to read" represented as a "sentence" or some such general name. This doesn't work. It
is not a normal sentence, and it should be named differently, like "vp-passing-object". If you do that, it is
possible to pass the marker to that structure. A normal sentence would not expect a marker. But a vp-passing-object would.

Note: at present the rule antecedent only takes a single argument, but I plan to change that in a short while.

    s(P1) -> what() np(E1) vp-passing-object(P1, E1)                                        // what-question(E1)
        np(E1) -> noun(E1)
            noun(E1) -> book(E1)                                                            // book(E1)          
        vp-passing-object(P1, E1) -> did() np(E3) vp-passing-subject-object(P1, E3, E1)
            np(E3) -> poss-pronoun(E3) noun(E3)
                poss-pronoun(E3) -> your(E3)                                                // your(E3) 
                noun(E3) -> teacher(E3)                                                     // teacher(E3) 
            vp-passing-subject-object(P1, E3, E1) -> ask(P2) vp-passing-object(P2, E1)      // ask(P1, E3, P2) 
                vp-passing-object(P3, E1) -> np(E2) to() read(P2)                           // read(P2, E2, E1)
                    np(E2) -> pronoun(E2)
                        pronoun(E2) -> you(E2)                                              // you(E2)  

This will expand the number of funny syntactic structures necessary, but this pays off in the way that the rules can be
written much simpler.

    Look ma, no features!

Principles

    A variable denotes a domain/semantic entity like a person or an event. So not a syntactic entity like a phrase or a state. 
    Relations connect entities.

CLE (p. 72) gives an example of a complex gap-threading challenge that requires a stack. Let's just see what my system
would do to this. Failing it does not completely disqualify it because this is a sought-for example.

    "Which babies were the toys easiest to take _ from _ ?"
    
    variant: "It was easiest to take the toys from which babies?"
    
    (Which babies) were (the toys) easiest (to take _ from _) ?
    
    s(P1) -> which() np(E1) vp-passing-prep(P1, E1)                                                 // what-question(E1)
        np(E1) -> baby(E1)                                                                          // baby(E1)
        vp-passing-prep(P1, E1) -> were() NP(E2) advp(P1) vp-passing-prep-object(P1, E1, E2)
            np(E2) -> toy(E2)                                                                       // toy(E2)
            advp(P1) -> adverb(P1)
                adverb(P1) -> easiest(P1)                                                           // easiest(P1)
            vp-passing-prep-object(P1, E1, E2) -> to() take() from()                                // take(P1, E2, E1)

So yes, it is possible, but indeed it introduces a "vp-passing-prep-object", and this makes you wonder how many
syntactic variants this will require in a complete grammar.

Can't I just replace "vp-passing-prep-object" by "vp-passing-2" to make this more general? The syntactic role of the
arguments never mattered before. I think so, but I am not sure.

    s(P1) -> which() np(E1) vp(P1, E1)                                                 // what-question(E1)
        np(E1) -> baby(E1)                                                             // baby(E1)
        vp(P1, E1) -> were() NP(E2) advp(P1) vp(P1, E1, E2)
            np(E2) -> toy(E2)                                                          // toy(E2)
            advp(P1) -> adverb(P1)
                adverb(P1) -> easiest(P1)                                              // easiest(P1)
            vp(P1, E1, E2) -> to() take() from()                                       // take(P1, E2, E1)

The words and phrases of the consequents constrain the use of the rules.

Perhaps it is possible that my system does not do well in restricting the possible sentences of a language. I.e. it's
possible that it would accept sentences that are not part of the language. But that's ok. The system is not a language
police. It presumes that the user just wants his job done and will create normal sentences. The job of the system is to
understand these. For the same reason I don't attach much value to things like agreement of number etc.

# 2020-02-21

In 7.5 "Memory" we read:

    To answer questions about past events, the BLOCKS programs remember selected parts of their subgoal tree. They do this by creating objects called events, and putting them on an EVENTLIST. The system does not remember the detailed series of specific steps like #MOVEHAND but keeps track of the larger goals like #PUTON and #STACKUP. The time of events is measured by a clock which starts at 0 and is incremented by 1 every time any motion occurs. ... MEMOREND puts information on the property list of the event name - the starting time, ending time, and reason for each event. The reason is the name of the event nearest up in the subgoal tree which is being remembered.
    
    A second kind of memory keeps track of the actual physical motions of objects, noting each time one is moved, and recording its name and the location it went to. This list can be used to establish where any object was at any time.    

Winograd uses this second order predicate TELL only once, and it is explained nowhere. I think I can assume it has no
special importance, and that I can ignore it here.

# 2020-02-20

    "Is at least one of them narrower than the one which I told you to pick up?"

This sentence contains a "hollow_non_finite_clause"

    "you to pick up"

see https://en.wikipedia.org/wiki/Non-finite_clause

This also means that this is the first sentence with a "gap": the entity "the one" is the object of the verb "pick up"
which is nested two levels deeper.

Is seems that my grammar is able to handle this sentence. That's very nice!

    quant(
        Q13, [some(Q13)], 
	    E18, [they(E18)], [
		    select_one(E18) 
		    quant(
		        Q14, [the(Q14)], 
			    E19, [tell(P10, S12) i(S12) pick_up(P11, S13, E19) you(S13)], [
			        narrower(E18, E19)])])]		

But I am not happy with the way "tell()" is handled here, even though it doesn't really play a role in this application.

# 2020-02-19

Whenever the system performs one of a small set of actions (MOVE, GRASP, UNGRASP) it could assert this relation with its
current timestamp. For example, at time 23 this relation could be added to a (which one?) database:

    GRASP(event:22, :b1)
    MOVE(event:23, 100, 50)
    UNGRASP(event:24, :b1)

    CAUSE_OF(event:22, event:18)

For now I can add all relations-with-events to the list so that includes relations like

    PICK_UP(event:415, :b6)

I should be able to figure out which block was "the one I told you to pick up". But what if there were multiple blocks
picked up?

Note that this is not Long Term Memory. This is Short Term Memory. Maybe it is not important to take the latest instance
of a picked up block.

    quant(the, one O,
        tell(:i, :you, pick_up(Ev, O))

"the one THAT i told you to pick up"

# 2020-02-18

Winograd uses an EVENTLIST structure to hold events. For each event it stores an event id, the main relation, and the
"reason": a reference to the event that caused it.

Understanding Natural Language has a special section on this sentence: 8.1.12 (Using clauses as objects), but is quite
obscure.

How do I infer from these questions that they make use of this event list? Well for one, they have the word "you" in
them and the sentence is in past tense.

tell(E1, S2) tense(E1, past) pick_up(S1, O1)

I might use inference

    tell(A, B) tense(A, past) -> tense(B, past)

This makes "pick up" into "picked up" (past tense)

    pick_up(E1, O1)

where E1 is not `:now` may be resolved from a data source containing action relations like this. This data source may be
filled automatically by the system each time some action occurs. Basic actions named by Winograd are GRASP, UNGRASP, and
MOVE.

# 2020-02-16

With all the preparatory work, NOT was quite straightforward!

Next up:

    Q: "Is at least one of them narrower than the one which I told you to pick up?"
    A: "Yes, the red cube"

This one's quite a beast! A quick analysis:

- a yes/no question that expects the found object(s) to be mentioned in the response
- a complex quantifier "at least one of", which can however be rephrased as "some"
- a reference to a group of entities mentioned before, which requires an adjustment to the anaphora queue
- a new concept "narrower than" that needs to be defined
- a reference to an object in an specific earlier event

That last one, especially, is very hard. It requires a queue of past events that you can refer to by describing them.
Then there is the "I told you" which is an implicit reference to the fact that the user (I) uttered a command.

Other sentences referring to earlier events:

    Had you touched any pyramid before you put the green one on the little cube?
    When did you pick it up?
    Why did you do that?
    How did you do it?
    How many objects did you touch while you were doing it?
    What did the red cube support before you started to clean it off?
    Have you picked up superblock since we began?
    Why did you drop it?

# 2020-02-15

I wanted to finish the story in the previous entry with "Moving the grammar rule np() -> she() up a few lines did the
trick", but in fact it didn't :/

I have rewritten the Earley parser. It now extracts all trees, and does this correctly, and efficiently.

Only now I completely understand the algorithm, and can I mentally reproduce a parsing process. It really is an awesome
algorithm, it gets better every year.

--

I made a new release: 1.7

Next up: question number 6:

    "How many blocks are not in the box?"

Negation!

Other sentences with "not":

* 14: I own block which are not red, but I don't own anything which supports a pyramid
* 35: is there anything which is bigger than every pyramid but is not as wide as the thing that supports it?

# 2020-02-11

This morning I had the following interaction with my system:

    Q: How many children has Madonna?
    A: She has 4 children
	Q: How old is she?
	A: I don't know 

This was my first attempt at anaphora resolution with DBpedia. I looked at what went wrong. I knew the system could
calculate Madonna's birthday so that was not it.

Turned out that the system did not use "she" as a pronoun, but first tried to use it as a name. "She" is the name of
many things: a magazine, several songs and movies, but these are all not persons, so no match. Until it found

    She: 7th century Chinese ruler of Qi  ( http://dbpedia.org/page/She_(Qi) )

The system went on happily, looking up the birth date and death date of this young ancient ruler. Unfortunately his(!)
birth date is unknown, and this lead to the result of "unknown".

# 2020-02-09

I managed to do anaphora resolution on both the lines 5 and 3 of the blocks world. I add entries to the anaphora queue
on completing the quant solving process and I attempt to select the entries from the anaphora queue first when in the
quant solving process.

It works, but it will need some work to make it work good. But it works. It actually works!

# 2020-02-08

I did. An id now looks like this:

    `person:123`

or when the entity type is not relevant:

    `:123`

I though about leaving out the : when there is no relevant entity type, but I decided not to because of the identifier:

    `http://dbpedia.org/class/yago/WikicatCountries`

this identifier happens to have a colon. I _could_ use another symbol as a separator, of course, but the problem would
remain. Forcing the colon has the added benefit that it makes you think about the entity type you might need. The
identifier now looks thus

    `:http://dbpedia.org/class/yago/WikicatCountries`

# 2020-02-07

I am making a release for the highlight that I am able to join the data of multiple databases.

---

Anaphora resolution:

It is a good idea to put the entity type in the `id`. That way I don't need to have relations present in order to
determine the entity type.

# 2020-02-03

The key cabinet is gone. It is finally replaced by a simple Binding, but a Binding that holds shared ids.

# 2020-02-01

Implementing shared ids. The key cabinet now maps variables to shared ids. The fact bases store the shared ids.

But where do I convert shared ids to db ids and the other way round? In the fact base or elsewhere?

For now I will do it in the general area, not in the fact base.

- Remove key cabinet altogether?
- Place named entities' ids in a binding?
- Add entity type to id term?

# 2020-01-29

I will interrupt my quest to answer question 5 to make a release that handles linked databases. I will make a new
integration test that handles the following question:

    Which poem wrote the grandfather of Charles Darwin?

Charles Darwin and the link to his grandfather are stored in db1. Erasmus Darwin and his poem is stored in db2. Both
databases are not directly linked, but there will be shared id tables.

# 2020-01-27

Shared-id == meta-id.

I will first un-use the key cabinet, and when that works, I will give it the new function of providing the shared-id -
db-id mapping.

I must make a test for the multiple-database query.

# 2020-01-25

What if a mapping between databases would exist? What would it look like? Preferably there would be a shared identity.

Of course many entities already have public id. Books have ISBN numbers. People have email addresses and
Social security numbers. Maybe they can be used in a certain domain, maybe not. But that's up to the designer of the
system. It is possible in some cases.

If I have 4 databases, A B C D, I could map A to B, A to C, A to D, B to C, B to D, C to D, but I could also map A to
shared and from shared to B.

    db 1: 8119              =>  johnbrown-1
    db 2: 23                =>  johnbrown-1

When I find an id in database A I can then do two things:

A) I could fetch the id's in A B etc. and place them in the key cabinet.
B) I could fetch the shared id and assign it to the query variable

In both cases I would need to find the database specific ID just before querying.

I prefer B. I would not need a key cabinet any more.

In most cases of course there is only one database. In this case the shared id is identical to the database id.

When and how should you create the mapping? Can it be done on-the-fly or must it be done periodically?

The mapping can be made as follows: suppose db1 has person fields `[id, email, name]` and db2 has fields `[id2, last
name, first name, email]` then the mapping should be created off-line by going through the persons in db1 and matching
them to the persons in db2, via heuristics or hard identities. The result would be mapping table for each database.

Even if not all db's have the entity, there must still be a shared id.

What would change?

- the key cabinet can go
- the id-term is a shared id; which can default to a database id, when there is only a single database
- in the configuration I would need to know which entity types have shared ids, and for each database a mapping table
- when a database is queried, the mapping from shared id to db id must be made; the response must be mapped again to the
    shared id

# 2020-01-24

Do I want a meta-id? (an entity that links one or more actual database id's)

Do I want to extend the notion of the id-term to meta-id? That if you create an id-term, that you are really creating a
meta-id with an initial single db id?

Can this replace the key cabinet?

meta-id: {db1: 15, db2: 18}

Matching meta's:

    {} + {db1: 16} => {db1: 16}
    {db1: 15} + {db1: 16} => no 
    {db1: 15} + {db2: 16} => {db1: 15, db2: 16}
    {db1: 15, db2: 4} + {db2: 4 db3: 108} => {db1: 15, db2: 4, db3: 108}

I can still use the key cabinet. The id term will then be a key from the cabinet. The cabinet holds the database ids;
possibly along with the entity type.

But how do these meta's behave?

When the sentence holds the name of an entity, the system can ask the user to clarify which of the entities, possibly in
different databases, is meant. A meta id will have one or more db ids.

But what happens further in the processing of the sentence? As long as a variable is not matched with a database id, it
will not be bound to a meta id, or perhaps only to an empty meta id. Once there's a database match, the db id will be
added to the meta id.

After that, the meta id will not gather any more ids. It is not possible that {db1: 13} will be used in `parent({db1:
13}, E2)` for db2. Because we do know which db2 entity matches 13 in db1, and we cannot leave it out entirely, because
the first argument _is known_, only not for db2, and we don't want any values from db2 as if the argument hadn't been
bound yet, for this would yield disallowed values. It would also not be possible to find new values for E2 in db2. So
this jeopardizes the idea that the system is capable of linking separate databases.

A way out would be a mapping function that links the entities of separate databases, but there is no such thing.

# 2020-01-23

In order to resolve anaphoric references I need to store the id's of earlier entities. The point is that the id's at the
moment are not unique even in a database, in the case of number ID's. An ID may belong to different tables. I can do two
things:

- store the table name in the id
- keep track of entity types during the processing.

Wait, I can look up the entity type for each relation, so this is not a problem.

---

I just noticed that my solver binds variables to ids of specific databases. Actually it is meant to bind to some
database-independent ids, that have links to ids in separate databases. This will be a problem when entities will be
found in multiple databases, as I intend to do. This is what I built the key cabinet for.

# 2020-01-22

I made the dialog context much simpler and more extendable by removing the fact base.

# 2020-01-18

I thought about anaphora resolution (handling pronouns like "it", and phrases like "the block", that refer to entities
that were just previously mentioned in the conversation).

This is how I intent to deal with it:

I will create a queue of entity-references that represent recently mentioned things in the conversation. This queue
consists of key-cabinet entries. Each entry holds the id's of the entity in one or more knowledge bases.

Whenever an entity is "consciously" processed by the system, it will be added to the queue. The queue will have a
certain maximum size and the last references will be removed when the queue is full.

In my solution I will not try to "detect" anaphoric references in noun phrases. I will treat anaphoric references just
like normal noun phrases. But I will make a change.

The critical structure is found in the processing of quantifications (quants). This process first fetches the id's of
all entities in the range. Then it processes the scope. And finally it filters the result through the quantifier.

The addition I make is in the loading of the range. When the range is specified by the word "he", for example, its sense
consists just of "male(E1)". This means that the entities considered would be all male persons. I will not load all
entities and filter them with the keys from the anaphoric queue. In stead, I will first attempt to restrict the range
with the available keys in the queue.

An example:

Anaphoric queue

    [{db=db2, id = 9591, type = block}]
    [{db=db1, id = 312, type = person} {db=db2, id = 111, type = person}]
    [{db=db1, id = 8, type = person} {db=db2, id = 9012, type = person}]
    [{db=db1, id = 31001, type = block}]

Input: When did the red man marry?

    when() 
    quant( 
        the()               // quantifier 
        red(E1) man(E1)     // range
        marry(E1, E2))      // scope

When the quant is processed, the processor will take the range

    "man(E)" 

and take the first entry of the anaphoric queue

    [{db=db2, id = 9591, type = block}]

since this doesn't match (type = block, not person), it tries the next

    [{db=db1, id = 312, type = person} {db=db2, id = 111, type = person}]

This gives a match.

Only if no match was found using the queue, will the range be matched with the knowledge bases without a key constraint.

This is the basic idea; I expect there will need to be made some adjustments to make this work.

# 2020-01-11

I put a new version online http://patrickvanbergen.com/dbpedia/app/ that allows case insensitive names. This will reduce
the number of failed queries quite a lot.

Did you know?

    Who is <name>?  --> Asks for a description 
    Who is <description>? --> Asks for a name 

# 2020-01-05

I made case-insensitive names possible, and at the same time checking the database in the parsing process. Introduced
s-selection to help finding the names. s-selection restricts predicate arguments and this in turn narrows the search
space for proper nouns in the database.

# 2020-01-02

Happy new year! 

I am introducing semantics in the parsing process, because I need some semantics to determine the type of the entity in
a name.

I want to use the relationizer that I already have for this, but it is too much linked to the nodes that I generate
after the parse is done.

Now I just had an interesting idea: what if I do the sense building as part of the chart states. That way, when the
parse is done, I just need to filter out the widest complete states and I will have the complete sense ready, without
having to create a tree and then parse that tree.
