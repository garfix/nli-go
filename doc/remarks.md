# 2019-09-06

    Find a block which is taller than the one you are holding and put it into the box.
    
    [
        seq(S5, 
            P5, [
                command(P5) 
                quant(
                    Q5, [number(Q5, 1)], 
                    E5, [
                        block(E5) 
                        quant(
                            Q6, [isa(Q6, the)], 
                            E6, [hold(P7, E7, E6) you(E7)], 
                            [taller(E5, E6)])
                    ], 
                    [find(P5, E5)]
                )], 
            P6, [
                command(P6) 
                quant(
                    Q7, [isa(Q7, the)], 
                    E9, [box(E9)], 
                    [put_into(P6, E8, E9)]) 
                it(E8)
            ])
    ]
  
# 2019-09-03

Again I need to solve the problem of quantification. By removing the extra QuantifierScoper step (it will once resurface in another form!) I have to think of a new way to do quantification right there in the grammar rules. 

adjp(E1) => taller() than() np(E2),            taller(sem(0), sem(2))

This does not provide a solution for the quantifiers, and it required a strange adaption to quantification code.

It occurred to me that It would be possible to place the quant() outside of the predicate it belongs to like this:

adjp(E1) => taller() than() np(E2),            quant(Q, sem(Q), R, sem(3), [ taller(E1, E2) ]) ERROR!
adjp(E1) => taller() than() np(E2),            quant(Q, sem(np/qp), R, sem(3), [ taller(E1, E2) ]) ERROR!

This doesn't work either. The quantifiers are problematic.

What if I included the qp() in the formula?

adjp(E1) => taller() than() np(E2),            quant(Q1, [isa(E1, a)], E2, sem(4), [ taller(E1, E2) ])
adjp(E1) => taller() than() qp(Q1) np(E2),     quant(Q1, sem(3), E2, sem(4), [ taller(E1, E2) ])

That is a proper solution! But its a bummer that I need two rules for every np(), in stead of one. This is intolerable.

What about:

adjp(E1) => taller() than() qp(Q1) nbar(E2),   quant(Q1, sem(3), E2, sem(4), [ taller(E1, E2) ])
qp(E1) => nil,                                 isa(E1, a)

The last rule represents an implicit existential quantifier when no determiner is present.

Awesome.

Except that the parser will have to be adjusted to allow for nil-constituents.

Alternatives?

adjp(E1) => taller() than() np(E2),     taller(E1, E2)
np(E1) => qp(Q1) nbar(E2),              quant(Q1, sem(1), E2, sem(2), sem(0)) ERROR!

No that is wrong, because there may be multiple np's in a predicate, and it messes with the bottom-to-top interpretation of semantics.

-- a night's sleep --

No, that conclusion is too quick. Checking Montague-like semantic interpretations (eg Natural Language Processing for Prolog Programmers, p.209) teaches just such a thing. In this routine, the last argument (sem(0), or: sem(parent), or: sense(parent)) "is awaiting the scope", that will be filled in by the parent.

So:

adjp(E1) => taller() than() np(E2),     taller(E1, E2)
np(E1) => qp(Q1) nbar(E2),              quant(Q1, sem(1), E2, sem(2), sem(parent))

How is this handled? When the relationizer processes the adjp(E1) it checks the child semantics.
When one of these child relations has 'sem(parent)' as argument, then there is

P parent relation
C child relation (with quant)

* the argument 'scope' in the quant of C is replaced by the current sem of P.
* the the sem of P is replaced by this quant 

My approach has some aspects of the Montague and the feature unification approaches, but it is much less verbose. Most standard semantic inheritance from child to parent is done implicitly, whereas these other approaches name each inheritance explicitly. I aim for simpleness. And even though my approach can still scare away the beginning programmer, I believe an experienced programmer will be able to understand it.   

The last rules have the advantage that simple 'exists' quantifiers can remain implicit, allowing for query optimization. Also, the complexity is isolated to a single place, that of the 'quant' definition.

# 2019-08-30

I just did it. I rewrote all grammars and removed the generic -> domain-specific transformations. With a few adjustments I removed an entire step in the analysis and made writing grammars much simpler.

I removed the difference between generic semantics and domain specific semantics. The generic semantics were very syntax oriented and the process of finding syntactic structures was hard enough, and then on top came the step of mapping these semantic structures to other semantic structures. Which required quite a lot of specific knowledge of the application, which made it quite hard to learn.

I also found that the generic->DS step was for the most part about fixing problems that the first step caused, or didn't fix. With the new "semantic grammar" I can create the "domain specific" semantics in one step.

I changed my opinion on semantic grammars. First I thought they were bad solutions that could never scale up. Now I see them as a completely different way of looking at sentences. It is more "human like", like the way you see the meaning of a sentence directly from the form, without seeing "syntactic structures". Syntactic structures are useful to determine if sentences belong to a language, but they are _too free_ and allow many forms that do not exist (like "pick down": verb particle). The new grammar can restrict this right away.

It is still possible, if wanted, to create multiple lexicons and grammars, general ones and domain specific ones.

All of this is quite a revelation! 

The number of rules will increase a lot, but in this age this is not a problem for the Earley algorithm.

# 2019-08-13

I _want_ to add raw words to the grammar rules.

Here are some thought about rewriting the grammar for the last blocks-world grammar:
    
    find a block which is taller than the one you are holding and put it in a box
    
    s(P1) -> imperative(P1)
    imperative(P1) -> imperative_clause(P1)
    imperative_clause(P1) -> imperative_clause(P1) conjunction(C1) imperative_clause(P1)
    
    imperative(P1) -> 'find' np(E1)                 do_find(P1) object(P1, E1)
    
    np(S1) -> qp(Q1) nbar(R1),                      quantification(Q1, R1, S1) 
    
    adjp(A1) -> adj(A) 'than' np(E1)    // taller than np
    
    np(E1) -> np(E1) relative_clause(P1)
    
    relative_clause(P1) -> 'which' copula(C1) adjp(A1)              // which is taller
    
    relative_clause(P1) -> 'that' np(E1) aux_be(C1) gerund(P1)
    relative_clause(P1) -> np(E1) aux_be(C1) gerund(P1)             // you are holding
    
    
    alternatives:
    
    // taller(A, B) can be part of the grammar directly (!)
    
    // E1 which is taller than E2
    np(E1) -> np(E1) 'which' copula 'taller' 'than' np(E1)       taller(E1, E2)
    
    // E1 which is / taller than E2
    np(E1) -> np(E1) 'which' copula comparative_phrase(E1)
    comparative_phrase(E1) -> 'taller' 'than' np(E2)             taller(E1, E2)

The take away is that syntax rules must facilitate the creation of the semantics.

Well-chosen rules make the transformation from global to domain specific unnecessary. 

# 2019-08-12

The nested quantifiers cause so much trouble that I left this project for a month.

This gave me some time to think. And read about other systems. Systems that use PHRAN, for example. PHRAN is a parser that parses directly to domain specific semantics. It works really really well. The parsing problems seem to go away like snow in the sun. But isn't this "just" a semantic parser? Isn't this cheating? Isn't there a catch that common sentences are other sentences are harder than necessary?

The problem I face is that I now have several nested structures (Quantifier Phrases and Conjunction Phrases), that I break up, modify, and reassemble, in order to make them domain specific and apply the right quantifier scopes. The assembling part now proves ever more difficult. Too difficult for a system that aims to be simple.

What if I parse directly to domain specific semantics? What if I allow the syntax rules to be domain specific? Would this be right or wrong?

I could have generic parse rules and domain specific parse rules. The latter could have higher priority given that we know that a certain domain is active. It would even allow for domain specific constructs like <employer> in the syntax rules, as PHRAN does.

I would not now mind so much about this any more. Practically because most important is that a system just works and is easy to configure. Theoretically because even humans may have domain specific parsing rules. Just as I don't see humans using transformational grammar, they probably don't convert a context free representation to a domain specific representation.  

# 2019-07-06

Back to the previous sentence. I used to treat one "the" as a number_of() in stead of a quantification(). I set this straight.
It turned out that the quantifier scoper was not up to the task and needed to be rewritten.

But now the seq() sequence is part of a scope and no longer top-level.
Also I found out that the sentence "find a block" has a quant in its range (!) The algorithm allows only quants in scopes.

So, back to the drawing board.

# 2019-07-01

The objects in the blocks world are identified by the combination of form, size and color. 
The box can be identified by "the box", since there is only one.
A block may be identified by "the large red block" since there are several blocks, and even several red ones.
There are two medium sized green blocks.
Pyramids can be identified by color alone.

I am now working on "which one?" as a response to "grasp the pyramid". Currently options are:

0: b2
1: b4
2: b5

In this case there color is different, so a good response could be

0: the green one
1: the blue one
2: the red one

It would provide sufficient information in the most efficient way. But it is hard to do.
Is is not even enough to find discerning attributes. A response like

0: the one with x position 100
1: the one with x position 150
2: the one with x position 200

would have the same information, but would still be less desirable to a human. Some attributes are more characteristic than others.

I could name them. 
    
    name(`b2`, 'the green pyramid')
    
Or for the moment I could just bail out and say: "I don't understand which one you mean.", almost like the SHRDLU response.    

# 2019-06-30

Back to "grasp the pyramid". I had solved this with

    determiner(E1, D1) isa(D1, the) => number_of(1, E1);
    
but I already knew it was no good. It results in grasping all pyramids and then counting if it was 1 picked up. Terrible really.

"the" is also not a simple quantifier. It is a determiner. It refers to either an object in the scene, or to an object in the dialog context.

And whereas "it" may refer to the latest subject in the dialog context, "the" may have a complete description attached. It seems terribly complicated! 

I checked Winograd and he uses "the" only for things in the scene of which there is only one (the table, the box), and for vert complex constructions (the one which I told you to pick up).

So for now I can treat "the" as a special kind of quantifier that applies only to the scene (not the dialog context) and means 

    the only instance of a group
    
If the group contains more than one, NLI-GO can respond with "I do not know which one you mean", as does SHRDLU. 

# 2019-06-27

Next question: "What does the box contain?"

For this question it is important that the "put into" has yielded an "contains" relation. Further this is the first answer that requires non-trivial generation:

"The blue pyramid and the blue block"

# 2019-06-23

"Find a block" succeeded! This is quite a milestone! 

I took a shortcut in that the active subject of the sentence was stored in the database by an explicit relation, this should be done by the framework.

Sparql query results are now cached to reduce dependency on DBPedia (I got "Not allowed" responses for too many queries) and to make tests faster.

# 2019-06-16

I added another nested structure: sequence.

I am now noticing that relations who have no connection to the rest of the structure (the relation set) fall outside the scope of the nested structures.

Thus, in the transformation phase from generic to domain specific, the connection must remain intact.

    So not:
    
    isa(P1, hold) subject(P1, S) object(P1, O) isa(S, you) => grasping(O);
    
    but
    
    isa(P1, hold) subject(P1, S) object(P1, O) isa(S, you) => grasping(P1, O);

The `P1` connects `grasping` with the rest of the structure and allows it to be nested along with the other relations.
`P1` is an event variable. I may also need it later when dealing with events in the past.

====

Let's talk about Stanford universal dependencies.

When I say "the block is taller than the pyramid"

these dependencies will have you do something like

    NN block -> relcl -> JJR taller -> nmod -> NN pyramid
    
and 

    NN pyramid -> case -> IN than
    
In relations

    block(A) relcl(A, C) taller(C) nmod(c, B) block(B) case(B, than)     
        
What I don't like is that there is a link to a pyramid, an entity, and this entity has a relation to a case.

I don't think that entities should have references to cases. To me this is very counter-intuitive.        
    
I want to write it like this

block(A) mod(A, M) isa(M, taller) mod(M, N) isa(N, than) ref(N, B) block(B) 

# 2019-06-14

The sentence has two imperative clauses and one sentence root. This changes the way I dealt with it before.

Make sure the sequences can be nested, just like other conjunctions.

# 2019-06-13

I just ignored "and". The left-to-right translation of words into meaning would take care of the order of the two clauses in the sentence.
It did not work. The optimizer changes the order of the relations and disturbs the order of the clauses.

"and" needs to be explicitly represented. And it _should be_. "And", at the sentence level, in imperative sentences, is not just syntactic glue. 
It implies order. It means, B should be executed after A. 

While the logical AND is commutative, the imperative AND is not.

Winograd has a whole section on the semantics of AND.

It needs to get its own "nested structure", the second after "quant".

# 2019-06-11

The optimizer had a complexity of n! and it started to show for 20 relations and several knowledge bases. I rewrote it and now it's back to normal.

# 2019-06-08

    "Find a block which is taller than the one you are holding"
    
Parse tree:  

    [s 
        [s_imperative 
            [s_imperative 
                [s_imperative 
                    [vp 
                        [vbar 
                            [vgp [verb Find]] 
                            [np 
                                [qp [quantifier a]] 
                                [nbar 
                                    [nbar [noun block]] 
                                    [sbar 
                                        [wh_determiner which] 
                                        [aux_copula is] 
                                        [adjp 
                                            [comparative_adjective taller] 
                                            [than_clause 
                                                [subordinating_conjunction than] 
                                                [gerund 
                                                    [np 
                                                        [dp [determiner the]] 
                                                        [nbar [noun one]]                                                    ] 
                                                    [np [pronoun you]] 
                                                    [aux_verb are] 
                                                    [gerund holding]]]]]]]]]] 
    
relations:
    
    root(P6) 
    isa(P6, find)                       // find
    object(P6, E7) 
        quantification(E7, Q6, R6) 
            number(Q6, 1)               // a 
            isa(R6, block)              // block
            mod(R6, S8) 
                subject(S8, S9) 
                    isa(S9, which)      // which
                isa(S8, taller)         // taller
                mod(S8, P8) 
                    case(P8, C6) 
                        isa(C6, than)   // than
                    object(P8, O5) 
                        determiner(O5, D6) 
                            isa(D6, the)// the 
                        number(O5, 1)   // one
                    subject(P8, S10)
                        isa(P8, hold)   // holding
                        isa(S10, you)   // you
                
I found the nbar -> nbar sbar rewrite to be much too free. It created ambiguity that wasn't necessary. I created some more restricted rules.

A modifier of a np is another np; its's not a preposition. The preposition is the case of the modifier.              

# 2019-06-03

Challenge of the day:

    "Find a block which is taller than the one you are holding and put it into the box."
    
Wow! Som many new things! Let's try to list them:

- A "which" modifier
- The concept "taller" and the comparison "taller than"
- The discourse object You
- A conjunction ("and") at the sentence level, joining two verbs.
- Anaphora (the pronoun "it"), that refers to a noun in the same sentence.   

"Find a block", once processed, should place the block found into the dialog context, where it can be picked up by the second part of the sentence.  

# 2019-06-02

Next question in the SHRDLU conversation: "Grasp the pyramid". 

The problem here is the simple word "the". It is combined with a noun ("pyramid") and forms a determiner phrase (DP).

The referent of a DP is an object in the dialog context (DC). 

I *have* a dialog context, but it just remembers the answers the user gave on previous questions. And it is not an actual knowledge base that can be accessed by the solver.

Then, I need the system to understand that "the pyramid" refers to information in that knowledge base. 

Finally, it would be nice if the system could answer "I don't understand which pyramid you mean" like SHRDLU does. The reason being that the pyramid has not been mentioned before in the dialog. But I might settle for a "Not OK" :) 

===

It occurred to me that "the pyramid" may mean something else. If there is only one pyramid in the scene, one may very well refer to it as "the pyramid". So if a search for pyramids in the scene results in 1 pyramid, then this is "the pyramid". However if I was talking about some specific pyramid and I say "the pyramid", then it is clear that I mean the pyramid from the dialog context.

I think actually that this second meaning is what Winograd means. I may be able to make this without too much work.   

# 2019-05-30

"Pick up a big red block" is almost done! I skipped the definition of "big" for now.

I added numbers as quantifiers. So I can now say "Pick up 1 block" and it will limit its actions to 1.

# 2019-05-29

ProblemSolver::SolveRelationSet() now solves a set in most-efficient order. Sometimes however, you need it to be executed in the order of the set.
When the sub-goals of a rule are executed, for example.

I have now split up these goals and feed them one-by-one to SolveRelationSet(), but I'd rather have a SolveRelationSetInOrder() function that I could call.

=====

Is there a difference in the number 5 between these two sentences?

* Put 5 blocks on the table.
* Are there 5 blocks on the table?

The first 5 is about 5 separate actions. For exactly five times, put a block on the table.
The second 5 is about counting. Find all blocks on the table. Count them.

# 2019-05-26

Earlier I described the command execution process as:

A command
* Is a relation set with *command()* and with one or more "command predicate" that ends with ! (like pick_up!() put_down!() ).
* It is recognized by the command() relation.
* It is executed as follows:
    * Find the command predicates (ending in !)
    * For each of the command predicates:
        * Bind the arguments using the input relation set without the command predicates
        * Execute (bind) the command predicate
        * Pass the bound variables to the next command predicate

But as I am implementing the command, without doing anything special, I get

    [[an(E5)]@shrdlu, [big(E5)]@shrdlu, [block(E5)]@shrdlu, [command()]@system-relations, [do_pick_up(E5)]@rules, [red(E5)]@shrdlu]

This is just about what I specified before. Before calling do_pick_up(E5), E5 must be determined. By placing `do_pick_up(E5)` at the end, this is just what happens, automatically

    [[an(E5)]@shrdlu, [big(E5)]@shrdlu, [block(E5)]@shrdlu, [command()]@system-relations, [red(E5)]@shrdlu, [do_pick_up(E5)]@rules]

So I don't have to do anything extra! Let's see how this works out.

I managed to place do_pick_up() at the end by creating stats for the other predicates (A stats-less predicate will always be executed last). This is not a proper solution, but it will do for now.

===

Now working on 'assert'.

    assert([grasping(X)])

The arguments (relations) are on a domain level, not db level. The available knowledge bases will be asked to accept this information.

In order to apply for a certain assert, a knowledge base should be able to:

 * Allow asserts / retracts
 * Handle the predicates in the assert

 Only fact bases can allow asserts for now.

# 2019-05-23

Sometimes I get some 405 (Method not allowed) responses from DBPedia. I think they tell me I am crossing the fair use limit.

If the number of queries to DBPedia proves to be too high, I might do some caching.
I could query all triples with a given verb at once, and store the results locally.
After that, whenever I need results with that verb, I can use the cache, for some time.

For example: select part of the Foaf names

    select ?a, ?b where { ?a <http://xmlns.com/foaf/0.1/name> ?b } offset 10000 limit 5000

===

I think I'll go for `do_pick_up()` rather than `pick_up!()`. The ! is a nice touch, but I don't like the change in allowed names it requires.

# 2019-05-22

I fixed the example question "which countries have population above 10000000".

I have to remark that I needed to change the entities.json entries for "name" to

{
  "person": {
    "name": "[person_name(Id, Name)]",
    "knownby": {
      "description": "[description(Id, Value)]"
    }
  },
  "country": {
    "name": "[country_name(Id, Name)]",
    "knownby": {
      "label": "[label(Id, Value)]",
      "founding_date": "[founding_date(Id, Value)]"
    }
  }
}

Person name was first "[name(Id, Name, fullname)]".
Then I changed it to "[name(Id, Name, fullname) person(Id)]", so the name "Iran" would be recognized as the name of a country, not a person.
But there was a problem with "person", it resolved to

    person(E) => type(E, `http://dbpedia.org/ontology/Person`);

and type had a very low size

    "type": {"size": 100, "distinctValues": [100000, 100] }

This meant that "person" would be placed first in the list of execution, and that meant that all persons would be loaded (!)
Hence the change to "[person_name(Id, Name)]" which maps to

    person_name(A, N) => birth_name(A, N) type(A, `http://dbpedia.org/ontology/Person`);
    person_name(A, N) => foaf_name(A, N) type(A, `http://dbpedia.org/ontology/Person`);

This is the desired order. It is also more specific, because a search for a country does need a search for birth_name:

    country_name(A, N) => foaf_name(A, N) type(A, `http://dbpedia.org/class/yago/WikicatCountries`);

# 2019-02-02

Winograd has multiple "theorems" for the same predicate:

TC-PICKUP, TCT-PICKUP, TCTE-PICKUP

and these are just different versions of commands.

Notice the difference in:

    Pick up the red block.

and

    the red block I told you to pick up.

The first "pick up" must perform an action, the second one must not perform an action. It just describes an action performed earlier.

Each "action predicate" can be fulfilled by other action predicates, and by description predicates.

Maybe I should make a distinction between these predicates and make action predicates look like this:

    do_pick_up()
    pick_up!()

May be I should change

    == pick up as a command ==
    root(P1) isa(P1, pick) modifier(P, Pt) isa(Pt, up) object(P, O) => pick_up(O);

to

    == pick up as a command ==
    root(P1) isa(P1, pick) modifier(P, Pt) isa(Pt, up) object(P, O) => do_pick_up(O);

Only root predicates can be commands.

And, since there can only be one command in the input, there is no need to specify it,
    so I can leave out

    action: goal(do_pick_up(E1)),

It is also important that I think about declaratives, i.e. "The red block is small", this could use an action like this:

        condition: declaration(),
        some_results: {
            answer: result(true)
        },

A declarative sentence must be asserted in whole.

===

To summarize:

A question
 * Is a relation set with *question()* and one of the question relations (what(), who() etc).
 * It is recognized by the question() relation
 * It is executed by binding its variables.

A declaration
* Is a relation set with *declaration()*.
* It is recognized by the declaration() relation.
* It is executed by *asserting* all of its relations.

A command
* Is a relation set with *command()* and with one or more "command predicate" that ends with ! (like pick_up!() put_down!() ).
* It is recognized by the command() relation.
* It is executed as follows:
    * Find the command predicates (ending in !)
    * For each of the command predicates:
        * Bind the arguments using the input relation set without the command predicates
        * Execute (bind) the command predicate
        * Pass the bound variables to the next command predicate

Two new system predicates are introduced: assert() and erase(). Both take a relation set as their sole argument.

===

which countries have population above 10000000
how many countries have population above 10000000

I managed to implement the second question, which is very impressive! But:

- some countries have populationCount, others have populationCensus, some both
- The limit of 100 results is not enough for all places with a population count
- DBpedia considers SEPA a country

So the result is not correct.

# 2019-01-31

What would the question solution look like?

        condition: question() yes_no() married_to(A, B),
        action: find(A, B)
        no_results: {
            answer: result(false)
        },

# 2019-01-30

To perform an action, it is necessary to mention the main command. Thus far I have:

    root(P1) isa(P1, pick) modifier(P, Pt) isa(Pt, up) object(P, O) => pick_up(O);

    == Pick up X ==
    {
        condition: command() pick_up(E1),
        action: goal(pick_up(E1)),
        no_results: {
            answer: dont_know()
        },
        some_results: {
            answer: canned('OK')
        }
    }

    pick_up(E1) :- at(E1, X, Y, Z) move_hand(X, Y, Z) grasp(X) raise_hand();
    move_hand(X, Y, Z) :- assert(at(`hand`, X, Y, Z));
    grasp(X) :- assert(grasping(X));
    raise_hand(X) :- at(`hand`, X Y Z1) add(Z1, 1000, Z2) move_hand(X Y Z2);

This is how this would execute:

- the condition of the solution is correct and of the same form as the questions.
- the relations in the input cannot be processed in any order, like in a question
- goal() accepts one relation of the input relation set
- goal() first evaluates the arguments of the goal relation: here it is just E1
- once the values for each of the arguments are found, pick_up() is evaluated, bound with the values just found
- pick_up() is evaluated just like any other relation set
- the results of goal(pick_up()) will be the bound variables, or empty set if it failed

I noticed that 'pick up' as a command may need to be modelled differently from 'pick_up' as a declarative predicate.

    (DEFPROP TC-PICKUP
         (THCONSE (X (WHY (EV)) EV)
              (#PICKUP $?X)
              (MEMORY)
              (THGOAL (#GRASP $?X) (THUSE TC-GRASP))
              (THGOAL (#RAISEHAND) (THNODB) (THUSE TC-RAISEHAND))
              (MEMOREND (#PICKUP $?EV $?X)))
         THEOREM)


# 2019-01-27

I am starting to think about the SHRDLU demo. Since Winograd's work is brilliant, I will merely try to mimic it, and not try to do it better.

I'll just handle Winograd's sample sentences one by one. First sentence:

    Pick up a large red ball.

My system does not know how to do anything, nor how to change something in a knowledge base. Both are necessary here.

The system now knows that pick_up() is a command. It does not yet know what to do with it. This may be a start:

    pick_up(X):
    - FIND[ X ]
    - GOAL[ grasp(X) ]

Both find() and hold() would be actions (or plans). Find() would be the process I nave use up until now.
Hold() would be a new action. In order to hold(X), the system would need to grasp(X) and then move(X, hold) where self is some temporary storage location.

The action grasp() would make changes to the database.

If I would rewrite Winograd's PLANNER code

    (DEFTHEOREM THEOREM3
        (THCONSE (X Y Z) (#PUT $?X $?Y))
        (THGOAL (#ON $?X $?Z))
        (THERASE (#ON $?X $?Z))
        (THASSERT (#ON $?X $?Y))
    )

into my own words, it would be

    PLAN put(X, Y) {
        FIND [ on(X, Z) ]
        ERASE [ on(X, Z) ]
        ASSERT[ on(X, Y) ]
    }

It is also possible to write, in Prolog

I don't understand the use of the term THEOREM in this context, and PLAN is, I think exactly what it is.
THEOREM3 is the name of the plan, and it may be used in a explanatory session, but since it is just a arbitrary identifier, it doesn't explain anything.

    put(X, Y) :- find(on(X, Z)) erase(on(X, Z)) assert(on(X, Y));

Thinking about my own form:

    == Pick up X ==
    {
        condition: command(P) pick_up(P, E1),
        action: goal(pick_up(E1)),
        no_results: {
            answer: dont_know()
        },
        some_results: {
            answer: canned('OK')
        }
    }

    pick_up(X) :- grasp(X) raise_hand();
        grasp(X) :- assert(grasping(X));
        raise_hand(X) :- at(`hand`, X Y Z1) add(Z1, 1000, Z2) move_hand(X Y Z2);

goal: resolve variables
find: resolve variables
erase: offer the relation set to each database capable of writing; the database must then remove the relations
assert: offer the relation set to each database capable of writing; the database must then add the relations

# 2019-01-22

I am now logging queries on the dbpedia demo site. This way I get to know how the application is used and what "my users" want.

They definitively want to type human names without capitals, and to just type their last names. I deal with that later.

Yesterday someone asked about the capital of Iraq. That was interesting. Apparently DBpedia has more than one entry for the country Iraq, so the user needs to disambiguate despite the fact that there is only one country at the moment.

I forgot to mention the result of the question my colleague at work asked:

    Who married Kim Kardashian?

The answer of the app was:

    Kanye West, Kris Humphries and The Underdogs married her

Apparently "The Underdogs" is the production team of Kim Kardashian, and someone listed it under "spouse". Funny, but it makes you wonder about the quality of dbpedia.

Back to capitals. "Iran" is not only a (or actually 3) country, but also a person. When someone now asks

    What is the capital of Iran?

He gets the answer: Which one? the Brazilian Footballer?

That's obviously silly. The question implies countries, and I will now try to attempt to extend the system with entity types.
For this I also need to change the order of the domain specific relation phase and entity recognition.

# 2019-01-13

I already fixed the quantification part. Very happy with this! Thing have gotten much more simple.

# 2019-01-12

I found out that quantification only needs to be done for "quantifiers". These exclude numerals.

This means that I can simplify things:

- quantification() will only be applied to dp's with 'every', 'all', 'none'.
- numerals will be treated as simple modifiers
- quantification() will not be processed in the relationizer step (step part will be removed)
- quantification scoping will do the part that was earlier done by the relationizer

I must also change:

- grammar will not contain sentence structures, but be restricted to root()
- syntactic relations will be made to look like the ones from Stanford Parser Universal Dependencies

--- this morning's brainstorm:

how many persons have more than 3 children?

have_child(P, C) person(P) child(C)

act(how_many, P) more_than(C, 3)

resultset

P = 5 C = 17
P = 5 C = 18

P = 6 C = 19
P = 6 C = 20
P = 6 C = 21
P = 6 C = 22

P = 8 C = 31
P = 8 C = 32
P = 8 C = 33
P = 8 C = 34

P = 9 C = 42

Dus het antwoord is 2

count(P) more_than(C, 3)

how many persons have more than 3 children with 2 friends?

friend(C, F) more_than(F, 2)

count(P) more_than(C, 3) more_than(F, 2)

group_by(P, C, F) <- hoe meer dependent, hoe meer naar achteren
having more_than(C, 3) more_than(F, 2)
select count(P)

P = 6 C = 19 F = 102
P = 6 C = 19 F = 103
P = 6 C = 20 F = 108
P = 6 C = 20 F = 109
P = 6 C = 20 F = 110
P = 6 C = 21 F = 120
P = 6 C = 22 F = 121

is at least one of them narrower than the one i told you to pick up?

narrower_than(a, b)
at_least(a, 1)

does every parent have 3 children?

have_child(a, b)
every(a) count(b, 3)

group by a
having count(b) = 3

does at least one parent have 3 children?

have_child(a, b) parent(a) child(b)
at_least(a, 1) count(b, 3)

group by a
having count(b) = 3
count(a) >= 1

which persons have two sons and three daughters?

have_son(p, s) have_daughter(p, d) count(d, 3) count(s, 2)

have_child(p, s) have_child(p, d) count(d, 3) count(s, 2)

p = 1 s = 11 d = 21
p = 1 s = 11 d = 22
p = 1 s = 12 d = 34
p = 1 s = 12 d = 35
p = 1 s = 12 d = 36

	all
	every
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

does every person have two sons and three daughters?

person(p) have_son(p, s) have_daughter(p, d) count(d, 3) count(s, 2) every(p)

every(p) person(p)                          <- maak een lijst met alle personen
	count(d, 3) have_daughter(p, d)         <- filter de personen met 3 dochters
	count(s, 2) have_son(p, s)				<- filter de personen met 3 zonen

quant
	- range (person)
	- quantifier (every)
	- scoped relations

alleen bepaalde ranges hebben moeten controleren of de actuele set overeenkomt met de totale set: every (x uit y), half (1+ uit y)
misschien is het mogelijk syntactisch onderscheid te maken? (part words)
dus: de range mag meestal leeg zijn

alleen bij part-of quants is het verschil tussen de range en de scoped relations van belang

de aggregation relations zijn built-in: all() count() more_than() less_than()

1) heuristieken om de quant, de range en de quantifier te bepalen
quant:
	- de aanwezigheid van een quantifier relatie zorgt voor een quant
quantifier:
	- een van de ingebouwde relaties
range:
	- alleen nodig voor all() half() ...?
	- zoek naar een relatie met de quantifier variable als enige argument

2) we kunnen ook bepalen dat de developer in de generic2domainspec transformatie moet zorgen voor de quantifications

determiner(E, D) -> quantification

nee dat kan niet

3) in plaats van vooraf alle entiteiten uit de range op te halen ... houdt bij welke entiteiten ge-evalueerd worden (ook die worden afgewezen)

4) forceer dat part-of quantifiers gebonden worden aan een entities-type

	isa(p, person) determiner(p, d) isa(d, all) -> every(d, person)

5) aparte variabele

																				lees dit als: de entiteit (Q1) wordt gequantificeerd door D1 over de range R1
	{ rule: np(E1) -> dp(Q1) nbar(R1),                                           sense: new_quantifier(E1, Q1, R1) }

	person(E1) have_son(E1, s) have_daughter(E1, d) count(d, 3) count(s, 2) new_quantifier(E1, Q1, R1)

Na de scope quantification worden de variabelen weer geunificeerd: E1, R1 -> E1

---

! Er bestaat een probleem dat sommige aantallen moeten worden geaggregeerd uit de database en andere aantallen zijn direct te vinden in de database.

Probleem: de transformations werken niet meer goed als je twee verschillende variabelen voor eenzelfde entiteit gebruikt

{ rule: np(E1) -> dp(Q1) nbar(R1),                                           sense: new_quantification(E1, Q1, R1) }
{ rule: np(E1) -> dp(Q1) nbar(E1/R1),                                           sense: new_quantification(E1, Q1, R1) }

root(S5) subject(S5, R5) object(S5, E5) new_quantification(R5, Q5, R6) specification(Q5, A5) isa(A5, how) isa(Q5, many) isa(R6, child) isa(S5, have) name(E5, 'Lord', 1) name(E5, 'Byron', 2)

isa(P1, have) isa(S, child) subject(P1, S) object(P1, O) => have_child(O, S);

6) is every/all niet gewoon een bijzondere uitzondering?

{ rule: np(E1) -> alle_phrase(Q1) nbar(R1),                                           sense: new_quantification(E1, Q1, R1) }

Bij 'every' heb je ook te maken met twee sets entiteiten: de hele groep en de subset.
Dat heb je niet bij 'meer dan 3'.

Bij 'every' is het ook altijd zo dat de entiteiten in de database afzonderlijk zijn opgeslagen, bij 'meer dan 1' hoeft dat niet.

Ik moet hier een naam voor hebben.

<https://en.wikipedia.org/wiki/Quantifier_(linguistics)>

Wow. Natuurlijke getallen zijn helemaal geen quantifiers(!)

<https://en.wikipedia.org/wiki/Numeral_(linguistics)>

Ok, dan moet ik quantification alleen gebruiken voor quantifiers. :)

Does every parent have 4 children?

isa(O, child) isa(P1, have) subject(P1, S) object(P1, O) quantifier(S, Q, R) isa(R, parent)

isa(O, child) isa(P1, have) subject(P1, S) object(P1, O) => have_child(P1, S, O);

dus dat gaat goed

How many children had Lord Byron ?


# 2019-01-11

I want to get rid of the quantification() in the syntactic phase. All relations must be freely modifiable until the solution phase.

quantification() was introduced in this phase (that I presumed to be semantic before) because it was the only place where I could safely gather the range and the quantifier relations.

# 2019-01-09

The relations that are produced by the relationizer should be considered "syntactic relations". Thus far I have considered them as semantic.

The difference is subtle. What matters to me is that some of the relations that are produced are not semantic. And these need to be stripped in the next step.

---

With the previous change the three representations will be:

- Syntactic Relations
- Application Semantics
- Database Relations

---

I looked at Dependency Grammars once again. Apparently DG always uses machine learning to learn the rules for creating the parse tree.

My grammar uses rewrite rules from Phrase Structure grammar, and syntactic relations from Dependency Grammar. I use the commonly used rules for VP, NP etc.
But the rules at the sentence level differ. I will make use of DG style to represent it, because it matches closer to the semantic relations I need.
