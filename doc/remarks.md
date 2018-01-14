# 2018-01-14

I released http://patrickvanbergen.com/dbpedia/app/ and mentioned it on Twitter. It is very unimpressive for a demo, but I had to release _something_ so that I can show 
other people what I am doing. It will be nice to gradually improve upon it. 

# 2018-01-09

This picture is just to clarify things to myself. A, B, C, and D are range variables in nested scopes.

A1 B1 C1 D1
A1 B1 C1 D2
A1 B1 C2 D1
A1 B1 C2 D3
A1 B2 C1 D1
A1 B2 C1 D4
A1 B2 C2 D1
A1 B2 C2 D5

When evaluating scopes, each scope is evaluated only once, to restrict the number of queries.
As you can see D1 and D2 occur multiple times in the result set. 
Suppose the quantifiers are: B = 2, C = 2, D = 2
This can be implemented as follows:
At quant C, all quants in C's scope are evaluated, which means D.
For all unique values C, the number of distinct values D should equal the quantifier of D, i.e. 2.
=> This does not hold, since below C2 are D2, D3, D4, and D5.
This means that the scopes can only be evaluated after the outermost scopes are evaluated, which means: after the query is done.

After the outermost relation set has been solved, the quantifiers of its quants should be evaluated. This affects the resulting bindings.

I think I am going to go with a special construct in SolveQuant. It keeps track of a global quant level variable.
This is not pretty, but it works. When the quant is done, and this is the outermost level (0), the quantifiers are evaluated.
I will chose this solution, because the function solveQuant() will solve the quant _fully_ and will not be dependent on some
outside function call. Furthermore, this attempt may not even succeed, and I don't want to rewrite my code too much for a failed
attempt. I will place a todo to be resolved later, if the attempt should succeed.

# 2018-01-07

I have done it that way, but I have introduced an IF / THEN construct 

    IF isa(E2, child) THEN subject(S5, E1) object(S5, E2) isa(S5, have) => has_child(E1, E2)
    
because this is much clearer (I think) and I didn't want to reuse any of the brackets for this purpose, because they already signify other things.

I am now running into the problem that in scoping, the range, which is always evaluated first, yields too many values.

If the range is isa(E, child), E resolves to all persons (who are children). This is too much for any real-world database.

So I am now thinking about combining the Range with the Scope. But i wonder why I didn't think of that before. I knew it crossed my mind. What was the reason?

---

I think this is it. A quant has three aspects, a Range, a Quantifier and a Scope. For example: dogs, 3 or more, have a bone. The relation between these aspects is thus:

The range (dogs) is the domain of discourse. And if that was all it could be just a part of the Scope ( => dog(X) has_bone(X) ), which limits the domain as well.

The quantifier is a check. Does the scope yield 3 entities? Does the scope variable yield 3 possible values?

The range then, is not necessary in most cases. To check if the scope yields 3 entities, it is not necessary to know that these entities are dogs. If they were cats, 
it would be just as well.

But when the Quantifier is _all_, or _most_, the scope becomes important. Because the number of dogs is different from that of cats. 
When the Quantifier is _all_, the number of results in the scope must match the number of dogs. COUNT(DISTINCT S1) = COUNT(R) where R is just the query isa(R, dog) 

The number of different values for R may be large, and even very large, but its number just needs to be counted. The individual values are not used for further processing.

The range relation set may be copied to extend the scope relation set, but it must remain independent as well.  

# 2018-01-03

I think I have found a better idea. Neither of the proposed alternatives from yesterday where really great. Also, allowing relation rewrites
without minding scope borders seems like a bad way to go.

What I am now thinking about is this

    [isa(E2, child)] subject(S5, E1) object(S5, E2) isa(S5, have) => has_child(E1, E2)

Which means: rewrite 'subject(S5, E1) object(S5, E2) isa(S5, have)' to 'has_child(E1, E2)',
if 'isa(E2, child)' occurs in the sentence. 'isa(E2, child)' is unaltered by the rewrite and may live in any scope of the sentence.

# 2018-01-02

Hi there! Happy 2018! I used to scope Quantifier and Range relations in the Relationize phase, quite early in the process.
But the problem was that it got in the way of generic 2 domain specific conversions, they became too complex.
Doing the scoping later (in the scoping phase where the quant is formed) is problematic too because collecting all relations
that contain the range variable results in too many relations. This is because the variable is used in other relations
higher up the parse tree as well. I will try to visualize:

    np(E1)       posession(E1, E2)
    |
    nbar(E1)
    |
    dp(D1) nbar(E1)

D1 forms the quantifier; E1 forms the range. The quantifier variable is fine. The range variable is not only bound to relations
below the nbar, but above it as well (in the example: possession(E1, E2))

It is hard to figure out which relations belong to the range and which relations do not. The solution I am now using is based on
the heuristic that each nbar relation forms a specialization. So I will only use the variable on the right side of the specification()

    specification(E1, X)

and all relations that are connected, directly or indirectly to this variable.
