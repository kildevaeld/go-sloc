package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"sort"
	"text/tabwriter"

	"github.com/kildevaeld/go-sloc/sloc"
)

type LData []LResult

func (d LData) Len() int { return len(d) }

func (d LData) Less(i, j int) bool {
	if d[i].CodeLines == d[j].CodeLines {
		return d[i].Name > d[j].Name
	}
	return d[i].CodeLines > d[j].CodeLines
}

func (d LData) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type LResult struct {
	Name         string
	FileCount    int
	CodeLines    int
	CommentLines int
	BlankLines   int
	TotalLines   int
}

func (r *LResult) Add(a LResult) {
	r.FileCount += a.FileCount
	r.CodeLines += a.CodeLines
	r.CommentLines += a.CommentLines
	r.BlankLines += a.BlankLines
	r.TotalLines += a.TotalLines
}

func printJSON(info map[string]*sloc.Stats) {
	bs, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}

func printInfo(info map[string]*sloc.Stats) {
	w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Language\tFiles\tCode\tComment\tBlank\tTotal\t")
	d := LData([]LResult{})
	total := &LResult{}
	total.Name = "Total"
	for n, i := range info {
		r := LResult{n, i.FileCount, i.CodeLines, i.CommentLines, i.BlankLines, i.TotalLines}
		d = append(d, r)
		total.Add(r)
	}
	d = append(d, *total)
	sort.Sort(d)
	for _, i := range d {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\t\n", i.Name, i.FileCount, i.CodeLines, i.CommentLines, i.BlankLines, i.TotalLines)
	}

	w.Flush()
}

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	useJson    = flag.Bool("json", false, "JSON-format output")
	version    = flag.Bool("V", false, "display version info and exit")
)

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("sloc %s\n", sloc.VERSION)
		return
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			return
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, `.`)
	}

	sloc := &sloc.SlocCounter{}

	for _, n := range args {
		sloc.Add(n)
		//add(n)
	}

	/*for _, f := range files {
		handleFile(f)
	}*/

	info := sloc.Sloc()

	if *useJson {
		printJSON(info)
	} else {
		printInfo(info)
	}
}
