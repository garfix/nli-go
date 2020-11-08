## Names, or proper nouns

To parse a name that consists of n words (here: 2 or 3), use this grammar

    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1) proper_noun(N1)}
    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1)}

    { rule: np(E1) -> proper_noun_group(E1) }

The "proper_noun" is what counts. It is a reserved word that matches to anything. The string of proper_nouns will be
passed, with separating spaces to the NameResolver, which will try to match it against any of the knowledge bases.

The syntax scheme is as follows. Here's an example:

It starts with a sense from a rule; here: marry(P1, E1, E2)

    { rule: np_comp4(P1) -> np(E1) marry(P1) to(P1) np(E2),     sense: marry(P1, E1, E2) }

It introduces two noun phrases (E1 and E2). You define their types in predicates.relation:

    marry(event, person, person)

Both noun phrases are persons. This type of definition is called s-selection (semantic selection of predicate
arguments).

You also need to tell the system how to find the name of such a sort. This is done in entities.json:

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