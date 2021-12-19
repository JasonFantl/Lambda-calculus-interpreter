
Name will only contain letters and numbers

Function names will be capatalized

variable names will be lowercase


## Grammar:

<program> ::= [these are seperated by new lines]
    | <fname> [<var> <var >...] = <exp>
    | <exp>

<exp> ::=
    | <var>
    | \<var> . <exp>
    | ( <exp> )
    | ( <exp> <exp> )
    | ( <fname> [<exp> <exp> ...] )

## Later:
use precedence to remove unneeded pranthesis
alpha transform (is this needed?)
beta reduction
