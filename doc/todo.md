# TODO

## Releases

### 1: simple full-circle nli

* language: english
* question types: yes/no, who
* real database access (MySQL)
* a few simple questions
* simple raw database response
* working example

### 2: omschrijving


# Known issues

## Syntax

- Permanent goal: improve the grammar; extend it with new phrases, make it more precise. I think there's such a thing as an NLI-English grammar that exists of grammar rules commonly used when talking to a computer. It's a small subset of full English grammar, with an emphasis on questions.
- Must be able to write whword in place of whword(); but wait, maybe we need multiple variables as well?
- The syntax is not very useful. Both too little and too much extra tokens.
- Must have a proper grammar syntax => relations
- Earley parser
- is het misschien nodig om predicates en constants te namespacen? Eigenlijk is de predicate al een namespace

## Domains

The domain tests are not a goal in themselves, but only help to make up test cases.

- Block's world

## Parsing

- Cannot handle left-recursive rules. Change the parser or restrict the rules?

## Answering

- The way to find an answer to a question is insufficient.

## Generation

- Generate simple responses
