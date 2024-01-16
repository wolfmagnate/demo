package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
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
type Vertex struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func readVertex(jsonPath string) ([][]Vertex, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("error opening JSON file: %w", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %w", err)
	}

	var vertices [][]Vertex
	if err := json.Unmarshal(byteValue, &vertices); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return vertices, nil
}

var rightLines, _ = readVertex("rightLine.json")

var leftLines, _ = readVertex("leftLine.json")

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
			s[i][x] = termenv.TrueColor.Color("#FFFFFF")
		}
	}
	b := make([][]termenv.Color, height)
	for i := 0; i < height; i++ {
		b[i] = make([]termenv.Color, width)
		for x := 0; x < width; x++ {
			b[i][x] = termenv.TrueColor.Color("#696969")
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
		m.Ratio += 0.03
		m.Ratio += m.Ratio / 12
		if m.Ratio >= 1 {
			m.Ratio = 0
			m.AnimationType = Point
		}
	case Point:
		m.Ratio += 0.07
		if m.Ratio >= 1 {
			m.Ratio = 0
			m.AnimationType = Open
		}
	case Open:
		m.Ratio += 0.01
		m.Ratio += m.Ratio / 20
		if m.Ratio >= 10 {
			m.Ratio = 0
			m.AnimationType = Dark
		}
	}
	return m
}

func renderPoints(v1, v2 Vertex) []Vertex {
	diffX := int(math.Abs(float64(v2.X - v1.X)))
	diffY := int(math.Abs(float64(v2.Y - v1.Y)))
	dirX := 0
	dirY := 0
	if diffX > 0 {
		dirX = (v2.X - v1.X) / diffX
	}
	if diffY > 0 {
		dirY = (v2.Y - v1.Y) / diffY
	}
	tmp := Vertex{
		X: v1.X,
		Y: v1.Y,
	}
	result := make([]Vertex, 0)
	for tmp.X != v2.X || tmp.Y != v2.Y {
		result = append(result, Vertex{
			tmp.X,
			tmp.Y,
		})
		tmp.X += dirX
		tmp.Y += dirY
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

func (m *SlideModel) setLeftBackground(border int, color string) {
	for y := 0; y < height; y++ {
		for x := 0; x <= border; x++ {
			m.background[y][x] = termenv.TrueColor.Color(color)
		}
	}
}
func (m *SlideModel) setRightBackground(border int, color string) {
	for y := 0; y < height; y++ {
		for x := border; x < width; x++ {
			m.background[y][x] = termenv.TrueColor.Color(color)
		}
	}
}

func (m *SlideModel) renderLines(ratio float64) {
	offset := int(math.Round(width / 2 * ratio))
	m.clearLeft(width/2 - offset)
	m.setLeftBackground(width/2-offset-1, "#252525")
	for _, line := range leftLines {
		for i := 0; i < len(line)-1; i++ {
			ps := renderPoints(line[i], line[i+1])
			for _, p := range ps {
				if p.X <= offset {
					continue
				}
				m.chars[p.Y][p.X-offset] = "█"
			}
		}
	}
	m.clearRight(width/2 + 1 + offset)
	m.setRightBackground(width/2-1+offset, "#252525")
	for _, line := range rightLines {
		for i := 0; i < len(line)-1; i++ {
			ps := renderPoints(line[i], line[i+1])
			for _, p := range ps {
				if p.X+offset >= width {
					continue
				}
				m.chars[p.Y][p.X+offset] = "█"
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
			m.foreground[y][x] = termenv.TrueColor.Color("#FFFFFF")
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
			m.foreground[y][x] = termenv.TrueColor.Color("#FFFFFF")
		}
	}
	for _, line := range leftLines {
		m.changeStyleAtPoint(line, ratio)
	}
	for _, line := range rightLines {
		m.changeStyleAtPoint(line, ratio)
	}
}

func (m *SlideModel) changeStyleAtPoint(line []Vertex, ratio float64) {
	linePoints := make([]Vertex, 0)
	for i := 0; i < len(line)-1; i++ {
		ps := renderPoints(line[i], line[i+1])
		linePoints = append(linePoints, ps...)
	}
	index := int(float64(len(linePoints)) * ratio)
	if index >= len(linePoints) {
		return
	}
	target := linePoints[index]
	m.foreground[target.Y][target.X] = termenv.TrueColor.Color("#00ff7f")
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func lerpColor(startColor, endColor [3]int, t float64) string {
	red := int(lerp(float64(startColor[0]), float64(endColor[0]), t))
	green := int(lerp(float64(startColor[1]), float64(endColor[1]), t))
	blue := int(lerp(float64(startColor[2]), float64(endColor[2]), t))

	return fmt.Sprintf("#%02x%02x%02x", red, green, blue)
}
func (m *SlideModel) changeStyleLine(line []Vertex, ratio float64, offset int) {
	linePoints := make([]Vertex, 0)
	for i := 0; i < len(line)-1; i++ {
		ps := renderPoints(line[i], line[i+1])
		linePoints = append(linePoints, ps...)
	}
	index := int(float64(len(linePoints)) * ratio)
	if index >= len(linePoints) {
		index = len(linePoints) - 1
	}
	endColor := [3]int{168, 168, 255} // #a8a8ff
	startColor := [3]int{0, 255, 127} // #00ff7f

	for i := 0; i <= index; i++ {
		target := linePoints[i]
		if target.X+offset < 0 || target.X+offset >= width {
			continue
		}
		t := float64(i) / float64(len(linePoints)-1)
		color := lerpColor(startColor, endColor, t)
		m.foreground[target.Y][target.X+offset] = termenv.TrueColor.Color(color)
	}
}

func (m *SlideModel) renderCore(core string, offset, startColumn, startRow int) {
	lines := strings.Split(core, "\n")

	for i, line := range lines {
		for j, c := range strings.Split(line, "") {
			column := j + startColumn + offset
			row := i + startRow
			if column < 0 || column >= width {
				continue
			}
			if row < 0 || row >= width {
				continue
			}
			m.chars[row][column] = string(c)
		}
	}
}

func (m *SlideModel) renderCoreColor(core, color string, offset, startColumn, startRow int) {
	lines := strings.Split(core, "\n")

	for i, line := range lines {
		for j, _ := range strings.Split(line, "") {
			column := j + startColumn + offset
			row := i + startRow
			if column < 0 || column >= width {
				continue
			}
			if row < 0 || row >= width {
				continue
			}
			m.foreground[row][column] = termenv.TrueColor.Color(color)
		}
	}
}

func (m *SlideModel) renderCenter(ratio float64) {
	offset := int(math.Round(width / 2 * ratio))

	m.renderCore(leftCore, -offset, 82, height/2-1)

	m.renderCore(rightCore, offset, 84, height/2-1)
}

func (m *SlideModel) renderCenterColor(ratio float64) {
	offset := int(math.Round(width / 2 * ratio))

	m.renderCoreColor(leftCore, "#7FFF7F", -offset, 82, height/2-1)

	m.renderCoreColor(rightCore, "#7FFF7F", offset, 84, height/2-1)
}

func (m *SlideModel) renderLogo() {
	logo := getLogo()
	offset := 10
	m.renderCore(logo, -offset, 50, height/2-5)
}

func (m *SlideModel) renderLogoColor() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m.background[y][x] = termenv.TrueColor.Color("#FF99CC")
		}
	}
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
		m.renderCenterColor(0)
	case Light:
		m.renderLines(0)
		m.renderCenter(0)
		m.renderLineColor(m.Ratio)
	case Open:
		m.clearAll()
		m.renderLogo()
		m.renderLogoColor()
		m.renderLines(m.Ratio)
		m.renderCenter(m.Ratio)
		m.renderLineColorWithOffset(5*m.Ratio, m.Ratio)
		m.renderCenterColor(m.Ratio)
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
