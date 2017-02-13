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

    object(Q, O) specification(O, S) isa(S, many) specification(S, T) isa(T, how) -> numberOf(O, N) focus(Q, N)

This forms

    question(Q) isa(Q, have) subject(Q, S) name(S, 'Janice', fullName) object(Q, O) isa(O, child) numberOf(O, N)

Can we execute this? No we have to combine "have" with "child"

    isa(Q, have) subject(Q, S) object(Q, O) isa(O, child) -> child(S, O)

This gives us

    question(Q) child(S, O) name(S, 'Janice', fullName) numberOf(O, N) focus(Q, N)

Can we execute this? Yes, after child() and name() are processed, there are 3 possible value for O left. Processing numberOf() fills N with 3.

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

 Second order predicates like numberOf() and max() act on result sets.

 ====

 When I started programming this I came across the problem that for some questions you have multiple answers. Can we handle these?

 Who were Mary's children?

    answer: name(C, N)


## 2017-02-02

There are several reasons why quantifier-constructs (exists, numberOf) should not be added to the lexicon:

 * the word itself is not always enough to determine the quantifier ('how many': the combination of these words means 'numberOf').
 * expressions can always give surface expressions another meaning than is apparent from the words. (every now and again, worth every penny)
 * some quantifiers cannot be deduced from the words alone and must be added later on ('not a lot', 'very little (people voted for Hillary)')

Trying some things:

How many children had Beatrice?

    solution: [
        condition: act(interrogation) focus(O) child(S, O)
        plan: numberOf(child(S, O), N)
        answer: numberOfAnswer(N)
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

    act(interrogation) focus(N) numberOf(hasChild(P, C), N) name(P, 'Janice')

New solution

    solutions: [
        condition: act(interrogation) focus(N) numberOf(hasChild(P, C), N),              // a question, about the number of children
        prep: gender(P, G),                                                              // look up the gender of the parent
        answer: gender(P, G) hasChild(P, C) numberOfAnswer(N);                           // "she, has children, number"
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

    isa(O, child) specification(O, S) isa(S, many) specification(S, T) isa(T, how) -> act(interrogation) focus(N) numberOf(O, isa(O, child), N)
    isa(Q, have) subject(Q, S) object(Q, O) isa(O, child) -> child(S, O)

By 'solution' I mean the matching of a question to an answer

    solution: [
        condition: act(interrogation) focus(N) numberOf(O, isa(O, child), N) child(S, O)
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

## 2016

Als je wilt dat de representatie een Horn clause repr is, moet je NOT en OR expliciet noemen.
Maar is het wel mogelijk om deze in de eerste parse op te nemen, of zijn
Maar is het wel nodig om ze op te nemen? Je kunt de meeste determiners ook niet opnemen.
En je neemt modale elementen (ik dacht dat ..) ook niet op
Ok, maar daarmee is je representatie echt NIET logisch te noemen
Niet alleen komen EN en OF niet overeen met hun logische equivalenten en is pragmatische interpretatie mogelijk,
    ook is keihard weergegeven NIET te beperkt, omdat er ook MISSCHIEN en NAUWELIJKS bestaan.
FOPC is in zijn algemeenheid gewoon te beperkt, en er is geen goed alternatief.

Er is een probleem met left-recursion in de simpele parser NP :- NP VP

een agent

agent: {
    grammar: {
        rules ...
    }
    lexicon: [
        entries ...
    ]
}

een lexicon op zichzelf:

lexicon: [
    {
        form: ..
        pos: ..
        sense: ..
    }
    {
        form: ..
        pos: ..
        sense: ..
    }
]




		predication(S1, marry)
		object(S1, E2)
		subject(S1, who)
		name(E1, 'Kurt Cobain')

Ik maak 'grammatical_subject' nu het predicaat dat aangeeft waar de hoofdzin is. Een predicatie-object is niet aanwezig in het domein-specifieke resultaat. 
Dit grammatical_subject geeft ook aan of de zin actief is of passief.
