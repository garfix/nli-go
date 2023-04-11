# Grammar predicates

Predicates that are primarily used as part of the sense of grammar rules.

## Do

Finds the entities specified by a `quant`, assign each of them in turn to a variable and execute `Scope`.

Does not continue to find entities after the quantifier has succeeded.
Fails if the number of entities that pass `Scope` is **less than** the same as specified by the quantifier of `Quant`.

    go:do(Quant ..., Scope)

* `Quant`: a quant
* `Scope`: a relation set

Check [quantification](quantification.md) for more information.

## Check

Find all entities specified by `Quant`.

Fails if the number of entities that pass `Scope` is not the same as specified by the quantifier of `Quant`.

    go:check(Quant, Scope)

* `Quant`: a quant
* `Scope`: a relation set

Check [quantification](quantification.md) for more information.

## Quant to list

Creates a new quant, based on an existing quant, but extended with an order function. If the original quant already had an order, it will be replaced.

    go:quant_ordered_list(Quant, &OrderFunction, List)

* `Quant`: a `quant` relation
* `OrderFunction`: a reference to a rule that functions as an order function
* `List`: a variable (to contain a list)

If the quant is complex and contains sub-quants; then these will be ordered by the `OrderFunction` as well

    Example:

The order relation takes two entities and returns a negative number, 0, or a positive number. negative when E1 goes before E2, 0 when E1 has the same order position as E2, positive when E1 goes after E2.

    by_easiness(E1, E2, R) :- if_then_else( cleartop(E1), unify(R, 1), unify(R, 0) );

    go:quant_ordered_list(Quant, &by_easyness, List)

## Quant match

Matches `Quant` against the entities in `List`.

    go:quant_match(Quant, List)

* `Quant`: a `quant` relation
* `List`: a list

## Back reference

The system will try to resolve E1 with the entities from the anaphora queue. It will check E1's type against the types of entities in the queue. It will also check if the value of the entity in the queue matches relation set `D`.

    go:back_reference(E1, D)

* `E1`: a variable
* `D`: a relation set

This allows you express "him", like this:

    go:back_reference(E1, gender(E1, male))

The recentless processed `person` entities will be processed, and the `gender()` check makes sure these persons are male.

## Definite reference

A definite reference checks not only in the anaphora queue, but in the databases as well.

If more than one entity matches, a remark is returned to the user: "I don't understand which one you mean"

    go:definite_reference(E1, D)

* `E1`: a variable
* `D`: a relation set

## Sortal back reference

A back reference to locate a concrete sort for an entity. For example: "one" in "Put a small one" should refer to "block" in the recent conversation (as stored in the anaphora queue).

    go:sortal_back_reference(E1)

* `E1`: a variable

This function goes though the anaphora queue to find the most recently used sort (for example: "block"); then finds the relation set that belongs to this sort (from the "Entity" field in `sort-properties.yml`) and executes this relation set (with its Id variable replaced by E1 in this example). The result of this relation set is the result of this function.

In this example `back_sortal_reference` will find the sort of E1, and resolve this into a concrete relation set (like `block(E1)`):

    { rule: noun(E1) -> 'one',                                             sense: go:back_sortal_reference(E1) }

"Find one" will thus be interpreted as "Find a block", if `one` refers to `block` in recents discourse.

More on this in [anaphora-resolution](../anaphora-resolution.md)

## Created canned response

Creates a canned response, based on a canned template, and one or more arguments, and places it in Output

    go:create_canned(Output, template, Argument)

* `Ouput`: a string
* `template`: an atom
* `Argument`: a string

## Translate

Translated a canned text to a new text in a given locale.

    go:translate(Source, Locale, Translation)

* `Source`: the source text
* `Locale`: the locale (i.e. en_US)
* `Translation`: the outpur translation

Use `go:locale(Locale)` to get the current locale.

