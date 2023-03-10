# Code

In the blocks world, episodic memory is implemented in two ways, according to Winograd's suggestion.

## Events

The application stores the actions it deems important in the database. These actions are stored as events, using a predication, an event id, and a parent event id. The parent event is the "cause" of the event, and events form a hierarchy.

This is the function that executes and persist an action:

    persist_action(EventId, ParentEventId, Action) :-
        increase_time()
        time(T1)
        go:assert(parent_event(EventId, ParentEventId))
        go:assert(start_time(EventId, T1))

        go:call(Action)

        time(T2)
        go:assert(end_time(EventId, T2))
    ;

`increase_time()` increments the internal timer by 1. The time is retrieved by `time(T1)`. For an event, its `parent_event` and its `start_time` and `end_time` are stored. `go:call` executes the action.

`persist_action` is called like this (a "pick up" action)

    do_pick_up(ParentEventId, E1) :-
        go:uuid(EventId, event)
        persist_action(EventId, ParentEventId,
            
            go:assert(pick_up(EventId, `:shrdlu`, E1))

            top_center(E1, X, Y, Z)
            phys_move_hand(X, Y, Z)
            phys_grasp(E1)
            phys_raise_hand(_)
        )
    ;

Note that a parent event id is passed to the function. It creates an event id and it stores (asserts) the predication of the event `pick_up(EventId, `:shrdlu`, E1)`

The database thus contains these relations:

    time(T)
    parent_event(EventId, ParentEventId)
    start_time(EventId, T1)
    end_time(EventId, T2)
    pick_up(EventId, Subj, Obj), put_down(EventId, Subj, Obj), put_on(EventId, Subj, Obj1, Obj2), etc

The `time` relation is dropped and added each time that time increases. The other relations are only added.

Note that past events have the same predication as current events. The event id allows it to be modified by time stamps.

"Why" questions are then handled by looking up the parent id of an event. "Why did you clear off that cube?" is solved by looking up the historical `clear_off` relation, followed by a look up of `parent_event`.

"How" questions are handled by enumerating the child events of a given event. "How did you do it?" is solved by resolving "it" to an event (it's id), followed by a look up of `parent_event` to find its child events.

"When" questions are handled by finding the top event of an event. "When did you pick it up?" first locates the pick up event, then traverses its event tree to the top. Then describes that top level event.

## Locations

To handle the question "What did the red cube support before you started to clean it off?" another kind of memory is needed. The application stores the past locations of all objects. Each period in which an object is located at one place is a location. An object initially has a location starting at time 0 and ending at the end of times (here: 1000000).

    location(start1, `block:small-red`, 100, 100, 0)
    start_time(start1, 0) 
    end_time(start1, 1000000)

The function that physically moves an object modifies the latest event, and adds a new one:

    phys_move_object(E1, X, Y, Z) :-
        ...

    /* modify latest location event of E1 */
    location(PrevEvent, E1, _, _, _)
    end_time(PrevEvent, 1000000)
    go:retract(end_time(PrevEvent, 1000000))
    go:assert(end_time(PrevEvent, T1))


    /* add a new location event */
    go:uuid(EventId2, event)
    go:assert(location(EventId2, E1, X, Y, Z))
    go:assert(start_time(EventId2, T1))
    go:assert(end_time(EventId2, 1000000))

After two relocations, the locations times may look something like this:

    0 - 20
    20 - 25
    25 - 28
    28 - 1000000

Past locations are stored like events, so that they have `start_time` and an `end_time` just like the action events. By doing this, these two kinds of events can be compared with predicates like `dom:before(EventId1, EventId2)`.

"What did the red cube support before you started to clean it off?" is then handled by finding the events for `support` and `clear_off`, and comparing them with `before`. There is a special `support` procedure that deduces the support-relation from the locations of objects.

    support(P1, A, B) :-
        location(P1, A, X, Y, Z)
        check_on_top_of(B, A)
    ;
