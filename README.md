# Introduction

This is an argument parser which can parse a certain formatted string to an argument array. For example:
```
echo -n "Hello \"world" ==> [echo, -n, Hello "world]
```

# Constructor
* NewArgParser  
```go
NewArgParser(quotes string, splitors string, escaper rune, lenient bool)   
```

create an arguments parser; any of the `splitors` inside a pair of identical `quotes` is treated as normal character; 
the `escaper` can escape a following rune; we can have more than one pair of `quotes` as well as `splitors`, yet only a pair of 
**identical** `quotes` are a integra scope of quote and all characters inside this scope will be treated as text except for the 
`escaper`; if the `lenient` is true, this parser will not report any error if any of the scopes is not finished( *meaning that 
something does not end properly* ) after the parsing process is over.  

* GeneralArgParser  

the general purpose `lenient` arguments parser with `"` and `'` as its `quotes`, `\t` as its `splitors`, `\` as its escaper.

# Member functions

* Parse  
```go
Parse(text string)
```
parse the given text
