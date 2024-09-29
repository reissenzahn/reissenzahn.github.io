---
title: "JavaScript"
date: "2024-08-03"
tags: ["Programming"]
---

## Hello World
```html
<!DOCTYPE html>
<html>
<body>
  <script>
    console.log('Hello, JavaScript!')
  </script>

  <script src="hello.js"></script>
</body>
</html>
```

## Strict Mode
```js
// enable strict mode for a script (may only be preceded by comments)
'use strict';

// enable strict mode for a function
function foo() {
  'use strict';
  // ...
}

// some modern feature (like classes) enable strict mode implicitly
```

## Variables
```js
// declare a variable
let message = 'hello';

// declare multiple variables
let user = 'John', age = 25, message = 'Hello';

// declare a constant (cannot be re-assigned)
const limit = 10;

// an assignment evaluates to the assigned value
let x;
(x = 1) === 1;

// chained assignment
let a, b, c;
a = b = c = 1 + 2;
```

## Data Types
```js
// number (integer and floating-point)
const x = 123;
const y = 3.14;

// special numeric values
1/0 === Infinity
-1/0 === -Infinity
'three' * 2 === NaN

// BigInt represents an integer of arbitrary length
const m = 1234567890123456789012345678901234567890n;

// string
const s = 'hello';
const t = "hello";

// template literals
`${s}, ${1 + 2}!`;

// boolean
const y = true;
const n = 4 < 1;

// null represents "empty" or "unknown"
const u = null;

// undefined represents "value not assigned"
let v;
const w = undefined;

// symbol represents a unique identifier
const i = Symbol('id');

// typeof evaluates to the type name of its operand
typeof true // "boolean"
typeof null // "object" (misleading, kept for compatibility)
typeof alert // "function" (functions are objects but this can be convenient)

// explicit type conversions
String(true);
Number('123');
Boolean(1);

// '', 0, null, undefined and NaN are falsy while other values are truthy
Boolean(NaN);
```

## Operators
```js
// exponentiation
2 ** 3 === 8;

// string concatenation
'hello ' + 'world';

// other arithmetic operators perform numeric conversion automatically
'6' / '3' === 2;

// unary + converts its operand to a number
+true === 1;

// modify-and-assign operators
let n = 2;
n += 2;

// prefix increment evaluates to i + 1
++i;

// postfix increment evaluates to i
i++;

// the comma operator evaluates to the result of the last expression
(1 + 2, 3 + 4, 5 + 6) === 11;

// strings are compared in lexicographic order
'hello' > 'bye' === true;

// equality performs unexpected conversions
(0 == false) === true;

// strict equality does not perform type conversions
(0 === false) === false;

// a chain of || returns the first truthy value or the last one if no truthy value is found
null || 0 || 1 || undefined === 1;
undefined || null || 0 === 0;

// short-circuit evaluation occurs as || processes its arguments until the first truthy value is reached
x % 2 == 0 || alert('odd');

// a chain of && returns the first falsy value or the last value if no falsy value is found
1 && 2 && 0 && 3 === 0;
1 && 2 && 3 === 3;

// short-circuit evaluation occurs as && processes its arguments until the first falsy value is reached
x % 2 == 0 && alert('even');

// a double !! can be ued to convert a value to a boolean
!!null === false;

// nullish coalescing operator evaluates to its first operand it is not null/undefined and its second operand otherwise
firstName ?? lastName ?? nickName ?? 'Anonymous';
```

## Conditionals
```js
// if evaluates a condition expression and converts the result to a boolean
if (x > 0) {
  console.log('+');
} else if (x < 0) {
  console.log('-');
} else {
  console.log('0');
}

// these values are falsy; all other values are truthy
if (0 || '' || null || undefined || NaN) {
  console.log('never');
}

// conditional operator
const parity = x % 2 === 0 ? 'even' : 'odd';
```

## Looping
```js
let i = 0;
while (i < 3) {
  console.log(i);
  i++;
}

let i = 0;
do {
  console.log(i);
  i++;
} while (i < 3);

for (let i = 0; i < 3; i++) {
  console.log(i);
}

let i = 0;
while (true) {
  console.log(i);
  if (i >= 3) break;
}

for (let i = 0; i < 6; i++) {
  if (i % 2 == 0) continue;
  console.log(i);
}

// labels can be used to break out of nested loops
outer:
for (let i = 0; i < 3; i++) {
  for (let j = 0; j < 3; j++) {
    break outer;
  }
}
```


## Promises
```js
// callback-based style of asynchronous programming
function loadScript(src, callback) {
  let script = document.createElement('script');
  script.src = src;

  script.onload = () => callback(null, script);
  script.onerror = () => callback(new Error(`load error: "${src}"`));

  document.head.append(script);
}

loadScript('./script.js', (error, script) => {
  if (error) {
    // handle error
  } else {
    // script loaded successfully
  }
});

// create a promise
const promise = new Promise((resolve, reject) => {
  // the executor function is called when the promise is created

  // 
  setTimeout(() => resolve('done'), 1000);

  setTimeout(() => reject(new Error('oops')), 1000);
});

promise.then(
  (result) => alert(result),  // called when the promise is resolved
  (error) => alert(error),    // called when the promise is rejected
)

// only handle resolved promise
promise.then(alert);

// only handle rejected promise
promise.catch(alert)

// always called when the promise settles
promise.finally(() => console.log('cleaning up...'));

// handlers can be attached to settled promises in which case they run immediately
const promise = new Promise(resolve => resolve('done'));
promise.then(alert);

// a more involved example
function loadScript(src) {
  return new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.src = src;
    script.onload = () => resolve(script);
    script.onerror = () => reject(new Error(`load error: "${src}"`));
    document.head.append(script);
  });
}

// .then() always returns a promise and the rest of the chain waits until it settles; if a value is returned or an error is thrown then the promise settles with the value or error

// handlers can be chained together; every call to .then() returns a new promise
new Promise((resolve, reject) => {
  setTimeout(() => resolve(1), 1000);
}).then((result) => {
  alert(result);
  return result * 2;  // create a new promise resolved with the result 2
}).then((result) => {
  alert(result);
  return result * 2;
}).then((result) => {
  alert(result);
  return result * 2;
});

// a handler may also create and return a promise in which case further handlers wait until it settles
new Promise((resolve, reject) => {
  setTimeout(() => resolve(1), 1000);
}).then((result) => {
  alert(result);
  return new Promise((resolve, reject) => {
    setTimeout(() => resolve(result * 2), 1000);
  });
}).then((result) => {
  alert(result);
});

// to be precise, a handler may return any "thenable" object with a .then() method
class Thenable {
  constructor(num) {
    this.num = num;
  }

  // the .then() method is passed resolve and reject arguments similar to an executor; the result is passed further down the chain
  then(resolve, reject) {
    alert(resolve);
    setTimeout(() => resolve(this.num * 2), 1000);
  }
}

new Promise(resolve => resolve(1))
  .then(result => {
    return new Thenable(result);
  })
  .then(alert);

// fetch returns a promise that resolves with a response object
fetch('/article/promise-chaining/user.json')
  .then((response) => {
    // response.text() returns a new promise that resolves with the full response text when it loads
    return response.text();
  })
  .then((text) => {
    alert(text);
  });

// when a promise rejects, the control jumps to the closest rejection handler
fetch('https://no-such-server.blabla')
  .then(response => response.json())
  .catch(err => alert(err))

// if an error is thrown in an executor or handler function then it will be caught and treated as a rejection
new Promise((resolve, reject) => {
  throw new Error("Whoops!");
}).catch(alert);

new Promise((resolve, reject) => {
  resolve("ok");
}).then((result) => {
  throw new Error("Whoops!");
}).catch(alert);

// if an error is handled in a .catch() then execution continues to the next closest .then() handler
new Promise((resolve, reject) => {
  throw new Error("Whoops!");
}).catch(function(error) {
  alert("The error is handled, continue normally");
}).then(() => alert("Next successful handler runs"));

// a .catch() can also re-throw an exception for it to be handled by the next closest .catch()
new Promise((resolve, reject) => {
  throw new Error("Whoops!");
}).catch(function(error) { // (*)
  if (error instanceof URIError) {
    // handle it
  } else {
    alert("Can't handle such error");
    throw error; // throwing this or another error jumps to the next catch
  }
}).then(function() {
  /* doesn't run here */
}).catch(error => { // (**)
  alert(`The unknown error has occurred: ${error}`);
  // don't return anything => execution goes the normal way
});

// an unhandled rejection
new Promise(function() {
  noSuchFunction();
}).then(() => {});

// catch unhandled rejections
window.addEventListener('unhandledrejection', function(event) {
  alert(event.promise); // [object Promise] - the promise that generated the error
  alert(event.reason); // Error: Whoops! - the unhandled error object
});

// take an iterable of promises and return a new promise that resolves with an array of their results; the order of the resulting array is the same as in the source of promises
Promise.all([
  new Promise(resolve => setTimeout(() => resolve(1), 3000)),
  new Promise(resolve => setTimeout(() => resolve(2), 2000)),
  new Promise(resolve => setTimeout(() => resolve(3), 1000)),
]).then(alert);

// if any of the promises is rejected, then the returned promise rejects with that error
Promise.all([
  new Promise((resolve, reject) => setTimeout(() => resolve(1), 1000)),
  new Promise((resolve, reject) => setTimeout(() => reject(new Error("Whoops!")), 2000)),
  new Promise((resolve, reject) => setTimeout(() => resolve(3), 3000))
]).catch(alert); // Error: Whoops!

// if any of the passed objects are not a promise then they will be passed to the resulting array as is
Promise.all([
  new Promise((resolve, reject) => {
    setTimeout(() => resolve(1), 1000)
  }),
  2,
  3
]).then(alert);

// wait for all promises to settle and return an array of result objects: {status:"fulfilled", value:result} or {status:"rejected", reason:error}
let urls = [
  'https://api.github.com/users/iliakan',
  'https://api.github.com/users/remy',
  'https://no-such-url'
];

Promise.allSettled(urls.map(url => fetch(url)))
  .then(results => {
    results.forEach((result, num) => {
      if (result.status == "fulfilled") {
        alert(`${urls[num]}: ${result.value.status}`);
      }
      if (result.status == "rejected") {
        alert(`${urls[num]}: ${result.reason}`);
      }
    });
  });

// wait for only the first settled promise and get its result (or error).
Promise.race([
  new Promise((resolve, reject) => setTimeout(() => resolve(1), 1000)),
  new Promise((resolve, reject) => setTimeout(() => reject(new Error("Whoops!")), 2000)),
  new Promise((resolve, reject) => setTimeout(() => resolve(3), 3000))
]).then(alert);

// waits only for the first fulfilled promise and gets its result; if all the promises are rejected then the returned promise is rejected with an AggregateError that contains all promise errors
Promise.any([
  new Promise((resolve, reject) => setTimeout(() => reject(new Error("Whoops!")), 1000)),
  new Promise((resolve, reject) => setTimeout(() => resolve(1), 2000)),
  new Promise((resolve, reject) => setTimeout(() => resolve(3), 3000))
]).then(alert);

// creates a resolved promise with the result value
Promise.resolve(123);

// create a rejected promise with an error
Promise.reject(new Error('oops'));

// promisification is the conversion of a function that accepts a callback into a function that returns a promise
function promisify(f) {
  return function (...args) {
    return new Promise((resolve, reject) => {
      args.push((err, result) => {
        if (err) {
          reject(err);
        } else {
          resolve(result);
        }
      });
      f.call(this, ...args);
    });
  };
}

// promise handlers are always asynchronous (even when a promise is immediately resolved)
const promise = Promise.resolve();
promise.then(() => alert("promise done!"));
alert("code finished"); // this alert shows first

// when a promise settles, its handlers are put into the microtask queue; tasks are only executed from the queue when nothing else is running

// an unhandled rejection occurs when a promise error is not handled at the end of the microtask queue; an unhandledrejection event is generated when the queue is complete and the engine examines promises to see if any of them is in the "rejected" state
let promise = Promise.reject(new Error("Promise Failed!")); // 1. Error: Promise Failed!
setTimeout(() => promise.catch(err => alert('caught')), 1000); // 2. caught
window.addEventListener('unhandledrejection', event => alert(event.reason));
```

- A promise has two internal properties:
  - The state is initially "pending" and changes to "fulfilled" if resolve() is called or "rejected" when reject() is called.
  - The result is initially undefined then changes to the value resolve() or reject() is called with
- A promise that is either resolved or rejected is called "settled".
- The executor should call only one resolve or one reject; all further calls are ignored.

## Async/Await
```js
// an async function always returns a promise; other values are wrapped in a resolved promise automatically
async function f() {
  return 1;
}
f().then(alert); // 1

// the await keyword can only be used within async functions and waits until a promise settles and returns its result
async function f() {
  const promise = new Promise((resolve, reject) => {
    setTimeout(() => resolve("done!"), 1000)
  });
  const result = await promise;
  alert(result); // "done!"
}

// await can be used at the top level inside a module, otherwise an anonymous async function can be used as a wrapper
(async () => {
  await promise;
})();

// await can be used with thenable objects
class Thenable {
  constructor(num) {
    this.num = num;
  }
  then(resolve, reject) {
    alert(resolve);
    setTimeout(() => resolve(this.num * 2), 1000); // (*)
  }
}

async function f() {
  const result = await new Thenable(1);
  alert(result);
}

// an async class method always returns a promise
class Waiter {
  async wait() {
    return await Promise.resolve(1);
  }
}

new Waiter().wait().then(alert);

// if a promise rejects then await will throw the resulting error
async function f() {
  try {
    await Promise.reject(new Error("Whoops!"));
  } catch(err) {
    alert(err); // TypeError: failed to fetch
  }
}

// if the error is not caught then the promise generated by the call to the async function becomes rejected
async function f() {
  let response = await fetch('http://no-such-url');
}

f().catch(alert);
```


## Generators
```js
// generator functions can yield multiple values, one after another, on-demand
function* gen() {
  yield 1;
  yield 2;
  return 3;
}

// a generator function returns a generator object when it is called
const g = gen();

// the next() method runs the execution until the nearest yield statement
const result = g.next();

// the result of a call to next has a value property containing the yielded value and a done property indicating if the generator has finished
console.log(g.next());
console.log(g.next());
console.log(g.next());

// subsequent calls to next() return {done: true}
console.log(g.next());

// generators are iterable; iteration ignores the last value when done=true
for (let value of g) {
  console.log(value);
}

console.log([0, ...gen()]);

// we can use a generator function for iteration by providing it as Symbol.iterator
const range = {
  from: 1,
  to: 5,

  *[Symbol.iterator]() { // shorthand for [Symbol.iterator]: function*()
    for(let value = this.from; value <= this.to; value++) {
      yield value;
    }
  }
};

alert([...range]); // 1,2,3,4,5

// a generator can yield value forever


// generator composition allows for embedding generators 
function* generateSequence(start, end) {
  for (let i = start; i <= end; i++) yield i;
}

function* generatePasswordCodes() {
  // delegates the execution to another generator
  yield* generateSequence(48, 57);
  yield* generateSequence(65, 90);
  yield* generateSequence(97, 122);
}

let str = '';
for(let code of generatePasswordCodes()) {
  str += String.fromCharCode(code);
}

// yield can pass a value inside the generator as well
function* gen() {
  const result = yield "2 + 2 = ?";
  alert(result);
}

const generator = gen();

// the first call to next() should always be made without an argument
const question = generator.next().value;

generator.next(4);

function* gen() {
  let ask1 = yield "2 + 2 = ?";
  alert(ask1); // 4

  let ask2 = yield "3 * 3 = ?"
  alert(ask2); // 9
}

let generator = gen();
alert( generator.next().value ); // "2 + 2 = ?"
alert( generator.next(4).value ); // "3 * 3 = ?"
alert( generator.next(9).done ); // true

// generator.throw() can be used to pass an error into a yield; the error is thrown in the line with that yield
function* gen() {
  try {
    let result = yield "2 + 2 = ?"; // (1)

    alert("The execution does not reach here, because the exception is thrown above");
  } catch(e) {
    alert(e); // shows the error
  }
}

let generator = gen();

let question = generator.next().value;

generator.throw(new Error("The answer is not found in my database")); // (2)


// if we don't catch the error then it falls through to the calling code
function* generate() {
  let result = yield "2 + 2 = ?"; // Error in this line
}

let generator = generate();

let question = generator.next().value;

try {
  generator.throw(new Error("The answer is not found in my database"));
} catch(e) {
  alert(e); // shows the error
}

// generator.return() finishes the generator execution and returns the given value
function* gen() {
  yield 1;
  yield 2;
  yield 3;
}

const g = gen();

g.next();        // { value: 1, done: false }
g.return('foo'); // { value: "foo", done: true }
g.next();        // { value: undefined, done: true }
```

## Objects
```js
// objects store properties consisting of a string key and an arbitrary value

// create an empty object with the Object constructor
new Object();

// create an empty object using object literal syntax
{};

// 
const user = {
  name: 'John',
  age: 30,
}

// get a property by key
console.log(user.name);

// assign a new property
user.isAdmin = true;

// remove a property
delete user.age;

// multi-word properties must be quoted
let user = {
  'likes birds': true,
};

// and accessed using square bracket notation
user["likes birds"] = true;
alert(user["likes birds"]); // true
delete user["likes birds"];

// square brackets allow for obtaining the property name as the result of an expression
let key = "likes birds";
user[key] = true;

// computed properties allow for obtaining the property name as the result of an expression for object literals
let fruit = prompt("Which fruit to buy?", "apple");
let bag = {
  [fruit]: 5,
  [fruit + 'Computers']: 5,
};
alert(bag.apple);

// there is a shorthand for making a property from a variable
function makeUser(name, age) {
  return {
    name, // same as name: name
    age,  // same as age: age
    // ...
  };
}

// property names can be strings or symbols; other types are automatically converted to strings
let obj = {
  0: "test" // same as "0": "test"
};

alert( obj["0"] ); // test
alert( obj[0] ); // test

// reading a non-existent property just returns undefined
let user = {};
user.noSuchProperty === undefined;

// check if a property name exists in an object
let user = { name: "John", age: undefined };
alert( "name" in user ); // true, user.name exists
alert( "age" in user ); // true, user.age exists
alert( "blabla" in user ); // false, user.blabla doesn't exist

// iterator over the keys of an object
let user = {
  name: "John",
  age: 30,
  isAdmin: true
};

for (let key in user) {
  alert( key );  // name, age, isAdmin
  alert( user[key] ); // John, 30, true
}

// keys are ordered as follows: integer properties are sorted (i.e. those that can be converted to-and-from an integer without change), others appear in creation order
let codes = {
  "49": "Germany",
  "41": "Switzerland",
  "44": "Great Britain",
  "1": "USA"
};

for (let code in codes) {
  alert(code); // 1, 41, 44, 49
}

let user = {
  name: "John",
  surname: "Smith"
};
user.age = 25; // add one more

for (let prop in user) {
  alert( prop ); // name, surname, age
}
```

## Statements
```js
// statements are delimited with a semicolon
alert('Hello'); alert('World');

// a line-break is usually also treated as a delimiter (automatic semicolon insertion)
alert('Hello')
alert('World')

// though this sometimes does not work
alert("There will be an error after this message")
[1, 2].forEach(alert)
```


## Functions
```js
// arrow functions provide a more concise syntax than function expressions
let sayHi = () => alert("Hello!");
const sum = (a, b) => a + b;

// parentheses can be omitted if there is only one parameter
const double = n => n * 2;

// multiline arrow functions
const sum = (a, b) => {
  let result = a + b;
  return result; // if we use curly braces, then we need an explicit "return"
};
```

## JSON
```js
let student = {
  name: 'John',
  age: 30,
  isAdmin: false,
  courses: ['html', 'css', 'js'],
  spouse: null
};

// serialize object to JSON
JSON.stringify(student);

// function properties, symbolic keys and properties that store undefined are ignored
let user = {
  sayHi() { // ignored
    alert("Hello");
  },
  [Symbol("id")]: 123, // ignored
  something: undefined // ignored
};

// specify properties to encode
let room = {
  number: 23
};

let meetup = {
  title: "Conference",
  participants: [{name: "John"}, {name: "Alice"}],
  place: room,
};

room.occupiedBy = meetup;

JSON.stringify(meetup, ['title', 'participants', 'place', 'name', 'number']);

// specify mapping function as replacer to be called for every key-value pair and return the replaced value or undefined if the value is to be skipped; the first call is made using a special wrapper object with "" as the key and the target object as the value
JSON.stringify(meetup, (key, value) => {
  return (key == 'occupiedBy') ? undefined : value;
});

// specify number of spaces to use for indentation
let user = {
  name: "John",
  age: 25,
  roles: {
    isAdmin: false,
    isEditor: true
  }
};

alert(JSON.stringify(user, null, 2));

// a toJSON method can be provided for JSON.stringify to call
let room = {
  number: 23,
  toJSON() {
    return this.number;
  }
};

let meetup = {
  title: "Conference",
  room
};

JSON.stringify(meetup)

// parse JSON into an object
let userData = '{ "name": "John", "age": 35, "isAdmin": false, "friends": [0,1,2,3] }';
let user = JSON.parse(userData);

// specify reviver function that will be called for each key-value pair and can transform the value
let str = '{"title":"Conference","date":"2017-11-30T12:00:00.000Z"}';
let meetup = JSON.parse(str, function(key, value) {
  if (key == 'date') return new Date(value);
  return value;
});

meetup.date.getDate()

let schedule = `{
  "meetups": [
    {"title":"Conference","date":"2017-11-30T12:00:00.000Z"},
    {"title":"Birthday","date":"2017-04-18T12:00:00.000Z"}
  ]
}`;

schedule = JSON.parse(schedule, function(key, value) {
  if (key == 'date') return new Date(value);
  return value;
});

alert( schedule.meetups[1].date.getDate() ); // works!
```


## Browser

```js
// display a dialogue with a message
alert('Hello!');

// display a dialogue with an input field
const input = prompt('What is your name?');

// display a dialogue with Ok/Cancel buttons
const input = confirm('OK?')

// set breakpoint when devtools are open
debugger;
```


## Rest Parameters & Spread
```js
// a rest parameter gathers the rest of the list of arguments into an array
function sumAll(...args) {
  let sum = 0;
  for (let arg of args) sum += arg;
  return sum;
}

sumAll(1, 2, 3);

// rest parameters must be at the end
function showName(firstName, lastName, ...titles) {
  alert(firstName + ' ' + lastName);

  alert(titles[0]);
  alert(titles[1]);
  alert(titles.length);
}

showName("Julius", "Caesar", "Consul", "Imperator");

// there is a special array-like object named arguments that contains all arguments by their index
function showName() {
  alert(arguments.length);
  alert(arguments[0]);
  alert(arguments[1]);

  for(let arg of arguments) alert(arg);
}

// arrow functions do not have this or arguments
function f() {
  let showArg = () => alert(arguments[0]);
  showArg();
}

// spread syntax expands an iterable object into a list of arguments
const arr = [3, 5, 1];
Math.max(...arr);

const arr1 = [1, -2, 3, 4];
const arr2 = [8, 3, -8, 1];
Math.max(...arr1, ...arr2);

Math.max(1, ...arr1, 2, ...arr2, 25);

// spread syntax can also be used to merge arrays
const arr = [3, 5, 1];
const arr2 = [8, 9, 15];

alert([0, ...arr, 2, ...arr2]);

// spread syntax works with any iterable
const str = "Hello";
[...str];

// this can be used to copy an array
const arr = [1, 2, 3];
const arrCopy = [...arr];

// and similarly for objects
const obj = { a: 1, b: 2, c: 3 };
const objCopy = { ...obj };
```

## Classes
```js
class User {
  constructor(name) {
    this.name = name;
  }

  sayHi() {
    alert(this.name);
  }
}

const user = new User("John");
user.sayHi();

// a class is a kind of function
alert(typeof User); // function

// the class syntax creates a function that becomes the result of the class declaration; the definition is taken from the constructor method

```


The function code is taken from the constructor method (assumed empty if we don’t write such method).  Stores class methods, such as sayHi, in User.prototype. After new User object is created, when we call its method, it’s taken from the prototype, just as described in the chapter F.prototype. So the object has access to class methods.


// class is a function
alert(typeof User); // function

// ...or, more precisely, the constructor method
alert(User === User.prototype.constructor); // true

// The methods are in User.prototype, e.g:
alert(User.prototype.sayHi); // the code of the sayHi method

// there are exactly two methods in the prototype
alert(Object.getOwnPropertyNames(User.prototype)); // constructor, sayHi
Not just a syntactic sugar
Sometimes people say that class is a “syntactic sugar” (syntax that is designed to make things easier to read, but doesn’t introduce anything new), because we could actually declare the same thing without using the class keyword at all:

// rewriting class User in pure functions

// 1. Create constructor function
function User(name) {
  this.name = name;
}
// a function prototype has "constructor" property by default,
// so we don't need to create it

// 2. Add the method to prototype
User.prototype.sayHi = function() {
  alert(this.name);
};

// Usage:
let user = new User("John");
user.sayHi();
The result of this definition is about the same. So, there are indeed reasons why class can be considered a syntactic sugar to define a constructor together with its prototype methods.

Still, there are important differences.

First, a function created by class is labelled by a special internal property [[IsClassConstructor]]: true. So it’s not entirely the same as creating it manually.

The language checks for that property in a variety of places. For example, unlike a regular function, it must be called with new:

class User {
  constructor() {}
}

alert(typeof User); // function
User(); // Error: Class constructor User cannot be invoked without 'new'
Also, a string representation of a class constructor in most JavaScript engines starts with the “class…”

class User {
  constructor() {}
}

alert(User); // class User { ... }
There are other differences, we’ll see them soon.

Class methods are non-enumerable. A class definition sets enumerable flag to false for all methods in the "prototype".

That’s good, because if we for..in over an object, we usually don’t want its class methods.

Classes always use strict. All code inside the class construct is automatically in strict mode.

Besides, class syntax brings many other features that we’ll explore later.

Class Expression
Just like functions, classes can be defined inside another expression, passed around, returned, assigned, etc.

Here’s an example of a class expression:

let User = class {
  sayHi() {
    alert("Hello");
  }
};
Similar to Named Function Expressions, class expressions may have a name.

If a class expression has a name, it’s visible inside the class only:

// "Named Class Expression"
// (no such term in the spec, but that's similar to Named Function Expression)
let User = class MyClass {
  sayHi() {
    alert(MyClass); // MyClass name is visible only inside the class
  }
};

new User().sayHi(); // works, shows MyClass definition

alert(MyClass); // error, MyClass name isn't visible outside of the class
We can even make classes dynamically “on-demand”, like this:

function makeClass(phrase) {
  // declare a class and return it
  return class {
    sayHi() {
      alert(phrase);
    }
  };
}

// Create a new class
let User = makeClass("Hello");

new User().sayHi(); // Hello
Getters/setters
Just like literal objects, classes may include getters/setters, computed properties etc.

Here’s an example for user.name implemented using get/set:

class User {

  constructor(name) {
    // invokes the setter
    this.name = name;
  }

  get name() {
    return this._name;
  }

  set name(value) {
    if (value.length < 4) {
      alert("Name is too short.");
      return;
    }
    this._name = value;
  }

}

let user = new User("John");
alert(user.name); // John

user = new User(""); // Name is too short.
Technically, such class declaration works by creating getters and setters in User.prototype.

Computed names […]
Here’s an example with a computed method name using brackets [...]:

class User {

  ['say' + 'Hi']() {
    alert("Hello");
  }

}

new User().sayHi();
Such features are easy to remember, as they resemble that of literal objects.

Class fields
Old browsers may need a polyfill
Class fields are a recent addition to the language.

Previously, our classes only had methods.

“Class fields” is a syntax that allows to add any properties.

For instance, let’s add name property to class User:

class User {
  name = "John";

  sayHi() {
    alert(`Hello, ${this.name}!`);
  }
}

new User().sayHi(); // Hello, John!
So, we just write " = " in the declaration, and that’s it.

The important difference of class fields is that they are set on individual objects, not User.prototype:

class User {
  name = "John";
}

let user = new User();
alert(user.name); // John
alert(User.prototype.name); // undefined
We can also assign values using more complex expressions and function calls:

class User {
  name = prompt("Name, please?", "John");
}

let user = new User();
alert(user.name); // John
Making bound methods with class fields
As demonstrated in the chapter Function binding functions in JavaScript have a dynamic this. It depends on the context of the call.

So if an object method is passed around and called in another context, this won’t be a reference to its object any more.

For instance, this code will show undefined:

class Button {
  constructor(value) {
    this.value = value;
  }

  click() {
    alert(this.value);
  }
}

let button = new Button("hello");

setTimeout(button.click, 1000); // undefined
The problem is called "losing this".

There are two approaches to fixing it, as discussed in the chapter Function binding:

Pass a wrapper-function, such as setTimeout(() => button.click(), 1000).
Bind the method to object, e.g. in the constructor.
Class fields provide another, quite elegant syntax:

class Button {
  constructor(value) {
    this.value = value;
  }
  click = () => {
    alert(this.value);
  }
}

let button = new Button("hello");

setTimeout(button.click, 1000); // hello
The class field click = () => {...} is created on a per-object basis, there’s a separate function for each Button object, with this inside it referencing that object. We can pass button.click around anywhere, and the value of this will always be correct.

That’s especially useful in browser environment, for event listeners.

Summary
The basic class syntax looks like this:

class MyClass {
  prop = value; // property

  constructor(...) { // constructor
    // ...
  }

  method(...) {} // method

  get something(...) {} // getter method
  set something(...) {} // setter method

  [Symbol.iterator]() {} // method with computed name (symbol here)
  // ...
}
MyClass is technically a function (the one that we provide as constructor), while methods, getters and setters are written to MyClass.prototype.

In the next chapters we’ll learn more about classes, including inheritance and other features.

Tasks
Rewrite to class
importance: 5
The Clock class (see the sandbox) is written in functional style. Rewrite it in the “class” syntax.

P.S. The clock ticks in the console, open it to see.

Open a sandbox for the task.

solution



## Resources
- https://javascript.info/
- https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference
- https://eslint.org/


The "switch" statement

Functions
Function expressions

Objects: the basics
Objects
Object references and copying
Garbage collection
Object methods, "this"
Constructor, operator "new"
Optional chaining '?.'
Symbol type
Object to primitive conversion

Data types
Methods of primitives
Numbers
Strings
Arrays
Array methods
Iterables
Map and Set
WeakMap and WeakSet
Object.keys, values, entries
Destructuring assignment
Date and time

Advanced working with functions
Recursion and stack
Variable scope, closure
The old "var"
Global object
Function object, NFE
The "new Function" syntax
Scheduling: setTimeout and setInterval
Decorators and forwarding, call/apply
Function binding
Arrow functions revisited

Object properties configuration
Property flags and descriptors
Property getters and setters
Prototypes, inheritance
Prototypal inheritance
F.prototype
Native prototypes
Prototype methods, objects without __proto__

Classes
Class inheritance
Static properties and methods
Private and protected properties and methods
Extending built-in classes
Class checking: "instanceof"
Mixins

Error handling, "try...catch"
Custom errors, extending Error

Async iteration and generators

Modules
Modules, introduction
Export and Import
Dynamic imports

Miscellaneous
Proxy and Reflect
Eval: run a code string
Currying
Reference Type
BigInt
Unicode, String internals
WeakRef and FinalizationRegistry



Browser: Document, Events, Interfaces
Learning how to manage the browser page: add elements, manipulate their size and position, dynamically create interfaces and interact with the visitor.

Document
Browser environment, specs
DOM tree
Walking the DOM
Searching: getElement*, querySelector*
Node properties: type, tag and contents
Attributes and properties
Modifying the document
Styles and classes
Element size and scrolling
Window sizes and scrolling
Coordinates
Introduction to Events
Introduction to browser events
Bubbling and capturing
Event delegation
Browser default actions
Dispatching custom events
UI Events
Mouse events
Moving the mouse: mouseover/out, mouseenter/leave
Drag'n'Drop with mouse events
Pointer events
Keyboard: keydown and keyup
Scrolling
Forms, controls
Form properties and methods
Focusing: focus/blur
Events: change, input, cut, copy, paste
Forms: event and method submit
Document and resource loading
Page: DOMContentLoaded, load, beforeunload, unload
Scripts: async, defer
Resource loading: onload and onerror

Miscellaneous
Mutation observer
Selection and Range
Event loop: microtasks and macrotasks
Additional articles
List of extra topics that assume you've covered the first two parts of tutorial. There is no clear hierarchy here, you can read articles in the order you want.
Frames and windows
Popups and window methods
Cross-window communication
The clickjacking attack

Binary data, files
ArrayBuffer, binary arrays
TextDecoder and TextEncoder
Blob
File and FileReader

Network requests
Fetch
FormData
Fetch: Download progress
Fetch: Abort
Fetch: Cross-Origin Requests
Fetch API
URL objects
XMLHttpRequest
Resumable file upload
Long polling
WebSocket
Server Sent Events

Storing data in the browser
Cookies, document.cookie
LocalStorage, sessionStorage
IndexedDB

