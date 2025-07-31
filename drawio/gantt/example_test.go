package gantt

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// ExampleParseGanttChart demonstrates how to parse a DrawIO Gantt chart
// and extract task information
func ExampleParseGanttChart() {
	// Read the DrawIO file
	data, err := os.ReadFile("../../diagrams/gantt/default.drawio")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Parse the XML
	var mxFile MxFile
	err = xml.Unmarshal(data, &mxFile)
	if err != nil {
		fmt.Printf("Error parsing XML: %v\n", err)
		return
	}

	// Extract task information from cells
	fmt.Println("Gantt Chart Tasks:")
	fmt.Println("==================")

	diagram := mxFile.Diagrams[0]
	for _, cell := range diagram.MxGraphModel.Root.Cells {
		// Look for cells that appear to be task names (have meaningful values and geometry)
		if cell.Value != "" && cell.MxGeometry != nil && isTaskCell(cell) {
			fmt.Printf("Task: %s\n", cell.Value)
		}
	}

	// Output:
	// Gantt Chart Tasks:
	// ==================
	// Task: Complete project execution
	// Task: Engineering
	// Task: Project examination
	// Task: Material specification
	// Task: Material ordering
	// Task: Equipment layouting
	// Task: Supervision and meetings
	// Task: Bill of works
	// Task: Workshop
	// Task: Project examination and material comparison
	// Task: Preparing distribution boards
	// Task: Mounting equipment
	// Task: Wiring
	// Task: Testing
	// Task: Packaging
	// Task: Field
	// Task: Field preparations and digging
	// Task: Cable laying
	// Task: Installation laying
	// Task: Mount distribution boards
	// Task: Wiring distribution boards
	// Task: Testing
}

// isTaskCell determines if a cell represents a task (not just a number or date)
func isTaskCell(cell MxCell) bool {
	value := strings.TrimSpace(cell.Value)
	
	// Skip empty values
	if value == "" {
		return false
	}
	
	// Skip single characters (calendar headers like M, T, W, etc.)
	if len(value) == 1 {
		return false
	}
	
	// Skip cells that are just numbers (task IDs)
	if len(value) <= 2 && isNumeric(value) {
		return false
	}
	
	// Skip cells that look like dates (contain dots and numbers, or month/day patterns)
	if strings.Contains(value, ".") && len(value) <= 10 {
		return false
	}
	
	// Skip date patterns like "7 May 12", "14 May 12", etc.
	if strings.Contains(strings.ToLower(value), "may") || 
	   strings.Contains(strings.ToLower(value), "apr") ||
	   strings.Contains(strings.ToLower(value), "jun") {
		return false
	}
	
	// Skip cells that are just duration (contain "day" or "days")
	if strings.Contains(strings.ToLower(value), "day") && len(value) <= 10 {
		return false
	}
	
	// Skip header cells
	if value == "Task Name" || value == "Duration" || value == "Start" || value == "Finish" {
		return false
	}
	
	// Skip date headers like "16 Apr 12", "23 Apr 12", etc.
	if strings.Contains(strings.ToLower(value), "apr") || 
	   strings.Contains(strings.ToLower(value), "may") || 
	   strings.Contains(strings.ToLower(value), "jun") {
		return false
	}
	
	// Only consider cells with meaningful length (likely task names)
	if len(value) < 4 {
		return false
	}
	
	return true
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}
