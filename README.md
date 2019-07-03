# Go Test Timer

`gott` finds your most time-consuming Go tests.

## Install

`go install go.coder.com/gott`

## Example

```
go test -v | gott

 ✔    TOTAL                                                           9.73s
 ✔    Test_getConfiguration                                           5.71s 17.96s
 ✔    --- Test_getConfiguration/GetSetupMode                                2.32s
 ✔    --- --- Test_getConfiguration/GetSetupMode/NotSetupMode         1.12s
 ✔    --- --- Test_getConfiguration/GetSetupMode/SetupMode            1.2s
 ✔    --- Test_getConfiguration/PostSetupMode                               2.86s
 ✔    --- --- Test_getConfiguration/PostSetupMode/Authorized          1.34s
 ✔    --- --- Test_getConfiguration/PostSetupMode/Unauthorized        1.52s
 ✔    --- Test_getConfiguration/SetupMode                                   12.78s
 ✔    --- --- Test_getConfiguration/SetupMode/BadConfiguration        2.29s
 ✔    --- --- Test_getConfiguration/SetupMode/Authorized              2.33s
 ✔    --- --- Test_getConfiguration/SetupMode/OnlyConfigurationRoutes 2.52s
 ✔    --- --- Test_getConfiguration/SetupMode/NoAuth                  2.65s
 ✔    --- --- Test_getConfiguration/SetupMode/ExitSetupMode           2.99s
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

`gott` provides how long it took for the test function to return followed by the time it
took for all children to return.

## Ranking Algorithm

`gott` orders test output by `max(testTook, childrenTook)`.

The longest test cases are at the bottom to reduce scrolling.