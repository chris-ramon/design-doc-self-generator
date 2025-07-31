# DrawIO Gantt Chart Parser

This package provides Go structs for parsing and working with DrawIO Gantt chart XML files.

## Overview

The package defines Go structs that correspond to the XML structure of DrawIO files, allowing you to:
- Parse existing DrawIO Gantt charts
- Extract task information, dates, and durations
- Create new DrawIO files programmatically
- Modify existing charts

## Structs

### Core Structs

- `MxFile`: Root element of a DrawIO file
- `Diagram`: Individual diagram within the file
- `MxGraphModel`: Graph model containing visual elements
- `Root`: Container for all cells
- `MxCell`: Individual cell representing tasks, labels, or visual elements
- `MxGeometry`: Position and size information for cells

## Usage

### Parsing an Existing DrawIO File

```go
package main

import (
    "encoding/xml"
    "fmt"
    "os"
    
    "your-module/drawio/gantt"
)

func main() {
    // Read the DrawIO file
    data, err := os.ReadFile("path/to/your/gantt.drawio")
    if err != nil {
        panic(err)
    }

    // Parse the XML
    var mxFile gantt.MxFile
    err = xml.Unmarshal(data, &mxFile)
    if err != nil {
        panic(err)
    }

    // Access the parsed data
    fmt.Printf("File version: %s\n", mxFile.Version)
    fmt.Printf("Number of diagrams: %d\n", len(mxFile.Diagrams))
    
    diagram := mxFile.Diagrams[0]
    fmt.Printf("Diagram name: %s\n", diagram.Name)
    fmt.Printf("Number of cells: %d\n", len(diagram.MxGraphModel.Root.Cells))
    
    // Extract tasks
    for _, cell := range diagram.MxGraphModel.Root.Cells {
        if cell.Value != "" && cell.MxGeometry != nil {
            fmt.Printf("Task: %s\n", cell.Value)
        }
    }
}
```

### Creating a New DrawIO File

```go
package main

import (
    "encoding/xml"
    "fmt"
    
    "your-module/drawio/gantt"
)

func main() {
    // Create a new DrawIO structure
    mxFile := gantt.MxFile{
        Host:    "Go",
        Agent:   "Go XML Generator",
        Version: "1.0",
        Diagrams: []gantt.Diagram{
            {
                Name: "Project-Gantt",
                ID:   "gantt-1",
                MxGraphModel: gantt.MxGraphModel{
                    Grid:     "1",
                    GridSize: "10",
                    Root: gantt.Root{
                        Cells: []gantt.MxCell{
                            {ID: "0"},
                            {ID: "1", Parent: "0"},
                            {
                                ID:     "task-1",
                                Value:  "Project Planning",
                                Style:  "strokeColor=#DEEDFF;fillColor=#ADC3D9",
                                Parent: "1",
                                Vertex: "1",
                                MxGeometry: &gantt.MxGeometry{
                                    X:      "100",
                                    Y:      "200",
                                    Width:  "200",
                                    Height: "20",
                                    As:     "geometry",
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    // Marshal to XML
    data, err := xml.MarshalIndent(mxFile, "", "  ")
    if err != nil {
        panic(err)
    }

    fmt.Println(string(data))
}
```

## Testing

Run the tests to verify the package works correctly:

```bash
go test -v
```

The tests include:
- Parsing the actual DrawIO file from `../../diagrams/gantt/default.drawio`
- Creating and marshaling new DrawIO structures
- Verifying round-trip XML parsing

## XML Structure

The DrawIO XML structure follows this hierarchy:

```
mxfile (root)
├── diagram
    └── mxGraphModel
        └── root
            └── mxCell (multiple)
                └── mxGeometry (optional)
```

Each `mxCell` can represent:
- Task names and descriptions
- Duration information
- Start and end dates
- Visual elements like bars and arrows
- Grouping containers

## Example Output

When parsing the included Gantt chart, you'll see tasks like:
- Complete project execution
- Engineering
- Project examination
- Material specification
- Workshop
- Field preparations and digging
- Testing and commissioning

Each task includes position, size, and styling information that can be used to reconstruct or modify the chart.
