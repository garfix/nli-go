# Considerations

## Exceptions

Building on strong negation it is possible to define rules with exceptions. This technique is taken from Answer Set Programming, an extension to Prolog.

The rule we want to create is "Birds fly". This is done by the traditional `fly(X) :- bird(X)`, but we append the clause that allows for exceptions: `not(-fly(X))`. The complete sentense thus says: "X can fly if X is a bird and there are no explicit instances that say that X cannot fly"  

    // birds fly (except for the ones that don't)
    fly(X) :- bird(X) not( -fly(X) );
    
The exceptions themselves are simple rules:    
    
    // pinguins can't fly
    -fly(X) :- penguin(X);
    
In this example the rule is posed positively and the exception negatively. But the reverse is also possible. The following example states that birds can't swim, except penguins.

    // birds don't swim (except the ones that do)
    -swim(X) :- bird(X) not( swim(X) );
    
    // penguins can swim
    swim(X) :- penguin(X);    

More about exceptions in Answer Set Programming in [this article](https://www.aaai.org/Papers/AAAI/2008/AAAI08-130.pdf): 

"Using Answer Set Programming and Lambda Calculus to Characterize Natural Language Sentences with Normatives and Exceptions"
-- Chi ta Baral and Juraj Dzifcak; Tran Cao Son (2008)

