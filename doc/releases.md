## 1.3: Simple DB-pedia queries


    name(A, F, firstName) name(A, L, lastName) join(N, ' ', F, L) => birth_name(A, N);
    name(A, F, firstName) name(A, S, secondName) name(A, L, lastName) join(N, ' ', F, S, L) => birth_name(A, N);
    name(A, N, fullName) => birth_name(A, N);

    name(A, F, firstName) name(A, L, lastName) join(N, ' ', F, L) => name(A, N);
    name(A, F, firstName) name(A, S, secondName) name(A, L, lastName) join(N, ' ', F, S, L) => name(A, N);
    name(A, N, fullName) => name(A, N);

    beter,
    - want de db bepaalt, alleen, want het formaat is waarop de naam wordt opgeslagen, het is geen DS ding
    maar
    - kb handles relation moet worden aangepast => ok
    - kb handles 1 SINGLE RELATION moet worden aangepast => ...
    - lost het probleem niet op: name() name() 2 namen


## 1.2: Command-line app "nli"

* An executable application with "answer" and "suggest subcommands"
* Use an existing javascript autosuggest line editor (Tag-it!) and create an example web app
* Build an example application from a configuration file
* Rebuild of log as a proper dependency and with productions

## 1.1: Quantifier Scoping

* handle scoped questions
    * One sentence with ALL and 2 as quantifiers
    * One sentence where the right quantifer outscopes the left
* examples from relationships
* new: parse tree as new step

## 1: simple full-circle nli (2017-02-28)

* language: english
* question types: yes/no, who, which, how many
* second order predicates, aggregations
* proper nouns
* real database access (MySQL)
* a few simple questions
* simple natural language responses
* working example
