Story Aggregator
=====

Usage: slush [options] url1 url2 ... urlN

This program aggregates meta information from mcclatchy-based articles and prints JSON to os.Stdout. 

### Options

`-template [path]`  
Pass a path to a template file and slush will pass all articles as an .Articles object. 

`-section`  
The section flag will pull all stories listed on a section page.

`-body`  
By default the JSON feed will not include the body. Adding this flag swaps that behavior
