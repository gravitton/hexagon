//lint:file-ignore U1000 the generator of every package shares this file, and none draws with all of it

package main

import (
	"cmp"
	"fmt"
	"strings"

	geom "github.com/gravitton/geometry"
)

const (
	mono      = "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace"
	sans      = "Helvetica Neue, Helvetica, Arial, sans-serif"
	pageWidth = 960.0
	margin    = 12.0
	gap       = 32.0
	captionAt = 24.0
	noteAt    = 15.0
	captionPt = 11.0
	notePt    = 10.0
	glyph     = 0.62
)

type point = geom.Point[float64]

type style struct {
	fill   string
	stroke string
	width  float64
	dashed bool
}

type font struct {
	size   float64
	family string
	anchor string
	halo   bool
}

type span struct {
	text  string
	color string
}

type canvas struct {
	theme   theme
	caption string
	note    string
	body    strings.Builder
	bounds  geom.Box[float64]
	drawn   bool
	view    geom.Box[float64]
}

func newCanvas(t theme, caption string) canvas {
	return canvas{theme: t, caption: caption}
}

func (c *canvas) clip(view geom.Box[float64]) {
	c.bounds, c.drawn, c.view = view, true, view
}

func (c *canvas) cover(at point) {
	if !c.view.IsZero() {
		return
	}

	bounds := geom.Bx(at, at)
	if c.drawn {
		bounds = c.bounds.Union(bounds)
	}
	c.bounds, c.drawn = bounds, true
}

func (c *canvas) width() float64 {
	return max(c.bounds.Width(), widthOf(c.caption, captionPt), widthOf(c.note, notePt))
}

func (c *canvas) height() float64 {
	return c.bounds.Height()
}

func (c *canvas) polygon(points []point, s style) {
	c.shape("polygon", points, s)
}

func (c *canvas) polyline(points []point, s style) {
	s.fill = "none"
	c.shape("polyline", points, s)
}

func (c *canvas) shape(element string, points []point, s style) {
	coordinates := make([]string, len(points))
	for i, at := range points {
		c.cover(at)
		coordinates[i] = fmt.Sprintf("%.1f,%.1f", at.X, at.Y)
	}

	fmt.Fprintf(&c.body, `<%s points="%s"%s stroke-linejoin="round"/>`+"\n", element, strings.Join(coordinates, " "), s)
}

func (c *canvas) line(a, b point, s style) {
	c.cover(a)
	c.cover(b)
	fmt.Fprintf(&c.body, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f"%s stroke-linecap="round"/>`+"\n", a.X, a.Y, b.X, b.Y, s)
}

func (c *canvas) circle(at point, radius float64, s style) {
	c.cover(at.AddXY(-radius, -radius))
	c.cover(at.AddXY(radius, radius))
	fmt.Fprintf(&c.body, `<circle cx="%.1f" cy="%.1f" r="%g"%s/>`+"\n", at.X, at.Y, radius, s)
}

func (c *canvas) arrow(from, to point, s style) {
	along := to.Subtract(from).Normalize()
	base := to.Add(along.Multiply(-4.5 * s.width))
	side := along.Normal().Multiply(2.25 * s.width)

	c.line(from, base, s)
	c.polygon([]point{to, base.Add(side), base.Add(side.Negate())}, style{fill: s.stroke})
}

func (c *canvas) text(at point, f font, spans ...span) {
	family, anchor := cmp.Or(f.family, mono), cmp.Or(f.anchor, "middle")
	var whole strings.Builder
	for _, s := range spans {
		whole.WriteString(s.text)
	}
	extent := widthOf(whole.String(), f.size)
	switch anchor {
	case "start":
		c.cover(at.AddXY(extent, 0))
	case "end":
		c.cover(at.AddXY(-extent, 0))
	case "middle":
		c.cover(at.AddXY(-extent/2, 0))
		c.cover(at.AddXY(extent/2, 0))
	default:
		panic(fmt.Sprintf("diagrams: unknown text anchor %q", anchor))
	}
	c.cover(at.AddXY(0, -f.size/2))
	c.cover(at.AddXY(0, f.size/2))

	halo := ""
	if f.halo {
		halo = fmt.Sprintf(` stroke="%s" stroke-width="4" stroke-linejoin="round" paint-order="stroke"`, c.theme.cell)
	}
	fmt.Fprintf(&c.body, `<text x="%.1f" y="%.1f" font-family="%s" font-size="%g" text-anchor="%s"%s`, at.X, at.Y+f.size*0.35, family, f.size, anchor, halo)
	if len(spans) == 1 {
		fmt.Fprintf(&c.body, ` fill="%s">%s</text>`+"\n", spans[0].color, spans[0].text)
		return
	}

	c.body.WriteString(">")
	for _, s := range spans {
		fmt.Fprintf(&c.body, `<tspan fill="%s">%s</tspan>`, s.color, s.text)
	}
	c.body.WriteString("</text>\n")
}

func (c *canvas) content(id int) string {
	if c.view.IsZero() {
		return c.body.String()
	}

	return fmt.Sprintf(`<clipPath id="view-%d"><rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="6"/></clipPath>`+"\n"+`<g clip-path="url(#view-%d)">`+"\n%s</g>\n",
		id, c.view.Min.X, c.view.Min.Y, c.view.Width(), c.view.Height(), id, c.body.String())
}

func (s style) String() string {
	fill := cmp.Or(s.fill, "none")
	if s.stroke == "" {
		return fmt.Sprintf(` fill="%s" stroke="none"`, fill)
	}

	dash := ""
	if s.dashed {
		dash = ` stroke-dasharray="5 4"`
	}

	return fmt.Sprintf(` fill="%s" stroke="%s" stroke-width="%g"%s`, fill, s.stroke, s.width, dash)
}

func render(t theme, panels []*panel) (string, error) {
	slots := slotsOf(panels)
	row, height := gap*float64(len(panels)-1), 0.0
	captioned, noted := false, false
	for i, p := range panels {
		row += slots[i]
		height = max(height, p.height())
		captioned = captioned || p.caption != ""
		noted = noted || p.note != ""
	}
	if row > pageWidth-margin*2 {
		return "", fmt.Errorf("a row of %.0f px is wider than %.0f px", row, pageWidth-margin*2)
	}

	total := height + margin*2
	if captioned {
		total += captionAt
	}
	if noted {
		total += noteAt
	}

	var svg strings.Builder
	fmt.Fprintf(&svg, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f" role="img">`+"\n", pageWidth, total, pageWidth, total)

	x := (pageWidth - row) / 2
	for i, p := range panels {
		left := x + (slots[i]-p.bounds.Width())/2 - p.bounds.Min.X
		top := margin + (height-p.height())/2 - p.bounds.Min.Y
		fmt.Fprintf(&svg, `<g transform="translate(%.1f %.1f)">`+"\n%s</g>\n", left, top, p.content(i))
		if p.caption != "" {
			fmt.Fprintf(&svg, `<text x="%.1f" y="%.1f" font-family="%s" font-size="%g" text-anchor="middle" fill="%s">%s</text>`+"\n",
				x+slots[i]/2, margin+height+captionAt-4, mono, captionPt, t.muted, p.caption)
		}
		if p.note != "" {
			fmt.Fprintf(&svg, `<text x="%.1f" y="%.1f" font-family="%s" font-size="%g" text-anchor="middle" fill="%s">%s</text>`+"\n",
				x+slots[i]/2, margin+height+captionAt+noteAt-4, sans, notePt, t.muted, p.note)
		}
		x += slots[i] + gap
	}
	svg.WriteString("</svg>\n")

	return svg.String(), nil
}

func slotsOf(panels []*panel) []float64 {
	widest, even := 0.0, true
	for _, p := range panels {
		widest = max(widest, p.width())
		even = even && p.bounds.Width() == panels[0].bounds.Width()
	}

	slots := make([]float64, len(panels))
	for i, p := range panels {
		slots[i] = p.width()
		if even && widest <= (pageWidth-margin*2-gap*float64(len(panels)-1))/float64(len(panels)) {
			slots[i] = widest
		}
	}

	return slots
}

func widthOf(text string, size float64) float64 {
	return float64(len([]rune(text))) * size * glyph
}
