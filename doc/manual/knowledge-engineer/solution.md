# Solutions

Each type of sentence has its own solution. A single solution can handle multiple sentences, but their top-level
structure is the same, they are processed the same way, and their response is the same, except for some variables.

Here is an example that contains all sections of a solution:

~~~
    {
        condition: question(_) how_many(B) have_child(A, B),
        transformations: 
            have_child(A, B) => have_n_children(A, Number);
        ,
        responses: 
            {
                condition: exists(),
                preparation: gender(A, Gender),
                answer: gender(A, Gender) have_child(A, C) count(C, Number)
            }
            {
                answer: dont_know()
            }
    }
~~~

The condition specifies what the input relations must look like for this solution to apply.

The transformations transform some relations from the input to the set that will actually be processed. This is kind of
"interpretation" on the part of the system. (When the user says A he actually means B). Transformations are optional.

A solution can have several responses, that depend on certain conditions. So each response has an optional condition.

Sometimes you need some extra information in the answer that was not retrieved in the question itself. For example you
may want to to answer "He had 2 children", whereas the gender of the subject was not part of the question. "preparation"
allows you to fetch this extra information.

"answer" are just a passive set of relations passed to the generator. They are not processed in any way. All variable bindings must already be available at this point.

A special function of "answer" is "make_and().

    preparation: make_and(A, And, R)

This construction allows you to create a nested structure of AND's, so that you can respond with

    John, Kale and Louis

## Canned responses

A canned response is just a literal text that may be used as an answer.

To use a canned response, use "canned()" in the answer of a solution, like this:

    {
        condition: question() who(B),
        responses: 
            {
                condition: exists(),
                answer: canned(D)
            }
            {
                answer: dont_know()
            }
    }

As you see the "answer" in the solution contains the single relation "canned()". When that happens, the contents of its variable will be used as the response.