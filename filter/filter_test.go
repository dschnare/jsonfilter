package filter

import (
	"testing"
	"os"
	"strings"
)

func TestFilterJsonText_noFilter(t *testing.T) {
	var (
		expectedFilter = "trim"
		expectedJson = map[string]interface{}{
			"a": "Hello World!",
			"size": float64(27),
			"counts": []interface{}{float64(2), float64(4), float64(56), float64(7)},
			"b": map[string]interface{}{
				"c": "This is a line of text",
				"age": float64(35),
				"d": []interface{}{
					"This is some text", 
					"So is this",
				},
			},
		}
	)

	if file,err := os.Open("./fixtures/dataset-1.json"); err == nil {
		filterRunner := func(command string, value string) (string, error) {
			if command != expectedFilter {
				t.Fatalf("Expected command to be '%v' got '%v'", expectedFilter, command)
			}
			return value,nil
		}
		if value,err := FilterJsonFromReaderWithFilterRunner(file, expectedFilter, filterRunner); err != nil {
			t.Fatalf("Expected no error :: %v", err.Error())
		} else {
			switch value.(type) {
			case map[string]interface{}:
				testValue(value, expectedJson, t)
			default: t.Fatalf("Unexpected return value from FilterJsonFromReader :: %v", value)
			}
		}

	} else {
		t.Fatalf("Failed to open fixture :: %v", err.Error())
	}
}

func TestFilterJsonText_uppercaseFilter(t *testing.T) {
	var (
		expectedFilter = "upper"
		expectedJson = map[string]interface{}{
			"a": "HELLO WORLD!",
			"size": float64(27),
			"counts": []interface{}{float64(2), float64(4), float64(56), float64(7)},
			"b": map[string]interface{}{
				"c": "THIS IS A LINE OF TEXT",
				"age": float64(35),
				"d": []interface{}{
					"THIS IS SOME TEXT",
					"SO IS THIS",
				},
			},
		}
	)

	if file,err := os.Open("./fixtures/dataset-1.json"); err == nil {
		filterRunner := func(command string, value string) (string, error) {
			if command != expectedFilter {
				t.Fatalf("Expected command to be '%v' got '%v'", expectedFilter, command)
			}
			return strings.ToUpper(value),nil
		}
		if value,err := FilterJsonFromReaderWithFilterRunner(file, expectedFilter, filterRunner); err != nil {
			t.Fatalf("Expected no error :: %v", err.Error())
		} else {
			switch value.(type) {
			case map[string]interface{}:
				testValue(value, expectedJson, t)
			default: t.Fatalf("Unexpected return value from FilterJsonFromReader :: %v", value)
			}
		}

	} else {
		t.Fatalf("Failed to open fixture :: %v", err.Error())
	}
}

func TestFilterJsonText_filterResolution(t *testing.T) {
	var (
		expectedJson = map[string]interface{}{
			"a": "HELLO WORLD!",
			"size": float64(27),
			"counts": []interface{}{float64(2), float64(4), float64(56), float64(7)},
			"b": map[string]interface{}{
				"c": "this is a line of text",
				"age": float64(35),
				"d": []interface{}{
					"This Is Some Text",
					"--So is this--",
				},
			},
		}
	)

	if file,err := os.Open("./fixtures/dataset-1.json"); err == nil {
		filterRunner := func(command string, value string) (string, error) {
			switch command {
			case "upper": return strings.ToUpper(value),nil
			case "lower": return strings.ToLower(value),nil
			case "title": return strings.Title(value),nil
			case "custom": return "--" + value + "--",nil
			default: t.Fatalf("Unexpected command :: %v", command)
			}
			return value,nil
		}
		if value,err := FilterJsonFromReaderWithFilterRunner(file, "./fixtures/filters.json", filterRunner); err != nil {
			t.Fatalf("Expected no error :: %v", err.Error())
		} else {
			switch value.(type) {
			case map[string]interface{}:
				testValue(value, expectedJson, t)
			default: t.Fatalf("Unexpected return value from FilterJsonFromReader :: %v", value)
			}
		}

	} else {
		t.Fatalf("Failed to open fixture :: %v", err.Error())
	}
}

func TestFilterJsonText_arrayFilter(t *testing.T) {
	var (
		expectedJson = map[string]interface{}{
			"a": "HELLO WORLD!",
			"size": float64(27),
			"counts": []interface{}{float64(2), float64(4), float64(56), float64(7)},
			"b": map[string]interface{}{
				"c": "this is a line of text",
				"age": float64(35),
				"d": []interface{}{
					"This Is Some Text",
					"So Is This",
					"SOME TEXT",
				},
			},
		}
	)

	if file,err := os.Open("./fixtures/dataset-3.json"); err == nil {
		filterRunner := func(command string, value string) (string, error) {
			switch command {
			case "upper": return strings.ToUpper(value),nil
			case "lower": return strings.ToLower(value),nil
			case "title": return strings.Title(value),nil
			case "custom": return "--" + value + "--",nil
			default: t.Fatalf("Unexpected command :: %v", command)
			}
			return value,nil
		}
		if value,err := FilterJsonFromReaderWithFilterRunner(file, "./fixtures/array-filter.json", filterRunner); err != nil {
			t.Fatalf("Expected no error :: %v", err.Error())
		} else {
			switch value.(type) {
			case map[string]interface{}:
				testValue(value, expectedJson, t)
			default: t.Fatalf("Unexpected return value from FilterJsonFromReader :: %v", value)
			}
		}

	} else {
		t.Fatalf("Failed to open fixture :: %v", err.Error())
	}
}

func TestFilterJsonText_objectFilter(t *testing.T) {
	var (
		expectedJson = map[string]interface{}{
			"a": "HELLO WORLD!",
			"size": float64(27),
			"counts": []interface{}{float64(2), float64(4), float64(56), float64(7)},
			"b": map[string]interface{}{
				"c": "this is a line of text",
				"age": float64(35),
				"d": []interface{}{
					"this is some text",
					"so is this",
				},
			},
		}
	)

	if file,err := os.Open("./fixtures/dataset-1.json"); err == nil {
		filterRunner := func(command string, value string) (string, error) {
			switch command {
			case "upper": return strings.ToUpper(value),nil
			case "lower": return strings.ToLower(value),nil
			case "title": return strings.Title(value),nil
			case "custom": return "--" + value + "--",nil
			default: t.Fatalf("Unexpected command :: %v", command)
			}
			return value,nil
		}
		if value,err := FilterJsonFromReaderWithFilterRunner(file, "./fixtures/object-filter.json", filterRunner); err != nil {
			t.Fatalf("Expected no error :: %v", err.Error())
		} else {
			switch value.(type) {
			case map[string]interface{}:
				testValue(value, expectedJson, t)
			default: t.Fatalf("Unexpected return value from FilterJsonFromReader :: %v", value)
			}
		}

	} else {
		t.Fatalf("Failed to open fixture :: %v", err.Error())
	}
}

func TestFilterJsonText_arrayFilterHomogenous(t *testing.T) {
	var (
		expectedJson = map[string]interface{}{
			"a": []interface{}{
				map[string]interface{}{
					"name": "DARREN",
					"age": float64(32),
				},
				map[string]interface{}{
					"name": "MAX",
					"age": float64(22),
				},
				map[string]interface{}{
					"name": "ZARA",
					"age": float64(19),
				},
			},
		}
	)

	if file,err := os.Open("./fixtures/dataset-2.json"); err == nil {
		filterRunner := func(command string, value string) (string, error) {
			switch command {
			case "upper": return strings.ToUpper(value),nil
			default: t.Fatalf("Unexpected command :: %v", command)
			}
			return value,nil
		}
		if value,err := FilterJsonFromReaderWithFilterRunner(file, "./fixtures/array-filter-homogenous.json", filterRunner); err != nil {
			t.Fatalf("Expected no error :: %v", err.Error())
		} else {
			switch value.(type) {
			case map[string]interface{}:
				testValue(value, expectedJson, t)
			default: t.Fatalf("Unexpected return value from FilterJsonFromReader :: %v", value)
			}
		}

	} else {
		t.Fatalf("Failed to open fixture :: %v", err.Error())
	}
}

func TestFilterJsonText_arrayFilterHomogenousObject(t *testing.T) {
	var (
		expectedJson = map[string]interface{}{
			"a": []interface{}{
				map[string]interface{}{
					"name": "darren",
					"age": float64(32),
				},
				map[string]interface{}{
					"name": "max",
					"age": float64(22),
				},
				map[string]interface{}{
					"name": "zara",
					"age": float64(19),
				},
			},
		}
	)

	if file,err := os.Open("./fixtures/dataset-2.json"); err == nil {
		filterRunner := func(command string, value string) (string, error) {
			switch command {
			case "lower": return strings.ToLower(value),nil
			default: t.Fatalf("Unexpected command :: %v", command)
			}
			return value,nil
		}
		if value,err := FilterJsonFromReaderWithFilterRunner(file, "./fixtures/array-filter-homogenous-object.json", filterRunner); err != nil {
			t.Fatalf("Expected no error :: %v", err.Error())
		} else {
			switch value.(type) {
			case map[string]interface{}:
				testValue(value, expectedJson, t)
			default: t.Fatalf("Unexpected return value from FilterJsonFromReader :: %v", value)
			}
		}

	} else {
		t.Fatalf("Failed to open fixture :: %v", err.Error())
	}
}

func TestFilterJsonText_commandLine(t *testing.T) {
	var (
		expectedJson = map[string]interface{}{
			"a": []interface{}{
				map[string]interface{}{
					"name": "DARREN",
					"age": float64(32),
				},
				map[string]interface{}{
					"name": "MAX",
					"age": float64(22),
				},
				map[string]interface{}{
					"name": "ZARA",
					"age": float64(19),
				},
			},
		}
	)

	if file,err := os.Open("./fixtures/dataset-2.json"); err == nil {
		if value,err := FilterJsonFromReader(file, "./fixtures/command-filter.json"); err != nil {
			t.Fatalf("Expected no error :: %v", err.Error())
		} else {
			switch value.(type) {
			case map[string]interface{}:
				testValue(value, expectedJson, t)
			default: t.Fatalf("Unexpected return value from FilterJsonFromReader :: %v", value)
			}
		}

	} else {
		t.Fatalf("Failed to open fixture :: %v", err.Error())
	}
}

// ---

func testValue(value interface{}, expected interface{}, t *testing.T) {
  switch expected.(type) {
  case bool:
  	b := expected.(bool)
  	if _b,ok := value.(bool); ok {
  		if b != _b {
  			t.Fatalf("expected bools to be equal %v : got %v", b, _b)
  		}
  	} else {
  		t.Fatalf("expected a bool type")
  	}
  case float64:
  	f := expected.(float64)
  	if _f,ok := value.(float64); ok {
  		if f != _f {
  			t.Fatalf("expected floats to be equal %v : got %v", f, _f)
  		}
  	} else {
  		t.Fatalf("expected a float64 type")
  	}
  case string:
    str := expected.(string)
    if _str,ok := value.(string); ok {
      if str != _str {
        t.Fatalf("expected strings to be equal '%v' : got '%v'", str, _str)
      }
    } else {
      t.Fatalf("expected a string type")
    }
  case map[string]interface{}:
    m := expected.(map[string]interface{})
    if _m,ok := value.(map[string]interface{}); ok {
      testMap(_m, m, t)
    } else {
      t.Fatalf("expected a map type %v", value)
    }
  case []interface{}:
    slice := expected.([]interface{})
    if _slice,ok := value.([]interface{}); ok {
      testSlice(_slice, slice, t)
    } else {
      t.Fatalf("expected a slice type :: %v", value)
    }
  }
}

func testMap(obj map[string]interface{}, expected map[string]interface{}, t *testing.T) {
  for k,v := range expected {
    if _v,ok := obj[k]; ok {
      testValue(_v, v, t)
    } else {
      t.Fatalf("expected map to have key '%v'", k)
    }
  }
}

func testSlice(slice []interface{}, expected []interface{}, t *testing.T) {
  if len(expected) != len(slice) {
    t.Fatalf("expected slices to be same length")
  }

  for k,v := range expected {
    testValue(slice[k], v, t)
  }
}