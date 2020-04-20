# Anaphora resolution

NLI-GO keeps track of the entities that it last processed in a structure called "anaphora queue". This queue consists of
the id's of these entities.

Whenever a quant is processed, the system will first try to resolve the range of the quant with each of
these ids filled in. When one of these ids gives a match, this id will be used as the range of the quantification.

In order to allow pronouns like "he", "she" and "it" in the input, you need to model pronouns in a way that reflects
their function as a quantification:

    { rule: pronoun(E1) -> it(E1),                                          sense: quant(_, [the(_)], E1, []) }

This basically says: by "it" I mean "the latest processed entity".