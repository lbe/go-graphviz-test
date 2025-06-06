package main

import (
	"bytes"
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type graphParams struct { // parameters defining width and depth of test digraphs
	m uint64
	n uint64
}

type edge struct { // edge in digraph
	from string
	to   string
}

type graphData struct { // slice of edges, can be extended for richer use
	edges []edge
}

// create children of node for test digraph
func getChildren(p string, n uint64) (c []string) {
	for i := uint64(1); i <= n; i++ {
		c = append(c, p+"."+strconv.FormatUint(i, 10))
	}
	return c
}

// create a digraph to use in testing
func createGraphData(p graphParams, d *graphData, i uint64, r string) {
	if i++; i < p.m {
		c := getChildren(r, p.n)
		for _, s := range c {
			d.edges = append(d.edges, edge{r, s})
			createGraphData(p, d, i, s)
		}
	}
}

// generate a set of test graphs of varying width and depth
func genGraphData(m, n uint64) map[graphParams]graphData {
	gd := make(map[graphParams]graphData)

	for i := uint64(2); i <= m; i++ {
		for j := uint64(2); j <= n; j++ {
			d := graphData{
				edges: []edge{},
			}
			createGraphData(graphParams{i, j}, &d, 0, "1")
			gd[graphParams{i, j}] = d
		}
	}

	return gd
}

// generate graphParams to be used for a specific test run
func getRandomGraphParams(minWidth, maxWidth, minDepth, maxDepth int) graphParams {
	return graphParams{
		m: uint64(rand.Intn(maxWidth-minWidth) + minWidth),
		n: uint64(rand.Intn(maxDepth-minDepth) + minDepth),
	}
}

var (
	gd     map[graphParams]graphData // map to hold the test digraphs
	wg     sync.WaitGroup            // wait group used for goroutines
	gMutex sync.Mutex                // mutex to insure exclusivity of graphviz operations
)

// write the buffer to a file, creating it if it doesn't exist
func writeBufToFile(buf *bytes.Buffer, fn string) {
	// Open the file for writing.  Create it if it doesn't exist,
	// truncate it if it does.
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write the buffer to the file.
	_, err = io.Copy(file, buf)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure all data is written to disk.
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

// convert a test digraph to a graphViz graph and generate the output as SVG
func createSvg(id string, p graphParams, file_type string, use_gmutex bool) {
	if use_gmutex {
		gMutex.Lock()
		defer gMutex.Unlock()
	}
	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()
	// set some defaults for the graph
	graph.SetLabel(id)
	graph.SetLabelLocation("t")
	graph.SetPageDir("TL")
	graph.SetRankDir("LR")
	graph.SetFontColor("black")

	nodes := make(map[string]*cgraph.Node) // map to track which nodes have been created

	for _, each := range gd[p].edges {
		if _, ok := nodes[each.from]; !ok { // create from node if not exists
			n, err := graph.CreateNodeByName(each.from)
			if err != nil {
				log.Fatal(err)
			}
			nodes[each.from] = n
		}
		if _, ok := nodes[each.to]; !ok { // create to node if not exists
			n, err := graph.CreateNodeByName(each.to)
			if err != nil {
				log.Fatal(err)
			}
			nodes[each.to] = n
		}
		// create the edge between the from and to nodes
		_, err := graph.CreateEdgeByName(each.from+"-"+each.to, nodes[each.from], nodes[each.to])
		if err != nil {
			log.Fatal(err)
		}
	}

	fn_out := path.Join("./"+file_type, id+"."+file_type) // generate output file
	switch file_type {
	case "svg":
		var bufSVG bytes.Buffer
		if err := g.Render(ctx, graph, graphviz.SVG, &bufSVG); err != nil {
			log.Println(p)
			log.Fatal(err)
		}
		writeBufToFile(&bufSVG, fn_out)
	case "dot":
		if err := g.RenderFilename(ctx, graph, graphviz.Format(graphviz.DOT), fn_out); err != nil {
			log.Println(p)
			log.Fatal(err)
		}
	default:
		log.Fatal("Unsupported file_type: " + file_type)
	}

	wg.Done()
}

func main() {
	var ct_graphs int64                               // number of graphviz calls
	var minWidth, maxWidth, minDepth, maxDepth uint64 // maximum width and depth of auto generated digraphs
	var use_goroutines bool
	var file_type string
	var use_gmutex bool

	flag.Int64Var(&ct_graphs, "ct", 1, "number of graphs to test per run")
	flag.Uint64Var(&minWidth, "minWidth", 2, "minimum width of test digraphs")
	flag.Uint64Var(&maxWidth, "maxWidth", 6, "maximum width of test digraphs")
	flag.Uint64Var(&minDepth, "minDepth", 2, "minimum depth of test digraphs")
	flag.Uint64Var(&maxDepth, "maxDepth", 6, "maximum depth of test digraphs")
	flag.StringVar(&file_type, "file_type", "svg", "file type output")
	flag.BoolVar(&use_goroutines, "use_goroutines", true, "use goroutines")
	flag.BoolVar(&use_gmutex, "use_gmutex", true, "use mutex to limit graphviz operations")

	flag.Parse()

	err := os.MkdirAll("./"+file_type, 0o755) // mkdir if not exists
	if err != nil {
		log.Fatal(err)
	}

	gd = genGraphData(maxWidth, maxDepth) // generate base graph data - a set of edges in a direct graph of varying sizes

	for i := int64(1); i <= ct_graphs; i++ {
		wg.Add(1)
		if use_goroutines {
			go createSvg(strconv.FormatInt(i, 10), getRandomGraphParams(int(minWidth), int(maxWidth), int(minDepth), int(maxDepth)), file_type, use_gmutex)
		} else {
			createSvg(strconv.FormatInt(i, 10), getRandomGraphParams(int(minWidth), int(maxWidth), int(minDepth), int(maxDepth)), file_type, use_gmutex)
		}
	}

	wg.Wait()
}
