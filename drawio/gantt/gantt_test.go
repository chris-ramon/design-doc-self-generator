package gantt

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParseDrawIOFile(t *testing.T) {
	// Get the path to the DrawIO file using runtime.Caller
	_, filename, _, _ := runtime.Caller(0)
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	drawioPath := filepath.Join(repoRoot, "diagrams", "gantt", "template", "basic.drawio")
	
	// Read the test DrawIO file
	data, err := os.ReadFile(drawioPath)
	if err != nil {
		t.Fatalf("Failed to read DrawIO file: %v", err)
	}

	// Parse the XML into our structs
	var mxFile MxFile
	err = xml.Unmarshal(data, &mxFile)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Verify basic structure
	if mxFile.Host == "" {
		t.Error("Expected host attribute to be set")
	}

	if mxFile.Version == "" {
		t.Error("Expected version attribute to be set")
	}

	if len(mxFile.Diagrams) == 0 {
		t.Error("Expected at least one diagram")
	}

	diagram := mxFile.Diagrams[0]
	if diagram.Name == "" {
		t.Error("Expected diagram name to be set")
	}

	if diagram.ID == "" {
		t.Error("Expected diagram ID to be set")
	}

	if len(diagram.MxGraphModel.Root.Cells) == 0 {
		t.Error("Expected at least one cell in the diagram")
	}

	// Verify we can find some expected cells
	foundCells := 0
	for _, cell := range diagram.MxGraphModel.Root.Cells {
		if cell.Value != "" {
			foundCells++
		}
	}

	if foundCells == 0 {
		t.Error("Expected to find cells with values")
	}

	t.Logf("Successfully parsed DrawIO file with %d cells", len(diagram.MxGraphModel.Root.Cells))
	t.Logf("Found %d cells with values", foundCells)
}

func TestMarshalDrawIOFile(t *testing.T) {
	// Create a simple test structure
	mxFile := MxFile{
		Host:    "test",
		Agent:   "test-agent",
		Version: "1.0",
		Diagrams: []Diagram{
			{
				Name: "Test-Page",
				ID:   "test-id",
				MxGraphModel: MxGraphModel{
					Grid:     "1",
					GridSize: "10",
					Root: Root{
						Cells: []MxCell{
							{
								ID: "0",
							},
							{
								ID:     "1",
								Parent: "0",
							},
							{
								ID:     "test-cell",
								Value:  "Test Task",
								Style:  "strokeColor=#DEEDFF;fillColor=#ADC3D9",
								Parent: "1",
								Vertex: "1",
								MxGeometry: &MxGeometry{
									X:      "100",
									Y:      "200",
									Width:  "120",
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
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	// Verify we can unmarshal it back
	var parsed MxFile
	err = xml.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal generated XML: %v", err)
	}

	// Basic verification
	if parsed.Host != mxFile.Host {
		t.Errorf("Expected host %s, got %s", mxFile.Host, parsed.Host)
	}

	if len(parsed.Diagrams) != 1 {
		t.Errorf("Expected 1 diagram, got %d", len(parsed.Diagrams))
	}

	if len(parsed.Diagrams[0].MxGraphModel.Root.Cells) != 3 {
		t.Errorf("Expected 3 cells, got %d", len(parsed.Diagrams[0].MxGraphModel.Root.Cells))
	}

	t.Logf("Successfully marshaled and unmarshaled test structure")
}
