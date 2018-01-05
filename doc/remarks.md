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
