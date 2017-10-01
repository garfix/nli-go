# HOW-TO

## Bidirectional relations

### The relation "spouse" is bidirectional, how do I deal with it?

You can add two lines to a .map file for a knowledge base:

    married_to(A, B) => spouse(A, B);
    married_to(A, B) => spouse(B, A);

or you can add a line to a rules file:

    married_to(A, B) :- married_to(B, A);
