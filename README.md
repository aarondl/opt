# Opt

**WARNING:** Pre-1.0 there could be significant changes to this. It's not clear
that separate sub-packages is the right way to build this package and a change in that
area could significantly break consumers.

Another generic optional package for Go. This one's focus is on ergonomics and
giving the right level of optionality for every situation.

The goal is to be able to create a Go struct or value that can interoperate with
other language datatypes which may omit a value completely, provide a null
for that value, or provide the value itself.

In addition to the above this package attempts to reasonably implement the
following useful interfaces for interacting with the wider Go ecosystem:

* For JSON interop: `json.Marshaller` & `json.Unmarshaller`
* For database/sql interop: `driver.Valuer` & `sql.Scanner`
* Other text encodings: `encoding.TextMarshaller` & `encoding.TextUnmarshaller`

### Note on sql.NullX types

The Go standard library provides a limited set of database/sql.NullX types.
These are helpful for databases but don't work well outside of that. For example
a JSONified struct with an sql.NullInt32 as of Go 1.18 still marshals as the
following: `{"age":{"Int32":5,"Valid":true}}` which is unexpected by a typical
JavaScript client. Furthermore they lack ergonomic helpers and constructors
which makes them grating to consume and to construct._

## Examples

Here are some examples of the API in action using the `omitnull.Val` type. This
type is the most flexible, and `null.Val` and `omit.Val` are both subsets of its
functionality since they are missing one possible states that it has.

```go
// Working with values that can be null or unset
var val omitnull.Val[int]
val.Set(5)  // set the value to 5
val.Null()  // set the value to null
val.Unset() // unset the value (omitted/undefined)

// Construct new values more directly
val = omitnull.From(5)
val = omitnull.FromPtr(&somePtr) // set to (null | value) based on ptr value

// Convert to other types
val.Ptr() // get a pointer back, will be nil if (null | unset).
omitnull.Map(val, func(i int) int { return i+1 }) // yea, it's here too :| :| :|

// Query state
val.Set(5)
val.IsSet()   // == true
val.IsNull()  // == false
val.IsUnset() // == false

// Fetch the value out with varying levels of safety
v, ok := val.Get() // returns (X, true) if the value is present
v := val.GetOr(6)  // returns 6 if no value is present
v := val.MustGet() // panics if no value is there

// Convert between opt types
val.MustGetNull() // returns null.Val, but lossy, panics if val == null
val.MustGetOmit() // returns omit.Val, but lossy, panics if val == omit
omitnull.FromNull(null.From(5))
omitnull.FromOmit(omit.From(5))

// Converting between incompatible types
o := omit.From(5)
n := null.From(6)
_ = null.From(o.MustGet()) // This panics when o == unset
_ = omit.From(n.MustGet()) // This panics when n == null
_ = null.FromBool(o.Get()) // This conflates null/omitted, technically incorrect
_ = omit.FromBool(n.Get()) // This conflates null/omitted, technically incorrect
```

## Permutations

This package provides a package for each permutation of the problem space.

* If you have a value that can be `null | value` use `null`
* If you have a value that can be `unset | value` (non-null, but omittable) then use `omit`
* If you have a value that can be any of `unset | null | value` use omitnull.

Consider the following Rust types to illustrate the point:

```rust
enum Nullable<T> {
	Null,
	Value(T)
}

enum Omittable<T> {
	Omitted,
	Value(T)
}

enum OmittableNullable<T> {
	Omitted,
	Null,
	Value(T)
}
```

Why do Omittable vs Nullable have to exist despite being seemingly identical?
While it's true that they are nearly identical the behavior is slightly
different. Imagine a scenario where you are reading JSON into a field of each
type, the `Omittable` should panic or error when fed a `null` value, because
`null` is not a valid value for it to hold. It can either be
unset/void/undefined/omitted or a valid `T` value. For this reason there are
separate types in this package for all 3 permutations.

## On Correctness

Why consider omitted? Why isn't `null | value` good enough?

In Go we have structs that contain fields like the following:

```go
type Example struct {
	Age int32
}
```

This is simple and straightforward, the integer is always present and it must
be in the range of valid values for an int32. Thanks to the Go zero-value
to avoid uninitialized memory this will always be the case. Consider these
examples:

```go
var e Example
e = Example{Age: 5} // Explicit set to 5
e = Example{}       // No value provided means 0-value takes over and age=0
```

This makes good sense for Go. However this breaks down a bit when interoperating
with other systems and as Go is positioned well as a network service language
this happens fairly frequently. The general approach in Go to create a JSON API
would have us unmarshalling json objects into structs. We use structs because
it's the easiest for Go folks to work with, much easier than `map[string]any`
everywhere. So the example you'll see in typical Go examples is similar to:

```go
var e Example
json.Unmarshal([]byte(`{"age":5}`), &e)
```

In the example above, everything behaves normally and we're happy. But what
happens if we do the following:

```go
var e Example
json.Unmarshal([]byte(`{"age": null}`), &e)
```

Unmarshal doesn't fail, and we don't really want it to. But we now have a
problem. Although the JSON did not specify a valid integer value we *do* have a
value in age. It's of course the 0-value since that's the default when e was
instantiated on the first line, and nothing overwrote it because `null` can't
fit in to an `int32`. But did the JSON contain a 0? No, it did not.

The general approach to handling this is to bring in `null` by either having
a custom type or by using a pointer. This is how it's done for most sql
solutions in Go today.

```go
type ExampleNull struct {
	Age *int32 // or sql.NullInt32
}
```

If we try the same code as above with this new struct:

```go
var e Example
json.Unmarshal([]byte(`{"age":null}`), &e)
```

Now `e.Age` will be `nil`. We correctly don't have a value here. But what about
this example:

```go
var e Example
json.Unmarshal([]byte(`{}`), &e)
```

`e.Age` is the exact same as before, `nil`, but is this correct? We have not
been provided a value at all, not a null, nor an integer. The question then is
does this loss of information have an impact?

This becomes a problem when interoperating with languages that have an explicit
way to have a struct or object value be undefined or unset. Go maps have this
property, a key either exists or it does not exist in the map, but a Go struct
field always exists. Despite this problem in order to avoid completely
unstructured data, gain a modicum of type safety and minimal validation it's
generally still recommended to use a struct when writing in Go.

Let's take a look at typescript's possible values the most permissive kind of
field inside an object:

```typescript
interface Example {
	age?: null | number; // same as: undefined | null | number
}
```

In Javascript/Typescript we can have three distinct value types inside the age
field. Contrast this with our Go solution where we only have two. Although this
should be sufficient to prove interoperability problems as
JSON/Javascript/Typescript can create an object/struct with missing keys that
clearly cannot be modeled well in Go with or without pointers, it's not clear
why mapping JSON's three values to the two values of Go structs is a problem,
continue to the next section to understand at least one context where it can
become an issue.

## API Contexts

From the previous section Go structs can hold two values if we use a pointer or
a null type, but what about the third case in Javascript (undefined / omitted).
This information is lost and coerced as `nil/null` in a typical Go scenario.

Consider a database which is storing a `User` with an age. Age is nullable
because the user may not want to store their age in our system.

In an API request, how might the user signal that they want to update the value
age to 5, remove the value, or change other values without affecting others?
This is easily achieved using a partial update and could look something like the
following API payloads:

```javascript
{"name": "hello", "age": 5}     // set name = "hello", set age = 5,    do not set other fields
{"name": "hello", "age": null}  // set name = "hello", set age = null, do not set other fields
{"name": "hello"}               // set name = "hello", do not set age, do not set other fields
```

The above works great, but how are we to know that the user wanted to set
something explicitly to `null`, or if they just omitted the value? In Go with
our two-value field types, that loss of information means we cannot know what
the user's intent was because each field has three possible values:

* omitted: don't update
* null: set to null
* value: set to value

There is of course the other option of always bringing the entire object out
and resaving all of its fields every time. These are whole object updates and
are less efficient, and also have a weakness around stale updates that requires
special handling.

```javascript
let obj = getObjFromAPI(); // returns {"name":"hello", "age":5}
obj.age = 6;               // mutate the object
saveObjToAPI(obj);         // save the object
```

A race can occur however, consider the scenario:

```javascript
// clientA
let obj = getObjFromAPI(); // returns {"name":"hello", "age":5}
obj.age = 6;               // mutate the object
saveObjToAPI(obj);         // save the object

// clientB
let obj = getObjFromAPI(); // returns {"name":"hello", "age":5}
obj.name = "hi";           // mutate the object
saveObjToAPI(obj);         // save the object
```

In the above example both clients receive the same object, and each one
overwrites a single field but now one of the updates will be clobbered. If
client A saves first, then B's update clobbers A's update and the change of
`age` to `6` is lost. If B saves first then A's update clobbers B's update
and the change of `name` to `"hi"` is lost.

Typically then you need to have some sort of version field on each object and
reject updates that do not have the latest version, but this is a decent amount
of additional work. Partial updates are also race-y but at a field level so it's
much less likely to produce a strange result, imagine the same scenario with
partial updates:

```javascript
// clientA
let obj = getObjFromAPI();     // returns {"name":"hello", "age":5}
updateObjectInAPI({"age": 6}); // update object

// clientB
let obj = getObjFromAPI();         // returns {"name":"hello", "age":5}
updateObjectInAPI({"name": "hi"}); // update object
```

In this case there's no conflict at all. Consider when there is a conflict when
both clients want to update the same field:

```javascript
// clientA
let obj = getObjFromAPI();     // returns {"name":"hello", "age":5}
updateObjectInAPI({"age": 6}); // update object

// clientB
let obj = getObjFromAPI();     // returns {"name":"hello", "age":5}
updateObjectInAPI({"age": 7}); // update object
```

Either way it still makes reasonable sense to a user when the object update
occurs in either order because the last user's update will remain which
is predictable and understandable.

Partial updates is a useful and efficient pattern for updating objects. Its also
quite convenient for clients who may not have the entire resource, but still
wish to update a subset of the object's fields. This is one common real-world
use case for the `omitnull.Val` type which can house all three field value
types.

# License

This code is licensed mostly with MIT but some files contain Go's BSD-3 Clause
as it borrows functions originally found in the standard library.
