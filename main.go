package main

import (
	"bufio"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var args bingoArgs
var allEntries []string

type bingoArgs struct {
	entriesPath        *string
	outputPath         *string
	imgPath            *string
	title              *string
	customTemplate     *bool
	customTemplatePath *string
}

type bingoModel struct {
	Entries   []string
	BingoName string
	ImgPath string
}

func loadArgs() {
	args.entriesPath = flag.String("Entries", "entries.txt", "Path to the file with the Entries.")
	args.outputPath = flag.String("output", "bingo_output.html", "Path to the output file.")
	args.imgPath = flag.String("image", "someImage.png", "Path to the free field image")
	args.title = flag.String("title", "Title of the bingo!", "Title of the generated html doc.")
	args.customTemplate = flag.Bool("custom", false, "Flag for turning on the custom template")
	args.customTemplatePath = flag.String("template-path", "template.gohtml", "Path to the custom template.")
	flag.Parse()
}

func genBingo() {
	var model bingoModel

	// Load the name
	model.BingoName = *args.title

	// Load image path for the model
	model.ImgPath = *args.imgPath

	// Load random entries from all of them
	model.Entries = randomFromLoaded()

	// Load template string content
	var tStr string

	if *args.customTemplate {
		templateByte, err := ioutil.ReadFile(*args.customTemplatePath)
		if err != nil {
			log.Fatal(err)
		}
		tStr = string(templateByte)
	} else {
		tStr = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .BingoName }}</title>
    <style>
        table, th, td {
            border: 1px solid black;
            background-color: white;
        }
        td {
            /*width: 160px;*/
            /*max-width: 160px;*/
            height: 160px;
            max-height: 160px;
        }
        img {
            width: 100%;
            height: 100%;
            max-width: 160px;
            margin-left: auto;
            margin-right: auto;
        }
    </style>
</head>
<body>
<h1>{{ .BingoName }}</h1>
<table>
    {{ $img := .ImgPath }}
    {{ range $ind, $val := .Entries}}
        {{ if mod $ind 5 }}
            <tr>
        {{ end }}

        {{ if eq $ind 12 }}
            <td onclick="changeColor(this)">
            <img src="{{ $img }}"></td>
        {{ else }}
            <td onclick="changeColor(this)">{{ $val }}</td>
        {{ end }}

        {{ if modWithStr (add $ind 1) 5 }}
            </tr>
        {{ end }}
    {{ end }}
</table>
<script>
    function changeColor(elem) {
        if (elem.style.backgroundColor === "")
            elem.style.backgroundColor = "white";
        if (elem.style.backgroundColor === "white") {
            elem.style.backgroundColor = "green";
        } else if (elem.style.backgroundColor.localeCompare("green") === 0) {
            elem.style.backgroundColor = "red";
        } else if (elem.style.backgroundColor.localeCompare("red") === 0) {
            elem.style.backgroundColor = "white";
        }
    }
</script>
</body>
</html>`
}

	// Create new template, define necessary functions and compile it
	t := template.New("bingo")
	t.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	t.Funcs(template.FuncMap{"modWithStr": func(i string, j int) bool {
		ii, err := strconv.Atoi(i)
		if err != nil {
			log.Fatal(err)
		}
		return ii%j == 0
	}})
	t.Funcs(template.FuncMap{"add": func(i, j int) string {return strconv.Itoa(i+j)}})
	_, err := t.Parse(tStr)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(*args.outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Execute the template and save the output to the specified file
	err = t.Execute(f, model)
	if err != nil {
		log.Fatal(err)
	}
}

func loadEntries() {
	file, err := os.Open(*args.entriesPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allEntries = append(allEntries, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func randomFromLoaded() (entries []string) {
	if len(allEntries) < 24 {
		log.Fatal("Not enough entries for bingo")
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	entered := make(map[int]bool)
	for i := 0; i < 25; i++ {
		if i == 12 {
			entries = append(entries, "")
			continue
		}
		var index int
		for true {
			index = r.Intn(len(allEntries))
			if _, ok := entered[index]; !ok {
				entered[index] = true
				break
			}
		}
		entries = append(entries, allEntries[index])
	}
	return
}

func main() {
	loadArgs()
	loadEntries()
	genBingo()
}
