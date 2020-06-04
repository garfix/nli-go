# Creating a grammar

Creating a custom grammar is not an easy thing to do. You need to be aware of domain knowledge and language syntax, but next to these you need to have some knowledge of how to map syntactic structures to semantic structures in the given formalism. The formalism here is [entity grammar](entity-grammar.md) and this document gives you examples of this mapping. 

I have constantly been struggling to get these rules right. These best practises may well change in the future, but at the moment these are the best ones I could come up with.

## Sentence level

At the sentence level there are three types of sentences: interrogative (questions), imperative (commands) and declaratives (assertions, teaching).

This then, is basically how one starts a grammar that needs these three types of sentences:

    { rule: s(S1) -> declarative(S1),                                       sense: intent(declare) }
    { rule: s(S1) -> imperative(S1),                                        sense: intent(command) }
    { rule: s(S1) -> interrogative(S1),                                     sense: intent(question) }

As in the rest of this document, you may not need all these rules. Just pick the ones you need and add others later. The art of creating a grammar is to use as few rules as you need, but no less ;)

Note that the head is formed by the build-in top-level node `s(S1)`. All sentences start with this node, and this node is _rewritten_ to one or more other nodes. Together they form a syntax tree. The `sense` you see is attached to the node that is formed by the head of the rule. `intent` relations are used by solutions to recognize certain types of sentences. When executed they always succeed.

It is important to note that a sentence can be interogative in a syntactic sense, while being a command in the semantic sense. "Can you please shut the door?" is not useful as a question. And this form of sentence must the rewritten as `imperative`, since the sense must be a `command`.

## Clauses

A sentence consists of at least one _clause_. It may also contain nested clauses. A clause is part of the sentence with a main predication. If you need nested clauses it may be useful to rewrite a sentence to a clause, so that you can reuse the clause in other parts of the sentence:

    { rule: interrogative(P1) -> interrogative_clause(P1) '?' } 
    { rule: imperative(P1) -> imperative_clause(P1)  '.' }
    
If you want punctuation marks at the end of the sentence to be optional, you can create two sentence level rewrites, like this:    

    { rule: imperative(P1) -> imperative_clause(P1)  '.' }
    { rule: imperative(P1) -> imperative_clause(P1) }
    
## The interrogative sentence

There are many types of questions. The type of question is determined by the way is handled and what the response should look like.

I don't think I found the ultimate way to represent questions, but what I'd like to do is to create a main rewrite rule for the fixed structure:

    { rule: interrogative_clause(P1) -> aux_do(_) do_clause(P1),                        sense: intent(yes_no) }
    
and then refine the changing part in a separate rule:

    { rule: do_clause(P1) -> np(E1) tv(P1, E1, E2) np(E2),                              sense: find([sem(1) sem(3)], sem(2)) }
    
The first rule states that these type of questions start with "do" and that they are yes/no questions, which means that their answer is a simple yes or no. The `intent(yes_no)` in the sense of the rule is a kind of tag that is used by the solution to recognize the type of sentence.

The seconde rule says that the changing part has a NP, a VP and another NP.

Check the example worlds for other examples of questions.        
    
## The imperative sentence

The basic structure of an imperative sentence (a command) is simply this:

    { rule: imperative_clause(P1) -> vp(P1) }    



## Verbs

