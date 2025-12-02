package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type ExerciseQuestion struct {
	id        int
	Exercise  int
	Question  string
	Errors    int
	DateAdded time.Time
	Comment   string
	Time      int
}

func GetExerciseQuestions(database *sqlx.DB) []ExerciseQuestion {
	data := GetExerciseQuestionsBySql(database, `select exercise, question, errors from english_exercise_questions order by date_added desc`)
	return data
}

// GetEnglishBookReadBySql Общая механика выборки из таблицы
func GetExerciseQuestionsBySql(database *sqlx.DB, sql string) []ExerciseQuestion {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	return scanExerciseQuestions(rows)
}

func scanExerciseQuestions(rows *sql.Rows) []ExerciseQuestion {
	var items []ExerciseQuestion
	for rows.Next() {
		p := ExerciseQuestion{}
		err := rows.Scan(&p.Exercise, &p.Question, &p.Errors)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}
