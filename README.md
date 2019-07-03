# Go Test Timer

`gott` is a command-line utility which finds your most time-consuming Go tests.

## Install

`go install go.coder.com/gott`

## Example

```
go test -v | gott

 ✔    --- TestA/A        20ms
 ✔    --- --- TestA/B/BB 10ms
 ✔    --- TestA/B        20ms
 ✔    TestA              40ms
 ✖    --- TestB/BigFail  40ms
 ✖    TestB              50ms
 ✖    TOTAL              368ms
```