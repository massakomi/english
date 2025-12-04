package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ExerciseQuestion struct {
	Id        int
	Exercise  int
	Question  string
	Errors    int
	DateAdded string
	Comment   sql.NullString
	Time      sql.NullInt64
}

func GetExerciseQuestions(database *sqlx.DB) []ExerciseQuestion {
	s := `select exercise, question, errors from english_exercise_questions order by date_added desc`
	data := GetExerciseQuestionsBySql(database, s, func(rows *sql.Rows, p *ExerciseQuestion) error {
		return rows.Scan(&p.Exercise, &p.Question, &p.Errors)
	})
	return data
}

func GetExerciseQuestionsByWhere(database *sqlx.DB, where string) []ExerciseQuestion {
	s := fmt.Sprintf(`select * from english_exercise_questions where %v order by date_added desc`, where)
	data := GetExerciseQuestionsBySql(database, s, func(rows *sql.Rows, p *ExerciseQuestion) error {
		return rows.Scan(&p.Id, &p.Exercise, &p.Question, &p.Errors, &p.DateAdded, &p.Comment, &p.Time)
	})
	return data
}

// GetEnglishBookReadBySql Общая механика выборки из таблицы
func GetExerciseQuestionsBySql(database *sqlx.DB, sql string, scanner func(rows *sql.Rows, p *ExerciseQuestion) error) []ExerciseQuestion {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []ExerciseQuestion
	for rows.Next() {
		p := ExerciseQuestion{}
		err := scanner(rows, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}
