package app

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/mergestat/timediff"
)

type List struct {
	Description   string
	CreatedAt     time.Time
	FormattedTime string
	Done          bool
}

type Todo struct {
	List []List
}

func NewTodo() *Todo {
	return &Todo{}
}

func (t *Todo) LoadCSV() {
	// Open the file with read/write permissions and create it if it doesn't exist
	file, err := os.OpenFile("todos.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		// If file cannot be opened or created, initialize an empty List and return
		t.List = []List{}
		return
	}

	// Ensure the file is closed when the function exits
	defer closeFile(file)

	// Lock the file to prevent concurrent access
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		// If locking fails, initialize an empty List and return
		t.List = []List{}
		return
	}

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all rows from the CSV file
	rows, err := reader.ReadAll()
	if err != nil {
		// If error in reading, initialize an empty List and return
		t.List = []List{}
		return
	}

	// Parse rows into the List slice
	var loadedList []List
	for i, row := range rows {
		if i == 0 {
			// Skip header
			continue
		}

		if len(row) < 3 {
			// Skip rows with insufficient columns
			continue
		}

		description := row[0]
		createdAt, err := time.Parse("02/01/2006 15:04:05", row[1])
		if err != nil {
			// Return if failed to parse time
			log.Fatal("Error parsing time :", err.Error())
		}

		done, err := strconv.ParseBool(row[2])
		if err != nil {
			// return if failed to parse boolean
			log.Fatal("Error parsing boolean :", err.Error())
		}

		// Add parsed row to the loaded list
		loadedList = append(loadedList, List{
			Description:   description,
			CreatedAt:     createdAt,
			FormattedTime: timediff.TimeDiff(createdAt),
			Done:          done,
		})
	}

	// Assign the loaded list to the Todo struct
	t.List = loadedList
}

func (t *Todo) saveToCSV() {
	// Open the file with write permissions and truncate it if it exists
	file, err := os.OpenFile("todos.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		// If file cannot be opened or created, return
		log.Fatal("Failed to open file :", err.Error())
	}

	// Ensure the file is closed when the function exits
	defer closeFile(file)

	// Lock the file to prevent concurrent access
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		// If locking fails, return
		log.Fatal("Failed to lock file :", err.Error())
	}

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Check if the file is empty
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal("Failed to get file stat :", err.Error())
	}

	if fileInfo.Size() == 0 {
		// If the file is empty, write the header
		if err := writer.Write([]string{"Description", "Created At", "Done"}); err != nil {
			log.Fatal("Failed to write header :", err.Error())
		}
	}

	// Write each List item as a row in the CSV file
	for _, list := range t.List {
		// Construct row
		row := []string{
			list.Description,
			list.CreatedAt.Format("02/01/2006 15:04:05"),
			strconv.FormatBool(list.Done),
		}

		// Write to csv
		if err := writer.Write(row); err != nil {
			log.Fatal("Failed to write to csv file :", err.Error())
		}
	}

	writer.Flush()
}

func (t *Todo) ShowTodos() {
	// Load list from csv file
	t.LoadCSV()

	// Initiate simpletable
	table := simpletable.New()

	// Create header
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "DESCRIPTION"},
			{Align: simpletable.AlignCenter, Text: "CREATED AT"},
			{Align: simpletable.AlignCenter, Text: "STATUS"},
		},
	}

	for i, list := range t.List {
		// Identify todo status by 'Done' property
		status := "Not Done"
		if list.Done {
			status = "Done"
		}

		// Build row
		row := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(i + 1)},
			{Text: list.Description},
			{Text: list.FormattedTime},
			{Align: simpletable.AlignCenter, Text: status},
		}

		// Append row
		table.Body.Cells = append(table.Body.Cells, row)
	}

	// Set table stle
	table.SetStyle(simpletable.StyleCompactLite)

	// Print to screen
	fmt.Println(table.String())
}

func (t *Todo) CompleteTodo(index int) {
	// Return error if index out of range
	if index < 0 && index > len(t.List) {
		log.Fatal("Todo not found")
	}

	// Fetch list by index
	list := t.List[index]

	// Update todo status
	list.Done = true

	// Update todo list
	t.List[index] = list

	// Save to csv file
	t.saveToCSV()

	// Load and show updated todos
	t.ShowTodos()
}

func (t *Todo) AddTodo(description string) {
	// Create new list
	newList := List{
		Description: description,
		CreatedAt:   time.Now(),
		Done:        false,
	}

	// Append list
	t.List = append(t.List, newList)

	// Save updated list to csv file
	t.saveToCSV()

	// List all todos
	t.ShowTodos()
}

func (t *Todo) EditTodo(index int, description string) {
	// Return error if index out of range
	if index < 0 && index > len(t.List) {
		log.Fatal("Todo not found")
	}

	// Fetch list by index
	list := t.List[index]

	// Update todo
	list.Done = true
	list.Description = description

	// Update todo list
	t.List[index] = list

	// Save to csv file
	t.saveToCSV()

	// Load and show updated todos
	t.ShowTodos()
}

func (t *Todo) DeleteTodo(index int) {
	// Return error if index out of range
	if index < 0 && index > len(t.List) {
		log.Fatal("Todo not found")
	}

	// Delete list by index
	t.List = append(t.List[:index], t.List[index+1:]...)

	// Save to csv file
	t.saveToCSV()

	// Load and show updated todos
	t.ShowTodos()
}

func (t *Todo) ImportTodos(filePath string, mode string) {
	importedLists := loadImportedFile(filePath)
	count := make(map[string]List)

	// Find unique list from existing lists
	for _, existingList := range t.List {
		lowerItem := strings.ToLower(existingList.Description)
		if _, exists := count[lowerItem]; !exists {
			count[lowerItem] = existingList
		}
	}

	// Find unique list from imported lists
	for _, importedList := range importedLists {
		lowerItem := strings.ToLower(importedList.Description)
		if _, exists := count[lowerItem]; !exists {
			count[lowerItem] = importedList
		}
	}

	// Overwrite existing lists with unique lists
	var uniqueList []List
	for _, list := range count {
		uniqueList = append(uniqueList, list)
	}

	t.List = uniqueList

	// If mode is replace-all then replace all existing lists with imported lists
	if mode == "replace-all" {
		t.List = importedLists
	}

	// Save updated todos to csv file
	t.saveToCSV()

	t.ShowTodos()
}

func loadImportedFile(filePath string) []List {
	// Open imported file
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to open imported file :", err.Error())
	}

	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all rows from the CSV file
	rows, err := reader.ReadAll()

	// Parse rows into the List slice
	var loadedList []List
	for i, row := range rows {
		if i == 0 {
			// Skip header
			continue
		}

		if len(row) < 3 {
			// Skip rows with insufficient columns
			continue
		}

		description := row[0]
		createdAt, err := time.Parse("02/01/2006 15:04:05", row[1])
		if err != nil {
			// Return if failed to parse time
			log.Fatal("Error parsing time :", err.Error())
		}

		done, err := strconv.ParseBool(row[2])
		if err != nil {
			// return if failed to parse boolean
			log.Fatal("Error parsing boolean :", err.Error())
		}

		// Add parsed row to the loaded list
		loadedList = append(loadedList, List{
			Description:   description,
			CreatedAt:     createdAt,
			FormattedTime: timediff.TimeDiff(createdAt),
			Done:          done,
		})
	}

	return loadedList
}

func closeFile(f *os.File) {
	// Unlock the file
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	// Close the file
	_ = f.Close()
}
