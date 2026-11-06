module github.com/gravitton/hexagon/docs/diagrams

go 1.27

require (
	github.com/gravitton/geometry v1.18.0
	github.com/gravitton/grid v1.3.0
	github.com/gravitton/hexagon v1.8.0
)

require github.com/gravitton/container v1.0.0 // indirect

replace github.com/gravitton/hexagon => ../..
