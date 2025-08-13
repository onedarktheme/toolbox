package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"strconv"
)

var emojiMap = map[string]string{
	"dark": "🌑",
	"frost": "🧊",
	"ember": "🔥",
	"wraith": "👻",
}

var roleMap = map[string]string{
	// Background stack
	"Base":        "bg0",
	"Mantle":      "bg1",
	"Crust":       "bg_d",

	// Foreground stack
	"Text":       "fg",
	"Subtext1":   "light_grey",
	"Subtext0":   "grey",
	"Overlay2":   "bg2",
	"Overlay1":   "bg3",
	"Overlay0":   "bg_blue",
	"Surface2":   "diff_text",
	"Surface1":   "diff_change",
	"Surface0":   "diff_add",

	// Accents
	"Red":        "red",
	"Orange":     "orange",
	"Yellow":     "yellow",
	"Green":      "green",
	"Cyan":       "cyan",
	"Blue":       "blue",
	"Purple":     "purple",
	"Dark Red":   "dark_red",
	"Dark Yellow":"dark_yellow",
	"Dark Purple":"dark_purple",
}

var roleOrder = []string {
	"Red",
	"Orange",
	"Yellow",
	"Green",
	"Cyan",
	"Blue",
	"Purple",
	"Dark Red",
	"Dark Yellow",
	"Dark Purple",
	"Text",
	"Subtext1",
	"Subtext0",
	"Overlay2",
	"Overlay1",
	"Overlay0",
	"Surface2",
	"Surface1",
	"Surface0",
	"Base",
	"Mantle",
	"Crust",
}

func printHelp() {
	helpText := `Usage: palette-to-html-table [options]

Options:
  -h, --help         Show this help message and exit
  -file <filename>   JSON file containing palettes (optional).
                     If omitted, the program reads JSON input from stdin by default.

Description:
  Reads a JSON palette definition and outputs an HTML table showing
  the colors with their hex, RGB, and HSL values.

Example:
  cat palettes.json | palette-to-html-table
  palette-to-html-table -file palettes.json
`
	fmt.Fprint(os.Stderr, helpText)
}

func hexToRGB(hex string) (r, g, b int, err error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		err = fmt.Errorf("invalid hex color: %s", hex)
		return
	}
	var ri, gi, bi uint64
	ri, err = strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return
	}
	gi, err = strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return
	}
	bi, err = strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return
	}
	r, g, b = int(ri), int(gi), int(bi)
	return
}

func rgbToHSL(r, g, b int) (h, s, l int) {
	fr := float64(r) / 255
	fg := float64(g) / 255
	fb := float64(b) / 255

	max := math.Max(fr, math.Max(fg, fb))
	min := math.Min(fr, math.Min(fg, fb))

	lum := (max + min) / 2

	var hue, sat float64
	if max == min {
		hue, sat = 0, 0
	} else {
		d := max - min
		if lum > 0.5 {
			sat = d / (2 - max - min)
		} else {
			sat = d / (max + min)
		}

		switch max {
		case fr:
			hue = (fg - fb) / d
			if fg < fb {
				hue += 6
			}
		case fg:
			hue = (fb - fr)/d + 2
		case fb:
			hue = (fr - fg)/d + 4
		}
		hue /= 6
	}

	h = int(math.Round(hue * 360))
	s = int(math.Round(sat * 100))
	l = int(math.Round(lum * 100))
	return
}

func generateHTML(paletteName string, colors map[string]string) (string, error) {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("<details><summary>%s %s</summary>\n", emojiMap[paletteName], strings.ToUpper(paletteName[:1]) + paletteName[1:]))
	builder.WriteString("<table>\n")
	builder.WriteString("\t<tr>\n")
	builder.WriteString("\t\t<th></th>\n")
	builder.WriteString("\t\t<th>Role</th>\n")
	builder.WriteString("\t\t<th>Hex</th>\n")
	builder.WriteString("\t\t<th>RGB</th>\n")
	builder.WriteString("\t\t<th>HSL</th>\n")
	builder.WriteString("\t</tr>\n")

	for _, role := range roleOrder {
		key, ok := roleMap[role]
		if !ok {
			continue
		}

		hex, ok := colors[key]
		if !ok {
			continue
		}

		r, g, b, err := hexToRGB(hex)
		if err != nil {
			return "", err
		}
		h, s, l := rgbToHSL(r, g, b)

		swatchPath := fmt.Sprintf("https://github.com/onedarktheme/onedark/blob/master/assets/palette/circles/%s-%s.png", paletteName, strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(role, " ", "-"), "_", "-")))

		builder.WriteString("\t<tr>\n")
		builder.WriteString(fmt.Sprintf("\t\t<td><img src=\"%s\" width=\"23\"/></td>\n", swatchPath))
		builder.WriteString(fmt.Sprintf("\t\t<td>%s</td>\n", role))
		builder.WriteString(fmt.Sprintf("\t\t<td><code>%s</code></td>\n", hex))
		builder.WriteString(fmt.Sprintf("\t\t<td><code>rgb(%d, %d, %d)</code></td>\n", r, g, b))
		builder.WriteString(fmt.Sprintf("\t\t<td><code>hsl(%d, %d%%, %d%%)</code></td>\n", h, s, l))
		builder.WriteString("\t</tr>\n")
	}

	builder.WriteString("</table>\n</details>\n")
	return builder.String(), nil
}

func main() {
	help := flag.Bool("help", false, "Show help")
	helpShort := flag.Bool("h", false, "Show help (shorthand)")
	inputFile := flag.String("file", "", "JSON file containing palettes (optional)")

	flag.Parse()

	if *help || *helpShort {
		printHelp()
		return
	}

	var input io.Reader
	if *inputFile != "" {
		f, err := os.Open(*inputFile)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer f.Close()
		input = f
	} else {
		input = os.Stdin
	}

	data, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	if len(data) == 0 {
		printHelp()
		return
	}

	palettes := map[string]map[string]string{}
	if err := json.Unmarshal(data, &palettes); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	for paletteName, colors := range palettes {
		html, err := generateHTML(paletteName, colors)
		if err != nil {
			log.Fatalf("Error generating HTML: %v", err)
		}
		fmt.Println(html)
	}
}