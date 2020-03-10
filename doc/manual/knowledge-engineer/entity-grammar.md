# Entity grammar

The author of this framework created this grammar because he needed a grammar that was both simple to use and powerful.
The grammar is not a traditional grammar that is aimed to distinguish between grammatical utterances and nonsensical
ones. Its aim is to provide a highway to a semantic representation of a sentence.

Let's start with a simple grammar that has only a single rule:

    { rule: s(P1) -> np(E1) vp(P1, E1) }

The syntactic categories (s, np, vp) are familiar to anyone who has come into contact with phrase structure grammars
before. They stand for "sentence", "noun phrase" and "verb phrase". The variables P1 and E1 stand for entities. Entities
play such a central role in this grammar that it is named after them. An entity can stand for anything: persons,
objects, concepts, and even predications (things you can say about something else). The names of the entity variables
start with a letter and this letter can be chosen by you, but it is good practise to make it represent what it stands
for: P = predication, E = any entity. The number is meant to distinguish between several variables of the same type.

We need some extra rules to make the grammar able to represent a complete sentence. The sentence here is:

~~~
Mary likes Jim
~~~

    { rule: s(P1) -> np(E1) vp(P1, E1) }
    { rule: vp(P1, E1) -> iv(P1, E1) }
    { rule: vp(P1, E1) -> tv(P1, E1, E2) np(E2) }
    { rule: np(E1) -> noun(E1) }

There's a lexicon to go with that:

    { form: 'Mary',             pos: noun(E1),          sense: name(E1, 'Mary') }
    { form: 'Jim',              pos: noun(E1),          sense: name(E1, 'Jim') }
    { form: 'likes',            pos: tv(P1, E1, E2),    sense: like(P1, E1, E2)) }

The rules with an arrow ( -> ) are rewrite rules. When a sentence is parsed, the parser starts with the category `s` and
attempts to build the sentence from this `s` by rewriting it with the right hand categories `np` and `vp`. Then, it
tries to rewrite the `np` with `noun` and the `vp` with either `iv` or `tv` followed by `np`. `iv` is short for
"intransitive verb" (a verb without an object), while `tv` stands for "transitive verb" (which does have an object).

The lexicon contains the words that form the terminals or leaves of the parse tree that is created by rewriting the
syntactic categories.
