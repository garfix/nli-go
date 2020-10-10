# Creating a grammar

Creating a custom grammar is not an easy thing to do. You need to be aware of domain knowledge and language syntax, but next to these you need to have some knowledge of how to map syntactic structures to semantic structures in the given formalism. The formalism here is [entity grammar](entity-grammar.md) and this document hopes to guide you in this process. 

I have constantly been struggling to get these rules right. These best practises may well change in the future, but at the moment these are the best ones I could come up with. While reading this over I notice how much this is still lacking, and some rules are just wrong, but it is a start, and surely better than nothing at all.

## Sentence level

At the sentence level there are three types of sentences: interrogative (questions), imperative (commands) and declaratives (assertions, teaching).

This then, is basically how one starts a grammar that needs these three types of sentences:

    { rule: s(S1) -> declarative(S1),                                       sense: intent(declare) }
    { rule: s(S1) -> imperative(S1),                                        sense: intent(command) }
    { rule: s(S1) -> interrogative(S1),                                     sense: intent(question) }

As in the rest of this document, you may not need all three rules. Just pick the ones you need and add the others later. The art of creating a grammar is to have as few rules as possible, but no less ;)

Note that the head is formed by the built-in top-level node `s(S1)`. All sentences start with this node, and this node is _rewritten_ to one or more other nodes. When all rules are applied, a syntax tree if formed. The `sense` you see is attached to the node that is formed by the head of the rule. An `intent` relation is used in a condition of one or more solutions to recognize certain types of sentences.

It is important to note that a sentence can be interrogative in a syntactic sense, while being imperative in the semantic sense. "Can you please shut the door?" is not useful as a question. And this form of sentence must the rewritten as `imperative`, since the sense must be a command.

## Clauses

A sentence consists of at least one main _clause_. It may also contain nested clauses. A clause is a part of the sentence with a verb. If you need nested clauses it may be useful to rewrite a sentence to a clause, so that you can reuse the clause in other parts of the sentence:

    { rule: interrogative(P1) -> interrogative_clause(P1) '?' } 
    { rule: imperative(P1) -> imperative_clause(P1)  '.' }
    
If you want punctuation marks at the end of the sentence to be optional, you can create two sentence level rewrites, like this:    

    { rule: imperative(P1) -> imperative_clause(P1)  '.' }
    { rule: imperative(P1) -> imperative_clause(P1) }
    
## The interrogative sentence

There are many types of questions. The type of question is determined by the way is handled and what the response should look like.

I don't think I found the ultimate way to represent questions, but what I'd like to do mostly is to create a main rewrite rule for the fixed structure:

    { rule: interrogative_clause(P1) -> aux_do(_) do_clause(P1),                        sense: intent(yes_no) }
    
and then refine the changing part in a separate rule:

    { rule: do_clause(P1) -> np(E1) tv(P1, E1, E2) np(E2),                              sense: quant_check($np1, quant_check( $np2, $tv)) }
    
The first rule states that these type of questions start with "do" and that they are yes/no questions, which means that their answer is a simple yes or no. The `intent(yes_no)` in the sense of the rule is a kind of tag that is used by the solution to recognize the type of sentence.

The seconde rule says that the changing part has a NP, a VP (a transitive verb) and another NP.

Check the example worlds for other examples of questions.        
    
## The imperative sentence

The basic structure of an imperative sentence (a command) is simply this:

    { rule: imperative_clause(P1) -> vp(P1) }
    
## The declarative sentence

A declarative sentence states something to be the case. The aim of the sentence is to teach the system something. This may be either a simple fact, or a rule. The predicate that is used for both these cases is `assert`. It adds information to a knowledge base or a rule base.

Let's start with simple declarative sentences. This one handles sentences like: "all red blocks are mine"

    { rule: declarative(P1) -> np(E1) copula(_) np(E2),         sense: assert( own(A, B) :- quant_check($np1, quant_check($np2, [equals(A, E2) equals(B, E1)]))) }        

A `copula` is a verb like "is" and "are" when it has no meaning of its own in the sentence. In the sense you can see the `assert` relation whose single argument is a rule. When this declarative is executed, the `assert` adds a rule to a rule base (the first rule base known to the system). 

If simple rules are all you need, you can leave it at that. But when you need to handle sentences with exceptions (using words like "but" or "except"), things become a lot more complicated, but not impossible.

To enable rules _with exceptions_ the construction becomes more complex. This is the start:

    { rule: declarative(P1) -> default_rule(P1) }
    { rule: declarative(P1) -> default_rule(P1) exception(P2) }

There are rules, and there are rules with exceptions.

    { rule: exception(P1) -> 'but' assertion(P1) }
    { rule: exception(P1) -> ',' exception(P1) }
    { rule: default_rule(P1) -> assertion(P1) }

This is how to implement "all red blocks are mine, but the blocks in the box are not mine", with an optional comma. The declarative is rewritten to a default rule with a possible exception. A default rule is a rule that explicitly states that it allows for exceptions. 

    { rule: default_rule(P1) -> np(E1) tv(P1, A, B) np(E2),                           sense: assert(
                                                                                        $tv :-
                                                                                        quant_check($np, quant_check($np2, equals(A, E1) equals(B, E2) not( -$tv )))) }

    { rule: assertion(P1) -> np(E1) dont(_) tv(P1, A, B) np(E2),                      sense: assert(
                                                                                        -$tv :-
                                                                                        quant_check($np1, quant_check($np2, equals(A, E1) equals(B, E2))) }

These last rules form the default rule and the exception. They are posed in a general way that allows for multiple application. The first says "NP verb NP", which handles clauses like "I like ice" and "you own the table". Note the sense: the head is `$tv`, which is the meaning of the second child in the rule (`tv(P1, A, B)`), and this head reoccurs later on in rule, in the form of  `not( -$tv )`. The meaning here is "I own all red blocks except for the ones that I don't own".

The second rule, the `assertion`, states in a general way "I don't own the blocks in the box". Note that the head of the rule to be taught is `-$tv`: a negative goal.

## Noun phrases

A noun phrase (NP) represents a quantified entity. An entity is a person, a block, a chair, a dragon, a profession; in short, anything that has an individual identity; and notably, has a unique ID in a database.

A quantified entity is a group of entities, where the size of the group is determined by a quantifier: `all`, `at least one`, `two ore more`.

All entities are quantified, even the ones like "I" and "people". While there is only one I, and always many people, it is necessary to state this explicitly.

The basic rewrite rule for `np` is this:

    { rule: np(E1) -> qp(_) nbar(E1),                                      sense: quant($qp, E1, $nbar) }
    
It says than an NP has a quantifier phrase and a proper noun phrase, or `nbar`. The `quant` combines the quantifier and the entity into a quantified entity, or `quant` for short.

While this is the basic rule, there are others, because sometimes the quantifier is implicit and needs to be made explicit:

    { rule: np(E1) -> nbar(E1),                                            sense: quant(quantifier(Result, Range, greater_than(Result, 0)), E1, $nbar) }

## Quantifier phrases

Here is how quantifiers are modelled:
    
    { rule: qp(_) -> quantifier(Result, Range),                            sense: quantifier(Result, Range, $quantifier) }
    { rule: quantifier(Result, Range) -> 'all',                            sense: equals(Result, Range) }
    { rule: quantifier(Result, Range) -> an(_),                            sense: greater_than(Result, 0) }
    { rule: quantifier(Result, Range) -> 'at' 'least' 'one' 'of',          sense: greater_than(Result, 0) }
    { rule: quantifier(Result, Range) -> number(N1),                       sense: equals(Result, N1) }
    { rule: quantifier(Result, Range) -> 'two',                            sense: equals(Result, 2) }
    
Since the quantifier `quantifier(Result, Range, greater_than(Result, 0))` (i.e. the existential quantifier) is very common, and may be needed in multiple places, you can use the atom `some` in its place. For example:

    quant(some, E1, $vp)
      
## Pronouns

A pronoun refers to an entity and can be represented thus:

    { rule: np(E1) -> pronoun(E1),                                         sense: quant(quantifier(Result, Range, greater_than(Result, 0)), E1, $pronoun) }
          
    { rule: pronoun(E1) -> 'you',                                          sense: you(E1) }
    { rule: pronoun(E1) -> 'i',                                            sense: i(E1) }
    { rule: pronoun(E1) -> 'it',                                           sense: back_reference(E1, []) }
    
## Proper nouns

Proper nouns, or names of things and people, are treated specially. They are not looked up in the grammar, but in the database. And this takes some extra work:

Here are the basics:

    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1) proper_noun(N1) }
    { rule: proper_noun_group(N1) -> proper_noun(N1) proper_noun(N1) }
    { rule: proper_noun_group(N1) -> proper_noun(N1) }

    { rule: np(E1) -> proper_noun_group(E1),                              sense: quant(quantifier(Result, Range, equals(Result, Range)), E1, []) }
    
Since this is an NP, it needs a quantifer; let's use the all-quantor.

The group exists of 1, 2 or 3 words, hence 3 rules. Do not use recursion. This is a heavy operation that uses much database access, so we don't want to try more words than necessary.

At the moment the `proper_noun_group` is processed, the parser has received top-down information about the entity type of the holder of the name. The parser knows that it is the name of a person. How? 

Let's have an example:

    { rule: nbar(E1) -> 'daughter' 'of' np(E2),                                sense: quant_check($np, has_daughter(E2, E1)) }
    
Here `np(E2)` will be rewritten to the name "Charles Babbage". The parser also sees that E2 is the first argument of the relation `has_daughter(E2, E1)`. And you can tell the system that the first argument of this relation is a person, by adding this line to the file "predicates.relation":

    has_daughter(person, person)
    
With this information, the parser knows that "Charles Babbage" is not the name of a book, but of a person. And it uses another file (entities.yml) to understand how to query the name in the knowledge base:

     person:
        name: person_name(Id, Name)
        knownby:
          description: description(Id, Value)
        
It uses the "name" property to find the relation needed to find out the id of the entity, given its name. "knownby" is used only for disambiguation. The system can ask the user "This Charles Babbage" or "That Charles Babbage"?                                 
  
## Determininer phrases

Determiners are words like "the" and "that" that hold a reference to an NP mentioned earlier. 

    { rule: np(E1) -> the(E1) nbar(E1),                                    sense: quant(quantifier(Result, Range, equals(Result, 1)), E1, definite_reference(E1, $nbar)) }      

## Verb phrases

There are transitive verbs (`tv`) and intransitive verbs (`iv`). Transitive verbs take a subject, an object and possibly an indirect object. Intransitive verbs just take a subject. Objects are mostly NP's, but complete clauses as well.

Here's a simple rewrite for a transitive verb phrase:

    { rule: vp(P1, E1) -> marry(P1) np(E2),                                sense: quant_check($np, marry(P1, E1, E2)) }
    
The `find` combines the `quant` or `quant`s with the sense of the verb. `find` iterates over all combinations of values for the quants.

Verb phrases can have different persons (first, second, third), pluralization, tenses (past, present, future), modals (can, should) and passivization. The word forms for the verb are different in these cases. Auxiliaries are used to specify past perfect. 

NLI-Go does not have a way to constrain agreement between NP and VP.

## Adverb phrases

Phrases that modify a verb are rarely used in an NLI. But it is possible.

The use of the variable `P1` that stands for the predication (the event, in most cases) itself. It allows you to modify the verb with adverbs ("give quickly")     

## Relative clauses

A relative clause modifies a noun phrase, like this:

    { rule: nbar(E1) -> noun(E1) relative_clause(E1) }
    { rule: relative_clause(E1) -> 'which' copula(C1) adjp(E1) }

This one allows for phrases like "(the cat) which is big", where "which is big" is the relative clause. Relative clauses can be quite large ("(the one) I told you to pick up").

## Adjective phrases

Adjective phrases (`adjp`) are phrases that modify a noun phrase.

### Attributive adjective phrases

An attributive adjective preceeds a noun phrase. There can be multiple adjective phrases, as in "the big red block".

    { rule: nbar(E1) -> adjp(E1) nbar(E1) }
    
### Predicative adjective phrases

Phrases like "is taller than" are also adjective, but in a predicative way.

    { rule: relative_clause(E1) -> 'which' copula(C1) adjp(E1) }
    { rule: adjp(E1) -> 'taller' 'than' np(E2),                            sense: quant_check($np, taller(E1, E2)) }
    
## Prepositional phrases

These phrases denote a relation between two noun phrases. Here's the rule for the preposition "in".

    { rule: pp(E1) -> 'in' np(E2),                                         sense: quant_check($np, contain(_, E2, E1)) }
    
## Numbers

Numbers form an open-ended group of words. 

This is now numbers are defined, using a regular expression.    

    { rule: number(N1) -> /^[0-9]+/ }
    
Here is how a number is used, as part of a quantifier: 

    { rule: quantifier(Result, Range) -> number(N1),                       sense: equals(Result, N1) }
           
In this example, the variable `N1` is replaced by the number in the relationization phase. So the resulting sense is

    equals(Result, 22)
    
given that the word "22" was used in the sentence.    

## Other open ended word forms

Other open-ended word forms can be treated with regular expression, like numbers. You may need to change the regular expression for the tokenizer to create separate tokens for this. You can set one with the field "tokenizer" in the config.json.

## Conjunctions

Here we just treat the conjunctions "and" and "or". They can be used on several levels of the sentence. In this example "and" connects two clauses:  

    { rule: imperative_clause(C) -> imperative_clause(P1) and(_) imperative_clause(P2),         sense: and($imperative_clause, $imperative_clause) }    

## Negations

A negation introduces the relation `not()` and can also occur on several levels of the sentence. Here's an adjective: 

    { rule: adjp(E1) -> 'not' adjective(E1),                               sense: not($adjective)) }

## Extraposition

Long term dependencies (extraposition) can be implemented using extra variables passed to the head.

"FillerStack_test" provides an example.

