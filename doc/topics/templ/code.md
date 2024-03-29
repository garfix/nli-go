# Code

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

To handle the question "What did the red cube support before you started to clean it off?" another kind of memory is needed. The application stores the past locations of all objects. The function that physically moves objects records their location as location events:

    phys_move_object(E1, X, Y, Z) :-
        ...

        go:uuid(EventId2, event)
        go:assert(location(EventId2, E1, X, Y, Z))
        go:assert(start_time(EventId2, T1))
        go:assert(end_time(EventId2, T1))

        ...
    ;

Past locations are stored like events, so that they have `start_time` and an `end_time` just like the action events. By doing this, these two kinds of events can be compared with predicates like `dom:before(EventId1, EventId2)`.

"What did the red cube support before you started to clean it off?" is then handled by finding the events for `support` and `clear_off`, and comparing them with `before`. There is a special `support` procedure that deduces the support-relation from the locations of objects.

    support(P1, A, B) :-
        location(P1, A, X, Y, Z)
        check_on_top_of(B, A)
    ;

To have the initial locations available as well, these are added to the database:

    location(start, `block:small-red`, 100, 100, 0)

Here `start` is a special event that starts and ends at time 0.

