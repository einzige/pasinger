package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
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

	// Define HTML template with embedded CSS
	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Train Map</title>
    <style>
        table {
            border-collapse: collapse;
            width: 100%;
            height: 100%;
	    font-family: Arial;
        }
        td {
            height: 35px;
            text-align: left;
            vertical-align: bottom;
        }
        .background-white { background-color: white; }
        .background-light { background-color: #eee; }
        .background-dark { background-color: #ddd; }
        .border-left { border-left: 2px solid black; }
        .border-bottom { border-bottom: 2px solid black; }
        .destination { display: flex; flex-direction: column; justify-content: center; padding-left: .25rem }
        .destination.is-highlighted { padding: 0; margin: 0; display: block }
        .destination-text { font-size: 1rem; }
        .destination.is-highlighted .destination-text { font-weight: bold; background: black; color: white; display: inline-block; padding: .2rem .5rem; font-size: 1.7rem }
        .destination-time { margin-top: 5px; display: flex; align-items: center }
        .destination.is-highlighted .destination-time { display: none }
        .destination-departure, .destination-arrival { font-size: .8rem; color: #555 }
	.destination-arrival { margin-left: .2rem; border: 1px dotted black; color: black; padding: .2rem; border-radius: 5px }
    </style>
</head>
<body>
    <table>
        {{range $y := .Rows}}
        <tr>
            {{range $x := $.Cols}}
                {{with index $.CellMap $y $x}}
                    <td class="{{join .Classes " "}}">
                        {{if eq .Type "destination"}}
                        <div class="destination {{if .Highlighted}}is-highlighted{{end}}">
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
		CellMap map[int]map[int]Cell
		Rows    []int
		Cols    []int
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
		CellMap: cellMap,
		Rows:    rows,
		Cols:    cols,
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
