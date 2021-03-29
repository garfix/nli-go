# Processing

Processing a request consists of these phases:

* Dialog Context: check if the input is just an answer to a question by the system
* Tokenization: from raw text to a string of tokens
* Parsing: from tokens to parse tree with attached generic relations
* Relationizer: combine the relations to a single set with unified variables
* Answering:
    * Conditions: match the relations to the conditions of a solution
    * Evaluation: find answer bindings by evaluating the relations
        * Rule sets are used to expand a relation (as a goal) into a sequence of subgoals (also relations)
        * In memory fact bases are used to look up simple facts about a domain
        * Database fact bases are wrappers around a database (for now: MySql) to read simple records
    * Preparation: find bindings needed to answer the question
    * Answer: create relations that hold the sense of the answer
* Generation: create an array of words from the generic relations
* Surface realization: concatenate the words to raw text
* Store dialog context

Note that there are two types of representation that are expressed by relations:

* Domain Specific: This is the interpretive step. This is the level of reasoning of a domain. Domain specific rules can be used with multiple databases.
* Database: Database relations are optimized for storage.

I will describe the components that make all of this possible.

### Dialog context

The system starts by checking if the user input is not a question, but rather an answer to a question posed earlier by the system.

The system may need to ask the user a question if his/her question is ambiguous and needs clarification.

The system compares user input with 'option' relations in its 'dialog-context' in memory database.

If there is a match, the input is stored as 'answer_open_question' in the dialog context, and the system proceeds to replace the user input (the answer) with the dialog context's 'original_input'.

For example:

* User: Who married Lord Byron?
* System: Which one? [a] the englishman [b] the american // stores 'Who married Lord Byron' as 'original_input', and 'a' and 'b' as 'option's
* User: a
* System: // stores a as the 'answer_open_question', but continues to use 'original_input' as the input the the rest of the request

If there is no match, the input is not taken to be a user answer to a system question, but as a new user question.

Any open options are now discarded.

### Tokenizer

The tokenizer splits a raw line of text in "words" or tokens. A token is either a string of letters/digits, or any other character. All whitespace is discarded.

For example, the sentence

    How many children had Lord Byron?

is split into

    How,many,children,had,Lord,Byron,?

### Parser

The parser is an Earley parser that turns an sequence of tokens into a parse tree with relation attachments.

Earley parsers are efficient and allow for left-recursive grammars, which is really very comfortable.

To use the parser, you need to define a grammar.

Each grammar entry contains a rule and, optionally, a sense. The rule is the syntactic part, but extended with entity variables. The sense consists of the relations that are created when parsing the sentence. The entity variables of the syntax reappear in the relations. Let me give you an example of how this works. When the following rewrite rule has been completed (example clause: John could marry Elsa)

    { rule: np_comp2(E1) -> child(E1) have(_) np(E2),                       sense: go:quant_check($np, dom:have_child(E2, E1)) }

    { rule: interrogative(P1) -> 'how' 'many' np_comp2(E1),                 sense: go:intent(how_many, E1) }

The example sentence (from Tokenizer paragraph) yields the following parse tree:

    [s 
        [sInterrogative 
            [whWord How] 
            [quantifier many] 
            [nbar 
                [noun children]
            ] 
            [auxVerb had]
            [np 
                [properNoun 
                    [firstName Lord]
                    [lastName Byron]
                ]
            ] 
            [questionMark ?]
        ]
    ]

### Relationizer

When the whole tree is parsed all relations will be connected in a relational model. I just call this a relation set.
 Here's the relation set for our sample sentence:

    go_quant_check(go_quant(go_quantifier(R9, R10, go_equals(R9, R10)), E8, none), dom_have_child(E8, E7)) go_intent(how_many, E7) go_intent(question)

### Answerer

The answerer turns a question into an answer. To do this, it goes through the following steps:

#### Find a solution

Each question requires a specific type of answer. To answer a question, a solution must be found. A solution looks like this

    condition: act(question, howMany) child(A, B),
    transformations: []
    responses: [
        {
            condition: go:exists(),
            preparation: gender(B, G) go:count(N, A),
            answer: gender(B, G) go:count(C, N) have_child(B, C)
        }
        {
            answer: dont_know()
        }
    ]

The first solution whose condition matches the question will be used.

Transformations is a list of transformations on the input relations. It is an interpretation of the input sentence. For example:
"How many children has A?" can mean: "get all children and count them" or it can mean "get the number-of-children attribute"

If the condition returns no results, the relation set from "no_results" will be used to phrase the result. Otherwise, the answer of "some_results" will be used.

"preparation" is a relation set that will be solved by the system just to prepare answering the question. Notice that preparation may contain aggregate predicates (i.e. number_of).

"answer" does not connect to any knowledge base. It just formats resulting bindings.

#### Execute the question

Now the question is "executed" as if it were a program. The result of this execution are variable bindings, like this:

    [
        { E1: 1, E2: 4 },
        { E1: 1, E2: 5 },
    ]

This result has two bindings. The first binds E1 to 1 and E2 to 4. In case you are wondering what 1 and 4 are, they are primary keys from relations in the database. More general, they are entity identifiers of some sort.

To evaluate a question, three sources of information may be inspected: in-memory fact bases, rule bases, and databases.

##### In-memory fact base

An in-memory fact base looks like this

    marriages(1, 4, '1815')
    marriages(6, 8, '1889')
    parent(2, 1)
    parent(6, 9)
    person(1, 'Lord Byron', 'M', '1788')
    person(2, 'Lady Lovelace', 'F', '1815')

It's a simple relational database. It can be used to test things, and to store additional information not present in the actual database.

##### Rule base

A rule base holds rules like this:

    siblings(A, B) :- parent(C, A) parent(C, B);

It looks like Prolog, and that's because it behaves like it. Whenever a relation is executed that matches the head of such a rule, the engine enters the rule and executes the tail relations as sub-goals.

Rule bases can be used to make inferences on the information of the database.

##### Databases

To use a database, you must tell the engine how a relation maps to one or more relations in the database. Here's an example

    dom:married_to(A, B) :- marriages(A, B, _);
    dom:name(A, N) :- person(A, N, _, _);
    dom:parent(P, C) :- parent(P, C);
    dom:child(C, P) :- parent(P, C);
    dom:gender(A, male) :- person(A, _, 'M', _);
    dom:gender(A, female) :- person(A, _, 'F', _);

In this example there's just a single relation at both the left (domain) and the right (database) side of the =>, but there could be more. It's a n:m mapping.

#### Preparation

After the question is executed, we have a set of bindings. These bindings are then bound to a sequence of relations called the preparation.

The preparation is meant to collect some more information needed to create the answer.

    preparation: gender(B, G) count(N, A),

In this example, the engine executes 'gender' because the gender is needed in the answer ('He ...'). count() is an aggregate function used to collect the number of children for the answer. This function is performed on the binding set. The different occurrences of A are counted and stored in variable N of all bindings.

#### Answer

The resulting bindings are then bound to the relations of the answer part of the solution, to create the answer.

### Generator

The generator generates a sequence of words based on the generic relations, using a generation grammar. These are different from the ones used for parsing, because there are some differences between parsing and generating sentences.

Here's part of the generation grammar

    rule: s(P) -> np(E) vp(P),                                                  condition: subject(P, E);
    rule: s(C) -> np(P1) comma(C) s(P2),                                        condition: conjunction(C, P1, P2) conjunction(P2, _, _);
    rule: s(C) -> np(P1) conjunction(C) np(P2),                                 condition: conjunction(C, P1, P2);

Generation starts from an s() clause. The first rule that matches condition is used. Next, the consequent of the rule
are used to generate the rest, all the way down, until words can be matched.

### Surface Representation

Finally the words are concatenated by spaces, except for comma's and periods. And the first letter is capitalized.

### Store Dialog Context

The dialog context is persisted in a JSON file.
