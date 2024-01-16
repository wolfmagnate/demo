package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
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
	Progress
	Horizontal
	Loopback
)

// slideModel はスライドのモデルを表す構造体です。
type SlideModel struct {
	AnimationType AnimationType // AnimationType は列挙型のように振る舞います。
	Ratio         float64       // Ratio は比率を表すfloat型です。
	ratio         float64
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
		ratio:         0.0,
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
		m.ratio = m.Ratio
	case Point:
		m.Ratio += 0.07
		if m.Ratio >= 1 {
			m.Ratio = 0
			m.AnimationType = Open
		}
		m.ratio = m.Ratio
	case Open:
		m.Ratio += 0.05
		if m.Ratio >= 1 {
			m.Ratio = 0
			m.AnimationType = Progress
		}
		m.ratio = Ease3(m.Ratio)
	case Progress:
		m.Ratio += 0.04
		if m.Ratio >= 1 {
			m.Ratio = 0
			m.AnimationType = Horizontal
		}
		m.ratio = Ease1(m.Ratio)
	case Horizontal:
		m.Ratio += 0.02
		if m.Ratio >= 1 {
			m.Ratio = 0
			m.AnimationType = Loopback
		}
		m.ratio = Ease2(m.Ratio)
	case Loopback:
		m.Ratio += 0.04
		if m.Ratio >= 1.1 {
			m.Ratio = 0
			m.AnimationType = Dark
		}
		m.ratio = Ease2(m.Ratio)
		m.ratio = m.Ratio
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

// parseColor は文字列形式のカラーコードを解析し、RGB値を表すintのスライスを返します。
func parseColor(hexColor string) ([3]int, error) {
	if len(hexColor) != 7 || hexColor[0] != '#' {
		return [3]int{}, fmt.Errorf("invalid color format")
	}
	red, err := strconv.ParseInt(hexColor[1:3], 16, 64)
	if err != nil {
		return [3]int{}, fmt.Errorf("error parsing red component: %v", err)
	}
	green, err := strconv.ParseInt(hexColor[3:5], 16, 64)
	if err != nil {
		return [3]int{}, fmt.Errorf("error parsing green component: %v", err)
	}
	blue, err := strconv.ParseInt(hexColor[5:7], 16, 64)
	if err != nil {
		return [3]int{}, fmt.Errorf("error parsing blue component: %v", err)
	}
	return [3]int{int(red), int(green), int(blue)}, nil
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

func (m *SlideModel) renderLogoBackgroundColor() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m.background[y][x] = termenv.TrueColor.Color("#FF99CC")
		}
	}
}

func (m *SlideModel) renderLogoColor(ratio float64, color1, color2 string, colorOffset float64) {
	logo := getLogo()
	lines := strings.Split(logo, "\n")
	offset := -10
	startColumn := 50
	startRow := height/2 - 5

	rgb1, _ := parseColor(color1)
	rgb2, _ := parseColor(color2)

	for i, line := range lines {
		chars := strings.Split(line, "")
		for j, _ := range chars {
			column := j + startColumn + offset
			row := i + startRow
			if column < 0 || column >= width {
				continue
			}
			if row < 0 || row >= width {
				continue
			}
			colorRatio := float64(j) / float64((len(chars) - 1))
			if colorRatio > ratio {
				m.foreground[row][column] = termenv.TrueColor.Color("#ffffff")
				continue
			}
			color := lerpColor(rgb1, rgb2, colorRatio)
			m.foreground[row][column] = termenv.TrueColor.Color(color)
		}
	}
}

func clamp(v float64) float64 {
	if v > 1 {
		return 1
	}
	if v < 0 {
		return 0
	}
	return v
}

func (m *SlideModel) renderHorizontalHeader(ratio float64) {
	width1, width2, width3 := getHorizontalWidths(ratio)

	leftBars := strings.Split(barLeft, "\n")
	rightBars := strings.Split(barRight, "\n")
	barWidth := 10
	for y := 0; y < 6; y++ {
		cnt := 0
		l := strings.Split(rightBars[y%6], "")

		for x := width1; x < width1+barWidth; x++ {
			if x >= width || x < 0 {
				continue
			}
			m.chars[y][x] = string(l[cnt])
			cnt++
		}
	}
	for y := 30; y < 35; y++ {
		cnt := 0
		l := strings.Split(rightBars[y%6], "")
		for x := width1; x < width1+barWidth; x++ {
			if x >= width || x < 0 {
				continue
			}
			m.chars[y][x] = string(l[cnt])
			cnt++
		}
	}
	for y := 6; y < 12; y++ {
		cnt := 0
		l := strings.Split(leftBars[y%6], "")
		for x := width - width2 - barWidth; x < width1-width2; x++ {
			if x >= width || x < 0 {
				continue
			}
			m.chars[y][x] = string(l[cnt])
			cnt++
		}
	}
	for y := 24; y < 30; y++ {
		cnt := 0
		l := strings.Split(leftBars[y%6], "")
		for x := width - width2 - barWidth; x < width1-width2; x++ {
			if x >= width || x < 0 {
				continue
			}
			m.chars[y][x] = string(l[cnt])
			cnt++
		}
	}
	if width3 == 0 {
		return
	}
	for y := 12; y < 18; y++ {
		cnt := 0
		l := strings.Split(rightBars[y%6], "")
		for x := width3; x < width3+barWidth; x++ {
			if x >= width || x < 0 {
				continue
			}
			m.chars[y][x] = string(l[cnt])
			cnt++
		}
	}
	for y := 18; y < 24; y++ {
		cnt := 0
		l := strings.Split(rightBars[y%6], "")
		for x := width3; x < width3+barWidth; x++ {
			if x >= width || x < 0 {
				continue
			}
			m.chars[y][x] = string(l[cnt])
			cnt++
		}
	}
}

func (m *SlideModel) renderHoritontalLine(ratio float64) {
	// まずはそれぞれの線の比率を計算する
	width1, width2, width3 := getHorizontalWidths(ratio)

	// 左から右へ
	// y=0,1,...5とy=30,31,...35
	// x=0,1,...width1
	for y := 0; y < 6; y++ {
		for x := 0; x < width1; x++ {
			m.chars[y][x] = "█"
		}
	}
	for y := 30; y < 35; y++ {
		for x := 0; x < width1; x++ {
			m.chars[y][x] = "█"
		}
	}

	for y := 6; y < 12; y++ {
		for x := width - 1; x >= width-width2; x-- {
			m.chars[y][x] = "█"
		}
	}
	for y := 24; y < 30; y++ {
		for x := width - 1; x >= width-width2; x-- {
			m.chars[y][x] = "█"
		}
	}

	for y := 12; y < 18; y++ {
		for x := 0; x < width3; x++ {
			m.chars[y][x] = "█"
		}
	}
	for y := 18; y < 24; y++ {
		for x := 0; x < width3; x++ {
			m.chars[y][x] = "█"
		}
	}
}

func getHorizontalWidths(ratio float64) (int, int, int) {
	var ratio1, ratio2, ratio3 = 0.0, 0.0, 0.0
	if ratio > 0.33 {
		ratio1 = 1.0
	} else {
		ratio1 = ratio * 3
	}
	if ratio <= 0.33 {
		ratio2 = 0.0
	} else if ratio > 0.66 {
		ratio2 = 1.0
	} else {
		ratio2 = (ratio - 0.33) * 3
	}
	if ratio <= 0.66 {
		ratio3 = 0.0
	} else {
		ratio3 = (ratio - 0.66) * 3
	}
	ratio1 = clamp(ratio1)
	ratio2 = clamp(ratio2)
	ratio3 = clamp(ratio3)
	width1 := int(width * ratio1)
	width2 := int(width * ratio2)
	width3 := int(width * ratio3)
	return width1, width2, width3
}

func getColorFromDistance(distance int, color1, color2, color3, color4 string) string {
	c1, _ := parseColor(color1)
	c2, _ := parseColor(color2)
	c3, _ := parseColor(color3)
	c4, _ := parseColor(color4)
	if distance < 170 {
		return lerpColor(c1, c2, float64(distance)/170)
	} else if distance < 340 {
		return lerpColor(c2, c3, float64(distance-170)/170)
	} else if distance < 510 {
		return lerpColor(c3, c4, float64(distance-340)/170)
	} else {
		return color4
	}
}

func calcDistance(b Vertex) int {
	b.Y = b.Y / 6
	bSum := 0
	if b.Y >= 3 {
		b.Y %= 3
		bSum += (2 - b.Y) * width
		if b.Y%2 == 0 {
			bSum += b.X
		} else {
			bSum += width - b.X - 1
		}
	} else {
		b.Y %= 3
		bSum += b.Y * width
		if b.Y%2 == 0 {
			bSum += b.X
		} else {
			bSum += width - b.X - 1
		}
	}
	return bSum
}

func calcHeadVertex(ratio float64) (Vertex, Vertex) {
	if ratio < 0.33 {
		w := int(width * ratio * 3)
		return Vertex{
				X: w,
				Y: 0,
			},
			Vertex{
				X: w,
				Y: 30,
			}
	} else if ratio < 0.66 {
		w := int(width * (ratio - 0.33) * 3)
		return Vertex{
				X: width - w,
				Y: 6,
			},
			Vertex{
				X: width - w,
				Y: 24,
			}
	} else {
		w := int(width * (ratio - 0.66) * 3)
		return Vertex{
				X: w,
				Y: 12,
			},
			Vertex{
				X: w,
				Y: 18,
			}
	}
}

func (m *SlideModel) renderHoritontalLineColor(ratio float64, color1, color2, color3, color4 string) {
	// まずはそれぞれの線の比率を計算する
	upper, lower := calcHeadVertex(ratio)
	upperDistance := calcDistance(upper)
	lowerDistance := calcDistance(lower)
	width1, width2, width3 := getHorizontalWidths(ratio)
	m.renderColorWithDistance(ratio, upperDistance, lowerDistance, width1, width2, width3, color1, color2, color3, color4)
}

func (m *SlideModel) renderLoopBackColor(ratio float64, color1, color2, color3, color4 string) {
	upper, lower := calcHeadVertex(ratio)
	upperDistance := calcDistance(upper)
	lowerDistance := calcDistance(lower)
	m.renderColorWithDistance(ratio, upperDistance+3*width, lowerDistance+3*width, width, width, width, color1, color2, color3, color4)
}

func (m *SlideModel) renderColorWithDistance(ratio float64, upperDistance, lowerDistance, width1, width2, width3 int, color1, color2, color3, color4 string) {
	for y := 0; y < 6; y++ {
		for x := 0; x < width1; x++ {
			d := calcDistance(Vertex{
				X: x, Y: y,
			})
			m.foreground[y][x] = termenv.TrueColor.Color(getColorFromDistance(upperDistance-d, color1, color2, color3, color4))
		}
	}
	for y := 30; y < 35; y++ {
		for x := 0; x < width1; x++ {
			d := calcDistance(Vertex{
				X: x, Y: y,
			})
			m.foreground[y][x] = termenv.TrueColor.Color(getColorFromDistance(lowerDistance-d, color1, color2, color3, color4))
		}
	}

	for y := 6; y < 12; y++ {
		for x := width - 1; x >= width-width2; x-- {
			d := calcDistance(Vertex{
				X: x, Y: y,
			})
			m.foreground[y][x] = termenv.TrueColor.Color(getColorFromDistance(upperDistance-d, color1, color2, color3, color4))
		}
	}
	for y := 24; y < 30; y++ {
		for x := width - 1; x >= width-width2; x-- {
			d := calcDistance(Vertex{
				X: x, Y: y,
			})
			m.foreground[y][x] = termenv.TrueColor.Color(getColorFromDistance(lowerDistance-d, color1, color2, color3, color4))
		}
	}

	for y := 12; y < 18; y++ {
		for x := 0; x < width3; x++ {
			d := calcDistance(Vertex{
				X: x, Y: y,
			})
			m.foreground[y][x] = termenv.TrueColor.Color(getColorFromDistance(upperDistance-d, color1, color2, color3, color4))
		}
	}
	for y := 18; y < 24; y++ {
		for x := 0; x < width3; x++ {
			d := calcDistance(Vertex{
				X: x, Y: y,
			})
			m.foreground[y][x] = termenv.TrueColor.Color(getColorFromDistance(lowerDistance-d, color1, color2, color3, color4))
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
		m.renderLogoBackgroundColor()
		m.renderLines(m.ratio)
		m.renderCenter(m.ratio)
		m.renderLineColorWithOffset(5*m.ratio, m.ratio)
		m.renderCenterColor(m.ratio)
	case Progress:
		m.clearAll()
		m.renderLogo()
		m.renderLogoBackgroundColor()
		m.renderLogoColor(m.ratio, "#8eff8e", "#7fffff", 2*m.ratio)
	case Horizontal:
		m.clearAll()
		m.renderLogo()
		m.renderLogoBackgroundColor()
		m.renderLogoColor(1, "#8eff8e", "#7fffff", 2*m.ratio)
		// 幅6の線をしたから引いていく
		m.renderHorizontalHeader(m.ratio)
		m.renderHoritontalLine(m.ratio)
		m.renderHoritontalLineColor(m.ratio, "#ffff74", "#7fff7f", "#7fbfff", "#252525")
	case Loopback:
		m.renderLoopBackColor(m.ratio, "#ffff74", "#7fff7f", "#7fbfff", "#252525")

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
