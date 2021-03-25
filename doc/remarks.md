## 2021-03-25

Up until now, `go:list_foreach` only fails if all children fail. I need it to fail under certain circumstances. So I'm thinking about a 

    go:cancel()

That breaks the loop and fails.

## 2021-03-23

The movements in the blocks world demo are now animated. The fact that the hand moved from A to B directly now became obvious, since it moved the blocks through other blocks. The hand now first moves up, before moving the block to another location.

I dropped the idea of using `go:action()` to communicate changes to the outside world. All interaction is now done via `wait_for()` relations.

The animation also reveals another problem: when the system builds a stack, it first decides on a location, then builds it. When building the first block, it may need to place the objects on top of it in some location. And it chooses the exact location where the stack should be. Later, the rest of the stack is still placed there. A solution could be to exclude this intended location from free space. 

## 2021-03-15

I must now create a more general communication protocol between the system (as a server) and the client (for example website).

Currently we have: asking something, getting a response text, getting response options, passing the option selection back to the server. Now we must add actions: the system tells the client what to do; the client notifies the system when its done.

Whereas before we could do most of the interaction in plain text; this is now becoming awkward. Its better to make the intention explicit.

    User: tell('How old is Lord Byron')
    System: select(uuid, 'Which one?', ['a', 'b', 'c'])
    User: pick(uuid, 'b')
    User: tell('Pick up a block')
    System: move_block(uuid, `block:red`, 100, 100, 0)
    User: done(uuid)

This would work, but it also requires to conversion from commands to actions on the side of the system. Can we do without? Can we make create relations that are directly executable?

    User: assert(goal(answer('How old is Lord Byron'))
    System: user_select(uuid, 'Which one?', ['a', 'b', 'c'], Selection)
    User: assert(user_select(uuid, 'Which one?', ['a', 'b', 'c'], 'b'))
    User: assert(goal(answer('Pick up a block'))
    System: user_move_block(uuid, `block:red`, 100, 100, 0)
    User: assert(user_move_block(uuid, `block:red`, 100, 100, 0))

(The relations the system sends are wrapped in `wait_for` relations).

This way, the system treats the user as just another knowledge source.

To create just two relation levels:

    User: go:tell('How old is Lord Byron')
    System: go:user_select(uuid, 'Which one?', ['a', 'b', 'c'], Selection)
    User: go:assert(go:user_select(uuid, 'Which one?', ['a', 'b', 'c'], 'b'))
    User: go:tell('Pick up a block')
    System: dom:move_object(uuid, `block:red`, 100, 100, 0)
    User: go:assert(dom:move_object(uuid, `block:red`, 100, 100, 0))

## 2021-03-06

The async rewrite is complete!

## 2021-03-05

All predicates that can be asserted / retracted are now indexed individually for facts and rules.

I created a write.yml file for rule bases. Now you can specify which rule base will receive the new rule, explicitly.

===

I will use `wait_for` instead of `ask` because it's a bit more general.

## 2021-03-01

Still working on making everything asynchronous. It's a lot of work, but its doable.

I am tackling some other issues in the process. Such as binding. I was never happy with the fact that the responsibility of binding was distributed over the code base. I am able to centralize it (in the Process and ProcessRunner classes). I just rewrote the binding code for rules as well, and it has become rediculously simple. As it should be.

Also, I am working on getting rid of scoped variable bindings.

Even though all relations are now executed asynchronously, this has not affected exection time. In fact, in the end it will be faster than before.

But today I want to talk about clarification questions. These will also change drastically. When the user asks "How old was Lord Byron" and the system needs to ask "Which one, A, B, or C?" the system would need to ask the user and restart the complete question from scratch. This changes. The system will halt until the user has answered the question, but when he/she does, the system will be able to continue where it left off.

I have thought of several approaches, but this is the one I'm currently excited about:

When the system needs to ask the user a question, it will just run

    go:ask(
        go:which_one(['George', 'Jack', 'Bob'], SelectionIndex)
    )

The `go:ask()` executes its child. When it succeeds, halt is successful. If it fails, it will not fail, but it will pause the process. Pausing the process can happen by sending a "processing instruction", which I already use for `go:break()`.

Next time the process is started, it will just run the `ask()` relation again, until the answer is made by the user.

The UI will ask the system for open `go:which_one()` relations (the last one on the process stack) and create facts for them:

    go:assert(go:which_one(['George', 'Jack', 'Bob'], 2))

Next time the process asks `go:which_one(...)` it will match this fact. This way, the answer will remain available to the system, and other knowledge sources could be able to answer the question as well.

## 2021-01-23

Save all goals and give them a status (i.e. "complete"), or remove them when they are finished?

## 2021-01-21

Starting one of the biggest rewrites of this program. Making everything asynchronous.

First problem: not all relations must be executed asynchonously. If I want to fetch all active goals at the start of a Run(), this can't wait. I need them now.

Relations that are executed in the scope of a goal are asynchronous; relations that are executed otherwise are immediate, synchronous.

## 2021-01-17

In order to save the state of a plan in execution, I want to create a custom call stack. Basically such a stack looks like this:

goal
|
goal
|
goal

The reason there is a call stack (or: sub goal stack) is that _some_ functions, like `list_foreach` perform a series of functions as part of their execution.

The advantage of creating a custom call-stack is that you can serialize and store is, and pick it up at some later time. And this is of course what I want. I want the system to be aware of its own state, and I want to be able to pause a process at any time, and I want to run the application for a limited amount of time, and be able to continue where it left off. If I can create this, I can use it for other processes too. Processes with steps that require user interaction.

This is one of those times where you'd wish that you had taken this decision earlier, because now you need to rewrite large parts of the code, and you had the idea in the back of your mind all the time. (Then again, there are _many_ ideas in the back of my head ;)

But things are not that simple. For one, the bindings need to be preserved.

    goal A / bindings 1
    |
    goal B / bindings 2
    |
    goal C / bindings 3

And then there's the thing that a goal consists of a relation set. Only one relation of this set can be executed at once, so we need to keep track of this

    relation * relation / bindings1
    |
    * relation relation / bindings2
    |
    relation relation * relation / bindings 3

The `*` denotes the current relation that is being executed. The bindings are now bound to the relation that is currently being executed.

So far, nothing special. But what exactly should happen when a function like `list_foreach` is executed? This function binds a variable to all values in a list, and calls a `scope` relation set, with this binding. How would that work? It can't call each `scope` immediately, like before. Now it must place `scope` on the call stack and leave its execution to the system (to its "fetch-decode-execution" cycle). But, hey!, that means that once `scope` is finished, `list_foreach` should be called again, to tell it that `scope` has finished, and that it should continue where it left off.

Yes, this is pretty complicated, but I see no other way. The relation that is being executed must get a state, and this state will be stored in the call stack.

    relation * relation / bindings1 / state 1
    |
    * relation relation / bindings2 / state 2
    |
    relation relation * relation / bindings 3 / state 3

As an example, let's go through the `list_foreach` and see what this means. In this example the values of the list are all added up. (S is a local, rewrtitable, variable).

First call:

    * go:list_foreach(List, ElementVar, Scope) / [S: 0, { List: [15, 27, 31], ListVar: X, Scope: go:add(S, X, S)}] / {}

`list_foreach` places the scope on the stack:

    * go:list_foreach(List, ElementVar, Scope) / [S: 0, { List: [15, 27, 31], ListVar: X, Scope: go:add(S, X, S)}] / { Index: 0 }
    |
    * go:add(S, X, S) / [{S: 0, X: 15}] / {}

Note that `list_foreach` has changed its state. Now it says: `{ Index: 0 }`

`add` quickly returns and the system pops it off the stack. `S` has become 15. The system still sees `list_foreach` as the current relation, and calls it again, but this time with a state. The result binding of `scope` is passed to `list_foreach` as well. 

    * go:list_foreach(List, ElementVar, Scope) / [S: 0, { List: [15, 27, 31], ListVar: X, Scope: go:add(S, X, S)}] / { Index: 0 }; child = [{S: 15, ...}]

`list_foreach` places the next scope on the stack:

    * go:list_foreach(List, ElementVar, Scope) / [S: 15, { List: [15, 27, 31], ListVar: X, Scope: go:add(S, X, S)}] / { Index: 1 }
    |
    * go:add(S, X, S) / [{S: 15, X: 15}] / {}

and when `add` returns, `list_foreach` is called for the third time and places the next scope on the stack:

    * go:list_foreach(List, ElementVar, Scope) / [S: 15, { List: [15, 27, 31], ListVar: X, Scope: go:add(S, X, S)}] / { Index: 2 }
    |
    * go:add(S, X, S) / [{S: 52, X: 15}] / {}

when `list_foreach` is called for the fourth time, is is done spawning child relations, and returns with new binding. This binding is passed to the next relation on the same stack frame. The state is reset to `{}`. 

## 2021-01-09

I noticed that the original SHRDLU demo used animations when moving objects around in the scene. And since this is an important part of the demo, I need to duplicate this. 

I noticed that SHRDLU only says "OK" when it is done moving all objects. This means, probably, that it moves objects in the process of answering the question. I could have NLI-GO do that too, but this is not what I want. I don't want the system to be unresponsive while it is executing the task. As for a simple reason why: you would not even be able to tell the system to stop.

This means that I will separate the answerer from the goal execution engine. They will need to be activated separately. At the same time, both will be modules of the same system, and share some of the resources. This is how I image this to work:

The answerer may assert `goal` relations:

    go:assert(go:goal(GoalSet, incomplete, `id`, none))

This means: spawn a goal (specified by GoalSet, a relation set), which is as yet unfinished ("incomplete"), has an id (`id`) and no parent id (`none`).

I will create a system command: 

    run(maxTime)

that tries to advance the incomplete goals by executing their goal sets. The `maxTime` argument may tell the system how long it should run maximally. For now it is not important. But it can be used to run this process each minute in a cron job or something.

A goal set, or any subgoal it invokes, may create actions. An action is a signal to external modules that something needs to be done.

    go:assert(go:action(ActionSet, <actionType>, incomplete, `id`))

The goal executor will halt until the action is completed. The action type can be used by external modules to listen to specific actions.

Finally, I will add another system command to notify the system that an action has been performed, or has failed: 

    updateActionState(actionId, newState) 

The web application then serves as the not-so-physical representation of the robot, that performs actions by creating animations. When the user enters a line of text, the web app will call `answer`. When this is done, it will call `run` to advance the system, and then `query` to ask for the pending actions. It will perform these actions by creating animations, and then call `updateActionState`, followed by another `run`, as long as the action query gives more results.

Isn't this a awful lot of work just to add some animation? Yes, if it were just for the animations, this much work would be unwarranted. But this is also a preparation for the agent that is necessary to perform much more complex tasks, and to create the goal hierarchy that is needed to answer the following 15 questions of the SHRDLU demo interaction. 

## 2021-01-05

I created a stack trace! It is output when the answerer finds no results. It is an automation of what I have been doing all along to fix a program. 

An object called CallStack keeps track of the functions being called recursively. When a function gives no results, a copy of the stack at that time is frozen.

Here's an example stack trace that is created when I try to fit a block into the box that doesn't fit:

    Stack trace
    10. go_greater_than_equals(W$6, Width$25)
        {ColIndex$5:4, Index$6:4, Line$6:1000, VerLines$2:[600, 640, 800, 840, 1000], W$6:0, Width$25:200, X1$18:1000}&{A1$1:0, A1$2:-1, B1$1:200, B1$2:-1, StartY$1:200, StartY$2:600, Success$11:true, Success$12:true, Success$14:true, Success$15:true, Success$3:true, Success$4:true, Success$7:true, Success$8:true}
    
    9. go_list_foreach(VerLines$2, Index$6, Line$6, go_subtract(Line$6, X1$18, W$6) go_greater_than_equals(W$6, Width$25) go_subtract(Index$6, ColIndex$5, ColSpan$2) go_break())
       {ColIndex$5:4, VerLines$2:[600, 640, 800, 840, 1000], Width$25:200, X1$18:1000}&{A1$1:0, A1$2:-1, B1$1:200, B1$2:-1, StartY$1:200, StartY$2:600, Success$11:true, Success$12:true, Success$14:true, Success$15:true, Success$3:true, Success$4:true, Success$7:true, Success$8:true}
    
    8. dom_find_span(Width$25, VerLines$2, ColIndex$5, ColSpan$2)
       {ColIndex$5:4, HorLines$2:[600, 640, 840, 940, 1000], Length$25:300, VerLines$2:[600, 640, 800, 840, 1000], Width$25:200, X$21:1000}&{A1$1:0, A1$2:-1, B1$1:200, B1$2:-1, StartY$1:200, StartY$2:600, Success$11:true, Success$12:true, Success$14:true, Success$15:true, Success$3:true, Success$4:true, Success$7:true, Success$8:true}
    
    7. go_list_foreach(VerLines$2, ColIndex$5, X$21, go_list_get(HorLines$2, 0, StartY$2) dom_find_span(Width$25, VerLines$2, ColIndex$5, ColSpan$2) go_add(ColIndex$5, ColSpan$2, V1$6) go_subtract(V1$6, 1, ColEnd$2) go_list_foreach(HorLines$2, LineIndex$2, Y2$6, go_greater_than(LineIndex$2, 0) go_subtract(LineIndex$2, 1, RowIndex$2) go_if_then_else(dom_span_free(ColIndex$5, ColEnd$2, RowIndex$2, fixed), go_subtract(Y2$6, StartY$2, SpanLength$2) go_greater_than_equals(SpanLength$2, Length$25) go_let(A1$2, X$21) go_let(B1$2, StartY$2) go_break(), go_let(StartY$2, Y2$6))) go_not_equals(A1$2, -1) go_break())
       {E5:`block:big-red`, HorLines$2:[600, 640, 840, 940, 1000], Length$25:300, VerLines$2:[600, 640, 800, 840, 1000], Width$25:200}&{A1$1:0, A1$2:-1, B1$1:200, B1$2:-1, StartY$1:200, StartY$2:600, Success$11:true, Success$12:true, Success$14:true, Success$15:true, Success$3:true, Success$4:true, Success$7:true, Success$8:true}
    
    6. dom_do_find_free_position(E5, fixed, HorLines$2, VerLines$2, X$18, Y$17)
       {BoundX1$2:600, BoundX2$2:1000, BoundY1$2:600, BoundY2$2:1000, E5:`block:big-red`, E6:`box:box`, HorLines$2:[600, 640, 840, 940, 1000], Objects$2:[`pyramid:blue`, `block:blue`], VerLines$2:[600, 640, 800, 840, 1000]}&{A1$1:0, B1$1:200, StartY$1:200, Success$3:true}
    
    5. dom_do_find_free_space(E6, E5, X$18, Y$17)
       {E5:`block:big-red`, E6:`box:box`, Z$26:0}&{A1$1:0, B1$1:200, StartY$1:200, Success$3:true}
    
    4. dom_do_put_in(E5, E6)
       {E5:`block:big-red`, E6:`box:box`}&{A1$1:0, B1$1:200, StartY$1:200, Success$3:true}
    
    3. dom_do_put_in_smart(S, E5, E6)
       {E5:`block:big-red`, E6:`box:box`}
    
    2. go_quant_foreach(go_quant(some, E6, go_definite_reference(E6, dom_box(E6))), dom_do_put_in_smart(S, E5, E6))
       {E5:`block:big-red`}
    
    1. go_quant_foreach(go_quant(some, E5, go_definite_reference(E5, dom_red(E5) dom_block(E5))), go_quant_foreach(go_quant(some, E6, go_definite_reference(E6, dom_box(E6))), dom_do_put_in_smart(S, E5, E6)))
       {}

This is very useful! It makes debugging a lot easier, I expect. What you also see here is that some constructs are very heavy and this shows up in the call stack. So this is an indication that I should keep them simple. 

## 2021-01-02

Very happy with the fact that it was possible to create the typical oblique projection of SHRDLU with THREE.JS. The demo is coming along fine. I am now working on the hand, which a largely ignored up to now, but which has such an important role in the demo. It hit on me a few days ago that the hand actually moves smoothly in the SHRDLU demo. In the database it just switches to a new position, so this is an important difference that I will need to find a solution for.

The hand picks up a block in the center, and when the hand moves up, the block should follow. The block is not at the same position as the hand of course, it must be translated.

## 2021-01-03

Since I now have an interactive blocks world demo, I can interact with it immediately and try different sentences. They all break down in terrible ways! I have seen cubes floating in space, cubes taking up the same space, execution times of 5 seconds... Wow. This thing is not robust by any means.

However, I fixed the first problems I encountered. And I made an interesting innovation: I added an extra rules layer of "physics". This layer of rules is responsible for ensuring that all relations (contain, support, cleartop) stay intact, whatever you do. All objects that move now go trough a single function `phys_move_object`. This function has as input just the position of the object. It breaks up any existing relations the object might have, and then rebuilds them, just by looking at the position of the object. I was afraid that this would be very expensive, and it is not cheap, but the enormous advantage is that you can now move objects around without thinking about the relations _at all_. They have been fixed in a special layer once and for all.
