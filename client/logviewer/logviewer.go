package logviewer

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// LogViewTemplate is the web code for viewing logs
var LogViewTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Log Viewer</title>
    <style>
        
        tr, th {
            border-bottom: 2px solid black;
        }
        td {
            margin-right: 3px

        }
    </style>
</head>
<body>
    <main>
        <h2>cmdctrl Log Viewer</h2>
        <table>
            <tr>
                <th>Date/Time</th>
                <th>Log</th>
            </tr>
            {{ range . }}
            <tr>
                <td>{{ .Date }}</td>
                <td>{{ .Log }}</td>
            </tr>
            {{ end }}
        </table>
    </main>
</body>
</html>
`

// Options represents options for running the logviewer
type Options struct {
	LogFile string
}

type logEntries struct {
	Date string
	Log  string
}

var t *template.Template
var fopt Options

// RunLogViewer will run the logviewer
func RunLogViewer(opt Options) error {
	fopt = opt
	var err error
	t, err = template.New("logviewer").Parse(LogViewTemplate)
	if err != nil {
		return err
	}
	s := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handle),
	}
	fmt.Println("You can now access your logs at http://localhost:8080")
	log.Fatal(s.ListenAndServe())
	return nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile(fopt.LogFile)
	if err != nil {
		return
	}
	items := []logEntries{}
	for _, curr := range strings.Split(string(f), "\n") {
		items = append(items, logEntries{
			Date: strings.Split(string(curr), " ")[0],
			Log:  strings.Join(strings.Split(string(curr), " ")[1:], " "),
		})
	}
	t.Execute(w, items)
}
