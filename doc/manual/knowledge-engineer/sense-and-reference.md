# Sense and reference

Since Gottlob Frege's "Über Sinn und Bedeutung" (1892) logicians divide semantics in _sense_ and _reference_. [Wikipedia](https://en.wikipedia.org/wiki/Sense_and_reference) 

The sense of a word or phrase is its logical (relational) construct. The reference of a word is a thing, or things, in the world. Furthermore, the reference must be real (as in tree) and not imaginary (as in unicorn).

## Database reference

When dealing with databases we can make another distinction. An entry in the database may or may not represent a thing in the real world. In any case it is different from it. And a natural language database application deals mainly with this separate layer, not with the objects in the world directly.

Let's call this new layer "database reference" for lack of a better word. We can then draw this diagram:

    word, phrase, sentence
    |
    sense
    |
    database reference
    |
    reference

The reference of a word doesn't really play a role in NLI-GO. The database _is_ the world. Database references are what matters. This fact is exacerbated by the constraint that a _reference_ must be real (not imaginary), and it is perfectly possible to build a database on Lord of the Rings characters, none of which have references (that I know of). 

## Proper nouns

Proper nouns may be said to have a sense (as some kind of description), but in NLI-GO the proper name only has a database reference.

    proper noun
    |
    database reference

The id of an entity can either be found by looking up its name in the database, or by converting a logical construct (sense) into a set of database queries.

## Relations

The relations that form the sense of a phrase or a sentence are grounded in two ways. First there is the mapping to database relations. Second there are built-in functions whose working is preprogrammed.

## The database reference of a sentence

Frege also introduced the idea that the reference of a sentence is its truth value. NLI-GO does not use this idea, because it is not useful in this domain.

Reference as truth value is a very scientific use of the sentence. The use is to find out the truth about the world. It applies well to declarative sentences ("Albert Einstein was born in 1879"), but not to questions ("When was Einstein born?") and commands ("Pick up a red block").

The reference of a sentence in NLI-GO is not defined. But we can speak of the database reference of a sentence. This would be the set of bindings that results from executing/answering the sentence, combined with the set of external actions that was instigated by the system.

