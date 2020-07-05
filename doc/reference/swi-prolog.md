# SWI Prolog predicates

Some interesting predicates I might use.

A "list" is an ordered collection of items, that may contain duplicates. 
A "set" is an unordered collections of items, which do not contain duplicates.
A "bag" is an unordered collection of items that may contain duplicates.

## forall

For all alternative bindings of Cond, Action can be proven.

`forall/2` does not change any variable bindings.

    forall(:Cond, :Action)
    
Example use

    forall(Generator, SideEffect)     
    
https://www.swi-prolog.org/pldoc/man?section=forall2    

## foreach

True if conjunction of results is true. Unlike forall/2, which runs a failure-driven loop that proves Goal for each solution of Generator, foreach/2 creates a conjunction. Each member of the conjunction is a copy of Goal, where the variables it shares with Generator are filled with the values from the corresponding solution.

The implementation executes forall/2 if Goal does not contain any variables that are not shared with Generator.
    
    foreach(:Generator, :Goal)
    
Example use

    ?- foreach(between(1,4,X), dif(X,Y)), Y = 5.
    Y = 5.
    ?- foreach(between(1,4,X), dif(X,Y)), Y = 3.
    false.    

https://www.swi-prolog.org/pldoc/doc_for?object=foreach/2    

## findall

Create a list of the instantiations Template gets successively on backtracking over Goal and unify the result with Bag. 
    
    findall(+Template, :Goal, -Bag)
    
As `findall/3`, but returns the result as the difference list Bag-Tail. The 3-argument version is defined as    
    
    findall(+Template, :Goal, -Bag, +Tail)

## findnsols

As `findall/3` and `findall/4`, but generates at most N solutions.

    findnsols(+N, @Template, :Goal, -List)
    findnsols(+N, @Template, :Goal, -List, ?Tail)

https://www.swi-prolog.org/pldoc/doc_for?object=findnsols/4

## maplist

True if Goal is successfully applied on all matching elements of the list.

    maplist(:Goal, ?List1)
    maplist(:Goal, ?List1, ?List2)
    maplist(:Goal, ?List1, ?List2, ?List3)
    maplist(:Goal, ?List1, ?List2, ?List3, ?List4)

https://www.swi-prolog.org/pldoc/man?predicate=maplist/2

## convlist

Similar to `maplist/3`, but elements for which call(Goal, ElemIn, _) fails are omitted from ListOut.

    convlist(:Goal, +ListIn, -ListOut)
    
https://www.swi-prolog.org/pldoc/doc_for?object=convlist/3    

## include

Filter elements for which Goal succeeds. True if List2 contains those elements Xi of List1 for which call(Goal, Xi) succeeds.

    include(:Goal, +List1, ?List2)
    
https://www.swi-prolog.org/pldoc/doc_for?object=include/3

## exclude

Filter elements for which Goal fails. True if List2 contains those elements Xi of List1 for which call(Goal, Xi) fails.

    exclude(:Goal, +List1, ?List2)
    
https://www.swi-prolog.org/pldoc/doc_for?object=exclude/3        

## partition

Filter elements of List according to Pred. True if Included contains all elements for which call(Pred, X) succeeds and Excluded contains the remaining elements.    

    partition(:Pred, +List, ?Included, ?Excluded)

https://www.swi-prolog.org/pldoc/doc_for?object=partition/4

## bagof

Unify Bag with the alternatives of Template. If Goal has free variables besides the one sharing with Template, `bagof/3` will backtrack over the alternatives of these free variables, unifying Bag with the corresponding alternatives of Template. The construct +Var^Goal tells `bagof/3` not to bind Var in Goal. `bagof/3` fails if Goal has no solutions.    

    bagof(+Template, :Goal, -Bag)        
    
https://www.swi-prolog.org/pldoc/doc_for?object=bagof/3

## setof

Equivalent to `bagof/3`, but sorts the result using sort/2 to get a sorted list of alternatives without duplicates.

    setof(+Template, +Goal, -Set)
    
https://www.swi-prolog.org/pldoc/doc_for?object=setof/3

## aggregate

Aggregate bindings in Goal according to Template. The `aggregate/3` version performs `bagof/3` on Goal.

    aggregate(+Template, :Goal, -Result)
    
Aggregate bindings in Goal according to Template. The `aggregate/4` version performs `setof/3` on Goal.    
    
    aggregate(+Template, +Discriminator, :Goal, -Result)

https://www.swi-prolog.org/pldoc/man?section=aggregate

## aggregate_all

Aggregate bindings in Goal according to Template. The `aggregate_all/3` version performs `findall/3` on Goal. Note that this predicate fails if Template contains one or more of min(X), max(X), min(X,Witness) or max(X,Witness) and Goal has no solutions, i.e., the minimum and maximum of an empty set is undefined.
The Template values count, sum(X), max(X), min(X), max(X,W) and min(X,W) are processed incrementally rather than using `findall/3` and run in constant memory.

    aggregate_all(+Template, :Goal, -Result)
    
Aggregate bindings in Goal according to Template. The `aggregate/4` version performs `setof/3` on Goal.    
    
    aggregate(+Template, +Discriminator, :Goal, -Result)
    
https://www.swi-prolog.org/pldoc/man?predicate=aggregate_all/3
