
Name will only contain letters and numbers

Function names will be capitalized

variable names will be lowercase

## Easy Grammar:

<exp> ::=
    | <var>
    | \<var> . <exp>
    | ( <exp> <exp> )

## Grammar:
<program> ::= [these are separated by new lines]
    | def <fname> <exp>
    | <exp>

<exp> ::=
    | <var>
    | \<var> [<var> ... <var>] . <exp> [right associative]
    | ( <exp> )
    | <exp> <exp> [left associative]
    | <fname>

## Later:
named functions
semantics
alpha transform (is this needed?)
beta reduction
named function reductions 
vscode language extension