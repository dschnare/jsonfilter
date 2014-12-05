/*
Package filter implements a JSON filter that filters all string values found in JSON data.
A filter can be specified at the command line or as a JSON file. If a JSON file is specified then
each key can have a specific filter. Otherwise if filter is not a JSON file the filter will be used
for all string values found in the JSON data.

Each filter is a command that will be executed on the command line with the string value to be
filtered piped into stdin. The filtered value is expected to be piped to stdout.

For example, to convert lowercase characters to uppercase:

  "tr '[:lower:]' '[:upper:]'"

Optionally you can override how filters are run by calling *WithFilterRunner(). By default the filter runner
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

*/
package filter

import (
  "os"
  "fmt"
  "os/exec"
  "io"
  "bytes"
  "bufio"
  "strings"
  "strconv"
  "encoding/json"
)

// FilterRunner defines the function signature used to override how filters are run.
// Each filter runner function is passed the raw command to run and the string value to filter.
type FilterRunner func(command string, value string) (string, error)

type visitorFunc func(path string, value string) (string, error)

// FilterJsonFromText filters JSON data from text. The filter can either be a command
// or a path to a JSON file. If filter is a command then all string values found in the JSON data
// will be filtered using the command. Returns the unmarshalled JSON data with all string values filtered.
// See http://golang.org/pkg/encoding/json/#Unmarshal for more details on the value returned.
func FilterJsonFromText(jsonText string, filter string) (value interface{},  err error) {
  if value,err = readJsonFromText(jsonText); err == nil {
    err = doFilter(value, filter, nil)
  }
  return
}

// FilterJsonFromText filters JSON data from text using a custom filter runner. The filter can either be a command
// or a path to a JSON file. If filter is a command then all string values found in the JSON data
// will be filtered using the command. Returns the unmarshalled JSON data with all string values filtered.
// See http://golang.org/pkg/encoding/json/#Unmarshal for more details on the value returned.
func FilterJsonFromTextWithFilterRunner(jsonText string, filter string, filterRunner FilterRunner) (value interface{},  err error) {
  if value,err = readJsonFromText(jsonText); err == nil {
    err = doFilter(value, filter, filterRunner)
  }
  return
}

// FilterJsonFromText filters JSON data from a reader. The filter can either be a command
// or a path to a JSON file. If filter is a command then all string values found in the JSON data
// will be filtered using the command. Returns the unmarshalled JSON data with all string values filtered.
// See http://golang.org/pkg/encoding/json/#Unmarshal for more details on the value returned.
func FilterJsonFromReader(reader io.Reader, filter string) (value interface{}, err error) {
  if value,err = readJsonFromReader(reader); err == nil {
    err = doFilter(value, filter, nil)
  }
  return
}

// FilterJsonFromText filters JSON data from a reader using a custom filter runner. The filter can either be a command
// or a path to a JSON file. If filter is a command then all string values found in the JSON data
// will be filtered using the command. Returns the unmarshalled JSON data with all string values filtered.
// See http://golang.org/pkg/encoding/json/#Unmarshal for more details on the value returned.
func FilterJsonFromReaderWithFilterRunner(reader io.Reader, filter string, filterRunner FilterRunner) (value interface{}, err error) {
  if value,err = readJsonFromReader(reader); err == nil {
    err = doFilter(value, filter, filterRunner)
  }
  return
}

func doFilter(value interface{}, filter string, filterRunner FilterRunner) (err error) {
  var filters interface{}

  if filters,err = loadFilters(filter); err == nil {
    _,err = traverse(value, func (path string, value string) (string, error) {
      return doRunFilter(path, value, filters, filterRunner)
    })
  }

  return
}

func doRunFilter(path string, value string, filters interface{}, filterRunner FilterRunner) (result string, err error) {
  if command,ok := getFilterCommand(path, filters); ok {
    if filterRunner == nil {
      return commandLineFilterRunner(command, value)
    } else {
      return filterRunner(command, value)
    }
  } else {
    result = value
  }

  return
}

func commandLineFilterRunner(command string, value string) (result string, err error) {
  var out bytes.Buffer
  parts := strings.Split(command, " ")
  cmd := exec.Command(parts[0], parts[1:]...)
  cmd.Stdin = strings.NewReader(value)
  cmd.Stdout = &out

  if err = cmd.Run(); err == nil {
    result = out.String()
  }

  return
}

func getFilterCommand(path string, filters interface{}) (command string, found bool) {
  var filterCommand interface{}
  // Path will be of the form:
  // ['key']['key'][num]['key'][num]
  keys := strings.FieldsFunc(path, func (r rune) bool {
    return r == '[' || r == ']'
  })

  if filterCommand,found = getFilterCommandRec(keys, filters); found {
    command,found = filterCommand.(string)
  }

  return
}

func getFilterCommandRec(keys []string, filters interface{}) (interface{}, bool) {
  if len(keys) == 0 {
    return filters,true
  }

  key := strings.Trim(keys[0], "'")

  switch filters.(type) {
  case string: 
    return filters,true
  case map[string]interface{}:
    m := filters.(map[string]interface{})
    if v,ok := m[key]; ok {
      return getFilterCommandRec(keys[1:], v)
    } else {
      return nil,false
    }
  case []interface{}:
    s := filters.([]interface{})
    if len(s) == 1 {
      return getFilterCommandRec(keys[1:], s[0])
    } else if i,err := strconv.ParseInt(key, 10, 32); err == nil && i < int64(len(s)) && i >= 0 {
      return getFilterCommandRec(keys[1:], s[i])
    } else {
      return nil,false
    }
  default: return nil,false
  }
}

func readJsonFromFile(fileName string) (interface{}, error) {
  var err error
  var file *os.File

  if file,err = os.Open(fileName); err == nil {
    return readJsonFromReader(bufio.NewReader(file))
  }

  return nil,err
}

func readJsonFromText(jsonText string) (interface{}, error) {
  return readJsonFromReader(strings.NewReader(jsonText))
}

func readJsonFromReader(reader io.Reader) (value interface{}, err error) {
  decoder := json.NewDecoder(reader)

  if err = decoder.Decode(&value); err == io.EOF {
    err = nil
  }

  return
}

func loadFilters(filter string) (interface{}, error) {
  if strings.HasSuffix(filter, ".json") {
    return readJsonFromFile(filter)
  } else {
    return filter,nil
  }
}

func traverse(value interface{}, visit visitorFunc) (interface{}, error) {
  return traverseWithPath(value, "", visit)
}

func traverseWithPath(value interface{}, path string, visit visitorFunc) (interface{}, error) {
  switch value.(type) {
  case string: return visit(path, value.(string))
  case map[string]interface{}: return traverseMap(value.(map[string]interface{}), path, visit)
  case []interface{}: 
    slice := value.([]interface{})
    return traverseSlice(&slice, path, visit)
  }

  return value,nil
}

func traverseMap(m map[string]interface{}, path string, visit visitorFunc) (value interface{}, err error) {
  value = m
  for k,v := range m {
    if m[k],err = traverseWithPath(v, fmt.Sprintf("%s['%s']", path, k), visit); err != nil {
      break
    }
  }
  return
}

func traverseSlice(s *[]interface{}, path string, visit visitorFunc) (value interface{}, err error) {
  slice := *s
  value = slice
  for k,v := range slice {
    if slice[k],err = traverseWithPath(v, fmt.Sprintf("%s[%d]", path, k), visit); err != nil {
      break
    }
  }
  return
}