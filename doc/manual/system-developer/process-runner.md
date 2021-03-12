# The process runner

Each of the goals (`go:goal`) in the database is implemented by a process.

The process has a stack of frames, and each frame contains one or more relations, along with the active binding.

This explicit call stack enables all processes to be halted at any time, and continued later on.

Since a frame has multiple relations, there is a cursor that points to the active relation.

Each relation can be handled by multiple handlers. A handler may be a simple function, a solution function, a database lookup, a database modification, or a rule firing.

The process also keeps track of its mutable variables. Normale variables can be bound only once. Mutable variables (defined by `go:let`) can be bound multiple times, and hold only a single binding.

## Execution

The process runner runs the process, step by step. 

Each step, the active relation of last frame of the stack is inspected and executed. When this leads to a new stack frame, the step is done. Otherwise, the process is advanced. This means:

- the next handler is selected, and when done:
- the next binding is selected, and when done:
- the next relation is selected, and when done:
- the frame is popped

## Child frames

When a relation excutes, it may create a frame with child relations. In fact, it may need to create such a child frame multiple times, before it is done executing.

These child frames are not executed inline. The relation creates a child frame. The child frame is excuted, and the process runner _re-enters_ the relation.

In a single step, a relation may create only a single child frame.

There's two ways to do this: simple and complex.

The simple way is used when the number of child frames to be created is just one or two, or a simple array.

Start the relation with the following code

    cursor := messenger.GetCursor()
    index := cursor.GetState("index", 0)
    cursor.SetState("index", index + 1)

The `index` variable can then be used to distinguish between the first call (`index` == 0)

    // create the child frame
    messenger.CreateChildStackFrame(condition, mentalese.InitBindingSet(binding))

and the following calls (`index` > 0)

    // get the results of the processed child frame
    newBindings = cursor.GetChildFrameResultBindings()
    // and, if necessary, the next frame:
    messenger.CreateChildStackFrame(condition, mentalese.InitBindingSet(binding))

The complex way looks like this. At the beginning of the function, the childIndex is reset to 0.

    messenger.GetCursor().SetState("childIndex", 0)

If in any part of the processing of the relation, there's a need to call a child frame, do

    bindings, loading := messenger.ExecuteChildStackFrameAsync(set, bindings)

Then, if `loading` is true, quit everything and make sure that the top-level function of the relation, returns an empty binding.
This is the only thing that needs to be done to create an almost inline processing of the child frame.

## Processing instructions

Some relations, like `go:break` and `go:list_foreach` need to tell the process runner that they require some special treatment of the process flow.

For this reason, it can call `AddProcessInstruction`:

    messenger.AddProcessInstruction(mentalese.ProcessInstructionBreak, "")
