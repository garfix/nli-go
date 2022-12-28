# Considerations

## Episodic memory in the system or in the application?

It's unusual for a database to keep track of past states. Most databases just keep track of the current state of affairs.

Is episodic memory something that needs to be added to the system, as a special faculty, or can it be handled in the application and the database? And should it?

I found that it is possible to handle episodic memory in the application, with no changes to the system. I think it's good to do it in the application, because different applications need different types of memory. For SHRDLU, for example, it's enough to have an integer  time that increases with each action. Other applications will have different requirements.
