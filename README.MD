## TODO
beta reductions only work on free variables, still variable captures on bound variables. Exe: `>>:s AND T T`

# Use

You can load a file when starting the program. For example, to load `bools.l` you would type 

``` shell
go run main.go -file=bools.l
```
_______

While in the shell you can load a file. For example

``` shell
>>:l bools.l
```

_______

While in the shell you can have the interpreter show each of the steps taken to evaluate an expression. For example

``` shell
>>:s (\a b . b a) x \x . x
```
_______

While in the shell you can load a file. For example

``` shell
>>:l bools.l
```
_______

Exit the shell by either pressing `Ctrl-z` or running

``` shell
>>:e
```
_______
Otherwise you can evaluate statements, examples can be seen in the next section

# Syntax

The basic lambda calculus usually has the grammar

```
<exp> ::=
    | <var>
    | \<var> . <exp>
    | ( <exp> <exp> )
```

but we want something a little easier, so we add a few features. 

First, parenthesis are not necessary for application, which is left associative. So we can re-write `(( a b ) c)` as `a b c`, and `(a (b c))` as `a (b c)`.

Then we often have nested lambda functions, so we now allow for multiple inputs defined at once. So can rewrite `\ x . \ y . y x` as `\ x y . y x`.

We also would find it handy to name some functions for later re-use, so we can store lambda expressions like `def Two \f x . f (f x)`, where `def` indicates a definition is being made, `Two` is the name we use to store the expression, and the rest is the expression.

**Note:** All names for definitions must start with a capital letter, and all variable names must start with a lower case letter.

We now have the grammar:

```
<program> ::=
    | def <name> <exp>
    | <exp>

<exp> ::=
    | <var>
    | \<var> [<var> ... <var>] . <exp> [right associative]
    | ( <exp> )
    | <exp> <exp> [left associative]
    | <name>
```

A fun place to start is with booleans, load a file with the following (or evaluate it in the interpreter)

```
def T \a b . a
def F \a b . b
def OR \a b . a T b
def AND \a b . a b F
```

And run some evaluations

```
>>AND T F
>>T a b
>> AND (OR T F) T 
```