package main

import "strings"

var C = `◢████████◤◤    
◢███◤     ◥███ 
███            
██             
███            
███            
██             
███           ◥
 ███◤      ◢██ 
   █████████◤  `

var H = `██          
██          
██          
███████████ 
██       ██◤
██       ███
██       ███
██       ███
██       ███
██       ███`

var A = `            
            
            
██████████  
        ███ 
◢█████████  
◢██◤      ██
██        ██
██       ███
◥█████████◤ `
var R = `        
        
██      
████████
██◤     
██      
██      
██      
██      
██      `

var M = `                     
                      
 ██  ______   ______ 
◥██████████◤█████████
◥██       ██       ██
◥██       ██       ██
◥██       ██       ██
◥██       ██       ██
◥██       ██       ██
 ██       ██       ██`

// splitIntoLines は与えられた文字列を行に分割する関数です。
func splitIntoLines(s string) []string {
	return strings.Split(s, "\n")
}

// combineLines は複数の文字列の配列を受け取り、それぞれの行を結合して新しい文字列を生成する関数です。
func combineLines(lines ...[]string) string {
	var combinedLines []string

	// 最も長い配列の長さを取得します。
	maxLength := 0
	for _, l := range lines {
		if len(l) > maxLength {
			maxLength = len(l)
		}
	}

	// 各行を順番に結合します。
	for i := 0; i < maxLength; i++ {
		var combinedLine string
		for _, l := range lines {
			if i < len(l) {
				combinedLine += "   " + l[i]
			}
		}
		combinedLines = append(combinedLines, combinedLine)
	}

	return strings.Join(combinedLines, "\n")
}

func getLogo() string {

	// 各文字列を行に分割します。
	cLines := splitIntoLines(C)
	hLines := splitIntoLines(H)
	aLines := splitIntoLines(A)
	rLines := splitIntoLines(R)
	mLines := splitIntoLines(M)

	// 分割した行を結合して新しい文字列を生成します。
	combined := combineLines(cLines, hLines, aLines, rLines, mLines)
	return combined
}
