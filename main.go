package main

import (
	"log"
	"math"
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

type graphData struct {  // slice of edges, can be extended for richer use
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
func getRandomGraphParams(min, max int) graphParams {
	return graphParams{
		m: uint64(rand.Intn(max-min) + min),
		n: uint64(rand.Intn(max-min) + min),
	}
}

var (
	gd map[graphParams]graphData // map to hold the test digraphs
	wg sync.WaitGroup // wait group used for goroutines
)

// convert a test digraph to a graphViz graph and generate the output as SVG
func createSvg(id string, p graphParams) {
	g := graphviz.New()
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
			n, err := graph.CreateNode(each.from)
			if err != nil {
				log.Fatal(err)
			}
			nodes[each.from] = n
		}
		if _, ok := nodes[each.to]; !ok { // create to node if not exists
			n, err := graph.CreateNode(each.to)
			if err != nil {
				log.Fatal(err)
			}
			nodes[each.to] = n
		}
		// create the edge between the from and to nodes
		_, err := graph.CreateEdge(each.from+"-"+each.to, nodes[each.from], nodes[each.to])
		if err != nil {
			log.Fatal(err)
		}
	}

	fn_out := path.Join("./Svg", id+".svg") // generate output file
	if err := g.RenderFilename(graph, graphviz.SVG, fn_out); err != nil {
		log.Fatal(err)
	}

	wg.Done()
}

func main() {
	var ct_graphs int64 = int64(math.Pow(2, 2)) // number of graphviz calls

	err := os.MkdirAll("./Svg", 0o755) // mkdir if not exists
	if err != nil {
		log.Fatal(err)
	}

	gd = genGraphData(6, 6) // generate base graph data - a set of edges in a direct graph of varying sizes

	for i := int64(1); i <= ct_graphs; i++ {
		wg.Add(1)
		go createSvg(strconv.FormatInt(i, 10), getRandomGraphParams(2, 6))
	}

	wg.Wait()
}

func Pow(i1, i2 int) {
	panic("unimplemented")
}
