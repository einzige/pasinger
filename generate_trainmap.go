package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Cell struct {
	X           int      `json:"x"`
	Y           int      `json:"y"`
	Type        string   `json:"type"`
	Text        string   `json:"text"`
	Highlighted bool     `json:"highlighted"`
	Classes     []string `json:"classes"`
}

func main() {
	// Load JSON file
	filePath := "trainmap_cells_corrected.json" // Your JSON file path
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	// Read and parse the JSON data
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var cells []Cell
	json.Unmarshal(byteValue, &cells)

	// Create an HTML file to write output
	outputFile, err := os.Create("trainmap_table.html")
	if err != nil {
		fmt.Println("Error creating HTML file:", err)
		return
	}
	defer outputFile.Close()

	// Define the current time as a single value
	currentTime := time.Now().Format("15:04:05")

	// Define HTML template with refactored CSS
	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Train Map</title>
    <style>
        body {
            font-family: Monospace;
            background-color: #ddd;
            overflow: hidden;
            filter: grayscale(100%);
        }

        #root {
            width: 1200px;
            height: 820px;
            position: relative;
        }

        #logo {
          position: absolute;
          top: 10px;
          left: 10px;
          z-index: 10000;
          height: auto;
          width: 100px;
        }

        table {
            border-collapse: collapse !important;
            width: 100%;
            height: 429px;
            position: relative;
            z-index: 9999;
        }

        td {
            height: 35px;
            text-align: left;
            vertical-align: bottom;
            border: none;
        }

        /* Cell background colors */
        .background-white {
            background-color: white;
            border-color: white;
        }

        .background-light {
            background-color: #eee;
            border-color: #eee;
        }

        .background-dark {
            background-color: #ddd;
            border-color: #ddd;
        }

        /* Borders */
        .border-left {
            border-left: 2px solid black;
        }

        .border-bottom {
            border-bottom: 2px solid black;
        }

        .note {
          font-size: .8rem;
          color: #aaa;
        }

        /* Destination styling */
        .destination {
            position: relative;
            display: flex;
            flex-direction: column-reverse;
            justify-content: center;
            padding-left: .25rem;
        }

        .destination::after {
            content: '';
            width: 4px;
            height: 4px;
            background-color: black;
            position: absolute;
            bottom: -4px;    /* Place at the top corner */
            left: -4px;  /* Place at the right corner */
        }

        .align-right .destination {
            align-items: flex-end;
            padding-right: .2rem;
        }

        .align-right .destination::after {
            left: inherit;
            right: -4px;
        }

        .destination.is-highlighted {
            padding: 0;
            margin: 0;
            display: block;
        }

        .destination-text {
            font-size: 1rem;
        }

        .destination.is-highlighted .destination-text {
            font-weight: bold;
            background: black;
            color: white;
            display: inline-block;
            padding: .2rem .5rem;
            font-size: 1.3rem;
            border-radius: 0 20px 0 0;
            padding-right: 2rem;
        }

        .destination-time {
            margin-top: 5px;
            display: flex;
            align-items: center;
        }

        .current-time {
            font-size: .8rem;
            color: #666;
        }

        .destination.is-highlighted .destination-time {
            display: none;
        }

        .destination-departure,
        .destination-arrival {
            font-size: .8rem;
            color: #555;
        }

        .destination-arrival {
            margin-left: .2rem;
            border: 1px dotted black;
            color: black;
            padding: .2rem;
            border-radius: 5px;
        }

        /* Images */
        .weather-img {
            height: 391px;
            width: 782px;
            z-index: 0;
            position: relative;
            top: -83px;
            overflow: hidden;
        }

        .loader-img {
            height: 311px;
            width: 418px;
            z-index: 0;
            position: absolute;
            right: 0;
            bottom: 0;
            overflow: hidden;
        }
    </style>
</head>
<body>
<div id="root">
    <svg xmlns:svg="http://www.w3.org/2000/svg" xmlns="http://www.w3.org/2000/svg" version="1.0" viewBox="0 0 363 254" width="100" id="logo">
      <defs id="defs3626"/>
      <g transform="translate(467,-268.21933)" id="layer1">
        <path d="M -140.27413,269.0187 L -431.01692,269.0187 C -450.99122,269.0187 -467.52964,284.81353 -467.52964,305.22373 L -467.52964,486.78733 C -467.52964,507.22317 -450.99122,523.40261 -431.01692,523.40261 L -140.27413,523.40261 C -120.6588,523.40261 -104.12038,507.22317 -104.12038,486.78733 L -104.12038,305.22373 C -104.12038,284.81353 -120.6588,269.0187 -140.27413,269.0187" id="path375" style="fill:#ed1c24;fill-opacity:1;fill-rule:nonzero;stroke:none"/>
        <path d="M -130.65877,486.78733 C -130.65877,492.58219 -134.50491,497.19756 -140.27413,497.19756 L -431.01692,497.19756 C -436.76049,497.19756 -440.99125,492.58219 -440.99125,486.78733 L -440.99125,305.22373 C -440.99125,299.45451 -436.76049,294.83914 -431.01692,294.83914 L -140.27413,294.83914 C -134.50491,294.83914 -130.65877,299.45451 -130.65877,305.22373 L -130.65877,486.78733" id="path379" style="fill:#ffffff;fill-opacity:1;fill-rule:nonzero;stroke:none"/>
        <path d="M -191.81245,430.50584 C -191.81245,413.55717 -202.58165,409.68539 -218.3252,409.68539 L -232.96618,409.68539 L -232.96618,452.48014 L -218.73545,452.48014 C -204.47908,452.48014 -191.81245,447.86477 -191.81245,430.50584 z M -232.96618,378.86496 L -219.09443,378.86496 C -206.01754,378.86496 -196.42782,373.09575 -196.42782,358.81373 C -196.42782,342.99327 -208.73548,338.76251 -221.81237,338.76251 L -232.96618,338.76251 L -232.96618,378.86496 z M -207.96625,479.45443 L -272.55582,479.45443 L -272.55582,312.94207 L -204.8637,312.94207 C -172.94071,312.94207 -155.6587,326.04459 -155.6587,357.27528 C -155.6587,373.48036 -169.1202,385.01879 -183.73555,392.3521 C -163.32535,398.12132 -149.50487,410.48026 -149.50487,431.27507 C -149.50487,465.19806 -177.96633,479.45443 -207.96625,479.45443" id="path383" style="fill:#ed1c24;fill-opacity:1;fill-rule:nonzero;stroke:none"/>
        <path d="M -329.83791,398.89055 C -329.83791,365.35218 -333.32508,339.17276 -369.45319,339.17276 L -377.55573,339.17276 L -377.55573,452.48014 L -363.325,452.48014 C -342.17121,452.48014 -329.83791,435.53147 -329.83791,398.89055 z M -356.40194,479.45443 L -417.52998,479.45443 L -417.52998,312.94207 L -356.40194,312.94207 C -313.32514,312.94207 -290.27392,339.91635 -290.27392,395.42902 C -290.27392,443.63401 -305.63285,479.06981 -356.40194,479.45443" id="path387" style="fill:#ed1c24;fill-opacity:1;fill-rule:nonzero;stroke:none"/>
      </g>
    </svg>
    <table>
        {{range $y := .Rows}}
        <tr>
            {{range $x := $.Cols}}
                {{with index $.CellMap $y $x}}
                    <td class="{{join .Classes " "}}">
                        {{if eq .Type "destination"}}
                        <div class="destination {{if .Highlighted}}is-highlighted{{end}}">
                            {{if .Highlighted}}
                            <div class="current-time">{{$.CurrentTime}}</div>
                            {{end}}
                            <div class="destination-text">{{.Text}}</div>
                            <div class="destination-time">
                                <div class="destination-departure">--:--</div>
                                <div class="destination-arrival">--:--</div>
                            </div>
                        </div>
                        {{else}}
                        {{.Text}}
                        {{end}}
                    </td>
                {{else}}
                    <td></td>
                {{end}}
            {{end}}
        </tr>
        {{end}}
    </table>
    <img class="weather-img" src="https://www.yr.no/en/content/2-2867714/meteogram.svg" />
    <img id="loader" class="loader-img" src="https://picsum.photos/418/391?grayscale" />
</div>
</body>
</html>
`

	// Create a map of cell data based on X, Y coordinates
	cellMap := make(map[int]map[int]Cell)
	maxX, maxY := 0, 0

	for _, cell := range cells {
		if cellMap[cell.Y] == nil {
			cellMap[cell.Y] = make(map[int]Cell)
		}
		cellMap[cell.Y][cell.X] = cell

		// Track maximum X and Y to create a table size
		if cell.X > maxX {
			maxX = cell.X
		}
		if cell.Y > maxY {
			maxY = cell.Y
		}
	}

	// Define template data
	type TemplateData struct {
		CellMap     map[int]map[int]Cell
		Rows        []int
		Cols        []int
		CurrentTime string // Single current time value passed to the template
	}

	// Prepare data for template rendering
	rows := make([]int, maxY+1)
	cols := make([]int, maxX+1)
	for i := 0; i <= maxY; i++ {
		rows[i] = i
	}
	for i := 0; i <= maxX; i++ {
		cols[i] = i
	}

	data := TemplateData{
		CellMap:     cellMap,
		Rows:        rows,
		Cols:        cols,
		CurrentTime: currentTime, // Pass current time
	}

	// Parse and execute the template
	tmpl, err := template.New("trainmap").Funcs(template.FuncMap{
		"join": func(a []string, sep string) string {
			return strings.Join(a, sep)
		},
	}).Parse(htmlTemplate)

	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}
