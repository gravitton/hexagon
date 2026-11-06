package main

import (
	"fmt"

	geom "github.com/gravitton/geometry"
	hex "github.com/gravitton/hexagon"
)

var heroMap = []string{
	"...............",
	"..#.....#......",
	"..#.....#....G.",
	"..#..~~.#......",
	"....S~~........",
	"..#..~~...##...",
	"..#.......#....",
	"..........#....",
	"...............",
}

type world struct {
	ground []hex.Hex
	rough  map[hex.Hex]bool
	walls  map[hex.Hex]bool
	scout  hex.Hex
	goal   hex.Hex
}

func hero(t theme) []*panel {
	w := worldOf(heroMap)

	return []*panel{w.field(t), w.line(t), w.sector(t)}
}

func (w world) field(t theme) *panel {
	visible := set()
	for h := range w.scout.FieldOfViewSeq(6, contains(w.walls)) {
		visible[h] = true
	}

	p := w.backdrop(t, "s.FieldOfView(6, blocked)", func(h hex.Hex) string {
		switch {
		case visible[h]:
			return t.accentSoft
		case w.scout.DistanceTo(h) <= 6:
			return t.shade
		default:
			return ""
		}
	})
	p.cell(w.scout, t.accent)
	p.cell(w.goal, t.cool)

	return p
}

func (w world) line(t theme) *panel {
	seen := w.scout.HasLineOfSight(w.goal, contains(w.walls))
	p := w.backdrop(t, "s.HasLineOfSight(g, blocked)", func(hex.Hex) string {
		return ""
	})
	p.note = verdict(seen)
	for _, h := range w.scout.LineTo(w.goal) {
		p.outline(h, t.accent, 1.5)
	}
	if seen {
		p.segment(w.scout, w.goal, t.accent, false)
	} else {
		p.segment(w.scout, w.goal, t.muted, true)
	}
	p.cell(w.scout, t.accent)
	p.cell(w.goal, t.cool)

	return p
}

func (w world) sector(t theme) *panel {
	facing := w.scout.DirectionTo(w.goal)
	visible := set()
	for h := range w.scout.FieldOfViewSectorSeq(6, facing, 2, contains(w.walls)) {
		visible[h] = true
	}

	p := w.backdrop(t, "s.FieldOfViewSector(6, toward, 2, blocked)", func(h hex.Hex) string {
		switch {
		case visible[h]:
			return t.accentSoft
		case w.scout.SectorContains(6, facing, 2, h):
			return t.shade
		default:
			return ""
		}
	})
	p.cell(w.scout, t.accent)
	p.cell(w.goal, t.cool)

	return p
}

func (w world) backdrop(t theme, caption string, fill func(hex.Hex) string) *panel {
	p := newPanel(t, pointyTop, 10.5, caption)
	for _, h := range w.ground {
		switch {
		case w.walls[h]:
			p.cell(h, t.wall)
		case fill(h) != "":
			p.cell(h, fill(h))
		case w.rough[h]:
			p.cell(h, t.leafSoft)
		default:
			p.cell(h, t.cell)
		}
	}

	return p
}

func worldOf(rows []string) world {
	w := world{rough: set(), walls: set()}
	for y, row := range rows {
		for x, tile := range row {
			h := hex.OffsetOddR.From(geom.Pt(x, y))
			w.ground = append(w.ground, h)
			switch tile {
			case '#':
				w.walls[h] = true
			case '~':
				w.rough[h] = true
			case 'S':
				w.scout = h
			case 'G':
				w.goal = h
			case '.':
			default:
				panic(fmt.Sprintf("diagrams: unknown tile %q on the hero map", tile))
			}
		}
	}

	return w
}
