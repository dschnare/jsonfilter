PACKAGE DOCUMENTATION

package filter
    import "github.com/dschnare/jsonfilter/filter"

    Package filter implements a JSON filter that filters all string values
    found in JSON data. A filter can be specified at the command line or as
    a JSON file. If a JSON file is specified then each key can have a
    specific filter. Otherwise if filter is not a JSON file the filter will
    be used for all string values found in the JSON data.

    Each filter is a command that will be executed on the command line with
    the string value to be filtered piped into stdin. The filtered value is
    expected to be piped to stdout.

    For example, to convert lowercase characters to uppercase:

	"tr '[:lower:]' '[:upper:]'"

    Optionally you can override how filters are run by calling
    *WithFilterRunner(). By default the filter runner will run each filter
    as a command on the command line. You can use a custom filter runner to
    define a filter language of your own or use it to mock out a test.

    To specify unique filters for specific JSON paths you can use a JSON
    file.

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

    Using the same filter file to filter a more complex JSON data object
    will result in any string value found under the path "a" being filtered.

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

    You can also specify unique paths by defining a filter with a longer
    path.

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

    Paths work with arrays too. To specify a filter for all items in an
    array then your path must lead to an array with a single element. That
    single element is your filter to use for each array item found in the
    JSON data. The filter can either be a command or even another array or
    object.

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

    If you don't specify a filter for each item in the array then all string
    values that come after the last defined filter will not be filtered.

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

FUNCTIONS

func FilterJsonFromReader(reader io.Reader, filter string) (value interface{}, err error)
    FilterJsonFromText filters JSON data from a reader. The filter can
    either be a command or a path to a JSON file. If filter is a command
    then all string values found in the JSON data will be filtered using the
    command. Returns the unmarshalled JSON data with all string values
    filtered. See http://golang.org/pkg/encoding/json/#Unmarshal for more
    details on the value returned.

func FilterJsonFromReaderWithFilterRunner(reader io.Reader, filter string, filterRunner FilterRunner) (value interface{}, err error)
    FilterJsonFromText filters JSON data from a reader using a custom filter
    runner. The filter can either be a command or a path to a JSON file. If
    filter is a command then all string values found in the JSON data will
    be filtered using the command. Returns the unmarshalled JSON data with
    all string values filtered. See
    http://golang.org/pkg/encoding/json/#Unmarshal for more details on the
    value returned.

func FilterJsonFromText(jsonText string, filter string) (value interface{}, err error)
    FilterJsonFromText filters JSON data from text. The filter can either be
    a command or a path to a JSON file. If filter is a command then all
    string values found in the JSON data will be filtered using the command.
    Returns the unmarshalled JSON data with all string values filtered. See
    http://golang.org/pkg/encoding/json/#Unmarshal for more details on the
    value returned.

func FilterJsonFromTextWithFilterRunner(jsonText string, filter string, filterRunner FilterRunner) (value interface{}, err error)
    FilterJsonFromText filters JSON data from text using a custom filter
    runner. The filter can either be a command or a path to a JSON file. If
    filter is a command then all string values found in the JSON data will
    be filtered using the command. Returns the unmarshalled JSON data with
    all string values filtered. See
    http://golang.org/pkg/encoding/json/#Unmarshal for more details on the
    value returned.

TYPES

type FilterRunner func(command string, value string) (string, error)
    FilterRunner defines the function signature used to override how filters
    are run. Each filter runner function is passed the raw command to run
    and the string value to filter.

SUBDIRECTORIES

	fixtures

