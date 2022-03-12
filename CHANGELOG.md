# Changelog

## 1.15 client / server - 12-03-2022

* a permanent server keeps processes active
* keywords if/then/else/end 
* assignment via [X := n]
* expressions for common procedure calls
* use procedure-calls as arguments

## 1.14 the "why?" release - 15-01-2022

* translation of canned texts using CSV files
* ellipsis / path DSL
* discourse entities
* centering
* parsing multiple sentences on a line 

## 1.13 asynchronous execution - 26-03-2021

* asynchronous execution of mentalese code
* processing of input has been rewritten as a mentalese program  
* goals, actions
* explicit declaration of the rule base that accepts a new rule

## 1.12 morphological analyzer - 20-12-2020

* nli executable allows interactive session 
* insert / delete in MySQL
* relation tags {{ Variable }} to help generate compound sentences
* allow for left recursion in generation grammar
* local variables; let
* separated sort extraction from parsing 
* support for comparatives and superlatives
* morphological analyzer

## 1.11 modules - 13-09-2020

* namespaces; predicate prefixing
* modules
* sort hierarchy
* configuration from JSON to YAML
* configurable tokenizer expression
* data type `list` and list functions
* if_then_else

## 1.10 strong negation - 01-06-2020

* strong negation: the negation operator for relations
* the predicate `exec` to execute a shell command
* changed the order child sense order evaluation from post-order to pre-order

## 1.9 generalized quantifiers - 02-05-2020

* any relationset can now be used as a quantifier
* back_reference() definite_reference()

## 1.8 do/find - 21-04-2020

* do / find for different kinds of iterating over entities
* not()
* multiple entities per anaphora queue position
* allow string constants in grammar rewrite rules
* merged lexicon into grammar

## 1.7 Anaphora resolution - 15-02-2020

* Handles pronouns and other entity references
* Extract all parse trees

## 1.6 Database linking and case-insensitive proper nouns - 07-02-2020

* Support for queries that span multiple databases (with different entity ids)
* Support for case-insensitive proper nouns

## 1.5: Providing support for new DBpedia queries - 21-01-2019

* Only real quantifiers like ALL are handled with quantification; numbers are not
* syntactic relations are modelled after Stanford Parser Universal Dependencies
* introduce root()
* Start logging anonymous user interactions to get a feel of what types of questions the users of dbpedia test app pose en then support these types of questions
* checking entity types in predicate arguments for better name resolution
* Support for non-ASCII letters

Added DBpedia demo support for:

* "Who is X?"
* "When did X die?"
* the father of X
* "What is the capital of X?"
* "How many countries have population above 130000000?"

## 1.4: Interactive with dialog context - 15-12-2018

* When a question is about Lord Byron, and the database has two persons "Lord Byron" asks "which one"
* Dialog context to store the selected person by the user

## 1.3: Simple DB-pedia queries - 24-09-2017

* Domain - knowledge base mapping changes from 1:n to n:m
* Support for Sparql bases
* Intermediate results are logged
* Optimization phase using knowledge base statistics

## 1.2: Command-line app "nli" - 06-05-2017

* An executable application with "answer" and "suggest subcommands"
* Use an existing javascript autosuggest line editor (Tag-it!) and create an example web app
* Build an example application from a configuration file
* Rebuild of log as a proper dependency and with productions

## 1.1: Quantifier Scoping - 11-04-2017

* handle scoped questions
    * One sentence with ALL and 2 as quantifiers
    * One sentence where the right quantifer outscopes the left
* examples from relationships
* new: parse tree as new step

## 1: simple full-circle nli - 28-02-2017

* language: english
* question types: yes/no, who, which, how many
* second order predicates, aggregations
* proper nouns
* real database access (MySQL)
* a few simple questions
* simple natural language responses
* working example
