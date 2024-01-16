package main

import (
	"math"
	"strings"

	"github.com/muesli/termenv"
)

const leftCore = ` ▄▄
██ 
 ▀▀`
const rightCore = `▄▄ 
 ██
▀▀ `

// 線を頂点で持つ
type vertex struct {
	x int
	y int
}

var rightLines = [][]vertex{
	{
		vertex{
			x: 86,
			y: 15,
		},
		vertex{
			x: 86,
			y: 8,
		},
		vertex{
			x: 90,
			y: 4,
		},
		vertex{
			x: 90,
			y: 0,
		},
	},
	{
		vertex{
			x: 86,
			y: 19,
		},
		vertex{
			x: 86,
			y: 26,
		},
		vertex{
			x: 90,
			y: 30,
		},
		vertex{
			x: 90,
			y: 34,
		},
	},
	{
		vertex{
			x: 87,
			y: 15,
		},
		vertex{
			x: 88,
			y: 15,
		},
		vertex{
			x: 89,
			y: 14,
		},
		vertex{
			x: 90,
			y: 14,
		},
		vertex{
			x: 91,
			y: 13,
		},
		vertex{
			x: 92,
			y: 13,
		},
		vertex{
			x: 93,
			y: 12,
		},
		vertex{
			x: 94,
			y: 12,
		},
		vertex{
			x: 95,
			y: 11,
		},
		vertex{
			x: 96,
			y: 11,
		},
		vertex{
			x: 97,
			y: 10,
		},
		vertex{
			x: 98,
			y: 10,
		},
		vertex{
			x: 108,
			y: 0,
		},
	},
	{
		vertex{
			x: 87,
			y: 19,
		},
		vertex{
			x: 88,
			y: 19,
		},
		vertex{
			x: 89,
			y: 20,
		},
		vertex{
			x: 90,
			y: 20,
		},
		vertex{
			x: 91,
			y: 21,
		},
		vertex{
			x: 92,
			y: 21,
		},
		vertex{
			x: 93,
			y: 22,
		},
		vertex{
			x: 94,
			y: 22,
		},
		vertex{
			x: 95,
			y: 23,
		},
		vertex{
			x: 96,
			y: 23,
		},
		vertex{
			x: 97,
			y: 24,
		},
		vertex{
			x: 98,
			y: 24,
		},
		vertex{
			x: 108,
			y: 34,
		},
	},
	{
		vertex{
			x: 87,
			y: 16,
		},
		vertex{
			x: 93,
			y: 16,
		},
		vertex{
			x: 95,
			y: 14,
		},
		vertex{
			x: 105,
			y: 14,
		},
		vertex{
			x: 112,
			y: 7,
		},
		vertex{
			x: 121,
			y: 7,
		},
		vertex{
			x: 128,
			y: 0,
		},
	},
	{
		vertex{
			x: 87,
			y: 18,
		},
		vertex{
			x: 93,
			y: 18,
		},
		vertex{
			x: 95,
			y: 20,
		},
		vertex{
			x: 105,
			y: 20,
		},
		vertex{
			x: 112,
			y: 27,
		},
		vertex{
			x: 121,
			y: 27,
		},
		vertex{
			x: 128,
			y: 34,
		},
	},
}

var leftLines = [][]vertex{
	{
		vertex{
			x: 82,
			y: 15,
		},
		vertex{
			x: 82,
			y: 8,
		},
		vertex{
			x: 78,
			y: 4,
		},
		vertex{
			x: 78,
			y: 0,
		},
	},
	{
		vertex{
			x: 82,
			y: 19,
		},
		vertex{
			x: 82,
			y: 26,
		},
		vertex{
			x: 78,
			y: 30,
		},
		vertex{
			x: 78,
			y: 34,
		},
	},
	{
		vertex{
			x: 81,
			y: 15,
		},
		vertex{
			x: 80,
			y: 15,
		},
		vertex{
			x: 79,
			y: 14,
		},
		vertex{
			x: 78,
			y: 14,
		},
		vertex{
			x: 77,
			y: 13,
		},
		vertex{
			x: 76,
			y: 13,
		},
		vertex{
			x: 75,
			y: 12,
		},
		vertex{
			x: 74,
			y: 12,
		},
		vertex{
			x: 73,
			y: 11,
		},
		vertex{
			x: 72,
			y: 11,
		},
		vertex{
			x: 71,
			y: 10,
		},
		vertex{
			x: 70,
			y: 10,
		},
		vertex{
			x: 60,
			y: 0,
		},
	},
	{
		vertex{
			x: 81,
			y: 19,
		},
		vertex{
			x: 80,
			y: 19,
		},
		vertex{
			x: 79,
			y: 20,
		},
		vertex{
			x: 78,
			y: 20,
		},
		vertex{
			x: 77,
			y: 21,
		},
		vertex{
			x: 76,
			y: 21,
		},
		vertex{
			x: 75,
			y: 22,
		},
		vertex{
			x: 74,
			y: 22,
		},
		vertex{
			x: 73,
			y: 23,
		},
		vertex{
			x: 72,
			y: 23,
		},
		vertex{
			x: 71,
			y: 24,
		},
		vertex{
			x: 70,
			y: 24,
		},
		vertex{
			x: 60,
			y: 34,
		},
	},
	{
		vertex{
			x: 81,
			y: 16,
		},
		vertex{
			x: 75,
			y: 16,
		},
		vertex{
			x: 73,
			y: 14,
		},
		vertex{
			x: 63,
			y: 14,
		},
		vertex{
			x: 56,
			y: 7,
		},
		vertex{
			x: 47,
			y: 7,
		},
		vertex{
			x: 40,
			y: 0,
		},
	},
	{
		vertex{
			x: 81,
			y: 18,
		},
		vertex{
			x: 75,
			y: 18,
		},
		vertex{
			x: 73,
			y: 20,
		},
		vertex{
			x: 63,
			y: 20,
		},
		vertex{
			x: 56,
			y: 27,
		},
		vertex{
			x: 47,
			y: 27,
		},
		vertex{
			x: 40,
			y: 34,
		},
	},
}

const width = 170
const height = 35

// AnimationType はアニメーションの種類を定義する型です。
type AnimationType int

// AnimationType の許容される値を定義します。
const (
	Dark AnimationType = iota
	Point
	Light
	Open
)

// slideModel はスライドのモデルを表す構造体です。
type SlideModel struct {
	AnimationType AnimationType // AnimationType は列挙型のように振る舞います。
	Ratio         float64       // Ratio は比率を表すfloat型です。
	chars         [][]string
	foreground    [][]termenv.Color
	background    [][]termenv.Color
}

func Init() *SlideModel {
	c := make([][]string, height)
	for i := 0; i < height; i++ {
		c[i] = make([]string, width)
		for x := 0; x < width; x++ {
			c[i][x] = " "
		}
	}
	s := make([][]termenv.Color, height)
	for i := 0; i < height; i++ {
		s[i] = make([]termenv.Color, width)
		for x := 0; x < width; x++ {
			s[i][x] = termenv.ANSI.Color("#FFFFFF")
		}
	}
	b := make([][]termenv.Color, height)
	for i := 0; i < height; i++ {
		b[i] = make([]termenv.Color, width)
		for x := 0; x < width; x++ {
			b[i][x] = termenv.ANSI.Color("#696969")
		}
	}
	return &SlideModel{
		AnimationType: Dark,
		Ratio:         0.0,
		chars:         c,
		foreground:    s,
		background:    b,
	}
}

func (m *SlideModel) Update() *SlideModel {
	switch m.AnimationType {
	case Dark:
		m.Ratio += 0.1
		m.Ratio += m.Ratio / 12
		if m.Ratio >= 1 {
			m.Ratio = 0
			m.AnimationType = Point
		}
	case Point:
		m.Ratio += 0.1
		m.Ratio += m.Ratio / 12
		if m.Ratio >= 1 {
			m.Ratio = 0
			m.AnimationType = Open
		}
	case Open:
		m.Ratio = 1
		if m.Ratio >= 10 {
			m.Ratio = 0
			m.AnimationType = Dark
		}
	}
	return m
}

func renderPoints(v1, v2 vertex) []vertex {
	diffX := int(math.Abs(float64(v2.x - v1.x)))
	diffY := int(math.Abs(float64(v2.y - v1.y)))
	dirX := 0
	dirY := 0
	if diffX > 0 {
		dirX = (v2.x - v1.x) / diffX
	}
	if diffY > 0 {
		dirY = (v2.y - v1.y) / diffY
	}
	tmp := vertex{
		x: v1.x,
		y: v1.y,
	}
	result := make([]vertex, 0)
	for tmp.x != v2.x || tmp.y != v2.y {
		result = append(result, vertex{
			tmp.x,
			tmp.y,
		})
		tmp.x += dirX
		tmp.y += dirY
	}
	return result
}

func (m *SlideModel) clearAll() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m.chars[y][x] = " "
		}
	}
}
func (m *SlideModel) clearLeft(border int) {
	for y := 0; y < height; y++ {
		for x := 0; x <= border; x++ {
			m.chars[y][x] = " "
		}
	}
}
func (m *SlideModel) clearRight(border int) {
	for y := 0; y < height; y++ {
		for x := border; x < width; x++ {
			m.chars[y][x] = " "
		}
	}
}

func (m *SlideModel) renderLines(ratio float64) {
	offset := int(math.Round(width / 2 * ratio))
	m.clearLeft(width/2 - offset)
	for _, line := range leftLines {
		for i := 0; i < len(line)-1; i++ {
			ps := renderPoints(line[i], line[i+1])
			for _, p := range ps {
				if p.x <= offset {
					continue
				}
				m.chars[p.y][p.x-offset] = "█"
			}
		}
	}
	m.clearRight(width/2 + 1 + offset)
	for _, line := range rightLines {
		for i := 0; i < len(line)-1; i++ {
			ps := renderPoints(line[i], line[i+1])
			for _, p := range ps {
				if p.x+offset >= width {
					continue
				}
				m.chars[p.y][p.x+offset] = "█"
			}
		}
	}
}

func (m *SlideModel) renderLineColor(ratio float64) {
	m.renderLineColorWithOffset(ratio, 0)
}
func (m *SlideModel) renderLineColorWithOffset(ratio, offsetRatio float64) {
	offset := int(math.Round(width / 2 * offsetRatio))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m.foreground[y][x] = termenv.ANSI.Color("#FFFFFF")
		}
	}
	for _, line := range leftLines {
		m.changeStyleLine(line, ratio, -offset)
	}
	for _, line := range rightLines {
		m.changeStyleLine(line, ratio, offset)
	}
}

func (m *SlideModel) renderPointColor(ratio float64) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m.foreground[y][x] = termenv.ANSI.Color("#FFFFFF")
		}
	}
	for _, line := range leftLines {
		m.changeStyleAtPoint(line, ratio)
	}
	for _, line := range rightLines {
		m.changeStyleAtPoint(line, ratio)
	}
}

func (m *SlideModel) changeStyleAtPoint(line []vertex, ratio float64) {
	linePoints := make([]vertex, 0)
	for i := 0; i < len(line)-1; i++ {
		ps := renderPoints(line[i], line[i+1])
		linePoints = append(linePoints, ps...)
	}
	index := int(float64(len(linePoints)) * ratio)
	if index >= len(linePoints) {
		return
	}
	target := linePoints[index]
	m.foreground[target.y][target.x] = termenv.ANSI.Color("#00ff7f")
}

func (m *SlideModel) changeStyleLine(line []vertex, ratio float64, offset int) {
	linePoints := make([]vertex, 0)
	for i := 0; i < len(line)-1; i++ {
		ps := renderPoints(line[i], line[i+1])
		linePoints = append(linePoints, ps...)
	}
	index := int(float64(len(linePoints)) * ratio)
	if index >= len(linePoints) {
		index = len(linePoints) - 1
	}
	for i := 0; i <= index; i++ {
		target := linePoints[i]
		if target.x+offset < 0 || target.x+offset >= width {
			continue
		}
		m.foreground[target.y][target.x+offset] = termenv.ANSI.Color("#00ff7f")
	}
}

func (m *SlideModel) renderCore(core string, offset, startColumn int) {
	lines := strings.Split(core, "\n")

	for i, line := range lines {
		for j, c := range strings.Split(line, "") {
			column := j + startColumn + offset
			if column < 0 || column >= width {
				continue
			}
			m.chars[i+height/2-1][column] = string(c)
		}
	}
}

func (m *SlideModel) renderCenter(ratio float64) {
	offset := int(math.Round(width / 2 * ratio))

	m.renderCore(leftCore, -offset, 82)

	m.renderCore(rightCore, offset, 84)
}

func (m *SlideModel) renderLogo() {
	logo := getLogo()
	offset := 10
	m.renderCore(logo, -offset, 60)

}

func (m *SlideModel) View() string {
	switch m.AnimationType {
	case Dark:
		m.renderLines(0)
		m.renderCenter(0)
		m.renderLineColor(-1)
	case Point:
		m.renderLines(0)
		m.renderCenter(0)
		m.renderPointColor(m.Ratio)
	case Light:
		m.renderLines(0)
		m.renderCenter(0)
		m.renderLineColor(m.Ratio)
	case Open:
		m.clearAll()
		m.renderLogo()
		m.renderLines(m.Ratio)
		m.renderCenter(m.Ratio)
		m.renderLineColorWithOffset(5*m.Ratio, m.Ratio)
	}

	b := strings.Builder{}
	for i, v1 := range m.chars {
		for j, v2 := range v1 {
			b.WriteString(termenv.String(v2).Foreground(m.foreground[i][j]).Background(m.background[i][j]).String())
		}
		b.WriteString("\n")
	}
	return b.String()
}
