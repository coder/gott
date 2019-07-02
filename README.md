# Go Test Timer

`gott` is a command-line utility which finds your most time-consuming Go tests.

## Install

`go install go.coder.com/gott`

## Example

```
go test -v | gott

 FAIL   TOTAL        433ms
 FAIL   TestB         50ms
 FAIL   TestB/BigFail 40ms
 PASS   TestA         30ms
 PASS   TestA/A       20ms
 PASS   TestA/B       10ms
```