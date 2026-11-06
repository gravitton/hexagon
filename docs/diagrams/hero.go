package main

import (
	hex "github.com/gravitton/hexagon"
)

func hero(t theme) []*panel {
	walls := set(
		offset(4, 1), offset(4, 2), offset(5, 2),
		offset(9, 5), offset(10, 5), offset(10, 6), offset(9, 6),
		offset(11, 1), offset(12, 1),
		offset(5, 7), offset(6, 7),
		offset(13, 4),
		offset(17, 2), offset(18, 2), offset(18, 3),
		offset(20, 6), offset(21, 6),
	)
	scout := offset(7, 4)
	visible := set()
	for h := range scout.FieldOfViewSeq(7, contains(walls)) {
		visible[h] = true
	}
	target := scout
	for _, h := range scout.Ring(7) {
		if visible[h] && hex.OffsetOddR.To(h).X > hex.OffsetOddR.To(target).X {
			target = h
		}
	}

	ranger := offset(22, 3)
	reach := set(ranger.Range(3)...)

	p := newPanel(t, pointyTop, 15, "")
	for row := range 9 {
		for column := range 30 {
			h := offset(column, row)
			switch {
			case walls[h]:
				p.cell(h, t.wall)
			case visible[h]:
				p.cell(h, t.accentSoft)
			case scout.DistanceTo(h) <= 7:
				p.cell(h, t.shade)
			case reach[h]:
				p.cell(h, t.coolSoft)
			default:
				p.cell(h, t.cell)
			}
		}
	}
	for _, h := range scout.LineTo(target) {
		p.outline(h, t.accent, 2)
	}
	p.segment(scout, target, t.accent, false)
	for _, h := range ranger.Ring(3) {
		p.outline(h, t.cool, 1.5)
	}
	p.cell(scout, t.accent)
	p.cell(ranger, t.cool)
	p.cell(target, t.accent)

	return []*panel{p}
}
