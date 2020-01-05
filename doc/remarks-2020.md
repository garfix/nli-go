# 2020-01-05

I made case-insensitive names possible, and at the same time checking the database in the parsing process. Introduced
s-selection to help finding the names. s-selection restricts predicate arguments and this in turn narrows the search
space for proper nouns in the database.

# 2020-01-02

Happy new year! 

I am introducing semantics in the parsing process, because I need some semantics to determine the type of the entity in
a name.

I want to use the relationizer that I already have for this, but it is too much linked to the nodes that I generate
after the parse is done.

Now I just had an interesting idea: what if I do the sense building as part of the chart states. That way, when the
parse is done, I just need to filter out the widest complete states and I will have the complete sense ready, without
having to create a tree and then parse that tree.
