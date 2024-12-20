package app

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

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

func closeFile(f *os.File) {
	// Unlock the file
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	// Close the file
	_ = f.Close()
}
