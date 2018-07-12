来源: https://hackmd.io/s/BJ_wJtbOb

## Workshop

- general review of error, panic(), defer, recover().
- errors
  - sentinel errors let us know what’s happening in our app, e.g io.EOF or sql.ErrNoRows. They are interesting for control flow. The speaker uses the stdlib ones but does not create them in his code. Does not say the cons.
  - https://github.com/pkg/errors to wrap errors and see tracing data.
- concurrency:
  - `WaitGroup`: don’t do `wg.Add(1)` inside the goroutine its doing the `Done()`.
  - example of https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
  - buffered channels
    -  they can hide a design mistake, so be really careful with them.
    -  what size do you use? Why?
    -  why you don’t want to block the sender?
    -  what about backpressure (more writes in than writes out)? You can avoid memory overflow, but what limit you put? Why that limit?
    -  they can hide a message loss problem.
  - `select`: be careful of ranging with a select with an empty default branch, the CPU will suffer.
- testing
  - https://github.com/smartystreets/goconvey is a cool tool. (Note: I knew about its BDD testing lib which I don’t like, but the reactive testing tool is cool).
  - dependency management: dep
  - context: API, cancelation propagation, trees
  - build tools, ldflags
  - protobufs & gRPC

## Writing beautiful packages

- single interface method interfaces allow functions to implement the interface
  - e.g http.Handler and http.HandlerFunc
- https://rakyll.org/style-packages
- for libs, leave concurrency to the user
- learn from the stdlib
- the pkg name is part of the API, and naming is hard
- take advantadge of zero values and avoid constructors if you can - e.g if you have not many fields
- allow injecting http.Client
- avoid global state and init

## Concurrency patterns

- blocking channels
  - channel blocks when no data:
    - no receiver, for unbuffered buffers
    - no sender, for all channels
  - blocking channels are good for synchronizing goroutines
  - blocking can lead to
    - deadlocks
    - scaling problems: adding more blocking goroutines can lead to worse performance
- closing channels
  - closing a channel sends a special closed message to all the readers
  - send after close panics
  - closing twice panics
  - closed makes the reader receive two things: the zero value of the type of the channel, and false.
  - the receiver always knows if the channel is closed, the sender does not.
    - corolary: always close the channel from the receiving side, not from the sending side!
- select
  - order of cases does not matter
  - a default case exists that’s executed if the other cases are blocked
  - to make channels nonblocking, use time.After in select or default
- channels are streams of data; combining streams is powerful. Depending on the shape of the combination there are different patterns:
  - fan-out (1:N): select with writes to several channels sends to the first non-blocking channel
  - turn-out (N:M): select with multiple reads, to select with multiple writes.
  - quit channel for cancelation
- channel failures
  - deadlocks
  - memory copying and performance
  - pasing pointers and race conditions
  - caches are about sharing, a cache with channels is not a good idea, use mutex instead?
    - RWLocks can reduce the problem
    - multiple mutexes will cause deadlocks
- three shades of code
  - blocking: the program can be locked
  - lock free: at least one part of the program makes progress
  - wait free: all parts of the program make progress
- sync.Atomic ops are thread-safe because they are based on CPU instructions
- Spin Lock or Spinning CAS: Compare and Swap in a loop. That’s how mutex are implemented (?)

## embedding

https://github.com/stabbycutyou/embeddingtalk

- inheritance does not exist in Go
  - embedding is not better than inheritance, is something different to resolve a different problem
  - embedding is for composing interfaces and structs
  - behaviour over lineage
  - no base class - super class relationship
    - “is-a” vs “has-a”
  - method dispatching has some edge cases
- what the spec says:
  - https://golang.org/ref/spec#Struct_types
  - A field declared with a type but no explicit field name is an anonymous field (colloquially called an embedded field). Such a field type must be specified as a type name T or as a pointer to a non-interface type name *T, and T itself may not be a pointer type. The unqualified type name acts as the field name.
  - https://golang.org/ref/spec#Interface_types
  - An interface T may use a (possibly qualified) interface type name E in place of a method specification. This is called embedding interface E in T; it adds all (exported and non-exported) methods of E to the interface T.

  - selectors: https://golang.org/ref/spec#Selectors
    - shallowest depth concept
    - promoted fields: https://golang.org/ref/spec#Struct_types
- examples:
  - emdedding a struct with another with the same field name
    - the field promoted is the one from the embedding not from the embedded struct
      - if you want to access the embedded one you need to do something like embedding.embedded.field
  - embedding multiple structs
  - nest embeddings
  - you cannot embed twice the same thing
  - same method and field example https://github.com/StabbyCutyou/embeddingtalk/blob/master/5.multiple/main.go
  - embed interfaces in structs
    - abstract behaviours instead of concrete behaviours
    - he does not say it but it’s greate for test doubles where one just use/overwrite some of the methods of the interface
      - warning: explodes at runtime
  - viewmodels to avoid marshaling some model fields https://github.com/StabbyCutyou/embeddingtalk/blob/master/7.marshalling/1.viewmodel/main.go
  - extending generated code
    - embedding generating code makes easy to extend it without modifying the generated code so it can be generated over and over again.
  - promoted methods are only called in the original receiver [warning]
