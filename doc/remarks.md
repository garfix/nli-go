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
