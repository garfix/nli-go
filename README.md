# NLI-GO

## Welcome to the nli-go wiki!

nli-go is a library, written in Go, that provides a natural language interface to databases. It is in fact just my hobby-project, but I am trying hard to make this the best nli system ever created!

For years to come, this system is not easy to use, not present for production, and not very robust either. Still, if you really need an nli, it may be worth your trouble to learn it. I think it's pretty cool.

## Purpose

This library helps a developer to create a system that allow end-users to use plain English / French / German to interface with a database. That means that an end user can type a question like

>  How many children had Lord Byron?

and the library looks up the answer, and formats the result, also in natural language:

> He had 2 children.

Every part of the system is configurable.

## Main parts

The system consists of the following main parts:

### Tokenizer

The tokenizer splits a raw line of text in "words" or tokens. A token is either a string of letters/digits, or any other character. All whitespace is discarded.

For example, the sentence

> How many children had Lord Byron?

is split into

> How,many,children,had,Lord,Byron,?

### Parser

The parser is an Earley parser that turns an sequence of tokens into a set of relations. It also produces a parse tree, but that's just because it can. The parse tree is not used further down the pipeline.

Earley parsers are efficient and allow for left-recursive grammars, which is really very comfortable.

To use the parser, you need to define a grammar and a lexicon. Here's part of the lexicon:

    form: 'how',        pos: whWord,        sense: isa(E, how);
    form: 'many',       pos: adjective,     sense: isa(E, many);
    form: 'has',        pos: auxVerb,       sense: isa(E, have);
    form: 'marry',      pos: verb,          sense: isa(E, marry);
    form: /^[A-Z]/,     pos: firstName,     sense: name(E, Form, firstName);
    form: /^[A-Z]/,     pos: lastName,      sense: name(E, Form, lastName);
    form: 'are',        pos: auxVerb,       sense: isa(E, be);
    form: 'and',        pos: conjunction;
    form: 'children',   pos: noun,          sense: isa(E, child);
    form: '?',          pos: questionMark;

Each entry in the lexicon has a surface form (form), a part of speech (pos), and a set of relations (usually just one). The sense is kept close to the surface form. This is not accidental. I am for the lexicon to be completely generic. There should be just one lexicon for each language (or dialect if necessary). All semantic complexity should come later in the process. So the meaning of 'many' is 'many' and the meaning of 'children' is 'child' (the fact that a plural is meant will be added in a later release).

As you can see I allow regular expressions to match proper nouns (i.e. Byron).

The grammar looks a bit like this (just some parts of it)

    rule: sInterrogative(S1) -> whWord(W1) adjective(A1) nbar(E1) auxVerb(S1) np(E2) questionMark(),
    sense: question(S1, whQuestion) focus(E1) determiner(E1, A1) specifier(A1, W1) subject(S1, E1) object(S1, E2);

    rule: clause(S1) -> np(E1) vp(S1),                                         sense: object(S1, E1);
    rule: clause(S1) -> np(E1) modal(M) vp(S1),                                sense: subject(S1, E1) modality(S1, M);

    rule: nbar(E1) -> noun(E1);
    rule: nbar(E1) -> adjp(A1) nbar(E1),                                       sense: specification(E1, A1);

    rule: vgp(V1) -> verb(V1);
    rule: vgp(V1) -> verb(V1) particle(P1),                                    sense: modifier(V1, P1);
    rule: vgp(V1) -> modal(A1) verb(V1),                                       sense: modality(V1, A1);

Each grammar entry contains a rule and, optionally, a sense. The rule is the syntactic part, but extended with entity variables. The sense consists of the relations that are created when parsing the sentence. The entity variables of the syntax reappear in the relations. Let me give you an example of how this works. When the following rewrite rule has been completed (example clause: John could marry Elsa)

> clause(S1) -> np(E1) modal(M) vp(S1)

The following relations are created

> subject(S1, E1) modality(S1, M)

Because subject and modality use the same variable, S1, these relations are connected. When the whole tree is parsed all relations will be connected in a relational model. I just call this a relation set.
 Here's the relation set for our sample sentence:

> [question(S1, whQuestion) focus(E1) determiner(E1, A1) specifier(A1, W1) subject(S1, E1) object(S1, E2) isa(W1, how) isa(A1, many) isa(E1, child) isa(S1, have) name(E2, 'Lord', firstName)  name(E2, 'Byron', lastName)]

