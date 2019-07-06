# Go Test Timer

`gott` finds the most time-consuming tests in large suites.

## Install

`go install go.coder.com/gott`

## Example

```
go test -v | gott

 ✔    --- TestA/A        20ms
 ✔    --- --- TestA/B/BB 10ms
 ✔    --- TestA/B        20ms  10ms
 ✔    TestA              40ms  40ms
 ✖    --- TestB/BigFail  50ms
 ✖    TestB              10ms  50ms
 ✖    TOTAL              562ms
```

## Usage

```
Usage: gott [flags]

Parses go test verbose output and produces a list of tests sorted by how time consuming they are.

gott flags:
	-c	omit entries that take less than this much time	(0s)
	-p	pass through go test output	(false)
```

## Parallel Children

`go test` reports how long a test function takes to return. Test functions don't wait on their
parallel children, so it can be difficult to answer the question _"How long did this test and
all of its children take?"_.

In `gott`,

The first column of durations show how long the test functions took to return.

The second column shows long it took for all children to return.

## Ranking Algorithm

`gott` orders outputted tests by `max(testTook, childrenTook)`.

The longest test cases are at the bottom to reduce scrolling.