# HOW-TO

## Bidirectional relations

### The relation "spouse" is bidirectional, how do I deal with it?

You can add two lines to a .map file for a knowledge base:

    married_to(A, B) => spouse(A, B);
    married_to(A, B) => spouse(B, A);

or you can add a line to a rules file:

    married_to(A, B) :- married_to(B, A);

## Canned responses

A canned response is just a literal text that may be used as an answer.

To use a canned response, use "canned()" in the answer of a solution, like this:

    {
        condition: question() who(B),
        no_results: {
            answer: dont_know()
        },
        some_results: {
            preparation: long_description(B, D),
            answer: canned(D)
        }
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

