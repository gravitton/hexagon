package main

import (
	"math"
	"slices"
	"strconv"
	"strings"

	geom "github.com/gravitton/geometry"
	hex "github.com/gravitton/hexagon"
)

var figures = []figure{
	{"hero", hero},
	{"coordinates", coordinates},
	{"arithmetic", arithmetic},
	{"directions", directions},
	{"neighbors", neighbors},
	{"traversal", traversal},
	{"circles", circles},
	{"chunks", chunks},
	{"sectors", sectors},
	{"sight", sight},
	{"walk", walk},
	{"field", field},
	{"transform", transform},
	{"systems-pointy", systemsPointy},
	{"systems-flat", systemsFlat},
	{"rounding", rounding},
}

var origin = hex.Pt(0, 0)

func coordinates(t theme) []*panel {
	fills := []string{t.accent, t.accentSoft, t.cell}
	inks := []string{t.onAccent, t.ink, t.ink}

	axial := newPanel(t, pointyTop, 28, "hex.Pt(q, r)")
	cube := newPanel(t, pointyTop, 28, "h.QRS()")
	length := newPanel(t, pointyTop, 28, "h.Length()")
	for _, h := range origin.Range(2) {
		axial.cell(h, t.cell)
		axial.text(h, 0, font{size: 11}, span{strconv.Itoa(h.Q), t.accent}, span{",", t.muted}, span{strconv.Itoa(h.R), t.cool})

		cube.cell(h, t.cell)
		at := cube.center(h)
		cube.canvas.text(at.AddXY(10, -7), font{size: 10}, span{strconv.Itoa(h.Q), t.accent})
		cube.canvas.text(at.AddXY(0, 11), font{size: 10}, span{strconv.Itoa(h.R), t.cool})
		cube.canvas.text(at.AddXY(-10, -7), font{size: 10}, span{strconv.Itoa(h.S()), t.leaf})

		length.cell(h, fills[h.Length()])
		length.text(h, 0, font{size: 13}, span{strconv.Itoa(h.Length()), inks[h.Length()]})
	}

	return []*panel{axial, cube, length}
}

func arithmetic(t theme) []*panel {
	a, b := hex.Pt(2, -1), hex.Pt(-1, 2)

	add := backdrop(t, 18, origin, 3, "a.Add(b)")
	sum := a.Add(b)
	add.cell(origin, t.shade)
	add.cell(a, t.accentSoft)
	add.cell(b, t.coolSoft)
	add.cell(sum, t.leafSoft)
	add.pointer(origin, b, t.cool, true)
	add.pointer(origin, a, t.accent, false)
	add.pointer(a, sum, t.cool, false)
	add.label(a, "a", t.ink)
	add.label(b, "b", t.ink)
	add.label(sum, "+", t.ink)

	subtract := backdrop(t, 18, origin, 3, "a.Subtract(b)")
	difference := a.Subtract(b)
	subtract.cell(origin, t.shade)
	subtract.cell(a, t.accentSoft)
	subtract.cell(b, t.coolSoft)
	subtract.cell(difference, t.leafSoft)
	subtract.pointer(b, a, t.leaf, true)
	subtract.pointer(origin, difference, t.leaf, false)
	subtract.label(a, "a", t.ink)
	subtract.label(b, "b", t.ink)
	subtract.label(difference, "−", t.ink)

	multiply := backdrop(t, 18, origin, 4, "a.Multiply(2)")
	multiply.cell(origin, t.shade)
	for _, factor := range []int{2, -1} {
		product := a.Multiply(factor)
		multiply.cell(product, t.leafSoft)
		multiply.pointer(origin, product, t.leaf, false)
		multiply.label(product, "×"+strings.Replace(strconv.Itoa(factor), "-", "−", 1), t.ink)
	}
	multiply.cell(a, t.accentSoft)
	multiply.label(a, "a", t.ink)

	return []*panel{add, subtract, multiply}
}

func directions(t theme) []*panel {
	compasses := []struct {
		orientation orientation
		caption     string
		names       [6]string
	}{
		{pointyTop, "PointyTop aliases", [6]string{"East", "SouthEast", "SouthWest", "West", "NorthWest", "NorthEast"}},
		{flatTop, "FlatTop aliases", [6]string{"SouthEast", "South", "SouthWest", "NorthWest", "North", "NorthEast"}},
	}

	panels := make([]*panel, 0, len(compasses))
	for _, compass := range compasses {
		p := newPanel(t, compass.orientation, 38, compass.caption)
		p.cell(origin, t.accent)
		for i, direction := range hex.Directions() {
			h := origin.Neighbor(direction)
			p.cell(h, t.sectors[i])
			degrees := int(math.Round(direction.Angle() * 180 / math.Pi))
			p.text(h, -11, font{size: 12}, span{direction.String(), t.ink})
			p.text(h, 3, font{size: 9.5, family: sans}, span{compass.names[i], t.muted})
			p.text(h, 16, font{size: 9}, span{strconv.Itoa(i) + " · " + strconv.Itoa(degrees) + "°", t.muted})
		}
		panels = append(panels, p)
	}

	return panels
}

func neighbors(t theme) []*panel {
	adjacent := backdrop(t, 15, origin, 3, "h.Neighbors()")
	for i, h := range origin.Neighbors() {
		adjacent.cell(h, t.sectors[i])
		adjacent.label(h, strconv.Itoa(i), t.ink)
	}
	adjacent.cell(origin, t.accent)

	diagonal := backdrop(t, 15, origin, 3, "h.DiagonalNeighbors()")
	for i, h := range origin.DiagonalNeighbors() {
		diagonal.cell(h, t.sectors[i])
		diagonal.label(h, strconv.Itoa(i), t.ink)
	}
	diagonal.cell(origin, t.accent)

	toward := newPanel(t, pointyTop, 15, "h.DirectionTo(target)")
	for _, h := range origin.Spiral(4)[1:] {
		toward.cell(h, t.sectors[origin.DirectionTo(h)])
	}
	toward.cell(origin, t.accent)

	return []*panel{adjacent, diagonal, toward}
}

func traversal(t theme) []*panel {
	disc := backdrop(t, 15, origin, 3, "h.Range(2)")
	disc.cells(origin.Range(2), t.accentSoft)
	disc.cell(origin, t.accent)

	ring := backdrop(t, 15, origin, 3, "h.Ring(2)")
	ring.cells(origin.Ring(2), t.accentSoft)
	ring.path(append(origin.Ring(2), origin.Ring(2)[0]), t.accent)
	for i, h := range origin.Ring(2) {
		ring.label(h, strconv.Itoa(i), t.ink)
	}
	ring.cell(origin, t.accent)

	spiral := backdrop(t, 15, origin, 3, "h.Spiral(2)")
	spiral.cells(origin.Spiral(2), t.accentSoft)
	spiral.path(origin.Spiral(2), t.accent)
	for i, h := range origin.Spiral(2) {
		spiral.label(h, strconv.Itoa(i), t.ink)
	}

	a, b := hex.Pt(-3, 1), hex.Pt(3, -2)
	line := backdrop(t, 15, origin, 3, "a.LineTo(b)")
	line.cells(a.LineTo(b), t.accentSoft)
	line.segment(a, b, t.ink, false)
	line.cell(a, t.accent)
	line.cell(b, t.accent)
	line.label(a, "a", t.onAccent)
	line.label(b, "b", t.onAccent)

	return []*panel{disc, ring, spiral, line}
}

func circles(t theme) []*panel {
	exact := disc(t, math.Sqrt(3), "h.Circle(math.Sqrt(3))")
	exact.note = "the range of one and the diagonals"

	whole := disc(t, 4, "h.Circle(4)")
	whole.note = "the hexes of Range(4)"

	half := disc(t, 4.5, "h.Circle(4.5)")
	half.note = "Range(4) and twelve hexes of ring five"

	return []*panel{exact, whole, half}
}

func disc(t theme, radius float64, caption string) *panel {
	p := backdrop(t, 13, origin, 5, caption)
	p.cells(origin.Circle(radius), t.accentSoft)
	p.cell(origin, t.accent)
	p.label(origin, "h", t.onAccent)

	step := p.center(hex.Pt(1, 0)).Subtract(p.center(origin)).Length()
	p.circle(p.center(origin), radius*step, style{stroke: t.muted, width: 1.2, dashed: true})

	return p
}

func chunks(t theme) []*panel {
	h, view := hex.Pt(3, 0), hex.Pt(2, -1)
	chunk, center, local := h.Chunk(2), h.Chunk(2).ChunkCenter(2), h.ChunkLocal(2)

	holding := tiling(t, view, "h.Chunk(2)")
	for _, member := range center.Range(2) {
		holding.outline(member, t.accent, 2)
	}
	holding.cell(h, t.accent)
	holding.label(h, "h", t.onAccent)
	holding.note = "the chunk " + chunk.String()

	wrap := tiling(t, view, "h.ChunkLocal(2)")
	wrap.segment(center, h, t.accent, false)
	wrap.segment(origin, local, t.accent, false)
	wrap.cell(center, t.cool)
	wrap.label(center, "c", t.onAccent)
	wrap.cell(origin, t.cool)
	wrap.label(origin, "o", t.onAccent)
	wrap.cell(h, t.accent)
	wrap.label(h, "h", t.onAccent)
	wrap.cell(local, t.accent)
	wrap.label(local, "h", t.onAccent)
	wrap.note = local.String() + " from c, where h wraps onto the map around o"

	index := tiling(t, view, "h.ChunkIndex(2)")
	for _, cell := range view.Range(5) {
		if cell.ChunkIndex(2) == h.ChunkIndex(2) {
			index.outline(cell, t.accent, 2)
		}
		index.label(cell, strconv.Itoa(cell.ChunkIndex(2)), t.ink)
	}
	index.cell(h, t.accent)
	index.label(h, strconv.Itoa(h.ChunkIndex(2)), t.onAccent)
	index.note = "every chunk numbers its hexes alike"

	return []*panel{holding, wrap, index}
}

func tiling(t theme, view hex.Hex, caption string) *panel {
	p := newPanel(t, pointyTop, 14, caption)
	fills := [3]string{t.accentSoft, t.coolSoft, t.leafSoft}
	for _, h := range view.Range(5) {
		chunk := h.Chunk(2)
		p.cell(h, fills[geom.FloorMod(chunk.Q-chunk.R, 3)])
	}

	return p
}

func sectors(t theme) []*panel {
	cone := backdrop(t, 16, origin, 4, "h.Sector(4, hex.QPlus, 1)")
	cone.cells(origin.Sector(4, hex.QPlus, 1), t.accentSoft)
	for i, h := range origin.Sector(4, hex.QPlus, 1) {
		cone.label(h, strconv.Itoa(i), t.ink)
	}
	if target := hex.Pt(2, -3); origin.SectorContains(4, hex.QPlus, 1, target) {
		cone.outline(target, t.accent, 2)
	}
	cone.cell(origin, t.accent)
	cone.label(origin, "h", t.onAccent)

	wide := backdrop(t, 16, origin, 4, "h.Sector(4, hex.QPlus, 2)")
	wide.cells(origin.Sector(4, hex.QPlus, 2), t.accentSoft)
	wide.cell(origin, t.accent)
	wide.label(origin, "h", t.onAccent)

	first, second := origin.Sector(4, hex.QPlus, 1), origin.Sector(4, hex.SMinus, 1)
	border := backdrop(t, 16, origin, 4, "neighbors share a border")
	border.cells(first, t.accentSoft)
	border.cells(second, t.coolSoft)
	for _, h := range first {
		if slices.Contains(second, h) {
			border.cell(h, t.leafSoft)
		}
	}
	border.cell(origin, t.accent)
	border.label(origin, "h", t.onAccent)

	return []*panel{cone, wide, border}
}

func sight(t theme) []*panel {
	a, b := origin, hex.Pt(2, 2)
	cases := []struct {
		name     string
		blocking []hex.Hex
	}{
		{"through", []hex.Hex{hex.Pt(1, 1)}},
		{"one beside an edge", []hex.Hex{hex.Pt(1, 0)}},
		{"both beside an edge", []hex.Hex{hex.Pt(1, 0), hex.Pt(0, 1)}},
	}

	panels := make([]*panel, 0, len(cases)+1)
	for _, c := range cases {
		seen := a.HasLineOfSight(b, contains(set(c.blocking...)))
		p := backdrop(t, 18, hex.Pt(1, 1), 2, c.name+": "+verdict(seen))
		p.cells(c.blocking, t.wall)
		if seen {
			p.segment(a, b, t.accent, false)
		} else {
			p.segment(a, b, t.muted, true)
		}
		p.cell(a, t.accent)
		p.cell(b, t.accentSoft)
		p.label(a, "a", t.onAccent)
		p.label(b, "b", t.ink)
		panels = append(panels, p)
	}

	return panels
}

func walk(t theme) []*panel {
	cases := []struct {
		target hex.Hex
		wall   hex.Hex
	}{
		{hex.Pt(3, 1), hex.Pt(2, 0)},
		{hex.Pt(2, 2), hex.Pt(1, 0)},
		{hex.Pt(2, 2), hex.Pt(0, 1)},
	}

	panels := make([]*panel, 0, len(cases))
	for _, c := range cases {
		blocked := contains(set(c.wall))
		line := origin.LineTo(c.target)
		walked := !slices.ContainsFunc(line[1:len(line)-1], blocked)
		seen := origin.HasLineOfSight(c.target, blocked)

		p := backdrop(t, 18, hex.Pt(1, 1), 2, "walking LineTo: "+verdict(walked))
		p.note = "HasLineOfSight: " + verdict(seen)
		p.cell(c.wall, t.wall)
		for _, h := range line {
			p.outline(h, t.cool, 1.5)
		}
		if seen {
			p.segment(origin, c.target, t.accent, false)
		} else {
			p.segment(origin, c.target, t.muted, true)
		}
		p.cell(origin, t.accent)
		p.cell(c.target, t.accentSoft)
		p.label(origin, "a", t.onAccent)
		p.label(c.target, "b", t.ink)
		panels = append(panels, p)
	}

	return panels
}

func field(t theme) []*panel {
	source := origin
	walls := []hex.Hex{
		hex.Pt(1, -2), hex.Pt(2, -2), hex.Pt(2, -1),
		hex.Pt(-2, 0), hex.Pt(-2, 1),
		hex.Pt(0, 2), hex.Pt(1, 2),
		hex.Pt(-1, -2),
	}
	blocked := contains(set(walls...))

	candidates := []hex.Hex{hex.Pt(4, -3), hex.Pt(4, -1), hex.Pt(-4, 1), hex.Pt(-1, 4), hex.Pt(3, 1), hex.Pt(-2, -2), hex.Pt(1, -4)}
	seen := set()
	for _, c := range candidates {
		seen[c] = source.HasLineOfSight(c, blocked)
	}
	pairwise := backdrop(t, 14, origin, 5, "a.HasLineOfSight(c, blocked)")
	pairwise.cells(walls, t.wall)
	for _, c := range candidates {
		if seen[c] {
			pairwise.segment(source, c, t.accent, false)
		} else {
			pairwise.segment(source, c, t.muted, true)
		}
	}
	for _, c := range candidates {
		if seen[c] {
			pairwise.cell(c, t.accentSoft)
		} else {
			pairwise.cell(c, t.shade)
		}
	}
	pairwise.cell(source, t.accent)
	pairwise.label(source, "a", t.onAccent)

	visible := set()
	for h := range source.FieldOfViewSeq(5, blocked) {
		visible[h] = true
	}
	sweep := newPanel(t, pointyTop, 14, "a.FieldOfView(5, blocked)")
	for _, h := range origin.Range(5) {
		switch {
		case blocked(h):
			sweep.cell(h, t.wall)
		case visible[h]:
			sweep.cell(h, t.accentSoft)
		default:
			sweep.cell(h, t.shade)
		}
	}
	sweep.cell(source, t.accent)
	sweep.label(source, "a", t.onAccent)

	facing := set()
	for h := range source.FieldOfViewSectorSeq(5, hex.SMinus, 2, blocked) {
		facing[h] = true
	}
	cone := backdrop(t, 14, origin, 5, "a.FieldOfViewSector(5, hex.SMinus, 2, blocked)")
	for _, h := range source.Sector(5, hex.SMinus, 2) {
		switch {
		case blocked(h):
			cone.cell(h, t.wall)
		case facing[h]:
			cone.cell(h, t.accentSoft)
		default:
			cone.cell(h, t.shade)
		}
	}
	cone.cells(walls, t.wall)
	cone.cell(source, t.accent)
	cone.label(source, "a", t.onAccent)

	return []*panel{pairwise, sweep, cone}
}

func transform(t theme) []*panel {
	a, b := hex.Pt(3, -1), hex.Pt(1, 0)

	turn := backdrop(t, 16, origin, 3, "a.Turn(steps)")
	for steps := range 6 {
		turn.cell(a.Turn(steps), t.accentSoft)
		turn.label(a.Turn(steps), strconv.Itoa(steps), t.ink)
	}
	turn.cell(origin, t.shade)
	turn.cell(a, t.accent)
	turn.label(a, "a", t.onAccent)

	around := backdrop(t, 16, origin, 3, "a.TurnAround(b, steps)")
	for steps := range 6 {
		around.cell(a.TurnAround(b, steps), t.coolSoft)
		around.label(a.TurnAround(b, steps), strconv.Itoa(steps), t.ink)
	}
	around.cell(b, t.cool)
	around.label(b, "b", t.onAccent)
	around.cell(a, t.accent)
	around.label(a, "a", t.onAccent)

	reflect := backdrop(t, 16, origin, 3, "a.ReflectQ(), ReflectR(), ReflectS()")
	mirrors := []struct {
		name  string
		axis  hex.Hex
		image hex.Hex
	}{
		{"Q", hex.Pt(2, -1), a.ReflectQ()},
		{"R", hex.Pt(-1, 2), a.ReflectR()},
		{"S", hex.Pt(1, 1), a.ReflectS()},
	}
	for _, mirror := range mirrors {
		reflect.cell(mirror.image, t.leafSoft)
	}
	reflect.cell(a, t.accent)
	for _, mirror := range mirrors {
		reflect.axis(mirror.axis, 1.75, strings.ToLower(mirror.name))
	}
	for _, mirror := range mirrors {
		reflect.label(mirror.image, mirror.name, t.ink)
	}
	reflect.label(a, "a", t.onAccent)

	return []*panel{turn, around, reflect}
}

func systemsPointy(t theme) []*panel {
	return systems(t, pointyTop, []view{
		{hex.Axial, hex.OffsetOddR, "q, r as stored"},
		{hex.OffsetOddR, hex.OffsetOddR, "odd rows shift right"},
		{hex.OffsetEvenR, hex.OffsetEvenR, "even rows shift right"},
		{hex.DoubleWidth, hex.OffsetOddR, "columns doubled"},
	}, func(index geom.Point[int]) int {
		return index.Y
	})
}

func systemsFlat(t theme) []*panel {
	return systems(t, flatTop, []view{
		{hex.Axial, hex.OffsetOddQ, "q, r as stored"},
		{hex.OffsetOddQ, hex.OffsetOddQ, "odd columns shift down"},
		{hex.OffsetEvenQ, hex.OffsetEvenQ, "even columns shift down"},
		{hex.DoubleHeight, hex.OffsetOddQ, "rows doubled"},
	}, func(index geom.Point[int]) int {
		return index.X
	})
}

type view struct {
	system hex.CoordinateSystem
	patch  hex.CoordinateSystem
	note   string
}

func systems(t theme, o orientation, views []view, shifted func(geom.Point[int]) int) []*panel {
	panels := make([]*panel, 0, len(views))
	for _, v := range views {
		p := newPanel(t, o, 18, v.system.String())
		p.note = v.note
		for row := range 4 {
			for column := range 5 {
				index := geom.Pt(column, row)
				h := v.patch.From(index)
				fill := t.cell
				if (shifted(index)%2 == 1) == (v.patch == hex.OffsetOddR || v.patch == hex.OffsetOddQ) {
					fill = t.shade
				}
				label := v.system.To(h)
				p.cell(h, fill)
				p.text(h, 0, font{size: 9.5}, span{strconv.Itoa(label.X), t.accent}, span{",", t.muted}, span{strconv.Itoa(label.Y), t.cool})
			}
		}
		panels = append(panels, p)
	}

	return panels
}

func rounding(t theme) []*panel {
	f := hex.FracPt(1.4, -1.8)
	colors := []string{t.accent, t.cool, t.leaf}

	p := newPanel(t, pointyTop, 30, "f.Round()")
	area := origin.Range(2)
	for _, h := range area {
		p.cell(h, t.cell)
	}
	for y := p.bounds.Min.Y + 3; y < p.bounds.Max.Y; y += 6 {
		for x := p.bounds.Min.X + 3; x < p.bounds.Max.X; x += 6 {
			at := geom.Pt(x, y)
			h := hex.Axial.From(p.layout.IndexAt(at))
			if h.Length() <= 2 {
				p.circle(at, 1.3, style{fill: colors[geom.FloorMod(h.Q-h.R, 3)]})
			}
		}
	}
	for _, h := range area {
		p.outline(h, t.edge, 1)
	}

	at, rounded := fractional(p, f), p.center(f.Round())
	p.line(at, rounded, style{stroke: t.ink, width: 1.5, dashed: true})
	p.circle(rounded, 3, style{fill: t.ink})
	p.circle(at, 4, style{fill: t.ink})
	p.canvas.text(at.Add(at.Subtract(rounded).Normalize().Multiply(13)), font{size: 13}, span{"f", t.ink})

	return []*panel{p}
}

func fractional(p *panel, f hex.FractionalHex) geom.Point[float64] {
	zero := p.center(origin)
	q := p.center(hex.Pt(1, 0)).Subtract(zero)
	r := p.center(hex.Pt(0, 1)).Subtract(zero)

	return zero.Add(q.Multiply(f.Q)).Add(r.Multiply(f.R))
}

func set(hexes ...hex.Hex) map[hex.Hex]bool {
	s := make(map[hex.Hex]bool, len(hexes))
	for _, h := range hexes {
		s[h] = true
	}

	return s
}

func contains(s map[hex.Hex]bool) func(hex.Hex) bool {
	return func(h hex.Hex) bool {
		return s[h]
	}
}

func verdict(seen bool) string {
	if seen {
		return "seen"
	}

	return "hidden"
}
