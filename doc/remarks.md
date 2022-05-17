## 2022-05-15

Fixing the tests again. "Why?" now fails because ellipsis resolution copies part of the previous sentence that hadn't been subjected to anaphora resolution. 

## 2022-05-14

All discourse entities must have a sort. Even if it is only `entity`.

Discourse entities are created when a user sentence is read, or when an answer is created by the system. In both cases the entities and their sorts are added.

## 2022-05-13

I'm now storing discourse entity sorts in a separate store.

Problem is, entities from responses are stored as well. And these may be an array of entities. Therefore the sorts may be arrays as well.

## 2022-05-09

About sorts: I dont want a double declaration of sort relation `dom:person(E1)` and tag `go:agree(E1, sort, person)`.
In stead, I just want the relation. And when the relation is a sort, this can be specified in `predicates.sort`. The sort and the relation are equal. No need to separate them.

## 2022-05-09

"the green one" is both a reference and one-anaphora. The one-anaphora should be resolved first.

It also means that the resolvings should be processed immediately, because they are used right after. And not collected to be processed later.

"Had you touched any pyramid before you put the green one on the little cube?"

"one" does not refer to "any pyramid" because the system does not know it's a pyramid (sort pyramid). It can only deduce that from the id at the moment.

Time to introduce sort tags.

## 2022-05-08

Since A is the simplest solution, I only need to specify which part of the sentence is the subject.
I can do this by annotating it with a tag.
But even without tags there's a way to make the earlier parts of the sentence more likely referents than latter ones: 
add the entities of the sentence from right to left to the anaphora queue, so that the first entity is at the end (and hence treated first).

... it turns out that is was already the way it is implemented. There was another problem. "the small red block" had no binding in the dialog, for some reason.

The reason: "Does the small red block support the green pyramid?" restult in no bindings (the answer is "no"), and it is these result bindings that are stored.

That was not the only problem: the system stores the entities in the anaphora queue before the variables have been replaced. So the queue may contain variables that are not actually used.

The next sentence to consider also has a problem with the anaphora queue as it currently stands:

    Had you touched any pyramid before you put the green one on the little cube?

"one" refers to an earlier entity in the same sentence (any pyramid).

It's probably best to extend the anaphora queue while it is being used to resolve references. If it is extended before this resolve process, it contains items that should not be there yet. If it is extended afer this resolve process, it doesn't contain the items that are needed.

The reason that I'm now building an explicit A Q again, is that I need to be able to add entities to it, one by one.

## 2022-05-07

In this interaction:

    What does the box contain?
    The blue pyramid and the blue block

the two objects are bound to the same variable.

But when only one of them is referenced

    What is the pyramid supported by?

the variable in the new sentence should not be replaced by the one holding the 2 objects.

This is an interesting issue. The two objects should be referencable together ("move them out of the box"). But in our case it seems as if a new referent is created out of the existing referent set:

    E18: [blue-block1, blue-pyramid]
    =>
    E23: blue-pyramid

Solution: if a reference refers to a single entity from a referent group, keep the reference variable unchanged, but bind it single referent's value.

===

I now have the following problem:

    Does the small red block support the green pyramid?
    Yes
    Put the littlest pyramid on top of it
    OK

... puts the littlest pyramid on top of the green pyramid (itself!)

Why does "it" not refer to "the green pyramid"?

- A the preferred referent is the subject
- B the object of "put_on" is a steady surface (which pyramids are not)
- C "put_on" should have a constraint that one cannot put something on top of itself

B is sufficient in this case but not in other cases

## 2022-05-06

Should the anaphora queue contain only bound entities?

You need to be able to refer to an unbound entity. However, how would you check that it an entity from the queue is the same as the one examined?

If the entity from the queue has an id, use this to check if they match.
If the entity from the queue does not have an id, use agreement to check if they match.

## 2022-04-23

The SHRDLU dialog does not have any definite anaphoric references (!) (where "the green block" refers to the block that was mentioned in the previous sentence)

References will be replaced by tags:

- `go:back_reference(E1, none)` "it" "them" - 
  - `tag: reference(E1) agree(E1, sort, object) agree(E1, number, singular)`  
  - `tag: reference(E1) agree(E1, sort, person) agree(E1, number, plural)`
- `go:definite_reference(E1, $nbar)` "the" 
  - no tag, but all quants are tried to be resolved 
- `go:sortal_back_reference(E1)` "one" 
  - `tag: sortal_reference(E1)`
  - searches for a sort, then tries to resolve 

## 2022-04-20

I'm going to use the `tag` property of a syntax rule for the features. And since I don't use "feature structures" it's better that I use a different name.

What about "agree"?

    tag: agree(E1, number, plural)
    tag: agree(E1, gender, male)
    tag: agree(E1, sort, person)

A tag that "agrees" should be checked for conflicts.

## 2022-04-17

This new diagram shows how I'm going to handle things:

![entities](diagram/entities2.png)

I am only going to add a meta-level of Features. This level of information (or "second order") is not about the value of the entity; it is about the possible values of the entity, or the type of the value.

An important aspect of features is that they mustn't conflict. 

This way you can say `plural(e5)` (`e5` is an atom that represents the object hold by the variable `E5`) while variable `E5` is not bound yet, and may not be bound at all. "I saw a dog. It chased a cat." - which dog exactly is not known at this time. Still, we know the sort and the number of the entity. Important for references.

The anaphora resolution phase merges variables if one refers to the other. Their meta data is merged. If `E11` refers to `E5`, the occurrences of `E11` will be replaced by `E5`.

You can reason about meta data as like normal data, but remember that the entities are different. So if the system needs to check `animal(E1)`, it may deduce this from `dog(E1)`.

## 2022-04-16

It's still possible to do it (assign atoms to sentence variables). When a database is consulted, the associated db id will be used, and if there is none, a variable can be used to consult the db.

Using atoms for entities will make probably some things easier. But the full extent of it is hard to predict.

But one of the things that can be done is to place the features inside an actual dialog database. Which means these facts can be used from anywhere in the code.

## 2022-04-04

Assigning discourse elements to sentence variables may not be such a good idea. All variables will be bound from the start and can't be bound later.

Another idea would be to replace some sentence variables by ones from previous sentences, in an anaphora resolution step.

## 2022-04-03

Also:

    The evening star is planet Venus
    The evening star appears in the west
    The morning star is the evening star
    Is the morning star planet Venus?
    yes
    Does the moning star appear in the west?
    no (!)

and

    Ice is hard
    Steam is volatile
    *Ice is steam

===

I'm going to introduce the discourse entities level between sentence variables and shared ids:

![entities](diagram/entities.png)

Whereas the list of discourse entities now hold ids, they should hold actual shared entity objects. This discourse entity has attributes:

- attributes derived from productions (gender, number, etc)
- has 0..n ids (reference to zero, one, or multiple shared entities)

Every sentence variable should refer to a discourse entity (DE). The discourse entity can say "plural", while the database ids are not (yet) known.

When are sentence variables linked to discourse entities? As soon as the sentence is parsed. Each variable then links to a different DE.
At anaphora resolution time, the sentence entity may unlink from the original DE and link to the final DE.

In DRT "From discourse to logic" sentence variables are equated `x = y`. Simple to implement, but a hassle to compute. You need to check these assignments every time you do a lookup. 

DRT

    Jones ownes Ulysses. It fascinates him.
    x y u v
    Jones(x)
    Ulysses(y)
    x owns y
    u fascinates v
    u = y
    v = x

NLI-GO

    Jones ownes Ulysses. It fascinates him.
    x: DE1 y: DE2 u: DE3 v: DE4
    Jones(x)
    Ulysses(y)
    x owns y
    u fascinates v
    process "it" (reference)    u: DE2
    process "him" (reference)   v: DE1
    x: DE1 y: DE2 u: DE2 v: DE1

All entities in the sentence need to be resolved to discourse entities, not just anaphora. "The blue block" may well refer to an entity named before.

Is it possible to have a separate anaphora resolution phase?

## 2022-04-02

For the last days I tried to incorporate Frege's sentence into the system

    "the morning star is the evening star"

But I give up. With reason. For one, this type of sentence is very very rare. It only occurs when you find out that two things that you thought were different are actually the same. This is absolutely irrelevant to an nli system. To create a special layer of representation just for this case makes it unnecessarily complex. Second: there are three ways of dealing with it (that I could think of):

- whenever you bind a variable to morning-star, bind it to evening-star as well (and the other way around)
- collapse the different database representations of the two entities into one
- assert identity and extend all inferences with an extra check for identity `planet(E) :- identical(E, F) planet(F)`

Very very complicated, and for what?

## 2022-03-24

I am thinking about creating a replacement for tags: productions; relations that are added to the dialog context once a parse is complete.
A production like `sort(E1, person)` can then be used to restrict the name-lookup process. And a production like `gender(E1, male)` can restrict the anaphora resolution process.

I may need meta information about these productions:

- `type`: default: `entity`
- `number`: agreement (should not conflict) 
- `sort`: derivable agreement (man -> person)
- `determinate`: changable (false -> true, but not: true -> false) 

## 2022-03-21

I'll try to build a general approach to anaphora resolution.

## 2022-03-11

The server is complete.

## 2022-02-28

Another idea: both break and cancel immediately return with no bindings. This stops the flow for that binding. Other bindings continue. But the break will also add its latest binding as a child-frame-bindings to the loop-cursor. Return does the same to the scope cursor.

## 2022-02-27

Still rebuilding the server, and making the processes synchronous.

I stumbled on a problem that I nadn't realized before. What if the process currently had multiple bindings, and the program flow is different for these bindings: 

- one binding causes the execution of child procedures, while the other doesn't
- one binding breaks while the other doesn't

Before I hadn't thought about this problem, while it was there all along. Actually, by executing child procedures _inline_, as I am doing now, I solved the first problem. But I stumbled on it only because of the second problem.

At the moment, if one binding causes a break, all bindings break. 

Yesterday, for a few hours, I thought that this would collapse my entire process structure. But it doesn't have to be so catastrophic. What I need to do is to keep the state of the program flow, per binding. And treat bindings as program flows. Program flows that may have breaked, and they should store this. Once a flow has breaked, it should not execute more steps, until the root of the break (the loop or the procedure), is reached.  

## 2022-02-12

Rebuilding the server.

Each client can tell the process runner to start a process.
Processes can run in parallel, but processes consist of step and the step is executed in a mutexed area: no two steps are executed together.
A process can create new processes; they are managed by the process list.
Once a process is done, it is removed from the process list; the system is ready when there are no more processes.

Need more mutexes.

## 2022-01-31

Make sure the process keeps running so that the call-stack is not broken down; it can't be rebuilt.

todo:

- continue fixing all `CreateChildStackFrame` calls, until you run into a problem (then undo the last one)
- register with the process for `wait for` messages
- have `wait for` retry for 1 minute, make it notify the event listeners, and continue running the process
- when a `wait for` is caught, return the `wait for` message
- on `break`, `return`, etc, break down the call stack
- continue fixing all `CreateChildStackFrame` calls
- clean up the rubbish

## 2022-01-30

About "simply processing a relation set": it is not simply possible to execute it and disregard the process stack. The process stack serves functions like `break` and `return`. Functions that require the call stack.

I'm thinking of the following technique for calling child relation sets. When a function F needs to execute a child set, it:

- adds a callback function to the current function F in the current stack frame
- places the child set on the stack
- returns a binding constant that says: ignore me

When the child set is executed, the process runner returns to F, it notices the callback function, and calls this callback, in stead of the function (it resumes where it left off). Before calling the callback, it clears the callback.

When a `break` or `return` occurs, F will be deleted and the process runner will not execute the callback. 

===

One callback function would not have been so bad, but some functions have a large substructure of helper functions, each of which may call a child relation set (`quant_ordered_list` comes to mind). Basically you just want to call child sets; and there are only a few problems: `break`, `cancel`, `return`, `fail`. All of which break down some part of the stack. 

After the stack has been broken down, the functions that execute these stack frames need to stop, immediately. This can be done by try/catch, or in Go, with panic/recover. Each function is wrapped by a panic/recover, and this wrapper checks if more of the calling functions need to be stopped, and if so, panic again.

Each of the keywords that tear down part of the stack, calls "panic", with a special control flag. The recover that catches this panic then checks if the call-stack-depth matches the process stack depth; if not, it panics again. Until the call stack depth matches the process stack depth. 

## 2022-01-29

I also want to simplify the process runner. Processes are currently very robust, but also very hard to understand. I dread having to make changes to functions that process relation sets as part of their procedure.

I thought about it for a long time and found myself coming back to the idea that a process, once started, should just keep running until its done. This is by far conceptually easiest to understand. Of course there were reasons for processes to interrupt: a user needed to answer a clarification question, or the user interface should update its state (moving blocks).

I worked this out. When a process is started, it simply processes a relation set. When the relation set is processed, the process is finished. A relation itself may also start processing another relation set. In the new scheme it can do this inline and wait for the result.

What about the callbacks to the user? The `wait for` relation that before stopped the process, will now continue to execute, _for a short amount of time_. It will wait for a fact, check 10 times a second, for a period of, say a minute. When the minute is over withouyt success, the `wait for` fails. The `wait for` will notify the wait-for-listeners, so that they don't waste time polling for `wait for`s.

## 2022-01-28

I introduced the server. It's a service that listens on a port for incomming messages, and in the background keeps all systems alive. The system now stays alive for the entire dialog with the user.

## 2022-01-22

The interactive demo doesn't work anymore. I know what's wrong with it: the dialog context is being rebuild at times that it shouldn't and this disturbs the process. The reason is that the process is restarted many times during the ansering of a single sentence. And that's ok for the processes, but not for the dialog entities and the clause list.

This is the drop. The process was already too slow; it had to redo the same code multiple times; and the grammars had to be reloaded. I've had enought: I'm going to build a real service. I'll turn nli-go into a service that keeps running. Nice!

===

There will be an application that listens to a port for incoming JSON messages. 
A dispatcher listens for incoming messages. A message will have a dialog id.
When this id is new, a new dialog context is created, and started. Then the message is passed to the dialog context.

I found this example https://gist.github.com/miguelmota/301340db93de42b537df5588c1380863

send message to port:

    echo "Hello" | netcat localhost 3333


## 2022-01-16

Changed go:equals() and go:not_equals() to    

    [T1 == T2]
    [T1 != T2]

Changed numerical comparisons to

    [T1 > T2]
    [T1 >= T2]
    [T1 < T2]
    [T1 <= T2]

Changed assignment to

    [V1 := T1]

Changed operations (see below for explanation) to

    [T1 + T2]
    [T1 - T2]
    [T1 / T2]
    [T1 * T2]

=== 

I want to rewrite go:add(T1, T2, Result) and like procedures to 

    [T1 + T2]

a possibility:

    go:eval(Result, go:add(T1, T2, Result))

and have assignment recognize `go:eval()` as a special case

Other possibility

    go:add(T1, T2)

and have procedure calling recognize that the last argument is missing; let it create an internal variable; and work with its result. This way we can use procedures as functions with very little overhead.

I will make this because it's such an interesting idea. You can also use it with facts from the database, like this:

    age(pat, 51)
    age(sue, 49)

    older(A, B) :- [age(A) > age(B)];

It's hard to change the process runner to have it create a child frame for each of the arguments, so for the moment I will have it execute the arguments immediately.

===

I don't have access tp the number of arguments, so for the moment I will do it a bit different:

    older(A, B) :- [age(A, return) > age(B, rv)];

Where the atom `return` takes the place of the return value.

## 2022-01-15

I added some keywords to mentalese:

old:

    go:if_then_else(
        go:not(cleartop(E2)) do_find_free_space(E2, E1, X1, Y1),
        do_put_on_position(ParentEventId, E1, E2, X1, Y1),
        do_cleartop(ParentEventId, E2) do_put_on_center(ParentEventId, E1, E2)
    )

new:

    if go:not(cleartop(E2)) do_find_free_space(E2, E1, X1, Y1) then
        do_put_on_position(ParentEventId, E1, E2, X1, Y1)
    else
        do_cleartop(ParentEventId, E2) do_put_on_center(ParentEventId, E1, E2)
    end

I created the keywords `return`, `fail`, `break` and `cancel`.

I removed `let()` and introduced `[X = n]` for both `let()` and most forms of `unify()`. I kept `unify()` because there were some valid use cases for it.

===

The processes I created work well, but I don't want to impose them on any other developer, because he would go insane quite quickly. Also, they cause a lot of code to be run multiple times unnecessaily. So I'm thinking of using some sort of Promises (see javascript) or async/await.

One of the main points of the processes is that an application may be killed at any time, and the proces still knows its state. If can continue where it left off.

## 2022-01-14

This is what I did: update these discourse entity values while processing the sentence. For example in `addToQueue`? This would not _add to the queue_, but update the discourse entity variables.

This then finally allowed me to run all 25 interactions, and remove the anaphora queue. But it became very clear to me that anaphora resolution is far from solved. The list of open problems is large. There's a long way to go before this subject is solved to satisfaction. 

## 2022-01-12

Great article "The defeat of the Winograd Schema Challenge" https://arxiv.org/pdf/2201.02387v1.pdf

It refers to this article that may be interesting:

A. Kehler, L. Kertz, H. Rohde, J. L. Elman, Coherence and coreference revisited, Journal of semantics 25 (2008) 1–44.

from this, an example of subject preference

    a. Bush narrowly defeated Kerry, and special interests promptly began lobbying him. [¼Bush]
    b. Kerry was narrowly defeated by Bush, and special interests promptly began lobbying him. [¼Kerry]

    "The alternation found in examples (6a,b) can be used to argue for the existence of a grammatical subject preference."

also

    the grammatical role parallelism preference which favours referents that occupy the same grammatical role as the pronoun.

===

Maybe I can update these discourse entity values while processing the sentence? For example in `addToQueue`? This would not _add to the queue_, but update it.

The problem is that `back_reference()` needs access to non-local information.

Is it possible to do anaphora resolution as a separate step?

## 2022-01-11

Nested entities have less priority than root entities.

Discourse entity values are not present while creating the anaphora queue, since they are only added at the end...

## 2022-01-10

A discourse variable often contains multiple values. Up until now I used a simple binding for discourse variables, but they could not store multiple values per variable. I thought of creating a new structure, but then I thought: what if I create a new type of term type: entity ids? But then I saw I already had an array type: the list.

Combinding multiple bindings in a single binding with multiple lists is all fine. But what about extracting these lists into multiple bindings? This can lead very quickly to a large number of bindings in an ongoing dialog.

The discourse entities may contain lists, and that's ok. These lists are not actively used any more.

"Had you touched any pyramid before you put [the green one] on the little cube?" contains a within-clause reference. And it causes a problem with the new implementation.

## 2022-01-09

Whenever a definite reference is asked, the system creates a temporary anaphora queue. It is built like this:

    for each of the root clauses, from last to first
        add all entities in order of priority

## 2022-01-07

On 5 January I managed to run all 25 interactions without any errors (!)

I still have to replace the anaphora queue and the sentences list with the new clause list.

===

A single root clause can contain multiple clauses. They may have references within the root clause.

    John mended the vase that he broke.

Highlighting the centers in the last interactions of the dialog:

[] = new center
() = retain center

    H: Had [you] touched any pyramid before you put the green one on the little cube?     // subject has highest prio (syntax)
    S: Yes, [the green one]                                                               // the answer entity gets to be the center (auto)  
    H: When did you pick (it) up?                                                         // back_reference() is calculated and yields the entity (solve!)  
    S: While I was stacking up a large red block, a large green cube and the red cube     // has no answer entities, so keep (it) ! (auto)
    H: Why (did you pick (it) up)?                                                        // ellipsis resolution creates the same entity for 'it' (auto)  
    S: To get rid of (it)                                                                 // note: write grammar (lookup of center)

Problems to solve

- `back_reference` needs access to the active clause of current input
- `update center` needs the sentence to be solved because it needs the result of the `back_reference()`s   

Todo

- I need to do `go:dialog_add_root_clause(RootClauseTree, false)` before the `solve` phase so that `back_reference` can use the latest clause.
- Update the center at a later time (when?) create a new relation for it

## 2022-01-05

I must implement all entities with atoms. Currently they are variables, but it means that variables are used as values, and this is clumsy. Then there must be a mapping from these atoms to database ids.

## 2022-01-01

Perhaps there can be an extra processing step in `respond.rule` that locates the backward-looking center of the new sentence.
