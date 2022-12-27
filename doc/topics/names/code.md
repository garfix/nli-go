# Proper nouns - code

To parse a name that consists of n words (here: 1, 2 or 3), use this grammar

    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1) proper_noun(N1) }
    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1) }
    { rule: proper_noun_group(N1) -> proper_noun(N1) }

    { rule: np(E1) -> proper_noun_group(E1) }

The "proper_noun" is what counts. It is a reserved word that matches to anything. The string of proper_nouns will be
passed, with separating spaces to the `NameResolver`, which will try to match it against any of the knowledge bases.

The syntax scheme is as follows. Here's an example:

It starts with a sense from a rule; here: marry(P1, E1, E2)

    { rule: np_comp4(P1) -> np(E1) marry(P1) to(P1) np(E2),     sense: marry(P1, E1, E2),   tag: go:sort(E1, person) go:sort(E2, person) }

Notice the `go:sort` tags that tell the system what sorts the arguments have.

You also need to tell the system how to find the name of such a sort. This is done in `sort-properties.yml`:

    person:
        name: person_name(Id, Name)
        knownby:
            description: description(Id, Value)

This is how the system finds a database id to match a name in the sentence "Where was Michael Jackson born?".

Now when the parser parses a sentence, it will come across the marry() predicate. From this is learns to expect two
persons for the NP arguments. The np leads to a proper_noun_group and this leads to 2 or 3 proper_noun's. When it
contains 2 proper_noun's, the parser will take the first proper noun under consideration ("Michael") and add the next
one ("Jackson") to create "Michael Jackson". It will learn from `sort-properties.yml` how to query the knowledge bases. If the
name is ambiguous, the user will be asked to disambiguate. This is where the description is for. When the name is found,
the proper_noun_group rule matches. The rule does not create any relations, but the variable-id pair will be stored in a
Binding. This allows for different values in each of the knowledge bases.

