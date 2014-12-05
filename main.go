// Copyright 2014 Darren Schnare. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

/*
Jsonfilter implements a command line tool to filter string values found in JSON data.
Use `jsonfilter --help` for details about each command line argument.

This tool supports standard in piping of a JSON file contents. If
the `--output` argument is not specified then output will be piped to
standard out.

  jsonfilter "json to filter" | jsonfilter [help|/?]
    -filter="": The filter(s) to apply to the strings contained in the JSON file.
    -help=false: Show the help message.
    -output="": The output file to write to.
    -pretty=false: Print JSON result with indentation. (shorthand)
    -pretty-print=false: Print JSON result with indentation.
*/
package main

import (
  "flag"
  "os"
  "io"
  "bytes"
  "fmt"
  "bufio"
  "encoding/json"
  jsonfilter "github.com/dschnare/jsonfilter/filter"
)

var (
  jsontext string
  // Flags
  output string
  help bool
  filter string
  prettyPrint bool
)

func usage() {
  fmt.Fprintf(os.Stderr, "Usage: jsonfilter \"json to filter\" | jsonfilter [help|/?]\n")
  flag.PrintDefaults()
}

func init() {
  const (
    helpDefault = false
    helpUsage = "Show the help message."
    outputDefault = ""
    outputUsage = "The output file to write to."
    filterDefault = ""
    filterUsage = "The filter(s) to apply to the strings contained in the JSON file."
    prettyPrintDefault = false
    prettyPrintUsage = "Print JSON result with indentation."
  )

  flag.Usage = usage

  flag.BoolVar(&help, "help", helpDefault, helpUsage)

  flag.BoolVar(&prettyPrint, "pretty-print", prettyPrintDefault, prettyPrintUsage)
  flag.BoolVar(&prettyPrint, "pretty", prettyPrintDefault, prettyPrintUsage + " (shorthand)")

  flag.StringVar(&filter, "filter", filterDefault, filterUsage)

  flag.StringVar(&output, "output", outputDefault, outputUsage)

  flag.Parse()

  if help {
    flag.Usage()
    os.Exit(0)
  } else if len(flag.Args()) == 1 && (flag.Arg(0) == "/?" || flag.Arg(0) == "help") {
    flag.Usage()
    os.Exit(0)
  } else if len(flag.Args()) == 1 {
    jsontext = flag.Arg(0)
  } else {
    if isPiped(os.Stdin) {
      var err error
      if jsontext,err = readFile(os.Stdin); err != nil {
        fmt.Printf("Failed to read from stdin :: %v\n", err.Error())
        os.Exit(1)
      }
    } else {
      flag.Usage()
      os.Exit(1)
    }
  }

  if len(filter) == 0 {
    fmt.Println("Expected a filter to be specified.")
    flag.Usage()
    os.Exit(1)
  }
}

func main() {
  if len(jsontext) == 0 {
    return
  }

  if value,err := jsonfilter.FilterJsonFromText(jsontext, filter); err == nil {
    if writer,err := createWriter(); err == nil {
      if err := doWrite(writer, value); err != nil {
        panic(err)
      }
    } else {
      panic(err)
    }
  } else {
    panic(err)   
  }
}

func isPiped(file *os.File) bool {
  if info,err := file.Stat(); err == nil {
  return info.Mode() == os.ModeNamedPipe
  }
  return false
}

func doWrite(writer *bufio.Writer, value interface{}) (err error) {
  var(
    b []byte
  )

  if prettyPrint {
    if b,err = json.Marshal(value); err == nil {
      var out bytes.Buffer
      if err = json.Indent(&out, b, "", "  "); err == nil {
        if _,err = out.WriteTo(writer); err == nil {
          err = writer.Flush()
        }
      }
    }
  } else {
    jsonEncoder := json.NewEncoder(writer)
    if err = jsonEncoder.Encode(value); err == nil {
      err = writer.Flush()
    }
  }

  return
}

func createWriter() (*bufio.Writer, error) {
  var writer *bufio.Writer

  if len(output) == 0 || isPiped(os.Stdout) {
    writer = bufio.NewWriter(os.Stdout)
  } else if file,err := os.Create(output); err == nil {
    writer = bufio.NewWriter(file)
  } else {
    return nil,err
  }

  return writer,nil
}

func readFile(file *os.File) (text string, err error) {
  var (
    buf bytes.Buffer
    c byte
  )
  reader := bufio.NewReader(file)

  for err == nil {
    c,err = reader.ReadByte()

    if err == nil {
      buf.WriteByte(c)
    }
  }

  if err == io.EOF {
    err = nil
  }

  text = buf.String()

  return
}