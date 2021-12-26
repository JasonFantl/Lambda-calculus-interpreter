
Name will only contain letters and numbers

Function names will be capitalized

variable names will be lowercase

## Easy Grammar:

<exp> ::=
    | <var>
    | \<var> . <exp>
    | ( <exp> <exp> )

## Next Grammar:

<exp> ::=
    | <var>
    | \<var> . <exp>
    | <exp> <exp>
    | ( <exp> )

## With associativity
<exp> ::=
    | <var>
    | \<var> . <exp>    [right associative]
    | <exp> <exp>       [left associative]
    | ( <exp> )

## Next grammar:
<program> ::= [these are separated by new lines]
    | def <fname> \<var> [<var> ... <var>] . <exp>
    | <exp>

<exp> ::=
    | <var>
    | \<var> [<var> ... <var>] . <exp>
    | ( <exp> )
    | <exp> <exp>
    | <fname>

## Later:
named functions
semantics
alpha transform (is this needed?)
beta reduction
named function reductions 