Declare predications in the sortal hierarchy

    go:has_sort(predication, entity);
    go:has_sort(stative, predication);
    go:has_sort(dynamic, predication);

Attach sort-predication to predication variables using tags:

{ rule: tv(P1, E1, E2) -> 'supports',       tag: go:sort(P1, stative) }
{ rule: tv(P1, E1, E2) -> 'clear' 'off',    tag: go:sort(P1, dynamic) }

Example use of this sort in an application

    filter_event(Event) :- 
        go:atom(Event, A) go:sort(A, Sort)
        if  [Sort == dynamic] then
            go:not(location(Event, _, _, _, _))
        else 
            if [Sort == stative] then
                location(Event, _, _, _, _)
            end
        end
    ;
