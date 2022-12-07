# Considerations

## Commands are language too!

Linguistics has traditionally been concerned mainly with declarative sentences. These statements are the standard, and they are the only ones that occur in logic, since they have a truth value.

But commands are very important in NLP, of course. And it turns out that they have some characteristics that are not found in statements. Most importantly, a statement describes the world, and a command changes the world.

The difference is also important in quantification. The phrase "two blocks" in "Stack up two blocks" has the system select and act on two blocks exactly. In "Did you stack up two blocks?" has the system count all objects and checks if their number is two.

Verbs can have two very different sences. One is the command sense that makes the system stack up blocks. The other sense only selects if stacking up has been performed at some point. Compare 

    { rule: vp_imperative(P1) -> 'pick' 'up' np(E1),                               sense: go:do($np1, dom:do_pick_up_smart(E1)) }

with    

    { rule: vp(P1) -> np(E1) 'pick' np(E2) 'up',                                   sense: go:check($np1, go:check($np2, dom:pick_up(P1, E1, E2))) }

## The time of acknowledgement

SHRDLU only says "OK" after it has performed the action. That's because it performs the action before responding. I think it's better to respond immediately and to start the action only after that. And that's how it's built. NLI-GO is done with the request after it has sent the acknowledgement. After the "OK" it starts executing the task in a separate thread, which can take quite some time. For example, when stacking up some blocks.

