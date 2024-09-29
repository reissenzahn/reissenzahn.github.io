---
title: "Ruby"
date: "2024-08-03"
tags: ["Programming"]
---

## Commands

```bash
# check file for syntax errors
ruby -cw hello.rb

# execute file
ruby hello.rb

# execute literal script
ruby -e 'puts "Hello, Ruby!"'

# require a feature
ruby -rprofile

# start interactive interpreter
irb

# start interactive interpreter and require feature
irb -rrbconfig

# show ruby executable installation directory
> RbConfig::CONFIG["bindir"]

# execute clean_tmp task in the admin namespace
rake admin:clean_tmp

# list rake tasks
rake --tasks

# install a gem
gem install prawn

# uninstall a gem
gem uninstall prawn
```

## Variables

```ruby
# bind a local variable to an object
x = 1
msg = "hello"

# an assignment expression has the value of its right-hand side
puts x = 1

# local variable
local_variable = 3.14

# instance variables store a value within individual objects
@instance_variable = 7

# class variables store a value per class hierarchy
@@class_variable = "hello"

# global variable
$GLOBAL_VARIABLE

# assign a value to a variable if and only if the variable is not nil or false; evaluates to the value of the variable
@stack ||= []
```

## Constants

```rb
# constants
CONSTANT_VARIABLE = "A"

# a constant defined in a class can be referred to from inside instances of the class or class methods
class Ticket
  VENUES = ["Convention Center", "Fairgrounds", "Town Hall"]

  def initialize(venue, date)
    if VENUES.include?(venue)
      @venue = venue
    else
      raise ArgumentError, "Unknown venue #{venue}"
    end

    @date = date
end

# as well as from outside the class using constant lookup notation
puts Ticket::VENUES

# a warning will be emitted when a constant is re-assigned
A = 1
A = 2
```

## Strings

```ruby
# convert string to an integer
"100".to_i

# concatenate two strings
"hello" + " there"

# format string
"%.2f" % 5.50

# string interpolation
"Hello, #{name}"

# remove trailing newline
"hello\n".chomp

# split string on delimiter
year, month, day = "2024-01-01".split("-")
```

## I/O

```ruby
# print string
print "hello"

# add a trailing newline to the string if there is not one already
puts "hello"

# print inspect string
p x

# read a line of input
input = gets

# read a file
contents = File.read("temp.dat")

# write a file
f = File.new("temp.out", "w")
f.puts "hello"
f.close
```

## Conditional Execution

```rb
if condition
  # do something
end

# one line if statement
if x > 10 then puts x end

# semicolons can be used to mimic line breaks
if x > 10; puts x; end

# else
if condition
else
end

# else if
if condition1
elsif condition2
elsif condition3
end

n = gets.to_i
if n > 0
  puts "Positive"
elsif n < 0
  puts "Negative"
else
  puts "Zero"
end

# nil evaluates to false
if nil; puts "Ain't gonna happen"; end

# not if
unless x > 100
  puts "Small number!"
else
  puts "Big number!"
end

# only false and nil cause a conditional expression to evaluate as false
if x % 2 == 0
  puts "even"
else
  puts "odd"
end

# conditional modifier
puts "Big number!" if x > 100

puts "Small number!" unless x > 100


# if an if statement succeeds, the entire statement evaluates to whatever is represented by th code in the successful branch
x = 1
if x < 0
  "negative"
elsif x > 0
  "positive"
else
  "zero"
end

# an if statement that does not succeed will evaluate to nil
if false
elsif false
end


# the parser allocates space for a local variable when it encounters the sequence identifier = value, creating the variable but not performing an assignment
if false
  x = 1
end

p x # nil
p y # NameError


# performing an assignment in a conditional test can be useful when calling a method that returns nil on failure and some other value on success
name = "David A. Black"
if m = /la/.match(name)
  puts "Found a match!"
  print "Here’s the unmatched start of the string: "
  puts m.pre_match
  print "Here’s the unmatched end of the string: "
  puts m.post_match
else
  puts "No match"
end


# case statement finds the first match using the case equality method (===)
answer = gets.chomp
case answer
when "yes"
  exit
# more than one possible match in a single when clause
when "no", "n"
  puts "Okay"
else
  puts "Huh?"
end

# for any object that does not override it, === works the same as ==
class Ticket
  attr_accessor :venue, :date
  def initialize(venue, date)
    self.venue = venue
    self.date = date
  end
  def ===(other_ticket)
    self.venue == other_ticket.venue
  end
end
ticket1 = Ticket.new("Town Hall", "07/08/18")
ticket2 = Ticket.new("Conference Center", "07/08/18")
ticket3 = Ticket.new("Town Hall", "08/09/18")
puts "ticket1 is for an event at: #{ticket1.venue}."
case ticket1    
when ticket2
  puts "Same location as ticket2!"
when ticket3
  puts "Same location as ticket3!"
else
  puts "No match."
end

# omitting the test expression will match for the first true when condition
case
when user.first_name == "David", user.last_name == "Black"
  puts "You might be David Black."
when Time.now.wday == 5
  puts "You’re not David Black, but at least it’s Friday!"
else
  puts "You’re not David Black, and it’s not Friday."
end


# a case statement evaluates to the value returned by the matched clause or nil if the match fails
```

## Looping

```rb
# loop takes a code block and loops unconditionally
loop 
loop { puts "Looping forever!" }

loop do
  puts "Looping forever!"
end

n = 1
loop do
  n = n + 1
  break if n > 9
end

n = 1
loop do
  n = n + 1
  next unless n == 10
  break
end

# while and until loop conditionally
n = 1
while n < 11
  puts n
  n = n + 1
end

n = 1
begin
  puts n
  n = n + 1
end while n < 11

n = 1
until n > 10
  puts n
  n = n + 1
end


# while and until can be used in a modifier position in one-line statements
n = 1
n = n + 1 until n == 10

# this does not behave the same way as the post-positioned while and until
a = 1
a += 1 until true
puts a  # 1


# loop over an array of values
celsius = [0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100]
for c in celsius
  puts "#{c}\t#{Temperature.c2f(c)}"
end
```

## Iteration

```rb
# loop over iterator
[1, 2, 3].each do |x|
  puts x
end


def my_loop
  while true
    yield
  end
end

my_loop { puts "My-looping forever!" }


# the difference between the two ways of delimiting a code block is a difference in precedence: 
array = [1,2,3]
array.map {|n| n * 10 }
array.map do |n| n * 10 end

puts array.map {|n| n * 10 }  # puts(array.map {|n| n * 10})
puts array.map do |n| n * 10 end # puts(array.map) do |n| n * 10 end


# the times() method is an instance method of the Integer class that runs the code block n times for any integer n and returns n
5.times { puts "Writing this 5 times!" }

# when a method yields, it can yield one or more values which the block picks up through its parameters
5.times {|i| puts "I'm on iteration #{i}!"}

class Integer
  def my_times
    c = 0
    until c == self
      yield c
      c += 1
    end
    self
  end
end

ret = 5.my_times {|i| puts "I’m on iteration #{i}!" }
puts ret


# you run the each method on a collection object and each yields each item in the collection to your code block, one at a time; it returns the original array
array = [1,2,3,4,5]
array.each {|e| puts "The block just got handed #{e}." }

class Array
  def my_each
    c = 0
    until c == size
      yield self[c]
      c += 1
    end
    self
  end
end

array.my_each {|e| puts "The block just got handed #{e}." }


# map walks through an array one element at a time and yields each element to the code block; it returns a new array that contains the accumulated return values of the code block from the iterations
names = ["David", "Alan", "Black"]
names.map {|name| name.upcase }

class Array
  def my_map
    c = 0
    acc = []
    until c == size
      acc << yield self[c]
      c += 1
    end
    acc 
  end

  def my_map
    acc = []
    my_each {|e| acc << yield e }
    acc
  end
end

# parameter bindings are similar to method arguments
def block_args_unleashed
  yield(1,2,3,4,5)
end

block_args_unleashed do |a,b=1,*c,d,e|
  puts "Arguments:"
  p a,b,c,d,e
end

# there are three types of block variables available: local variables that exist already when the block is created, block parameters which are always block local and true block-locals which are listed after the semicolon and protect any same-named variables from the outer scope
def block_scope_demo
  x = 100
  1.times do
    puts x
    x = 200
  end
  puts x
end

def block_local_parameter
  x = 100
  [1,2,3].each do |x|
    puts "Parameter x is #{x}"
    x = x + 10
    puts "Reassigned to x in block; it’s now #{x}"
  end
  puts "Outer x is still #{x}"
end

celsius.each do |c;fahrenheit|
  fahrenheit = Temperature.c2f(2)
  puts "#{c}\t#{fahrenheit}"
end
```

## Exceptions

```rb
# an exception is an instance of the class Exception or a descendant of that class

# by default, the raise method raises a RuntimeError
begin
  result = 100 / 0
rescue
  puts "Divide by zero!"
end

# pinpoint the exception to trap
begin
  result = 100 / 0
rescue ZeroDivisionError
  puts "Divide by zero!"
end

# the beginning of a method or code block provides an implicit begin/end context
def open_user_file
  print "File to open: "
  filename = gets.chomp
  fh = File.open(filename)
  yield fh
  fh.close
  rescue
    puts "Couldn’t open your file!"
end

open_user_file do |filename|
  fh = File.open(filename) 
  yield fh
  fh.close
  rescue
    puts "Couldn’t open your file!"
end


# open an irb session for debugging
binding.irb


# the safe navigation operator allows for avoiding calling a method on nil; it will only call the next method if the receiver is not nil, otherwise the expression returns nil
class Roster
  attr_accessor :players
end

class Player
  attr_accessor :name, :position

  def initialize(name, position)
    @name = name
    @position = position
  end
end

moore = Player.new("Maya Moore", "Forward")
taurasi = Player.new("Diana Taurasi", "Guard")
tourney_roster1 = Roster.new
tourney_roster1.players = [moore, taurasi]

if tourney_roster2.players&.first&.position == "Forward"
  puts "Forward: #{tourney_roster1.players.first.name}"
end

tourney_roster2.players&.first
tourney_roster2.players&.first&.position


# raise an exception explicitly
raise ArgumentError, "I need a number under 10" unless x < 10

# raise a RuntimeError
raise "Problem!"


# assign a caught exception object to a variable
begin
  fussy_method(20)
rescue ArgumentError => e
  puts "That was not an acceptable number!"
  puts "Here’s the backtrace for this exception:"
  puts e.backtrace
  puts "And here’s the exception object’s message:"
  puts e.message
end


# re-raise an exception
begin
  fh = File.open(filename)
rescue => e
  logfile.puts("User tried to open #{filename}, #{Time.now}")
  logfile.puts("Exception: #{e.message}")
  raise
end


# the ensure clause is executed whether an exception is raised or not
def line_from_file(filename, substring)
  fh = File.open(filename)
  begin
    line = fh.gets
    raise ArgumentError unless line.include?(substring)
  rescue ArgumentError
    puts "Invalid line!"
    raise
  ensure
    fh.close
  end
  return line
end


# a new exception class can be created by inheriting from Exception or from a descendant class of Exception
class MyNewException < Exception
end

begin
  raise MyNewException, "some new kind of error has occurred!"
rescue MyNewException => e 
  puts "Just raised an exception: #{e}"
end

class InvalidLineError < StandardError
end
def line_from_file(filename, substring)
  fh = File.open(filename)
  line = fh.gets
  raise InvalidLineError unless line.include?(substring)
  return line
  rescue InvalidLineError
    puts "Invalid line!"
    raise
  ensure
    fh.close
end

# namespacing exceptions
module TextHandler
  class InvalidLineError < StandardError
  end
end
 
def line_from_file(filename, substring)
  fh = File.open(filename)
  line = fh.gets
  raise TextHandler::InvalidLineError unless line.include?(substring)
end
```

## Require

```ruby
# loading files

# load file from the current working directory or from the load path
load "hello.rb"

# load from a relative directory
load "../hello.rb"

# load a feature
require "./hello"

# load a feature searching relative to the current file
require_relative "hello"

# load a gem
require "bundler"

# show load path
puts $:
```

## Objects

```rb
# create a generic object
obj = Object.new

# define a singleton method for an object
def obj.talk
  puts "Hello"
end

# send a message to an object to call a method
obj.talk

# define a method that takes a parametere
def obj.c2f(c)
  c * 9.0 / 5 + 32
end

obj.c2f(100)

# create a more elaborate object
ticket = Object.new

def ticket.date
  "1903-01-02"
end

def ticket.venue
  "Town Hall"
end

def ticket.performer
  "Mark Twain"
end

def ticket.seat
  "Second Balcony, row J, seat 12"
end

def ticket.price
  5.50
end

# methods that return a boolean conventionally have names that end in a question mark
def ticket.available?
  # true and false are objects too
  false
end

# list object methods
obj.methods.sort

# every object has a unique identifier associated with it
obj.object_id

# ask an object if it responds to a message
obj.respond_to?("talk")

# send a message to an object
request = gets
obj.send(request)

# public_send cannot call private methods
obj.public_send(request)

# freeze an object to prevent mutations
"hello".freeze

# duplicate an object, unfreezing it if it is frozen
"hello".dup

# duplicate an object, without unfreezing it
"hello".clone

# freezing an object does not freeze the objects it contains
["one", "two", "three"][2].replace("four")
```


## Parameters

```rb
# variadic arguments are assigned to the parameter as an array
def multi_args(a, b, *rest)
  p a, b, rest
end

multi_args(1, 2, 3, 4, 5)

# assign default values to arguments
def default_args(a, b, c = 3)
  p a, b, c
end

default_args(1, 2)
default_args(1, 2, 4)

# values are assigned to as many variables as according to the priority: normal > default > variadic
def mixed_args(a, b = 1, *c, d, e)
  p a, b, c, d, e
end

mixed_args(1, 2, 3)
mixed_args(1, 2, 3, 4)
mixed_args(1, 2, 3, 4, 5)
```

## Classes

```rb
class Ticket
  # define an instance method
  def event
    "Not sure yet..."
  end
end

# instantiate a new object using the constructor
ticket = Ticket.new

# re-open a class
class C
  def x
  end
end

class C
  def y
  end
end

# instance variables are only visible to the object to which they belong
class Person
  def set_name(name)
    @name = name
  end

  def get_name
    @name
  end
end

person = Person.new
person.set_name("Joe")
person.get_name

class Ticket
  # initialize is executed each time an instance of the class is created
  def initialize(venue, date)
    @venue = venue
    @date = date
  end

  def venue
    @venue
  end

  def date
    @date
  end

  # define setter method
  def price=(price)
    if (price * 100).to_i == price * 100
      @price = price
    end
  end

  def price
    @price
  end
end

ticket = Ticket.new("Town Hall", "2013-11-12")

# equivalent to ticket.price=(63.00), expression evaluates to the value on the right-hand side
ticket.price = 63.00


# a class can be created by sending the new message to Class
my_class = Class.new

# that class can create instances of its own
my_instance = my_class.new

# we can also define methods on this anonymous class
c = Class.new do
  def say_hello
    puts "Hello!"
  end
end


# define a class method (i.e. a singleton method defined on a class)
def Ticket.most_expensive(*tickets)
  tickets.max_by(&:price)
end

highest = Ticket.most_expensive(t1, t2, t3)
```

## Attributes

```rb
class Ticket
  # create a method that reads and returns the value of the instance variable
  attr_reader :venue, :date, :price

  # def venue
  #   @venue
  # end

  # create a method that writes an instance variable
  attr_writer :price

  # def price=(price)
  #   @price = price
  # end

  # create reader and writer methods
  attr_accessor :performer

  def initialize(venue, date)
    @venue = venue
    @date = date
  end
end
```

## Inheritance

```rb
class Publication
  attr_accessor :publisher
end

# declare a subclass of another class
class Magazine < Publication
  attr_accessor :editor
end

mag = Magazine.new
mag.publisher = "David A. Black"
mag.editor = "Joe Leo"
puts mag.publisher, mag.editor

# every class is a descendant of Object
class C
end

class D < C
end

puts D.superclass # C
puts D.superclass.superclass # Object

# check whether an object has a given class either as its class or as one of its ancestor superclasses
mag.is_a?(Magazine)
mag.is_a?(Publication)

# the behavior of an object can deviate from those supplied by its class
mag = Magazine.new

def mag.wings
  puts "Flying away!"
end

mag.wings
```

## Symbols

```rb
# symbols are a naming or labeling facility
:hello
```


## Rake

```rb
namespace :admin do
  desc "Interactively delete all files in /tmp"

  task :clean_tmp do
    Dir["/tmp/*"].each do |f|
      next unless File.file?(f)

      print "Delete #{f}? "
      answer = $stdin.gets

      case answer
      when /^y/
        File.unlink(f)
      when /^q/
        break
      end
    end
  end
end


# define a task in the top-level namespace
task :clean_tmp do
  # ...
end

# define task in nested namespaces
namespace :admin do
  namespace :clean do
    task :tmp do
      # ...
    end
  end
end
```

## Modules

```rb
module MyFirstModule
  def ruby_version
    system("ruby -v")
  end
end

class ModuleTester
  include MyFirstModule
end

mt = ModuleTester.new
my.ruby_version


module Stacklike
  def stack
    @stack ||= []
  end

  def add_to_stack(obj)
    stack.push(obj)
  end

  def take_from_stack
    stack.pop
  end
end

class Stack
  include Stacklike
end

s = Stack.new
s.add_to_stack("one")
s.add_to_stack("two")
s.add_to_stack("three")

puts s.take_from_stack
puts s.stack

class Suitcase
end

class CargoHold
  include Stacklike

  def load(obj)
    puts "Loading: #{obj.object_id}"
    add_to_stack(obj)
  end

  def unload(obj)
    take_from_stack
  end
end

ch = CargoHold.new
ch.load(SuitCase.new)
ch.load(SuitCase.new)
puts ch.unload
```

## Method Lookup

```rb
module M
  def report
    puts "M: report()"
  end
end

class C
  include M
end

class D < C
end

obj = D.new
obj.report


# when a class mixes in two or more modules and more than one implements the method being searched for, the modules are searched in reverse order of inclusion
module M
  def report
    puts "M: report()"
  end
end

module N
  def report
    puts "N: report()"
  end
end

class C
  include M
  include N
end

c = C.new
c.report  # N: report()


module MeFirst
  def report
    puts "Hello from module!"
  end
end

class Person
  prepend MeFirst

  def report
    puts "Hello from class!"
  end
end

p = Person.new
p.report  # Hello from module!


# extend will make the methods of a module available as class methods
module Convertible
  def c2f(celsius)
    celsius * 9.0 / 5 + 32
  end

  def f2c(fahrenheit)
    (fahrenheit - 32) * 5 / 9.0
  end
end

class Thermometer
  extend Convertible
end

puts Temperature.c2f(100)
puts Temperature.f2c(212)


# the super keyword can be used to jump up to the next-highest definition in the method-lookup path of the method you are currently executing
module M
  def report
    puts "M: report()"
  end
end

class C
  include M

  def report
    puts "C: report()"
    super  # keep looking and find the next match
  end
end

c = C.new
c.report


class Bicycle
  attr_reader :gears, :wheels, :seats

  def initialize(gears = 1)
    @wheels = 2
    @seats = 1
    @gears = gears
  end
end

class Tandem < Bicycle
  def initialize(gears)
    # super automatically forwards the arguments that were passed to the method from which it is called if invoked with no argument list
    super
    @seats = 2
  end
end

# inspect this namespace and determine if a super method exists
class Bicycle
  attr_reader :gears, :wheels, :seats

  def initialize(gears = 1)
    @wheels = 2
    @seats = 1
    @gears = gears
  end

  def rent
    puts "Sorry but this model is sold out."
  end
end

class Tandem < Bicycle
  def initialize(gears)
    super
    @seats = 2
  end

  def rent
    puts "This bike is available!"
  end
end

# inspect method hierarchies with method and super_method
t = Tandem.new(1)
t.method(:rent).call
t.method(:rent).super_method.call


# override method_missing
o = Object.new

def o.method_missing(m, *args)
  puts "Unknown method: #{m}"
end

o.blah

# it is common to want to intercept an unrecognized message and decide whether to handle it or pass it along to the original method_missing
class Student
  def method_missing(m, *args)
    if m.to_s.start_with?("grade_for_")
      puts "You got an A in #{m.to_s.split("_").last.capitalize}!"
    else
      super
    end
  end
end


# a more extensive example
j = Person.new("John")
p = Person.new("Paul")
g = Person.new("George")
r = Person.new("Ringo")
j.has_friend(p)
j.has_friend(g)
g.has_friend(p)

r.has_hobby("rings")

Person.all_with_friends(p).each do |person|
  puts "#{person.name} is friends with #{p.name}"
end

Person.all_with_hobbies("rings").each do |person|
  puts "#{person.name} is into rings"
end

class Person
  PEOPLE = []

  attr_reader :name, :hobbies, :friends

  def initialize(name)
    @name = name
    @hobbies = []
    @friends = []

    PEOPLE << self
  end

  def has_hobby(hobby)
    @hobbies << hobby
  end

  def has_friend(friend)
    @friends << friend
  end

  def Person.method_missing(m, *args)
    # convert symbol to string
    method = m.to_s

    if method.start_with?("all_with_")
      attr = method[9..-1]
      if Person.public_method_defined?(attr)
        PEOPLE.find_all do |person|
          person.send(attr).include?(args[0])
        end
      else
        raise ArgumentError, "Unknown: #{attr}"
      end
    else
      super
    end
  end
end
```

## Nesting

```rb
# a class definition can be nested inside a module definition; this can be used to create namespaces
module Tools
  class Hammer
  end
end

# the double-colon constant lookup token is used to point the way to the name of a class nested inside a module
h = Tools::Hammer.new
```

## Self

```rb
# self at the top level is the object main
puts "top-level: #{self}"

class C
  # self inside a class definition is the class object itself
  puts "class: #{self}"

  # self inside a method is the object to which the message was sent
  def self.x
    # for a class method, self is the class object
    puts "class method: #{self}"
  end

  def m
    # for an instance method, self is the instance of the class
    puts "instance method: #{self}"
  end
end

c = C.new
c.m

# self in a singleton method is the object that owns the method
obj = Object.new
def obj.m
  puts "singleton method: #{self}"
end
obj.m


# there are a few different ways to define class methods
class C
  def C.x
  end
end

class C
  # self is C in the class definition
  def self.x
  end
end

class C
  class << self
    def x
    end

    def y
    end
  end
end


# if the receiver of the method is self, you can omit the receiver and the dot
class C
  def x
    puts "calling x()"
  end

  def y
    puts "calling y()"
    x
  end
end

c = C.new
c.y


# a given instance variable belongs to whatever object is the current object at thar point in the program
```

## Scope

```rb
# global variables are available everywhere; they never go out of scope
$gvar = "I'm a global!"

class C
  def examine_global
    puts $gvar
  end
end

c = C.new
c.examine_global


# there are a number of pre-initialized global variables

# name of the startup file for the currently running program
$0

# directories that make up the load path
$:

# process ID of the Ruby process
$$


# local scope
class C
  a = 5
  module M
    a = 4
    module N
      a = 3
      class D
        a = 2
        def show_a
          a = 1
          puts a
        end
        puts a
      end
      puts a
    end
    puts a
  end
  puts a
end
d = C::M::N::D.new
d.show_a


# constants have a kind of global visibility or reachability; as long as you know the path to a constant through the classes/modules in which it is nested, you can get to that constants
module M
  class C
    X = 2
    class D
      module N
        X = 1
      end
    end
  end
end

# constant lookup refers to the process of resolving a constant identifier
puts M::C::D::N::X
puts M::C::X

# constants are identified relative to the point of execution
module M
  class C
    class D
      module N
        X = 1
      end
    end
    puts D::N::X
  end
end

# a constant path separator can be prepended to the constant lookup to force an absolute constant path
class Violin
  class String
    attr_accessor :pitch

    def initialize(pitch)
      @pitch = pitch
    end
  end

  def history
    ::String.new(make + ", " + date)
  end
end


# class variables are shared between a class and instances of that class and are not visible to any other objects
class Car
  @@makes = []
  @@cars = {}
  @@total_count = 0

  attr_reader :make

  def self.total_count
    @@total_count
  end

  def self.add_make(make)
    unless @@makes.include?(make)
      @@makes << make
      @@cars[make] = 0
    end
  end

  def initialize(make)
    if @@makes.include?(make)
      puts "Creating a new #{make}!"
      @make = make
      @@cars[make] += 1
      @@total_count += 1
    else
      raise "No such make: #{make}."
    end
  end

  def make_mates
    @@cars[self.make]
  end
end

Car.add_make("Honda")
Car.add_make("Ford")

h = Car.new("Honda")
f = Car.new("Ford")
h2 = Car.new("Honda")

puts Car.total_count
puts h2.make_mates


# class variables are scoped to the class hierarchy
class Parent
  @@value = 100
end
class Child < Parent
  @@value = 200
end
class Parent
  puts @@value  # 200
end


# using instance variables of class objects to maintain per-class state is usually preferable to class variables
class Car
  @@makes = []
  @@cars = {}

  attr_reader :make

  def self.total_count
    @total_count ||= 0
  end

  def self.total_count=(n)
    @total_count = n
  end

  def self.add_make(make)
    unless @@makes.include?(make)
      @@makes << make
      @@cars[make] = 0
    end
  end

  def initialize(make)
    if @@makes.include?(make)
      puts "Creating a new #{make}!"
      @make = make
      @@cars[make] += 1
      self.class.total_count += 1
    else
      raise "No such make: #{make}."
    end
  end

  def make_mates
    @@cars[self.make]
  end
end

class Hybrid < Car
end

h3 = Hybrid.new("Honda")
f2 = Hybrid.new("Ford")

puts "There are #{Hybrid.total_count} hybrids on the road!"
```

## Method-access

```rb
class Cake
  def initialize(batter)
    @batter = batter
    @baked = true
  end
end

class Egg
end

class Flour
end

class Baker
  def bake_cake
    @batter = []
    pour_flour
    add_egg
    stir_batter
    return Cake.new(@batter)
  end

  # all instance methods defined below will be private
  private

  def pour_flour
    @batter.push(Flour.new)
  end

  def add_egg
    @batter.push(Egg.new)
  end

  def stir_batter
  end
end

# private can also take arguments as a list of of the methods to be made private
private :pour_flour, :add_egg, :stir_batter

# private setters must be called with self as the receiver
class Dog
  attr_reader :age, :dog_years

  def dog_years=(years)
    @dog_years = years
  end

  def age=(years)
    @age = years
    self.dog_years = years * 7
  end

  private :dog_years=
end

# protected methods
class C
  def initialize(n)
    @n = n
  end

  def n
    @n
  end

  def compare(c)
    if c.n > n
      puts "Smaller"
    else
      puts "Bigger"
    end
  end

  protected :n
end
```

## Top-level methods

```rb
# define a top-level method
def talk
  puts "Hello"
end

# a method defined at the top level is stored as a private instance method of the Object class

# equivalent
class Object
  private
  def talk
    puts "Hello"
  end
end

# as such, these methods can only be called on self, without an explicit receiver and these private instance methods can be called from anywhere as Object lies in the method-lookup path of every class
```


## Theory
- Objects
  - Ruby sees all data structures and values as objects.
  - Every object is capable of understanding a certain set of messages.
  - Each message that an object understands corresponds directly to a method.
  - Message sending is achieved via the dot operator which sends a message to its receiver.
  - Methods can take arguments which are also objects.
  - The default object self is always defined though which object is self changes.
  - A class defines the functionality of an object and every object is an instance of exactly one class.
  - Objects an acquire methods and behaviors that were not defined in their class.
  - An object is said to respond to a message if the object has a method defined whose name corresponds to the message.
- Load & Require
  - The load path is a list of directories that the interpreter search for files to load.
  - `load` always loads the file while while `require` does not reload a file that is already loaded.
  - `require` loads a feature which allows for treating Ruby extensions the same was as C extensions.
  - `require` is typically used to load extensions and libraries.
  - Typically, `require` is used for standard/installed library while `require_relative` is used for local files.
- Variables
  - Variables typically hold references to objects though some objects are stored in variables as immediate values: integers, symbols, true, false and nil.
  - Any object that is represented as an immediate value is always exactly the same object.
- Classes
  - Classes are objects and can respond to messages.
  - An attribute is a property of an object whose value can be read and/or written through the object via attribute reader and attribute writer methods.
  - A class can have only one superclass (single inheritance).
  - Every class is an instance of a class called `Class`.
  - In the absence of an explicit receiver, messages go to self.
  - In the topmost level of a class definition body, self is the class object itself.
  - Every object is a descendant of `Object` and `Object` inherits from `BasicObject`.
  - `BasicObject` provides a blank-slate object with a minimal set of instance methods.
  - Objects get their methods from their class, their superclasses (and their ancestors) and from their singleton methods.
  - The superclass of `Class` is `Module`.
- Modules
  - Modules are bundles of methods and constants but, unlike classes, do not have instances.
  - The functionality of a module can be added to a specific class or object.
  - `Class` is a subclass of the `Module` class so every class object is also a module object.
  - Modules get mixed in to classes using the `include` method or the `prepend` method such that the instances of the class have access to the instance methods defined in the module.
  - You can mix in more than one module.
  - `Object` mixes in the `Kernel` module which contains most of the fundamental methods.
- Method lookup
  - Classes and modules have methods; objects have the ability to traverse classes and modules in search of methods. 
  - To resolve a message into a method, an object looks for the method in this order:
    1. Modules prepended to its class, in reverse order of prepending
    2. Its class
    3. Modules included in its class, in reverse order of inclusion
    4. Modules prepended to its superclass
    5. It class's superclass
    6. Modules included in its superclass
    7. And so on up to `Object` (and its mix-in `Kernel`) and `BasicObject`
  - An object can see only one version of a method with a given name at any given time; if the method-lookup path includes two or more same-named methods, the first one encountered is executed
  - The `method_missing` method is called if a method cannot be found.
  - Defining a method a second time in a class or module overrides the previous definition.
  - When a class mixes in two or more modules and more than one implements the method being searched for, the modules are searched in reverse order of inclusion.
  - Including a module more than once has no effect.
  - If you `prepend` a module to a class, the object looks in that module first before it looks in the class.
  - Singleton methods lie in a special class: the object's singleton class.
  - The `super` keyword allows for navigating the lookup path explicitly.
  - `super` handles arguments as follows:
    - Called with no arguments, `super` automatically forwards the arguments that were passed to the method from which it is called.
    - Called with an empty argument list (`super()`), `super` sends no arguments.
    - Called with specific arguments (`super(a, b, c)`), `super` sends exactly those arguments.
  - The `Kernel` module provides an instance method called `method_missing` which is executed whenever an object receives a message that it does not know how to respond to (i.e. that does not match a method anywhere in its method-lookup path).
    - `method_missing` can be overridden to intercept these calls.
- self
  - `self` refers to the default or current object.
  - There is always one (and only one) current self object.
  - The current `self` object is determined as follows:
    - Top level: `main` (built-in top-level default object)
    - Class definition: The class object.
    - Module definition: The module object.
    - Top-level method definition: Whatever object is `self` when the method is called; top-level methods are available as private methods to all objects.
    - Instance-method definition in a class: An instance of the class responding to the method call.
    - Instance-method definition in a module: Instance of class that mixes in the module or an individual object extended with the module.
    - Singleton method on a specific object: the object that owns the method.
- Scope
  - Global variables are available everywhere.
  - At any given moment, your program is in a particular local scope. The main thing that changes from one local scope to another is your supply of local variables; you get a new supply when you leave a local scope.
    - The top level has its own local scope.
    - Every class or module definition block has its own local scope, even nested ones.
    - Every method definition has its own local scope.
- Method-access
  - Public is the default access level.
  - Private means that the method cannot be called with an explicit receiver, with the exception of setter methods which must have the receiver self.
  - Protected means you can call the method on an object x as long as the default object is an instance of the same class as x or of an ancestor or descendant class of x's class
    - This is typically desirable if you want one instance of a certain class to do something with another instance of its class.
    - A protected method is thus like a private method, but with an exemption for cases where the class of self and the class of the object having the method called on it are the same or related by inheritance.
- Iterators
  - A code block is a delimited set of program instructions written as part of the method call and available to be executed from the method.
  - A loose convention holds that one-line code blocks use curly braces and multiline blocks use do/end.
  - An iterator is a Ruby method that is called with a code block provided; the method can call the block using the yield keyword.
  - When the method yields to the block, the code in the block runs and then control returns to the method.
  - Yielding takes place while the method is still running. After the code block executes, control returns to the method at the statement immediately following the call to yield.
  - The code block is part of the syntax of the method call; it is not an argument.
  - A code block can take parameters.
  - A method can yield only if it is called with a code block.
  - Methods are not obliged to yield.
  - A block can return a value which comes back as the value returned from yield.


## Resources
- https://bundler.io/
- https://lostisland.github.io/faraday/#/
- http://whatisthor.com/
- https://pry.github.io/
