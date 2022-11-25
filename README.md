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

##Update Nov 26, 2022 1:30 PM CST

Based upon feedback from a post on this topic on reddit: [goroutines and goccy/go-graphviz package](https://www.reddit.com/r/golang/comments/z474g4/goroutines_and_goccygographviz_package/), I have added a mutex to insure that only one graphviz coroutine is active at one time.
This is controlled with the -use_gmutux command line switch.
This case with goroutines works.  Unfortunately, this limits the performance
run times similar to the no goroutines base.

I also realized that some of the dumps contains faults internal to the 
graphviz code handling layout for SVG generation.  I added the ability
to generate DOT output so that the graphviz layout code did not come
into play.  This is controlled with the file_type command line switch.

To assist with controlling the size of the graphs whent testing with SVG, 
I added minWidth and minDepth command line switches.

The full syntax is now:
```
/go-graphviz-test -h
Usage of ./go-graphviz-test:
  -ct int
    	number of graphs to test per run (default 1)
  -file_type string
    	file type output (default "svg")
  -maxDepth uint
    	maximum depth of test digraphs (default 6)
  -maxWidth uint
    	maximum width of test digraphs (default 6)
  -minDepth uint
    	minimum depth of test digraphs (default 2)
  -minWidth uint
    	minimum width of test digraphs (default 2)
  -use_gmutex
    	use mutex to limit graphviz operations (default true)
  -use_goroutines
    	use goroutines (default true)
```
Based upon my current understanding, go-graphviz is not threadsafe when
using lexically scoped variables such as used in `createSvg`