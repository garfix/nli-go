## Dialog context

## Context set

Sets a deictic center (a relation set) in the dialog context.

    go:context_set(DeicticCenter, Event, Modifier)

* DeicticCenter: a deictic center (the atom `time` or `place`)
* Event: a variable that represent an event that the modifier attaches to
* Modifier: a relation set that contains the variable `Event`

## Context extend

Extends an existing deictic center (a relation set) in the dialog context.

    go:context_extend(DeicticCenter, Event, Modifier)

* DeicticCenter a deictic center (the atom `time` or `place`)
* Event: a variable that represent an event that the modifier attaches to
* Modifier: a relation set that contains the variable `Event`

## Context clear

Removes a deictic center from the dialog context.

    go:context_clear(DeicticCenter)

* DeicticCenter a deictic center (the atom `time` or `place`)

## Context call

Executes a deictic center (a relation set) from the dialog context.

    go:context_call(DeicticCenter, Event)

* DeicticCenter a deictic center (the atom `time` or `place`)
* Event: a variable that represent an event that the modifier attaches to

## Dialog read bindings

Reads the entity ids of the discourse into the current request.

    go:dialog_write_bindings

## Dialog write bindings

Writes the entity ids of the current request back into the discourse.

    go:dialog_write_bindings
