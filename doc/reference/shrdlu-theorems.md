# Types of theorems in SHRDLU

With page numbers in "Understanding Natural Language"

## PLANNER

All functions that start with `TH` are PLANNER functions.

THGOAL
: will try to find an assertion in the data base, or prove it using other theorems (109)

    (THGOAL (#ON $?X $?Z))

THTBF
: to be followed? modifier of `THGOAL`. planner's way of saying "try anything you know which can help prove it". (110) advice that causes the evaluator to try all theorems whose consequent is of a form which matches the goal (i.e. a theorem with a consequent `($?Z :TURING)` would be tried, but one of the form `(#HAPPY $?Z)` or `(#FALLIBLE $?Y $?Z)` would not. 

    (THGOAL (#PERSUASIVE $?Y) (THTBF THTRUE))
    
THUSE
: modifier of `THGOAL`. gives advice on what other theorems to use and in what order (109)

    (THGOAL (#GRASP :B1) (THUSE TC-GRASP))

THFIND
: takes four pieces of information. When we use ALL, it looks for as many as it can find, and succeeds if it finds any. If we use an integer, it succeeds as soon as it finds that many, without looking for more. If we want to be more complex, we can tell it three things: (a) how many it needs to succeed (b) how many it needs to quit looking, and (c) whether to succeed or fail if it reaches the upper limit set in b. Thus if we want to find exactly three objects, we can use a parameter of (3 4 NIL), which means: "Don't succeed unless there are three, look for a fourth, but if you find it, fail". The second bit of information tells it what we want in the list it returns. For our purposes, this will always be the variable name of the object we are interested in. The third item is a list of variables to be used in the process. This acts much like an existential quantifier in the predicate calculus notation. The fourth item is the body of the statement. It is the body that must be satisfied of each object to be found.(111)

    (THFIND ALL $?X (X) THGOAL(#BLOCK $?X) THGOAL(#COLOR $?X RED))
    
THASSERT
: a function which, when evaluated, stores its argument in the data base of assertions or the database of theorems (which are cross-referenced to give the system efficient look-up capabilities) (113)

    (THASSERT (#HUMAN :TURING))
    
THERASE
: removes the assertion from the database.

    (THERASE (#ON $?X $?Y))    
    
THCONSE
: consequent. implication. This states that if we ever want to establish a goal of the form `(#FALLIBLE $?X)` we can do this by accomplishing the goal `(#HUMAN $?X)`, where, as before, the prefix characters $? indicate that X is a variable. (113)

    (DEFTHEOREM THEOREM1
        (THCONSE (X) (#FALLIBLE $?X) (THGOAL (#HUMAN $?X))))
        
THANTE
: antecedent. induction. The more conclusions we draw right at the time information is asserted, the easier proofs will be, since they will not have to deduce these consequences over and over again. The following theorem says that when we assert that X likes something, we should also assert `(#HUMAN $?X)` (116)

    (DEFTHEOREM THEOREM2
        (THANTE (X Y) (#LIKES $?X $?Y)
            (THASSERT (#HUMAN $?X))))
                    
THPROG
: program. like a function in an imperative programming language. Planner's equivalent of a LISP PROG, complete with GO statements, tags, RETURN, etc. Acts as an existential quantifier. (114)

    (THPROG (Y)
        (THGOAL (#FALLIBLE $?Y) (THTBF THTRUE)))        

THAMONG
: chooses its variable bindings from "among" a given list (136)

    (THFIND 3 $?X1 (X1) (THAMONG X1 (QUOTE(:B1 :B4 :B6 :B7))))
    
THFAIL
: causes not just that theorem but the entire goal to fail, regardless of what other theorems there are (143)

    (THCONSE (X1)
        (#OWN :FRIEND $?X1)
        (THGOAL (#IS $?X1 #BLOCK))
        (THGOAL (#COLOR $?X1 #RED))
        )
            (THFAIL THGOAL))

THCOND
: conditional? if-then? (122)

THPUTPROP
: ? (160)    

THAND
: logical AND (109)

THOR
: logical OR (109)

THNOT
: logical NOT

THTRUE
: logical TRUE

THFALSE
: logical FALSE
