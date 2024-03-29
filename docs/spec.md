# Taurine Spec

This document provides a basic outline of how taurine code is written. This project is a big WIP so errors are not going to be obvious...

## Etch statement

The `etch` statement is used to print a list of expressions to the screen.

```
etch "Hello World";
etch "I am ", 22, " years old."; // "I am 22 years old."
```

## Variables

To declare a variable, use the `var` keyword. The syntax is `var (type) symbol = value;`

```
var (num) myAge = 22;
var (str) myName = "Brayden";
etch myName, " is ", myAge, " years old."; // "Brayden is 22 years old."
```

The following data types are allowed:

| Type   | Description           |
|--------|-----------------------|
| `str`  | string                |
| `num`  | floating point number |
| `int`  | integer               |
| `bool` | boolean               |
| `arr`  | array                 |
| `obj`  | object                |

## Read statement

To read a string from stdin, use the `read` statement.

```
var (str) name;
read name, "Enter your name: ";
etch "Hello, ", name; // "Hello, <name>"
```

## Functions

You can declare a function with the `func` keyword. The syntaix is `func (returnType) functionName(type argName, type argName) {}`

```
func (num) factorial(num n) {
  if n == 1 {
    return 1;
  }
  return n * factorial(n - 1);
}
etch factorial(5); // "120.000000"
```

Functions can be assigned to variables.

```
var (func) fn = func (num) (num n) {
  return n * 2;
}
```

## Expression grouping

Expressions are grouped together with parethesis `()`.

```
var (num) x = 3 * 2 + 4;   // 10
var (num) y = 3 * (2 + 4); // 18
```

## Indexing

Accessing an array or string character at a given index can be done with `@`.

```
var (str) myName = "Brayden";
etch myName@1;        // "r"

var (arr) myArr = [10, 20, 30];
etch myArr@2; // "30"
```

## Objects

Javascript-style objects

```
var (obj) myObj = { hello: "world" };
var (obj) myObj = {
  hello: "world",
  x: 3,
}
```

## Dot notation

With the addition of objects will be the addition of dot notation

```
var (obj) myObj = { hello: "world" };
etch myObj.hello;
```

## Functions as expressions

Pass functions as arguments and assign them to variables.

## 

# COMING SOON

The below features have not been implemented yet.

## String interpolation

Interpolate variables in strings

```
var (str) myStr = "world";
var (str) helloWorld = "hello \(myStr)"
```

