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

lbe Nov 25, 2022
