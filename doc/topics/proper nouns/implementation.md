# Proper nouns - implementation

We can create rewrite rules to detect proper nouns of 1 to 3 words like this:

    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1) proper_noun(N1) }
    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1) }
    { rule: proper_noun_group(N1) -> proper_noun(N1) }

    { rule: np(E1) -> proper_noun_group(E1) }

Where `proper_noun` is a reserved category that matches any string without withespace. Because of this, this construct quickly produces multiple parse trees. Some of them are quite bizarre. Even dots and comma's can be part of them.

After parsing, the sorts of the entities are determined. Only then the system tries to find the name in the database. Since it knows the sort, it knows in which table to look.
