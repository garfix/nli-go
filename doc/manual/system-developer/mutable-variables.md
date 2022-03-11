# Mutable variables

Mutable variables are stored in the "scope"-typed cursor.

They are treated differently from normal variables:

- They are added to the binding in the `step`, and removed afterwards. They don't stay in a binding.
- They are added to the scope by `go:let()`: explicit assignment 
- They are added to the scope when a "simple"-typed cursor is completed (an implicit kind of assignment)
