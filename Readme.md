[![GoDoc](https://godoc.org/github.com/dschnare/jsonfilter?status.svg)](https://godoc.org/github.com/dschnare/jsonfilter)

# Overview

Jsonfilter exposes a Go packge and a command line tool for filter string values found in JSON data.

# Usage

	jsonfilter "json to filter" | jsonfilter [help|/?]
		-filter="": The filter(s) to apply to the strings contained in the JSON file.
		-help=false: Show the help message.
		-output="": The output file to write to.
		-pretty=false: Print JSON result with indentation. (shorthand)
		-pretty-print=false: Print JSON result with indentation.

Where `filter` can either be a command to use to filter all string values or a path to a JSON file.

If no JSON is specified as an argument then it is expected to be piped into stdin.

If no output file is specified as an argument then the output is piped to stdout.

# Filtering

A filter can be specified at the command line or as a JSON file. If a JSON file is specified then
each key can have a specific filter. Otherwise if filter is not a JSON file the filter will be used
for all string values found in the JSON data.

Each filter is a command that will be executed on the command line with the string value to be
filtered piped into stdin. The filtered value is expected to be piped to stdout.

For example, to convert lowercase characters to uppercase:

	"tr '[:lower:]' '[:upper:]'"

Optionally you can override how filters are run by calling **WithFilterRunner()**. By default the filter runner
will run each filter as a command on the command line. You can use a custom filter runner to define a filter
language of your own or use it to mock out a test.

To specify unique filters for specific JSON paths you can use a JSON file.

	// filter.json
	{
	"a": "tr '[:lower:]' '[:upper:]'"
	}

	// data.json
	{
	"a": "hi"
	}

	// result
	{
	"a": "HI"
	}

Using the same filter file to filter a more complex JSON data object will result
in any string value found under the path "a" being filtered.

	// data2.json
	{
	"a": ["yep", "nope"]
	}

	// result
	{
	"a": ["YEP", "NOPE"]
	}

	// ---

	// data3.json
	{
	"a": {"type": "message", "value": "hello world!"}
	}

	// result
	{
	"a": {"type": "MESSAGE", "value:" "HELLO WORLD!""}
	}


You can also specify unique paths by defining a filter with a longer path.

	// filter2.json
	{
	"a": {"value": "tr '[:lower:]' '[:upper:]'"}
	}

	// data3.json
	{
	"a": {"type": "message", "value": "hello world!"}
	}

	// result
	{
	"a": {"type": "message", "value:" "HELLO WORLD!""}
	}

Paths work with arrays too. To specify a filter for all items in an array
then your path must lead to an array with a single element. That single
element is your filter to use for each array item found in the JSON data.
The filter can either be a command or even another array or object.

	// filter3.json
	{
	"a": [{"value": "tr '[:lower:]' '[:upper:]'"}]
	}

	// data3.json
	{
	"a": [{"type": "message", "value": "hello world!"}]
	}

	// result
	{
	"a": [{"type": "message", "value:" "HELLO WORLD!""}]
	}

You can even specify a unique filter for each item in an array.

	// filter4.json
	{
	"a": ["tr '[:lower:]' '[:upper:]'", "tr '[:upper:]' '[:lower:]'"]
	}

	// data4.json
	{
	"a": ["hello world!", "APPLES"]
	}

	// result
	{
	"a": ["HELLO WORLD!", "apples"]
	}

If you don't specify a filter for each item in the array then all string values
that come after the last defined filter will not be filtered.

	// filter4.json
	{
	"a": ["tr '[:lower:]' '[:upper:]'", "tr '[:upper:]' '[:lower:]'"]
	}

	// data5.json
	{
	"a": ["hello world!", "APPLES", "This text will be left as-is."]
	}

	// result
	{
	"a": ["HELLO WORLD!", "apples", "This text will be left as-is."]
	}


# Packages

**[filter](http://godoc.org/github.com/dschnare/jsonfilter/filter)** - Package filter implements a JSON filter that filters all string values found in JSON data.

# Commands

**[jsonfilter](http://godoc.org/github.com/dschnare/jsonfilter)** - Jsonfilter implements a command line tool to filter string values found in JSON data.