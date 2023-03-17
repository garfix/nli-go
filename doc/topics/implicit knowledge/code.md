All relations that are specified as tags in syntax rules, are added to the dialog context.

For example:

    { rule: relation(E1, E2) -> 'husband',   sense: dom:has_husband(E1, E2),   tag: go:sort(E1, person) go:sort(E2, person) }

The tags of this rule are mainly used to resolve the sense of the sentence. 

They are also stored as facts in the dialog context database, that is available only within the active dialog. Because it's not possible to store variables in a database, the variables are turned into atoms.

    E1 -> e$23
    E2 -> e$24

    go:sort(e$23, person) 
    go:sort(e$24, person)

When an application rule needs this information to solve a problem, it can acces the information by turning the local variable into an atom and using this atom to access the data:

    go:sort(go:atom(A, Atom), Sort)

This (only) works, because the system replaces the formal variable of a function with the variable of the calling function. This way, the original variables of the sentence are available deep into the system.


