NLI-GO launches concurrent processes. There are three types of processes that co-occur within a single system:

- language
- robot
- simple

"simple" processes just query some data and don't interact with the user.

A "language process" is the process of responding to user input. It may interact with the user for clarification.

A "robot process" involves the execution of a plan by the robot.

The language process and the robot process are [critical sections](https://en.wikipedia.org/wiki/Critical_section). Each of them can be executed only once within one system at any moment. If this is violated, the robot performs two actions at once and this is clearly impossible. The language process would mess up the dialog context.

In order to prevent this, the process is given a "resource", like "language" or "robot". The process list refuses a new process that wants to use a resource that's already in use.
