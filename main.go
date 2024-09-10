package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Student struct {
	Name   string
	Grades map[string]float64
}

// main function drives the Student Grade System program.
// It provides options to manage student records, including adding, viewing, updating, and analyzing grades.
func main() {

	students := make(map[string]Student)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Student Grade System")
		fmt.Println("1. Add Student")
		fmt.Println("2. View All Students")
		fmt.Println("3. Search a specific student")
		fmt.Println("4. Update Student Grades")
		fmt.Println("5. Delete a student")
		fmt.Println("6. View Avg Grade")
		fmt.Println("7. Exit")
		fmt.Print("Choose an option: ")
		scanner.Scan()
		option := scanner.Text()
		switch option {
		case "1":
			addStudent(scanner, students)
		case "2":
			viewAllStudents(students)
		case "3":
			searchStudent(scanner, students)
		case "4":
			updateStudentGrades(scanner, students)
		case "5":
			deleteStudent(scanner, students)
		case "6":
			viewAvgGrade(scanner, students)
		case "7":
			fmt.Println("Exiting....")
			return
		default:
			fmt.Println("Invalid option. Please choose a valid option.")
		}
	}
}

// addStudent adds a new student to the system with their respective grades.
// Prompts the user for student details and grades, storing them in the map.
// Params: scan (*bufio.Scanner) - Scanner for input, students (map[string]Student) - map to store student data.
func addStudent(scan *bufio.Scanner, students map[string]Student) {
	fmt.Print("Enter student name: ")
	scan.Scan()
	name := strings.TrimSpace(scan.Text())
	subjects := make(map[string]float64)
	for {
		fmt.Print("Enter subject name or done: ")
		scan.Scan()
		subject := strings.TrimSpace(scan.Text())
		if subject == "done" {
			break
		}
		fmt.Print("Enter grade for ", subject, ": ")
		scan.Scan()
		grade, err := strconv.ParseFloat(strings.TrimSpace(scan.Text()), 64)
		if err != nil {
			fmt.Println("Invalid grade")
			continue
		}
		subjects[subject] = grade
	}

	students[name] = Student{Name: name, Grades: subjects}
	saveGradesToJson(students)
	fmt.Println("Student added successfully!")

}

// viewAllStudents displays all students and their grades.
// Iterates through the students map and prints each student's details.
// Params: students (map[string]Student) - map containing all student data.
func viewAllStudents(students map[string]Student) {
	for _, student := range students {
		fmt.Println("Name :", student.Name)
		for subject, grade := range student.Grades {
			fmt.Printf("  Subject: %s, Grade: %.2f\n", subject, grade)

		}
	}
}

// searchStudent finds and displays a specific student and their grades based on the provided name.
// Params: scan (*bufio.Scanner) - Scanner for input, students (map[string]Student) - map containing student data.
func searchStudent(scan *bufio.Scanner, students map[string]Student) {
	fmt.Print("Enter student name: ")
	scan.Scan()
	name := strings.TrimSpace(scan.Text())
	student, exists := students[name]
	if !exists {
		fmt.Println("Student Not found")
		return
	}
	fmt.Println("Name :", student.Name)
	for subject, grade := range student.Grades {
		fmt.Printf("  Subject: %s, Grade: %.2f\n", subject, grade)
	}
}

// updateStudentGrades updates the grades of an existing student.
// Allows the user to modify grades for specific subjects.
// Params: scan (*bufio.Scanner) - Scanner for input, students (map[string]Student) - map containing student data.
func updateStudentGrades(scan *bufio.Scanner, students map[string]Student) {
	fmt.Print("Enter student name: ")
	scan.Scan()
	name := strings.TrimSpace(scan.Text())
	student, exists := students[name]
	if !exists {
		fmt.Println("Student Not found")
		return
	}
	fmt.Println("Name :", student.Name)
	for {
		fmt.Print("Enter subject to update (or 'done' to finish): ")
		scan.Scan()
		subject := strings.TrimSpace(scan.Text())
		if subject == "done" {
			break
		}
		fmt.Print("Enter new grade: ")
		scan.Scan()
		gradeStr := strings.TrimSpace(scan.Text())
		grade, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			fmt.Println("Invalid grade")
			return
		}
		student.Grades[subject] = grade
		students[name] = student
		fmt.Println("Grade updated")

	}
	saveGradesToJson(students)
}

// deleteStudent removes a student from the system.
// Prompts for the student name and deletes the corresponding record if found.
// Params: scan (*bufio.Scanner) - Scanner for input, students (map[string]Student) - map containing student data.
func deleteStudent(scan *bufio.Scanner, students map[string]Student) {
	fmt.Print("Enter student name: ")
	scan.Scan()
	name := strings.TrimSpace(scan.Text())
	_, exists := students[name]
	if !exists {
		fmt.Println("Student Not found")
		return
	}
	delete(students, name)
	fmt.Println("Student deleted")
	saveGradesToJson(students)
}

// viewAvgGrade calculates and displays the average, highest, lowest, and median grades of a student.
// Provides an analysis of the student's grades across subjects.
// Params: scan (*bufio.Scanner) - Scanner for input, students (map[string]Student) - map containing student data.
func viewAvgGrade(scan *bufio.Scanner, students map[string]Student) {
	fmt.Print("Enter student name: ")
	scan.Scan()
	name := strings.TrimSpace(scan.Text())
	student, exists := students[name]
	if !exists {
		fmt.Println("Student Not found")
		return
	}
	var sum, highest, lowest, median float64
	var grades []float64
	var highestSubject, lowestSubject string

	for subject, grade := range student.Grades {
		sum += grade
		grades = append(grades, grade)
		if grade > highest || highest == 0 {
			highest = grade
			highestSubject = subject
		}
		if grade < lowest || lowest == 0 {
			lowest = grade
			lowestSubject = subject
		}
	}
	median = calculateMedian(grades)

	fmt.Println("avg grade of student is: ", sum/float64(len(student.Grades)))
	fmt.Printf("Highest Grade: %.2f in %s\n", highest, highestSubject)
	fmt.Printf("Lowest Grade: %.2f in %s\n", lowest, lowestSubject)
	fmt.Println("Median Gradent of student is: ", median)

}

// calculateMedian computes the median of a slice of grades.
// Sorts the grades and returns the middle value or the average of two middle values for even slices.
// Params: grades ([]float64) - slice of grades.
// Returns: float64 - median value of the grades.
func calculateMedian(grades []float64) float64 {
	n := len(grades)
	if n == 0 {
		return 0
	}
	sort.Float64s(grades)
	// If the number of grades is odd, return the middle element
	if n%2 != 0 {
		return grades[n/2]
	}
	return (grades[(n-1)/2] + grades[n/2]) / 2.0
}

// saveGradesToFile saves the student grades to a JSON file.
// Params: students (map[string]Student) - map containing all student data.
func saveGradesToJson(students map[string]Student) {
	jsonData, err := json.MarshalIndent(students, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	err = os.WriteFile("grades.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
	fmt.Println("Grades saved to grades.json successfully!")
}
