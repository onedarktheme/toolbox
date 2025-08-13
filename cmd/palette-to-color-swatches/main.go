package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"
	"fmt"
	"flag"
	"io"
	"io/ioutil"
	"log"
)

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
	helpText := `Usage: palette-to-color-swatches [options]

Options:
  -h, --help         Show this help message and exit
  -file <filename>   JSON file containing palettes (optional).
                     If omitted, the program reads JSON input from stdin by default.

Description:
  Reads a JSON palette definition and outputs 23x23 circle color swatches
  for each color in the palettes and saves them as PNGs in ./assets/

Example:
  cat palettes.json | palette-to-color-swatches
  palette-to-color-swatches -file palettes.json
`
	fmt.Fprint(os.Stderr, helpText)
}


func hexToRGBA(hex string) (color.RGBA, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color: %s", hex)
	}
	r64, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil { return color.RGBA{}, err }
	g64, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil { return color.RGBA{}, err }
	b64, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil { return color.RGBA{}, err }
	return color.RGBA{R: uint8(r64), G: uint8(g64), B: uint8(b64), A: 255}, nil
}

func drawCircle(img *image.RGBA, centerX, centerY, radius int, col color.Color) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x + y*y <= radius*radius {
				img.Set(centerX+x, centerY+y, col)
			}
		}
	}
}

func generateColorSwatches(paletteName string, colors map[string]string) {
	const size = 23
	const radius = 11

	// Ensure output directory exists
	outputDir := "assets/palette/circles"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	for _, role := range roleOrder {
		key, ok := roleMap[role]
		if !ok {
			continue
		}

		hex, ok := colors[key]
		if !ok {
			continue
		}

		img := image.NewRGBA(image.Rect(0, 0, size, size))
		draw.Draw(img, img.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
		swatchColor, _ := hexToRGBA(hex)

		drawCircle(img, size/2, size/2, radius, swatchColor)

		filePath := fmt.Sprintf("%s/%s-%s.png", outputDir, paletteName, strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(role, " ", "-"), "_", "-")))
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file %s: %v", filePath, err)
			continue
		}

		if err := png.Encode(file, img); err != nil {
			log.Printf("Failed to encode PNG for %s: %v", filePath, err)
		}

		file.Close()
	}
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
		generateColorSwatches(paletteName, colors)
	}
}