# Code

The relations of the response are constructed in the solution file. These are language agnostic.

The generation phase converts these relations into a line of natural language text. This requires a generation grammar.

## Grammar

The grammar contains rewrite rules to rewrite the top-level construct `s()` into leaf nodes.

Here's an example `s()` rewrite rule for a response like "Peter and John" 

    { rule: s(C) -> np(P1) 'and' np(P2),                                 condition: and(P1, P2) }

The rule says: the syntax tree of the response contains the nodes `np(P1) 'and' np(P2)` if the relation `and(
P1, P2)` is present in the response. The contents of the variables `P1` and `P2` is bound to the syntax tree nodes.

## Literal text

Sometimes it is needed to output the literal contents of a variable. When it holds the name of a person, for instance.
In this case you can use `text` to include the text directly.

    { rule: proper_noun(E1) -> text(Name),                                        condition: name(E1, Name) }

## Rule condition

In some cases it's more natural to provide the condition as part of the head of the rule. In this example numbers for months are mapped to names.

    { rule: month_name('01') -> 'January' }
    { rule: month_name('02') -> 'February' }
    { rule: month_name('03') -> 'March' }

The head of `month_name` contains a string here. When the generator needs to generate text for `month_name(C)`, `month_name('02')` matches if `C` is bound to `02`.
 