# Anaphora resolution

NLI-GO handles pronouns like "he", "she", "it", but also expressions like "the block" and "the red one". These
expressions are called commonly called anaphora. Some refer to recent user input and some to a recent system response.
It also handles references to entities mentioned earlier within the same sentence (intra-sentential).

## The anaphora queue

NLI-GO uses a structure called the anaphora queue to store the most recent references. The queue is located in the
dialog context and is stored on file.

The queue is simply a queue (first in first out) of entity references (id + entity type). New entities are stored on
front. When the queue becomes larger than 10 items, the ones at the end fall off.

## Pronoun grammar

Pronouns must have a very specific sense in the grammar, in order to be treated as quants.

Here's "it" for example:

    { rule: pronoun(E1) -> it(E1),                                          sense: quant(_, [ the(_) ], E1, [], sem(parent)) }

and here's "she"

    { rule: pronoun(E1) -> it(E1),                                          sense: quant(_, [ the(_) ], E1, [ female(E1) ], sem(parent)) }

## Adding entities

Any entity that is part of the range of a quant is added to the queue. This allows a user to refer to entities in the
same sentence and to entities in previous input sentences.

Any entity that is part of the result set of a question is added to the queue. This allows a user to refer to a previous
response.

## Handling anaphora

Anaphora is handled in the quant solver. A quant has a range, such as "blue block" and a quantifier like "the".

Normally the solver will fetch the ids of all entities in the range, and use each in turn to process the scope.

In the case where the quantifier is "the" the solver will first try the entities in the anaphora queue, before fetching
all entities in the range. When the queue contains the ids person:5 and person:13, it will first attempt to resolve

    blue block (with id person:5)

If this does not return 1 result, it will continue to the next

    blue block (with id person:13)

If this returns 1 result, person:13 is taken to be the range of the quant. If not, the rest of the ids in the queue are
tried. Only when all fail, will the solver just use the range without id binding.
