
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
