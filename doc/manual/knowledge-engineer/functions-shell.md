# Shell functions

These functions allow you to run command shell commands.

## exec

Executes a shell command with arguments.

    exec(Command, Args...)
    
* `Command`: a string
* `Args`: zero or more strings

## exec_response

Executes a shell command with arguments. The output is strored in a variable

    exec_response(Output, Command, Args...)
    
* `Output`: a variable    
* `Command`: a string
* `Args`: zero or more strings

