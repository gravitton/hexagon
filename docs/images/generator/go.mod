module github.com/gravitton/hexagon/docs/images/generator

go 1.27

require (
	github.com/gravitton/geometry v1.20.1
	github.com/gravitton/grid v1.6.1
	github.com/gravitton/hexagon v1.9.1
)

require github.com/gravitton/container v1.1.0 // indirect

replace github.com/gravitton/hexagon => ../../..
