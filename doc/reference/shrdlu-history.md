# SHRDLU history

Winograd did not take the easy road. If SHRDLU was to answer questions about the when, the why, and the how of its previous actions, it needed to have a history.

It was not common for a database to have a history 50 years ago, and it is still not common (and this is an understatement) for a database to keep a history. This is one of the reasons why some of the questions SHRDLU is able to answer still evokes so much awe today. 

So Winograd has created a history for his program. It is part of the content of the database, not of the structure of the database management system.

What does this history look like? In this page I will collect everything Winograd has said about it.

## What the history isn't

In 6.2.4 of Natural Language Understanding Winograd mentions that a history _could_ be built by adding a timestamp to each fact in the database:

    (#ON :B1 :B2 S111)

This relation would say: B1 is on B2 at time 111. 

While this representation creates a history, Winograd shortly explaings that it is not easy to work with, and that he had let this idea go.

## Event list

In 7.5 we read: The event list contains event-objects. It keeps track of the larger goals like #PUTON and #STACKUP.

The time of events is measured by a clock which starts at 0 and is incremented by 1 every time any motion occurs.

Theorems call MEMORY and MEMOREND; MEMOREND causes an event to be created, combining the original goal statement with a name (E1, E2, ...) 

MEMOREND puts information on the property list of the event name -- the starting time, ending time, and reason for each event. The reason is the name of the event nearest up in the subgoal tree which is being remembered. The reason for goals called by the linguistic part of the system is a special symbol meaning "because you asked me to."

MEMORY is called at the beginning of a theorem to establish the start time and declare that theorem as the "reason" for the subgoals it calls.

    Clock
        - timestamp [0..>

    Event
        - goal (i.e. #PUTON or #STACKUP)
        - name (E1, E2, ...)
        - object (:B1, :B6, ...)

    Event properties
        - event name (E1, E2, ...)
        - starting time (0, 1, ...)
        - ending time (0, 1, ...)
        - reason (E1, E2, ...)

Example event (see 8.1.5)

    (#PICKUP E23 :B5)

Events are only created for top-level goals. So #PICKUP is a goal when the user asks specifically to pick up a block. When picking up and putting down blocks is part of building a stack, this does not cause events. The fact that a a block was picked up can be deduced from the fact that it was put somewhere, and the theorem TCTE-PICKUP actually looks at a number of different types of events (like #PUTON en #PUTIN) to find all the occasions on whcih an object was really picked up.

## The physical motions list

7.5 continues: A second kind of memory keeps track of the actual physical motions of objects, noting each time one is moved, and recording its name and the location it went to.

    Motion
        - time (0, 1, ...)
        - name (:B1, :B6, ...)
        - location (X, Y, Z)
    
