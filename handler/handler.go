package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"test-crud-api/model"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "apitest"
	password = "mypass"
	dbname   = "mydb"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

var student model.Student
var students []model.Student

func createConn() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	err = db.Ping()
	CheckError(err)

	fmt.Println("Successfully connected to DB")

	return db
}

func createStudent(student model.Student) int64 {
	db := createConn()
	defer db.Close()

	var id int64

	insertQuery := `INSERT INTO students (name, surname, university, major) VALUES($1, $2, $3, $4) RETURNING id`
	err := db.QueryRow(insertQuery, student.Name, student.Surname, student.University, student.Major).Scan(&id)
	CheckError(err)

	return id
}

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewDecoder(r.Body).Decode(&student)
	CheckError(err)

	insertedID := createStudent(student)

	res := response{
		ID:      insertedID,
		Message: "Student created successfully",
	}
	json.NewEncoder(w).Encode(res)
}

func getStudent(id int64) (model.Student, error) {
	db := createConn()
	defer db.Close()

	selectQuery := `SELECT * FROM students WHERE id=$1`
	err := db.QueryRow(selectQuery, id).Scan(&student.ID, &student.Name, &student.Surname, &student.University, &student.Major)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows returned")
		return student, nil
	case nil:
		return student, nil
	}

	return student, err
}

func GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	val := mux.Vars(r)

	id, err := strconv.Atoi(val["id"])
	CheckError(err)

	student, err = getStudent(int64(id))
	CheckError(err)

	json.NewEncoder(w).Encode(student)
}

func getAllStudents() ([]model.Student, error) {
	db := createConn()
	defer db.Close()

	selectAllQuery := `SELECT * FROM students`

	rows, err := db.Query(selectAllQuery)
	CheckError(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&student.ID, &student.Name, &student.Surname, &student.University, &student.Major)
		CheckError(err)
		students = append(students, student)
	}
	return students, err
}

func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	students, err := getAllStudents()
	CheckError(err)

	json.NewEncoder(w).Encode(students)
}

func updateStudent(id int64) int64 {

	db := createConn()
	defer db.Close()

	updateQuery := `UPDATE students SET name=$2, surname=$3, university=$4, major=$5 WHERE id=$1`

	result, err := db.Exec(updateQuery, id, student.Name, student.Surname, student.University, student.Major)
	CheckError(err)

	updatedRow, err := result.RowsAffected()
	CheckError(err)

	fmt.Printf("Updated rows: %v", updatedRow)

	return updatedRow
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	val := mux.Vars(r)
	id, err := strconv.Atoi(val["id"])
	CheckError(err)

	err = json.NewDecoder(r.Body).Decode(&student)
	CheckError(err)

	updatedRows := updateStudent(int64(id))

	msg := fmt.Sprintf("The number of updated rows: %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func deleteStudent(id int64) int64 {

	db := createConn()
	defer db.Close()

	deleteQuery := `DELETE FROM students WHERE id=$1`

	result, err := db.Exec(deleteQuery, id)
	CheckError(err)

	deletedRow, e := result.RowsAffected()
	CheckError(e)

	fmt.Printf("Row deleted: %v", deletedRow)

	return deletedRow
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	val := mux.Vars(r)

	id, err := strconv.Atoi(val["id"])
	CheckError(err)

	deletedRow := deleteStudent(int64(id))

	msg := fmt.Sprintf("Deleted row: %v", deletedRow)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}
