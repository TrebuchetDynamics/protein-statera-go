package render

import (
	"fmt"
	"html"
	"sort"
	"strings"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/evidence"
)

// HTMLReport renders a static HTML evidence report.
func HTMLReport(report evidence.Report) string {
	var b strings.Builder
	fmt.Fprintln(&b, "<!doctype html>")
	fmt.Fprintln(&b, "<html lang=\"en\">")
	fmt.Fprintln(&b, "<head>")
	fmt.Fprintln(&b, "<meta charset=\"utf-8\">")
	fmt.Fprintf(&b, "<title>Protein Statera Go - %s</title>\n", html.EscapeString(report.StructureID))
	fmt.Fprintln(&b, "<style>body{font-family:system-ui,sans-serif;line-height:1.5;max-width:900px;margin:2rem auto;padding:0 1rem}table{border-collapse:collapse}td,th{border:1px solid #999;padding:.35rem .5rem;text-align:left}</style>")
	fmt.Fprintln(&b, "</head>")
	fmt.Fprintln(&b, "<body>")
	fmt.Fprintf(&b, "<h1>Protein Structure Report: %s</h1>\n", html.EscapeString(report.StructureID))
	fmt.Fprintf(&b, "<p>evidence_class=%s</p>\n", html.EscapeString(string(report.Evidence)))
	fmt.Fprintln(&b, "<h2>Metrics</h2>")
	fmt.Fprintln(&b, "<table><thead><tr><th>metric</th><th>value</th></tr></thead><tbody>")
	for _, key := range sortedFloatKeys(report.Metrics) {
		fmt.Fprintf(&b, "<tr><td>%s</td><td>%.2f</td></tr>\n", html.EscapeString(key), report.Metrics[key])
	}
	fmt.Fprintln(&b, "</tbody></table>")
	fmt.Fprintln(&b, "<h2>Confidence</h2>")
	fmt.Fprintln(&b, "<table><thead><tr><th>band</th><th>count</th></tr></thead><tbody>")
	for _, key := range sortedIntKeys(report.Confidence) {
		fmt.Fprintf(&b, "<tr><td>%s</td><td>%d</td></tr>\n", html.EscapeString(key), report.Confidence[key])
	}
	fmt.Fprintln(&b, "</tbody></table>")
	fmt.Fprintln(&b, "<h2>Interpretation</h2>")
	fmt.Fprintln(&b, "<ul>")
	for _, note := range report.Notes {
		fmt.Fprintf(&b, "<li>%s</li>\n", html.EscapeString(note))
	}
	fmt.Fprintln(&b, "</ul>")
	fmt.Fprintln(&b, "</body></html>")
	return b.String()
}

func sortedFloatKeys(values map[string]float64) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedIntKeys(values map[string]int) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
