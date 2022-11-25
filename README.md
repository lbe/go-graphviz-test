# go-graphviz-test
Generate SVGs of Graphs using go-graphiz pkg

I am new to Go.  I have created this repo to request help to solve a problem 
that has me banging my head against the wall!

The problem: my code crashes when I use goroutines when I have more that 2 
goroutines running.

The code: This code generates a set of 49 directed graphs (digraphs) and 
then based upon some settings in the code, processes one or more of these graphs
though the go package goccy/go-graphviz and saves the output as an SVG file.  

The base code works as expected when go routines are not used or when go 
routines are used and there is only 1 or 2 graphs generated.  Anything past 2 results
in a segmentation fault.

Usage for this program is as follows:
```
./go-graphviz-test -h
Usage of ./go-graphviz-test:
  -ct int
    	number of graphs to test per run (default 1)
  -maxDepth uint
    	maximum depth of test digraphs (default 6)
  -maxWidth uint
    	maximum width of test digraphs (default 6)
  -use_goroutines
    	use go routines (default true)
```

The following work without a problem
```
./go-graphviz-test -use_goroutines=false -ct 1
./go-graphviz-test -use_goroutines=false -ct 2
./go-graphviz-test -use_goroutines=false -ct 3
./go-graphviz-test -use_goroutines=false -ct 1000
./go-graphviz-test -use_goroutines=true -ct 1
./go-graphviz-test -use_goroutines=true -ct 2
```

The following fail with dumps when calling the `RenderFilename` on one or more goroutines
```
./go-graphviz-test -use_goroutines=true -ct 3
./go-graphviz-test -use_goroutines=true -ct 4
```

Looking at the dumps, I see many more go routines dumping that I expect.
I suspect that go is breaking my function createSvg down into more than just
1 go routine as I would expect.  I have scoured documentation on go routines
but have not found anything to confirm this or to suggest what I should do
differently.

Any advice will be appreciated.

lbe Nov 25, 2022
