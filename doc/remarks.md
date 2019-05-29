# 2019-05-29

ProblemSolver::SolveRelationSet() now solves a set in most-efficient order. Sometimes however, you need it to be executed in the order of the set.
When the sub-goals of a rule are executed, for example.

I have now split up these goals and feed them one-by-one to SolveRelationSet(), but I'd rather have a SolveRelationSetInOrder() function that I could call.

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

https://en.wikipedia.org/wiki/Quantifier_(linguistics)

Wow. Natuurlijke getallen zijn helemaal geen quantifiers(!)

https://en.wikipedia.org/wiki/Numeral_(linguistics)

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
