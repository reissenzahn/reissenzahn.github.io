---
title: "Kotlin"
date: "2024-08-03"
tags: ["Programming"]
---

https://play.kotlinlang.org

## Overview
- Kotlin is multi-paradigm statically-typed language with a pragmatic, concise and safe design and interoperability with Java.
- Code compiled with the Kotlin compiler depends on the Kotlin runtime library which contains the definitions of the Kotlin standard library classes and extensions to the standard Java APIs.

```kotlin
data class Person(val name: String, val age: Int? = null)

fun main(args: Array<String>) {
  val people = listOf(Person("Alice"), Person("Bob", age = 29));

  val eldest = people.maxBy { it.age ?: 0 }

  println("eldest: $eldest")
}
```

```
# compile and run source code
kotlinc hello.kt -include -runtime -d hello.jar
java -jar hello.jar

# start interactive shell
kotlinc
```


## Variables
```kotlin
// declare an immutable variable
val answer: Int

// initialize a variable
val answer: Int = 42

// the type of a variable can be inferred
val answer = 42

// a val variable must be initialized exactly once during the execution of the block in which it is defined; but it can be initialized to different values
val message: String
if (canPerformOperation()) {
  message = "Success"
} else {
  message = "Failed"
}

// initialize a mutable variable
var answer = 42

// the object that a val reference points to may itself be mutable
val languages = arrayListOf("Java")
languages.add("Kotlin")
```

## Functions
```kotlin
// declare a function with a block body
fun max(a: Int, b: Int): Int {
  return if (a > b) a else b
}

// declare a function with an expression body
fun max(a: Int, b: Int): Int = if (a > b) a else b

// the return type of an expression-body function can be inferred
fun max(a: Int, b: Int) = if (a > b) a else b

println(max(1, 2))
```

## Control Flow
```kotlin
// if is an expression with a result value
val max = if (a > b) a else b
```

## Strings
```kotlin
// string templates allow for referring to local variables in string literals
val name = "Kotlin"
println("Hello, $name!")

// the $ character can be escaped to prevent attempting to interpret a variable reference
println("\$x")

// string templates can also include more complex expressions, including nesting double quotes
println("1 + 1 = ${1 + 1}")
```

## Classes
```kotlin
// declare class with a single property
class Person(val name: String)
```

## Properties
```kotlin
// properties are a first-class language feature which replaces fields and accessor methods
class Person(
  // declare an immutable property (generates a field and a getter)
  val name: String,

  // declare a mutable property (generates a field, getter and setter)
  var isMarried: Boolean
)

val person = Person("Bob", false)
println(person.name)

person.isMarried = true
println(person.isMarried)

// custom property accessors can be defined
class Rectangle(val height: Int, val width: Int) {
  val isSquare: Boolean
    get() {  // property getter declaration
      return height == width
    }
}

val rectangle = Rectangle(41, 43)
println(rectangle.isSquare)
```

## Packages

Kotlin also has the concept of packages, similar to that in Java. Every Kotlin file can have a package statement at the beginning, and all declarations (classes, functions, and properties) defined in the file will be placed in that package. Declarations defined in other files can be used directly if they’re in the same package; they need to be imported if they’re in a different package. As in Java, import statements are placed at the beginning of the file and use the import keyword. Here’s an example of a source file showing the syntax for the package declaration and import statement.

package geometry.shapes
import java.util.Random

class Rectangle(val height: Int, val width: Int) {
val isSquare: Boolean
get() = height == width
}
fun createRandomRectangle(): Rectangle {
val random = Random()
return Rectangle(random.nextInt(), random.nextInt())
}

Kotlin doesn’t make a distinction between importing classes and functions, and it allows you to import any kind of declaration using the import keyword. You can import the top-level function by name.

package geometry.example
import geometry.shapes.createRandomRectangle
fun main(args: Array<String>) {
Listing 2.8 Putting a class and a function declaration in a package
Listing 2.9 Importing the function from another package
Package
declaration
Imports the standard
Java library class
Imports a function
by name
Classes and properties 27
println(createRandomRectangle().isSquare)
}
You can also import all declarations defined in a particular package by putting .*
after the package name. Note that this star import will make visible not only classes
defined in the package, but also top-level functions and properties. In listing 2.9, writing import geometry.shapes.* instead of the explicit import makes the code
compile correctly as well.
 In Java, you put your classes into a structure of files and directories that matches
the package structure. For example, if you have a package named shapes with several
classes, you need to put every class into a separate file with a matching name and store
those files in a directory also called shapes. Figure 2.2 shows how the geometry package and its subpackages could be organized in Java. Assume that the createRandomRectangle function is located in a separate class, RectangleUtil.
Figure 2.2 In Java, the directory hierarchy duplicates the package hierarchy.
In Kotlin, you can put multiple classes in the same file and choose any name for that
file. Kotlin also doesn’t impose any restrictions on the layout of source files on disk;
you can use any directory structure to organize your files. For instance, you can define
all the content of the package geometry.shapes in the file shapes.kt and place this
file in the geometry folder without creating a separate shapes folder (see figure 2.3).
Figure 2.3 Your package hierarchy doesn’t need to follow the directory hierarchy.
In most cases, however, it’s still a good practice to follow Java’s directory layout and to
organize source files into directories according to the package structure. Sticking to
that structure is especially important in projects where Kotlin is mixed with Java,
Prints “true”
incredibly rarely
geometry.example package
geometry.shapes package
Rectangle class
geometry.example package
geometry.shapes package
28 CHAPTER 2 Kotlin basics
because doing so lets you migrate the code gradually without introducing any surprises. But you shouldn’t hesitate to pull multiple classes into the same file, especially
if the classes are small (and in Kotlin, they often are).
 Now you know how programs are structured. Let’s move on with learning basic
concepts and look at control structures in Kotlin.
2.3 Representing and handling choices: enums and “when”
In this section, we’re going to talk about the when construct. It can be thought of as a
replacement for the switch construct in Java, but it’s more powerful and is used
more often. Along the way, we’ll give you an example of declaring enums in Kotlin
and discuss the concept of smart casts.
2.3.1 Declaring enum classes
Let’s start by adding some imaginary bright pictures to this serious book and looking
at an enum of colors.
enum class Color {
RED, ORANGE, YELLOW, GREEN, BLUE, INDIGO, VIOLET
}
This is a rare case when a Kotlin declaration uses more keywords than the corresponding Java one: enum class versus just enum in Java. In Kotlin, enum is a so-called soft
keyword: it has a special meaning when it comes before class, but you can use it as a
regular name in other places. On the other hand, class is still a keyword, and you’ll
continue to declare variables named clazz or aClass.
 Just as in Java, enums aren’t lists of values: you can declare properties and methods
on enum classes. Here’s how it works.
enum class Color(
val r: Int, val g: Int, val b: Int
) {
RED(255, 0, 0), ORANGE(255, 165, 0),
YELLOW(255, 255, 0), GREEN(0, 255, 0), BLUE(0, 0, 255),
INDIGO(75, 0, 130), VIOLET(238, 130, 238);
fun rgb() = (r * 256 + g) * 256 + b
}
>>> println(Color.BLUE.rgb())
255
Enum constants use the same constructor and property declaration syntax as you saw
earlier for regular classes. When you declare each enum constant, you need to provide
Listing 2.10 Declaring a simple enum class
Listing 2.11 Declaring an enum class with properties
Declares properties
of enum constants
Specifies
property
values
when each
constant is
created
The semicolon
here is required.
Defines a method
on the enum class
Representing and handling choices: enums and “when” 29
the property values for that constant. Note that this example shows the only place in the
Kotlin syntax where you’re required to use semicolons: if you define any methods in the
enum class, the semicolon separates the enum constant list from the method definitions. Now let’s see some cool ways to deal with enum constants in your code.
2.3.2 Using “when” to deal with enum classes
Do you remember how children use mnemonic phrases to memorize the colors of the
rainbow? Here’s one: “Richard Of York Gave Battle In Vain!” Imagine you need a
function that gives you a mnemonic for each color (and you don’t want to store this
information in the enum itself). In Java, you can use a switch statement for this. The
corresponding Kotlin construct is when.
 Like if, when is an expression that returns a value, so you can write a function
with an expression body, returning the when expression directly. When we talked
about functions at the beginning of the chapter, we promised an example of a multiline function with an expression body. Here’s such an example.
fun getMnemonic(color: Color) =
when (color) {
Color.RED -> "Richard"
Color.ORANGE -> "Of"
Color.YELLOW -> "York"
Color.GREEN -> "Gave"
Color.BLUE -> "Battle"
Color.INDIGO -> "In"
Color.VIOLET -> "Vain"
}
>>> println(getMnemonic(Color.BLUE))
Battle
The code finds the branch corresponding to the passed color value. Unlike in Java,
you don’t need to write break statements in each branch (a missing break is often a
cause for bugs in Java code). If a match is successful, only the corresponding branch is
executed. You can also combine multiple values in the same branch if you separate
them with commas.
fun getWarmth(color: Color) = when(color) {
Color.RED, Color.ORANGE, Color.YELLOW -> "warm"
Color.GREEN -> "neutral"
Color.BLUE, Color.INDIGO, Color.VIOLET -> "cold"
}
>>> println(getWarmth(Color.ORANGE))
warm
Listing 2.12 Using when for choosing the right enum value
Listing 2.13 Combining options in one when branch
Returns a “when”
expression directly
Returns the corresponding
string if the color equals
the enum constant
30 CHAPTER 2 Kotlin basics
These examples use enum constants by their full name, specifying the Color enum
class name. You can simplify the code by importing the constant values.
import ch02.colors.Color
import ch02.colors.Color.*
fun getWarmth(color: Color) = when(color) {
RED, ORANGE, YELLOW -> "warm"
GREEN -> "neutral"
BLUE, INDIGO, VIOLET -> "cold"
}
2.3.3 Using “when” with arbitrary objects
The when construct in Kotlin is more powerful than Java’s switch. Unlike switch,
which requires you to use constants (enum constants, strings, or number literals) as
branch conditions, when allows any objects. Let’s write a function that mixes two colors if they can be mixed in this small palette. You don’t have lots of options, and you
can easily enumerate them all.
fun mix(c1: Color, c2: Color) =
when (setOf(c1, c2)) {
setOf(RED, YELLOW) -> ORANGE
setOf(YELLOW, BLUE) -> GREEN
setOf(BLUE, VIOLET) -> INDIGO
else -> throw Exception("Dirty color")
}
>>> println(mix(BLUE, YELLOW))
GREEN
If colors c1 and c2 are RED and YELLOW (or vice versa), the result of mixing them is
ORANGE, and so on. To implement this, you use set comparison. The Kotlin standard
library contains a function setOf that creates a Set containing the objects specified
as its arguments. A set is a collection for which the order of items doesn’t matter; two
sets are equal if they contain the same items. Thus, if the sets setOf(c1, c2) and
setOf(RED, YELLOW) are equal, it means either c1 is RED and c2 is YELLOW, or vice
versa. This is exactly what you want to check.
 The when expression matches its argument against all branches in order until
some branch condition is satisfied. Thus setOf(c1, c2) is checked for equality: first
with setOf(RED, YELLOW) and then with other sets of colors, one after another. If
none of the other branch conditions is satisfied, the else branch is evaluated.
 Being able to use any expression as a when branch condition lets you write concise
and beautiful code in many cases. In this example, the condition is an equality check;
next you’ll see how the condition may be any Boolean expression.
Listing 2.14 Importing enum constants to access without qualifier
Listing 2.15 Using different objects in when branches
Imports the Color
class declared
in another package Explicitly imports
enum constants to
use them by names Uses imported
constants
by name
An argument of the “when” expression
can be any object. It’s checked for
equality with the branch conditions.
Enumerates pairs
of colors that
can be mixed
Executed if none of
the other branches
were matched
Representing and handling choices: enums and “when” 31
2.3.4 Using “when” without an argument
You may have noticed that listing 2.15 is somewhat inefficient. Every time you call this
function, it creates several Set instances that are used only to check whether two
given colors match the other two colors. Normally this isn’t an issue, but if the function is called often, it’s worth rewriting the code in a different way to avoid creating
garbage. You can do it by using the when expression without an argument. The code is
less readable, but that’s the price you often have to pay to achieve better performance.
fun mixOptimized(c1: Color, c2: Color) =
when {
(c1 == RED && c2 == YELLOW) ||
(c1 == YELLOW && c2 == RED) ->
ORANGE
(c1 == YELLOW && c2 == BLUE) ||
(c1 == BLUE && c2 == YELLOW) ->
GREEN
(c1 == BLUE && c2 == VIOLET) ||
(c1 == VIOLET && c2 == BLUE) ->
INDIGO
else -> throw Exception("Dirty color")
}
>>> println(mixOptimized(BLUE, YELLOW))
GREEN
If no argument is supplied for the when expression, the branch condition is any Boolean expression. The mixOptimized function does the same thing as mix did earlier.
Its advantage is that it doesn’t create any extra objects, but the cost is that it’s harder
to read.
 Let’s move on and look at examples of the when construct in which smart casts
come into play.
2.3.5 Smart casts: combining type checks and casts
As the example for this section, you’ll write a function that evaluates simple arithmetic
expressions like (1 + 2) + 4. The expressions will contain only one type of operation:
the sum of two numbers. Other arithmetic operations (subtraction, multiplication,
division) can be implemented in a similar way, and you can do that as an exercise.
 First, how do you encode the expressions? You store them in a tree-like structure,
where each node is either a sum (Sum) or a number (Num). Num is always a leaf node,
whereas a Sum node has two children: the arguments of the sum operation. The following listing shows a simple structure of classes used to encode the expressions: an
interface called Expr and two classes, Num and Sum, that implement it. Note that the
Expr interface doesn’t declare any methods; it’s used as a marker interface to provide
Listing 2.16 when without an argument
No argument
for “when”
32 CHAPTER 2 Kotlin basics
a common type for different kinds of expressions. To mark that a class implements an
interface, you use a colon (:) followed by the interface name:
interface Expr
class Num(val value: Int) : Expr
class Sum(val left: Expr, val right: Expr) : Expr
Sum stores the references to left and right arguments of type Expr; in this small example, they can be
either Num or Sum. To store the expression (1 + 2) + 4
mentioned earlier, you create an object
Sum(Sum(Num(1), Num(2)), Num(4)). Figure 2.4
shows its tree representation.
 Now let’s look at how to compute the value of an
expression. Evaluating the example expression should
return 7:
>>> println (eval(Sum(Sum(Num(1), Num(2)), Num (4))))
7
The Expr interface has two implementations, so you have to try two options in order
to evaluate a result value for an expression:
 If an expression is a number, you return the corresponding value.
 If it’s a sum, you have to evaluate the left and right expressions and return
their sum.
First we’ll look at this function written in the normal Java way, and then we’ll refactor
it to be written in a Kotlin style. In Java, you’d probably use a sequence of if statements to check the options, so let’s use the same approach in Kotlin.
fun eval(e: Expr): Int {
if (e is Num) {
val n = e as Num
return n.value
}
if (e is Sum) {
return eval(e.right) + eval(e.left)
}
throw IllegalArgumentException("Unknown expression")
}
Listing 2.17 Expression class hierarchy
Listing 2.18 Evaluating expressions with an if-cascade
Simple value object class with
one property, value,
implementing the Expr interface
The argument of a Sum
operation can be any Expr:
either Num or another Sum
Sum
Sum
Num(1)
Num(4)
Num(2)
Figure 2.4 A representation
of the expression Sum(SumNum(1), Num(2)), Num(4))
This explicit cast to
Num is redundant.
The variable e
is smart-cast.
Representing and handling choices: enums and “when” 33
>>> println(eval(Sum(Sum(Num(1), Num(2)), Num(4))))
7
In Kotlin, you check whether a variable is of a certain type by using an is check. If
you’ve programmed in C#, this notation should be familiar. The is check is similar to
instanceof in Java. But in Java, if you’ve checked that a variable has a certain type
and needs to access members of that type, you need to add an explicit cast following
the instanceof check. When the initial variable is used more than once, you often
store the cast result in a separate variable. In Kotlin, the compiler does this job for
you. If you check the variable for a certain type, you don’t need to cast it afterward;
you can use it as having the type you checked for. In effect, the compiler performs the
cast for you, and we call it a smart cast.
 In the eval function, after you check
whether the variable e has Num type, the
compiler interprets it as a Num variable.
You can then access the value property of
Num without an explicit cast: e.value.
The same goes for the right and left
properties of Sum: you write only e.right and e.left in the corresponding context.
In the IDE, these smart-cast values are emphasized with a background color, so it’s easy
to grasp that this value was checked beforehand. See figure 2.5.
 The smart cast works only if a variable couldn’t have changed after the is check.
When you’re using a smart cast with a property of a class, as in this example, the property has to be a val and it can’t have a custom accessor. Otherwise, it would not be
possible to verify that every access to the property would return the same value.
 An explicit cast to the specific type is expressed via the as keyword:
val n = e as Num
Now let’s look at how to refactor the eval function into a more idiomatic Kotlin style.
2.3.6 Refactoring: replacing “if” with “when”
How does if in Kotlin differ from if in Java? You have seen the difference already. At
the beginning of the chapter, you saw the if expression used in the context where Java
would have a ternary operator: if (a > b) a else b works like Java’s a > b ? a : b. In
Kotlin, there is no ternary operator, because, unlike in Java, the if expression returns
a value. That means you can rewrite the eval function to use the expression-body syntax, removing the return statement and the curly braces and using the if expression
as the function body instead.
fun eval(e: Expr): Int =
if (e is Num) {
e.value
} else if (e is Sum) {
Listing 2.19 Using if-expressions that return values
if (e is Sum) {
 return eval(e.right) + eval(e.left)
}
Figure 2.5 The IDE highlights smart casts with
a background color.
34 CHAPTER 2 Kotlin basics
eval(e.right) + eval(e.left)
} else {
throw IllegalArgumentException("Unknown expression")
}
>>> println(eval(Sum(Num(1), Num(2))))
3
The curly braces are optional if there’s only one expression in an if branch. If an if
branch is a block, the last expression is returned as a result.
 Let’s polish this code even more and rewrite it using when.
fun eval(e: Expr): Int =
when (e) {
is Num ->
e.value
is Sum ->
eval(e.right) + eval(e.left)
else ->
throw IllegalArgumentException("Unknown expression")
}
The when expression isn’t restricted to checking values for equality, which is what you
saw earlier. Here you use a different form of when branches, allowing you to check
the type of the when argument value. Just as in the if example in listing 2.19, the
type check applies a smart cast, so you can access members of Num and Sum without
extra casts.
 Compare the last two Kotlin versions of the eval function, and think about how
you can apply when as a replacement for sequences of if expressions in your own
code as well. When the branch logic is complicated, you can use a block expression as
a branch body. Let’s see how this works.
2.3.7 Blocks as branches of “if” and “when”
Both if and when can have blocks as branches. In this case, the last expression in the
block is the result. If you want to add logging to the example function, you can do so
in the block and return the last value as before.
fun evalWithLogging(e: Expr): Int =
when (e) {
is Num -> {
println("num: ${e.value}")
e.value
}
is Sum -> {
val left = evalWithLogging(e.left)
Listing 2.20 Using when instead of if-cascade
Listing 2.21 Using when with compound actions in branches
“when” branches that
Smart casts are check the argument type
applied here.
This is the last expression
in the block and is returned
if e is of type Num.
Iterating over things: “while” and “for” loops 35
val right = evalWithLogging(e.right)
println("sum: $left + $right")
left + right
}
else -> throw IllegalArgumentException("Unknown expression")
}
Now you can look at the logs printed by the evalWithLogging function and follow
the order of computation:
>>> println(evalWithLogging(Sum(Sum(Num(1), Num(2)), Num(4))))
num: 1
num: 2
sum: 1 + 2
num: 4
sum: 3 + 4
7
The rule “the last expression in a block is the result” holds in all cases where a block
can be used and a result is expected. As you’ll see at the end of this chapter, the same
rule works for the try body and catch clauses, and chapter 5 discusses its application
to lambda expressions. But as we mentioned in section 2.2, this rule doesn’t hold for
regular functions. A function can have either an expression body that can’t be a block
or a block body with explicit return statements inside.
 You’ve become acquainted with Kotlin ways to choose the right things among
many. Now it’s a good time to see how you can iterate over things.
2.4 Iterating over things: “while” and “for” loops
Of all the features discussed in this chapter, iteration in Kotlin is probably the most
similar to Java. The while loop is identical to the one in Java, so it deserves only a
brief mention in the beginning of this section. The for loop exists in only one form,
which is equivalent to Java’s for-each loop. It’s written for <item> in <elements>, as in C#. The most common application of this loop is iterating over collections, just as in Java. We’ll explore how it can cover other looping scenarios as well.
2.4.1 The “while” loop
Kotlin has while and do-while loops, and their syntax doesn’t differ from the corresponding loops in Java:
while (condition) {
/*...*/
}
do {
/*...*/
} while (condition)
Kotlin doesn’t bring anything new to these simple loops, so we won’t linger. Let’s
move on to discuss the various uses of the for loop.
This expression is returned
if e is of type Sum.
The body is executed while
the condition is true.
The body is executed for the first time
unconditionally. After that, it’s
executed while the condition is true.
36 CHAPTER 2 Kotlin basics
2.4.2 Iterating over numbers: ranges and progressions
As we just mentioned, in Kotlin there’s no regular Java for loop, where you initialize a
variable, update its value on every step through the loop, and exit the loop when the
value reaches a certain bound. To replace the most common use cases of such loops,
Kotlin uses the concepts of ranges.
 A range is essentially just an interval between two values, usually numbers: a start
and an end. You write it using the .. operator:
val oneToTen = 1..10
Note that ranges in Kotlin are closed or inclusive, meaning the second value is always a
part of the range.
 The most basic thing you can do with integer ranges is loop over all the values. If
you can iterate over all the values in a range, such a range is called a progression.
 Let’s use integer ranges to play the Fizz-Buzz game. It’s a nice way to survive a long
trip in a car and remember your forgotten division skills. Players take turns counting
incrementally, replacing any number divisible by three with the word fizz and any
number divisible by five with the word buzz. If a number is a multiple of both three
and five, you say “FizzBuzz.”
 The following listing prints the right answers for the numbers from 1 to 100. Note
how you check the possible conditions with a when expression without an argument.
fun fizzBuzz(i: Int) = when {
i % 15 == 0 -> "FizzBuzz "
i % 3 == 0 -> "Fizz "
i % 5 == 0 -> "Buzz "
else -> "$i "
}
>>> for (i in 1..100) {
... print(fizzBuzz(i))
... }
}
1 2 Fizz 4 Buzz Fizz 7 ...
Suppose you get tired of these rules after an hour of driving and want to complicate
things a bit. Let’s start counting backward from 100 and include only even numbers.
>>> for (i in 100 downTo 1 step 2) {
... print(fizzBuzz(i))
... }
Buzz 98 Fizz 94 92 FizzBuzz 88 ...
Listing 2.22 Using when to implement the Fizz-Buzz game
Listing 2.23 Iterating over a range with a step
If i is divisible by 15,
returns FizzBuzz. As
in Java, % is the
modulus operator.
If i is divisible
by 3,
returns Fizz If i is divisible
by 5, returns Buzz
Else returns the
number itself Iterates over the
integer range 1..100
Iterating over things: “while” and “for” loops 37
Now you’re iterating over a progression that has a step, which allows it to skip some
numbers. The step can also be negative, in which case the progression goes backward
rather than forward. In this example, 100 downTo 1 is a progression that goes backward (with step -1). Then step changes the absolute value of the step to 2 while keeping the direction (in effect, setting the step to -2).
 As we mentioned earlier, the .. syntax always creates a range that includes the end
point (the value to the right of ..). In many cases, it’s more convenient to iterate over
half-closed ranges, which don’t include the specified end point. To create such a
range, use the until function. For example, the loop for (x in 0 until size)
is equivalent to for (x in 0..size-1), but it expresses the idea somewhat more
clearly. Later, in section 3.4.3, you’ll learn more about the syntax for downTo, step,
and until in these examples.
 You can see how working with ranges and progressions helped you cope with the
advanced rules for the FizzBuzz game. Now let’s look at other examples that use the
for loop.
2.4.3 Iterating over maps
We’ve mentioned that the most common scenario of using a for ... in loop is iterating over a collection. This works exactly as it does in Java, so we won’t say much about
it. Let’s see how you can iterate over a map, instead.
 As an example, we’ll look at a small program that prints binary representations for
characters. You’ll store these binary representations in a map (just for illustrative purposes). The following code creates a map, fills it with binary representations of some
letters, and then prints the map’s contents.
val binaryReps = TreeMap<Char, String>()
for (c in 'A'..'F') {
val binary = Integer.toBinaryString(c.toInt())
binaryReps[c] = binary
}
for ((letter, binary) in binaryReps) {
println("$letter = $binary")
}
The .. syntax to create a range works not only for numbers, but also for characters.
Here you use it to iterate over all characters from A up to and including F.
 Listing 2.24 shows that the for loop allows you to unpack an element of a collection you’re iterating over (in this case, a collection of key/value pairs in the map). You
store the result of the unpacking in two separate variables: letter receives the key,
Listing 2.24 Initializing and iterating over a map
Uses TreeMap so
the keys are sorted
Iterates over the
characters from A to F
using a range of characters
Converts
ASCII code
to binary
Stores the value in a
map by the c key
Iterates over a map,
assigning the map key and
value to two variables
38 CHAPTER 2 Kotlin basics
and binary receives the value. Later, in section 7.4.1, you’ll find out more about this
unpacking syntax.
 Another nice trick used in listing 2.24 is the shorthand syntax for getting and
updating the values of a map by key. Instead of calling get and put, you can use
map[key] to read values and map[key] = value to set them. The code
binaryReps[c] = binary
is equivalent to its Java version:
binaryReps.put(c, binary)
The output is similar to the following (we’ve arranged it in two columns instead of one):
A = 1000001 D = 1000100
B = 1000010 E = 1000101
C = 1000011 F = 1000110
You can use the same unpacking syntax to iterate over a collection while keeping track
of the index of the current item. You don’t need to create a separate variable to store
the index and increment it by hand:
val list = arrayListOf("10", "11", "1001")
for ((index, element) in list.withIndex()) {
println("$index: $element")
}
The code prints what you expect:
0: 10
1: 11
2: 1001
We’ll dig into the whereabouts of withIndex in the next chapter.
 You’ve seen how you can use the in keyword to iterate over a range or a collection.
You can also use in to check whether a value belongs to the range or collection.
2.4.4 Using “in” to check collection and range membership
You use the in operator to check whether a value is in a range, or its opposite, !in, to
check whether a value isn’t in a range. Here’s how you can use in to check whether a
character belongs to a range of characters.
fun isLetter(c: Char) = c in 'a'..'z' || c in 'A'..'Z'
fun isNotDigit(c: Char) = c !in '0'..'9'
>>> println(isLetter('q'))
true
>>> println(isNotDigit('x'))
true
Listing 2.25 Checking range membership using in
Iterates over a collection
with an index
Exceptions in Kotlin 39
This technique for checking whether a character is a letter looks simple. Under the
hood, nothing tricky happens: you still check that the character’s code is somewhere
between the code of the first letter and the code of the last one. But this logic is concisely hidden in the implementation of the range classes in the standard library:
c in 'a'..'z'
The in and !in operators also work in when expressions.
fun recognize(c: Char) = when (c) {
in '0'..'9' -> "It's a digit!"
in 'a'..'z', in 'A'..'Z' -> "It's a letter!"
else -> "I don't know…"
}
>>> println(recognize('8'))
It's a digit!
Ranges aren’t restricted to characters, either. If you have any class that supports comparing instances (by implementing the java.lang.Comparable interface), you can
create ranges of objects of that type. If you have such a range, you can’t enumerate all
objects in the range. Think about it: can you, for example, enumerate all strings
between “Java” and “Kotlin”? No, you can’t. But you can still check whether another
object belongs to the range, using the in operator:
>>> println("Kotlin" in "Java".."Scala")
true
Note that the strings are compared alphabetically here, because that’s how the
String class implements the Comparable interface.
 The same in check works with collections as well:
>>> println("Kotlin" in setOf("Java", "Scala"))
false
Later, in section 7.3.2, you’ll see how to use ranges and progressions with your own
data types and what objects in general you can use in checks with.
 There’s one more group of Java statements we want to look at in this chapter: statements for dealing with exceptions.
2.5 Exceptions in Kotlin
Exception handling in Kotlin is similar to the way it’s done in Java and many other languages. A function can complete in a normal way or throw an exception if an error
occurs. The function caller can catch this exception and process it; if it doesn’t, the
exception is rethrown further up the stack.
Listing 2.26 Using in checks as when branches
Transforms to
 a <= c && c <= z
Checks whether the value is
You can in the range from 0 to 9
combine
multiple
ranges.
The same as “Java” <= “Kotlin”
&& “Kotlin” <= “Scala”
This set doesn’t contain
the string “Kotlin”.
40 CHAPTER 2 Kotlin basics
 The basic form for exception-handling statements in Kotlin is similar to Java’s. You
throw an exception in a non-surprising manner:
if (percentage !in 0..100) {
throw IllegalArgumentException(
"A percentage value must be between 0 and 100: $percentage")
}
As with all other classes, you don’t have to use the new keyword to create an instance
of the exception.
 Unlike in Java, in Kotlin the throw construct is an expression and can be used as a
part of other expressions:
val percentage =
if (number in 0..100)
number
else
throw IllegalArgumentException(
"A percentage value must be between 0 and 100: $number")
In this example, if the condition is satisfied, the program behaves correctly, and the
percentage variable is initialized with number. Otherwise, an exception is thrown,
and the variable isn’t initialized. We’ll discuss the technical details of throw as a part
of other expressions, in section 6.2.6.
2.5.1 “try”, “catch”, and “finally”
Just as in Java, you use the try construct with catch and finally clauses to handle
exceptions. You can see it in the following listing, which reads a line from the given
file, tries to parse it as a number, and returns either the number or null if the line
isn’t a valid number.
fun readNumber(reader: BufferedReader): Int? {
try {
val line = reader.readLine()
return Integer.parseInt(line)
}
catch (e: NumberFormatException) {
return null
}
finally {
reader.close()
}
}
>>> val reader = BufferedReader(StringReader("239"))
>>> println(readNumber(reader))
239
Listing 2.27 Using try as in Java
“throw” is an
expression.
You don’t have to explicitly
specify exceptions that can be
thrown from this function.
The exception type
is on the right.
“finally” works just
as it does in Java.
Exceptions in Kotlin 41
The biggest difference from Java is that the throws clause isn’t present in the code: if
you wrote this function in Java, you’d explicitly write throws IOException after the
function declaration. You’d need to do this because IOException is a checked exception. In Java, it’s an exception that needs to be handled explicitly. You have to declare
all checked exceptions that your function can throw, and if you call another function,
you need to handle its checked exceptions or declare that your function can throw
them, too.
 Just like many other modern JVM languages, Kotlin doesn’t differentiate between
checked and unchecked exceptions. You don’t specify the exceptions thrown by a
function, and you may or may not handle any exceptions. This design decision is
based on the practice of using checked exceptions in Java. Experience has shown that
the Java rules often require a lot of meaningless code to rethrow or ignore exceptions,
and the rules don’t consistently protect you from the errors that can happen.
 For example, in listing 2.27, NumberFormatException isn’t a checked exception.
Therefore, the Java compiler doesn’t force you to catch it, and you can easily see the
exception happen at runtime. This is unfortunate, because invalid input data is a common situation and should be handled gracefully. At the same time, the BufferedReader.close method can throw an IOException, which is a checked exception
and needs to be handled. Most programs can’t take any meaningful action if closing a
stream fails, so the code required to catch the exception from the close method is
boilerplate.
 What about Java 7’s try-with-resources? Kotlin doesn’t have any special syntax
for this; it’s implemented as a library function. In section 8.2.5, you’ll see how this is
possible.
2.5.2 “try” as an expression
To see another significant difference between Java and Kotlin, let’s modify the example a little. Let’s remove the finally section (because you’ve already seen how this
works) and add some code to print the number you read from the file.
fun readNumber(reader: BufferedReader) {
val number = try {
Integer.parseInt(reader.readLine())
} catch (e: NumberFormatException) {
return
}
println(number)
}
>>> val reader = BufferedReader(StringReader("not a number"))
>>> readNumber(reader)
Listing 2.28 Using try as an expression
Becomes the value of
the “try” expression
Nothing
is printed.
42 CHAPTER 2 Kotlin basics
The try keyword in Kotlin, just like if and when, introduces an expression, and you
can assign its value to a variable. Unlike with if, you always need to enclose the statement body in curly braces. Just as in other statements, if the body contains multiple
expressions, the value of the try expression as a whole is the value of the last expression.
 This example puts a return statement in the catch block, so the execution of the
function doesn’t continue after the catch block. If you want to continue execution,
the catch clause also needs to have a value, which will be the value of the last expression in it. Here’s how this works.
fun readNumber(reader: BufferedReader) {
val number = try {
Integer.parseInt(reader.readLine())
} catch (e: NumberFormatException) {
null
}
println(number)
}
>>> val reader = BufferedReader(StringReader("not a number"))
>>> readNumber(reader)
null
If the execution of a try code block behaves normally, the last expression in the block
is the result. If an exception is caught, the last expression in a corresponding catch
block is the result. In listing 2.29, the result value is null if a NumberFormatException is caught.
 At this point, if you’re impatient, you can start writing programs in Kotlin in a way
that’s similar to how you code in Java. As you read this book, you’ll continue to learn how
to change your habitual ways of thinking and use the full power of the new language.
2.6 Summary
 The fun keyword is used to declare a function. The val and var keywords
declare read-only and mutable variables, respectively.
 String templates help you avoid noisy string concatenation. Prefix a variable
name with $ or surround an expression with ${ } to have its value injected into
the string.
 Value-object classes are expressed in a concise way in Kotlin.
 The familiar if is now an expression with a return value.
 The when expression is analogous to switch in Java but is more powerful.
 You don’t have to cast a variable explicitly after checking that it has a certain
type: the compiler casts it for you automatically using a smart cast.
Listing 2.29 Returning a value in catch
This value is used when
no exception happens.
The null value is used
in case of an exception.
An exception is
thrown, so the function
prints “null”.
Summary 43
 The for, while, and do-while loops are similar to their counterparts in Java,
but the for loop is now more convenient, especially when you need to iterate
over a map or a collection with an index.
 The concise syntax 1..5 creates a range. Ranges and progressions allow Kotlin
to use a uniform syntax and set of abstractions in for loops and also work with
the in and !in operators that check whether a value belongs to a range.
 Exception handling in Kotlin is very similar to that in Java, except that Kotlin
doesn’t require you to declare the exceptions that can be thrown by a function.



