
Name will only contain letters and numbers

Function names will be capatalized

variable names will be lowercase


## First pass of grammar:

<var> ::= a single letter

<lexp> ::=
    | <var>
    | \<var> . <lexp>
    | ( <lexp> <lexp> )

## Later:
word variables
read from file
use precedence to remove unneeded pranthesis
named functions
alpha transform (is this needed?)
beta reduction
