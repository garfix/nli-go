## 2021-04-24

I added memory events in a number of strategic places. I displaced and renamed many of the blocks world files, to emphasize the role of the events.

## 2021-04-22

Not all relations cause memory events. Only the ones that are relevant for later recollection. These events are added explicitly.

This probably does not correspond to the way long term memory works in humans.

Myself, I would not have created a special event for "let go", but I will do it for the demo.

## 2021-04-20

Another thing about

    When did you pick it up?

"When", in this case, refers to the time frame "before you put the green one on the little cube". If you just look for the occurrences where the green pyramid was picked up, there are two. The question just referes to the first of these occurrences.

So I need to introduce a time frame for the dialog context, or "did you pick it up" itself must refer to the event "Had you touched" in the previous sentence. 

## 2021-04-19

Looking ahead to the coming interactions, it's a good idea to visualize the should goal structure that allows the system to answer the coming questions 25 - 30:

    G1: stack_up([`block:big-red`, `block:green-back`, `block:small-red`])
        Event1: put_on(`block:big-red`, `table:table`)
        Event2: let_go(`block:big-red`)
        Event3: put_on(`block:green-back`, `block:big-red`)
        Event4: let_go(`block:green-back`)
        Event5: put_on(`block:small-red`, `block:green-back`)
            Event6: clean_off(`block:small-red`)
                Event7: get_rid_of(`pyramid:green`)
                    Event8:  pick_up(`pyramid:green`)
        Event6: let_go(`block:small-red`)

To answer the questions:

    25: why? -> Event7
    26: why did you do that? -> Event6
    27: why did you clear of that block? -> Event5
    28: why did you do that? -> G1
    29: why did you do that? -> null ("Becaause you asked me to")
    30: how did you do it? -> Event1, Event2, Event3, Event4, Event5, Event6

I must now create a relation between goal-event and event-event.

## 2021-04-18

Some remarks about this interaction:

    Is there a large block behind a pyramid?
    Yes, three of them: a large red one, a large green cube and the blue one

More specifically: "a large red one". 

There are 5 blocks, 3 of which are cubes. The `type` of all blocks is `block`, but the `shape` of a cube is `cube`, and not (also) `block`.

Now, to determine if we need to add "large" to distinguish two blocks, we should look at the _type_ of the object, not the _shape_. There is only one red object whose _shape_ is block (and not cube). But there are two objects whose _type_ is block.

On the other hand, in

    Is at least one of them narrower than the one which I told you to pick up?
    Yes, the red cube

it is not necessary to say "small red cube", since there is only one red cube.

The point is that there's a difference between the presentation shape "cube" and all allowable shapes (a cube is also a block). 

To determine the accurate description of an object, you must first determine the presentation shape of the object (the `characteristic_shape` I called it) and then check how many objects have this shape, by not checking the characteristic shape of these objects, but by checking the regular shape.

If the characteristic shape of an object is `block`, you need to check how many blocks there are; and these should include cubes as well (even though, if you would ask to describe them, they would be called `cube`, not `block`).

So there's

- type: the database shape
- shape: any of the shapes and subshapes that an object may be said to have (an object can be both a cube and a block)
- characteristic shape: the most characteristic shape of the object (`cube` for a cubed block)

---

Interaction 24:

    H: When did you pick it up?
    C: While I was stacking up the red cube, a large red block and a large green cube.

Winograd's comment: "A time is described by naming the main goal which was being worked on."

Winograd also writes on naming events in 8.3.3: "Naming events is relatively straightforward". This is a comfort, but it also comes at a surprise: "straightforward"?

Let's do an initial analysis:

- The question asks for an event. The event is described by the clause "you pick it up"; where "it" refers to the green pyramid from the last answer.
- The answer describes the event by creating a clause from the predicate and the arguments

So what's straightforward about it? Winograd:

    With each event type (such as #PUTON or #STACKUP) we associate a small program which generates an English name for the event and combines it properly with the names of the objects involved.

Ok, this makes sense. In stead of creating English names we will create a relational representation, that will be turned into English by the generator.

This is enough to get started.

---

The "pick up" mentioned here is not an explicit command, like "pick up a big red block". This "pick up" event takes place every time an object is grasped.

## 2021-04-15

The answer to #23 should still be

    Yes, the green one

Not: "Yes, the green pyramid". This is because the question already includes the type of the block: "Had you touched any pyramid".

I also need this for the previous interaction "Is there a large block behind a pyramid??"

---

It looks like the active sentence is a knowledge source itself. "Is E1 a block?" This fact can only be found in the input sentence itself.

---

Another take at this, is to simply view the `sort` of a variable; if I need to know if `E5` is a block, I can check its sort. It's stored with the variable.

I can introduce a new function, like this:

    go:get_sort(E5, Sort)

Or, allowing inheritance:

    go:isa(E5, block)

This works because the sort of an entity is deduced from the relations used in the quants of the input sentence.

---

Unfortunately, the type of the entity is not only determined in the relationization process. The identifiers of the entities also contain their type. This leads to an overfitting of the word "one".

The meaning of the question "Had you touched any pyramid" should more explicitly contain the requested sort:

"Describe the entity, of which it is given that it is a pyramid (no need to repeat that)", such that ...

    describe(E1, Sort, Description)

## 2021-04-12

I added "slots" to the process entity. These are a sort of global variables needed to store some objects that need to be accessed from various places: like `locale` and the `sense` of a sentence. These slots can both be accessed from the outside (the Rule code) and the inside (the Go code).

I managed to complete #23 this way.

More about James Allen's "Natural Language Understanding": on page 367 he has a very interesting description of the difference between reflexive and irreflexive pronouns.

## 2021-04-11

I now have a case where I want to access information about the sentence I am currently processing. And I find that I don't have access to it any more.

By turning sentence processing in a more general relation processing, I lost the sentence as a first class citizen.

So I need to bring this back. I can't store sentence information in the nli database, because it is shared by other sentences (in the future nli-go will work on several sentences at once). Same thing for the dialog context.

What I can do is store sentence information in the `process` object.

---

Some more thought: a process can have multiple relational representations of the same sentence. So storing the sense with the process is not possible.

But what I can do is store the sense in a mutable variable, just like I did with the locale. This way, it won't be necessary to add another field to the process either.

Or, like I thought I did. Because I was storing the locale in a the nli-go memory base. I should change that as well. 

"Mutable variables" are then really turning into global (process scoped) variables. And I can use these variables "from the inside" as well. This is very convenient.



## 2021-04-10

This anaphora queue makes you wonder: what _is_ the best way to represent anaphoric relations? 

According to Winograd, SHRDLU looks inspects the current and previous sentences. This is probably very time-consuming. There is no separate structure to hold anaphoric information.

I just looked at James Allen's "Natural Language Understanding". He has a complete chapter on this. I should get into this a lot deeper. Later.

## 2021-04-09

If I just store the current variable with the entity in the queue, I can use this to determine if two entities share a relation in the current sentence. When the sentence is finished, I will remove the variables from the queue. This way I can determine from the presence of variables that an entity is mentioned in the present sentence or not.

## 2021-04-08

An error in the anphora queue generation was one of the problems causing int. 23 to fail. When I corrected it, int. 21 suddenly failed:

    Put the littlest pyramid on top of it

The system's interpretation now is: "Put the littlest pyramid on top of <the littlest pyramid>". This is because "it" refers to the last entry added to the anaphora queue, and this is now (as it should be) "the littlest pyramid".

Here is a picture showing the system has put the little pyramid on top of itself (it's no longer there, of course):

![Initial blocks world](archive/blocksworld6.png)

The reason why this interpretation is wrong, is not so simple. The interpretation corresponds to the word "itself", but why is that?

LFG would say that the antecedent of a reflexive pronouns like "itself" can only occur in the Minimal Complete Nucleus containing the pronoun.
(Syntax and Semantics, Dalrymple, p. 280)

This means, I think, that when the antecedent and the pronoun are part of the same relation, the pronoun is reflexive.

Or that a non-reflexive pronoun must not refer to the subject of the same relation.

Currently I am not storing this type of information in the anaphora queue, so that must change.

## 2021-04-05

SHRDLU uses a Time Semantic Structure to describe the interval that applies to a certain clause.

    TSS
        - tense (present, past)
        - progressive (yes, no)
        - start time (0, 1, ...; nil)
        - end time (0, 1, ...; nil)

So this structure is attached to each clause. I prefer to work with individual relations in stead of "object structures".

    dom:tense(P1, past)
    dom:progressive(P1)
    dom:start_time(P1, 23)

Current approach:

- Evaluate "you put the green one on the little cube" into a start- and end time
- Use the start time to create `dom:end_time(P1, End)`
- Create a normal representation for "Had you touched any pyramid" into the main event, and extend it with `dom:tense(P1, past)`
- Evaluate the main event with the function `dom:eval_in_time($event_description)`

---

I can store an past action in memory like this

    do_pick_up(E1) :-
        time(T1)
        // implementation
        time(T2)
        go:uuid(Id)
        go:assert(pick_up(Id, `:shrdlu`, E1))
        go:assert(start_time(Id, T1))
        go:assert(end_time(Id, T2));

`before` can be implemented as

    before(P1, P2) :-
        end_time(P1, End)
        start_time(P2, Start)
        go:less_than(End, Start);

New approach:

- Look up "you put the green one on the little cube" (as `put_on(P1, E1, E2)`) in the system's event assertions, this gives event ID values for `P1`
- Look up "Had you touched any pyramid" (as `touch(P1, E1, E2)`) in the assertions (gives no match) and use deduction to come from `touch` to `pick_up`
- Given the values of P1 and P2, use `before(P1, P2)`

## 2021-04-03

Interaction 23, shall we?

    H: Had you touched any pyramid before you put the green one on the little cube?
    C: Yes, the green one

"The system deduces from its memory what event is being referred to, and can relate other events to it in time. The analysis includes complex tenses of verb groups, and there are special facilities for keeping track of the various forms of irregular verbs like "have""

Interaction 23 to 33 form a new level of complexity. Daunting, but utterly fascinating. Only at the highest level have I an idea of how to deal with these. It will require a lot of thinking, trying, failing, and refactoring. But I have made all the preparations I could. So let's just get going.

First remarks on the interaction:

- The sentence is in past tense. It is about an action that was performed in some past period of time.
- It can be paraphrases as "Did x happen, before time T (described by y)"
- "you put the green one on the little cube" describes an action, so past actions must be stored in a form so that they can be recalled like this
- "Had you touched any pyramid" also describes an action in the past
- Another paraphrase: x(T1) y(T2) T1 < T2
- "The green one" refers to "any pyramid" earlier in the same sentence.

Earlier on I created the document "shrdlu-history.md" that collects information about SHRDLU's history.

Coming questions about moving things in the past:

    How many objects did you touch while you were doing it?
    What did the red cube support before you started to clean it off?
    Have you picked up superblock since we began?



## 2021-04-02

I added the function `go:translate(Source, Locale, Translation)` to translated canned texts to a chosen locale.

`go:locale(Locale)` can be used to find the current locale.

The translations can be stored in a simple CSV file in the grammar.

## 2021-04-01

The rule bases weren't stored in the session. They are now.

I placed the learnable 'own' rules in a separate rules base (blocks world), to avoid storing all rules in the session. Only a few rules can be learned, so only these will need to be stored in the session.

## 2021-03-31

Haha, I made it!

![Initial blocks world](archive/blocksworld5.png)

## 2021-03-30

Next interaction (22):

    H: How many things are not on top of green cubes?
    C: I'm not sure what you mean by "on top of" in the phrase "on top of green cubes",
       do you mean:
       1 - directly on the surface
       2 - anywhere on top of
    H: 2
    C: three of them

    "When there are several interpretations and none is clearly better, the system has to ask the user for more information. The paraphrases were included with the directory definition of "on", and the request for more information occurs automatically when the heuristics for resolving ambiguities aren't sufficient. A request like this can also occur when the meaning of a pronoun such as "it" or "they" is unclear."

The preposition "on" is ambiguous. Winograd has put this ambiguitity in the dictionary. I will add rules for it. The `wait_for` / `user_select` construction allows me to do this.

I will change the system's clarification to a simple "Did you mean?"

"Things" must mean "block" or "pyramid", but not "table" and "box". 

---

I got an answer, from the system. But the answer was "two of them", not "three of them". What happened? 

Apparently, in interaction 20 SHRDLU doesn't remove the long pyramid from the green block, but it _puts the small block next to it_, on the same block. The book doesn't have a picture of this situation, but it makes sense. This is why the long pyramid is placed in a corner of the block.

Also, I got myself into a left recursion with

    anywhere_on(A, B) :- anywhere_on(C, A) anywhere_on(C, B);

I fixed it.

## 2021-03-29

go:number_of() -> go:count()

Data in memory bases is now only persisted if a change has taken place. This saves 25% of the time to run the blocks world.

---

Agreement:
- number: the NP must agree in number with the VP (singular, plural)
- gender (male, female, neuter)
- tense (present, past, future)

It would be nice to have agreement. I haven't needed it yet, but it is necessary in some cases.
But is it? Yes, I agree that agreement is necessary for generation; but is it necessary for parsing?
It's a big deal in CLE; so let's see if they can come up with an example where agreement is necessary for a proper parse.

I haven't found one. And even if I could find a far-fetched example of a sentence that could only be parsed properly with agreement, is it worth it? Agreement requires a lot of extra work.

And what about generation? Yes, to generate a sentence one must remember that the number, the gender and the tense of the NP and the VP match. But is it necessary to have explicit mention of agreement in the rules? It is the semantics that leads the chosen rule, both for the NP and the VP. There's no need to have them match.

If I were to have agreement, its syntax would be an extra field next to 'rule' and 'sense': 'agree'. But I will wait until have a good case before I add it in.

    You might not need agreement.

## 2021-03-27

Relations are now used as messages, in the communication between the web client and the server. Relations can be converted to JSON effortlessly.

---

This sequence also goes wrong:

    Put a small red block into the box
    Will you please stack up both of the red blocks and either a green cube or a pyramid?

The system creates a stack of a green block with on it a single small red block. Now I think I know why: it somehow interprets "both of the red blocks" as "a small red block" from the anaphora queue.

## 2021-03-26

The rewrite I started on January 21 is done! Answering and execution are now separated. This is a major advancement and I will create a new release.

It is now possible to create a plan and execute that plan later, by creating a `goal` for it. This change shows itself most manifestly in the fact that the system answers "OK" directly after it has received the question, and before it starts executing the plan. This is an advancement even over SHRDLU, that says "OK" only after the execution is complete. To appreciate this difference it is only necessary to imagine that you want to tell the system to "Stop!" performing its current action. 

The working representation of the blocks after all interactions now looks like this:

![Initial blocks world](archive/blocksworld3.png)

And here's a first impression of the web demo I am working on.

    "Pick up a big red block."

![Initial blocks world](archive/blocksworld4.png)

It is still too brittle to release.

The system class now looks nothing like it did before. It used to contain all language processes. Now it just passes relational messages. 

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
