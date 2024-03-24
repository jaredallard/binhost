module github.com/jaredallard/binhost

go 1.22

// https://github.com/jaredallard/xz/commit/c20e88619f2c09907bb17bf0b18bbe2627ee570a
replace github.com/jamespfennell/xz => github.com/jaredallard/xz v0.0.0-20240323042956-c20e88619f2c

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/jamespfennell/xz v0.1.2
	gotest.tools/v3 v3.5.1
)

require github.com/google/go-cmp v0.5.9 // indirect
