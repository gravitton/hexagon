module github.com/gravitton/hexagon/docs/images/generator

go 1.27

require (
	github.com/gravitton/geometry v1.19.0
	github.com/gravitton/grid v1.4.0
	github.com/gravitton/hexagon v1.8.1
)

require github.com/gravitton/container v1.0.1 // indirect

replace github.com/gravitton/hexagon => ../../..
