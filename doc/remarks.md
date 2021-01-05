## 2021-01-05

I created a stack trace! It is output when the answerer finds no results. It is an automation of what I have been doing all along to fix a program. 

An object called CallStack keeps track of the functions being called recursively. When a function gives no results, a copy of the stack at that time is frozen.

Here's an example stack trace that is created when I try to fit a block into the box that doesn't fit:

    Stack trace
    10. go_greater_than_equals(W$6, Width$25)
        {ColIndex$5:4, Index$6:4, Line$6:1000, VerLines$2:[600, 640, 800, 840, 1000], W$6:0, Width$25:200, X1$18:1000}&{A1$1:0, A1$2:-1, B1$1:200, B1$2:-1, StartY$1:200, StartY$2:600, Success$11:true, Success$12:true, Success$14:true, Success$15:true, Success$3:true, Success$4:true, Success$7:true, Success$8:true}
    
    9. go_list_foreach(VerLines$2, Index$6, Line$6, go_subtract(Line$6, X1$18, W$6) go_greater_than_equals(W$6, Width$25) go_subtract(Index$6, ColIndex$5, ColSpan$2) go_break())
       {ColIndex$5:4, VerLines$2:[600, 640, 800, 840, 1000], Width$25:200, X1$18:1000}&{A1$1:0, A1$2:-1, B1$1:200, B1$2:-1, StartY$1:200, StartY$2:600, Success$11:true, Success$12:true, Success$14:true, Success$15:true, Success$3:true, Success$4:true, Success$7:true, Success$8:true}
    
    8. dom_find_span(Width$25, VerLines$2, ColIndex$5, ColSpan$2)
       {ColIndex$5:4, HorLines$2:[600, 640, 840, 940, 1000], Length$25:300, VerLines$2:[600, 640, 800, 840, 1000], Width$25:200, X$21:1000}&{A1$1:0, A1$2:-1, B1$1:200, B1$2:-1, StartY$1:200, StartY$2:600, Success$11:true, Success$12:true, Success$14:true, Success$15:true, Success$3:true, Success$4:true, Success$7:true, Success$8:true}
    
    7. go_list_foreach(VerLines$2, ColIndex$5, X$21, go_list_get(HorLines$2, 0, StartY$2) dom_find_span(Width$25, VerLines$2, ColIndex$5, ColSpan$2) go_add(ColIndex$5, ColSpan$2, V1$6) go_subtract(V1$6, 1, ColEnd$2) go_list_foreach(HorLines$2, LineIndex$2, Y2$6, go_greater_than(LineIndex$2, 0) go_subtract(LineIndex$2, 1, RowIndex$2) go_if_then_else(dom_span_free(ColIndex$5, ColEnd$2, RowIndex$2, fixed), go_subtract(Y2$6, StartY$2, SpanLength$2) go_greater_than_equals(SpanLength$2, Length$25) go_let(A1$2, X$21) go_let(B1$2, StartY$2) go_break(), go_let(StartY$2, Y2$6))) go_not_equals(A1$2, -1) go_break())
       {E5:`block:big-red`, HorLines$2:[600, 640, 840, 940, 1000], Length$25:300, VerLines$2:[600, 640, 800, 840, 1000], Width$25:200}&{A1$1:0, A1$2:-1, B1$1:200, B1$2:-1, StartY$1:200, StartY$2:600, Success$11:true, Success$12:true, Success$14:true, Success$15:true, Success$3:true, Success$4:true, Success$7:true, Success$8:true}
    
    6. dom_do_find_free_position(E5, fixed, HorLines$2, VerLines$2, X$18, Y$17)
       {BoundX1$2:600, BoundX2$2:1000, BoundY1$2:600, BoundY2$2:1000, E5:`block:big-red`, E6:`box:box`, HorLines$2:[600, 640, 840, 940, 1000], Objects$2:[`pyramid:blue`, `block:blue`], VerLines$2:[600, 640, 800, 840, 1000]}&{A1$1:0, B1$1:200, StartY$1:200, Success$3:true}
    
    5. dom_do_find_free_space(E6, E5, X$18, Y$17)
       {E5:`block:big-red`, E6:`box:box`, Z$26:0}&{A1$1:0, B1$1:200, StartY$1:200, Success$3:true}
    
    4. dom_do_put_in(E5, E6)
       {E5:`block:big-red`, E6:`box:box`}&{A1$1:0, B1$1:200, StartY$1:200, Success$3:true}
    
    3. dom_do_put_in_smart(S, E5, E6)
       {E5:`block:big-red`, E6:`box:box`}
    
    2. go_quant_foreach(go_quant(some, E6, go_definite_reference(E6, dom_box(E6))), dom_do_put_in_smart(S, E5, E6))
       {E5:`block:big-red`}
    
    1. go_quant_foreach(go_quant(some, E5, go_definite_reference(E5, dom_red(E5) dom_block(E5))), go_quant_foreach(go_quant(some, E6, go_definite_reference(E6, dom_box(E6))), dom_do_put_in_smart(S, E5, E6)))
       {}

This is very useful! It makes debugging a lot easier, I expect. What you also see here is that some constructs are very heavy and this shows up in the call stack. So this is an indication that I should keep them simple. 

## 2021-01-02

Very happy with the fact that it was possible to create the typical oblique projection of SHRDLU with THREE.JS. The demo is coming along fine. I am now working on the hand, which a largely ignored up to now, but which has such an important role in the demo. It hit on me a few days ago that the hand actually moves smoothly in the SHRDLU demo. In the database it just switches to a new position, so this is an important difference that I will need to find a solution for.

The hand picks up a block in the center, and when the hand moves up, the block should follow. The block is not at the same position as the hand of course, it must be translated.

## 2021-01-03

Since I now have an interactive blocks world demo, I can interact with it immediately and try different sentences. They all break down in terrible ways! I have seen cubes floating in space, cubes taking up the same space, execution times of 5 seconds... Wow. This thing is not robust by any means.

However, I fixed the first problems I encountered. And I made an interesting innovation: I added an extra rules layer of "physics". This layer of rules is responsible for ensuring that all relations (contain, support, cleartop) stay intact, whatever you do. All objects that move now go trough a single function `phys_move_object`. This function has as input just the position of the object. It breaks up any existing relations the object might have, and then rebuilds them, just by looking at the position of the object. I was afraid that this would be very expensive, and it is not cheap, but the enormous advantage is that you can now move objects around without thinking about the relations _at all_. They have been fixed in a special layer once and for all.
