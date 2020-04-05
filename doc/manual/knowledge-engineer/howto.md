# HOW-TO

## Names

To parse a name that consists of n words (here: 2 or 3), use this grammar

    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1) proper_noun(N1)}
    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1)}

    { rule: np(E1) -> proper_noun_group(E1) }

The "proper_noun" is what counts. It is a reserved word that matches to anything. The string of proper_nouns will be
passed, with separating spaces to the NameResolver, which will try to match it against any of the knowledge bases.

The syntax scheme is as follows. Its an example:

It starts with a sense from a rule; here: marry(P1, E1, E2)

    { rule: np_comp4(P1) -> np(E1) marry(P1) to(P1) np(E2),     sense: marry(P1, E1, E2) }

It introduces two noun phrases (E1 and E2). You define their types in predicates.json:

    "marry": { "entityTypes": ["event", "person", "person"] },

Both noun phrases are persons. This type of definition is called s-selection (semantic selection of predicate
arguments).

You also need to tell the system how to find the name of such an entity type. This is done in entities.json:

    "person": {
        "name": "person_name(Id, Name)",
        "knownby": {
          "description": "description(Id, Value)"
        }
    },

This is how the system finds a database id to match a name in the sentence "Where was Michael Jackson born?".

Now when the parser parses a sentence, it will come across the marry() predicate. From this is learns to expect two
persons for the NP arguments. The np leads to a proper_noun_group and this leads to 2 or 3 proper_noun's. When it
contains 2 proper_noun's, the parser will take the first proper noun under consideration ("Michael") and add the next
one ("Jackson") to create "Michael Jackson". It will learn from entities.json how to query the knowledge bases. If the
name is ambiguous, the user will be asked to disambiguate. This is where the description is for. When the name is found,
the proper_noun_group rule matches. The rule does not create any relations, but the variable-id pair will be stored in a
Binding. This allows for different values in each of the knowledge bases.

## Bidirectional relations

### The relation "spouse" is bidirectional, how do I deal with it?

You can add two lines to a .map file for a knowledge base:

    married_to(A, B) :- spouse(A, B);
    married_to(A, B) :- spouse(B, A);

or you can add a line to a rules file:

    married_to(A, B) :- married_to(B, A);

## Canned responses

A canned response is just a literal text that may be used as an answer.

To use a canned response, use "canned()" in the answer of a solution, like this:

    {
        condition: question() who(B),
        responses: [
            {
                condition: exists(),
                answer: canned(D)
            }
            {
                answer: dont_know()
            }
        ]
    }

As you see the "answer" in the solution contains the single relation "canned()". When that happens, the contents of its variable will be used as the response.

## Specify entity-types for predicate arguments

Create a file predicates.json, for example like this

    {
      "has_capital": {"entityTypes": ["country", "city"] }
    }

This file specifies the entity types of the arguments of the domain specific predicate "has_capital".

The entity types used here are the same as in the entities file.

Add the file to the config file.

    {
      "predicates": [
        "predicates.json"
      ]
    }

## Negation

It is possible to use "not" in a simple case.

Here's an example from the blocks world: "How many blocks are not in the box?"

"not" is modelled in the grammar:

    { rule: how_many_clause(E1) -> np(E1) copula() not() pp(E1),           sense: not(sem(4)) }

not() is a "nested structure" that wraps a relation set.

This set is specified in the example as "sem(4)". This means: the combined senses of all syntactic structures that were
linked to the fourth consequent (which is "pp(E1)").

A not() predicate can only be evaluated correctly when it is evaluated as part of a quant scope.