package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"canoe/internal/model"

	_ "github.com/lib/pq"
)

func ConnectPostgresDB() (*sql.DB, error) {
	connstring := "user=postgres dbname=postgres password='postgres' host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connstring)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func QueryQuestionCount(db *sql.DB) (int, error) {
	rows, err := db.Query("select count(*) from questions")
	if err != nil {
		return 0, err
	}

	var res = 0
	for rows.Next() {
		rows.Scan(&res)
	}
	return res, nil
}

func QueryWordsCount(db *sql.DB) (int, error) {
	rows, err := db.Query("select count(*) from parole")
	if err != nil {
		return 0, err
	}

	var res = 0
	for rows.Next() {
		rows.Scan(&res)
	}
	return res, nil
}

func QueryQuestionBySerial(db *sql.DB, id int) (model.Question, error) {
	var question model.Question
	var rawOption string
	sqlStr := fmt.Sprintf("select id, created_at, question_text, question_type, options, answer, class, serial from questions where serial = %d", id)
	err := db.QueryRow(sqlStr).Scan(&question.Id, &question.CreatedAt, &question.QuestionText, &question.QuestionType, &rawOption, &question.Answer, &question.Class, &question.Serial)
	if err != nil {
		return model.Question{}, err
	}

	var opt model.Option
	err = json.Unmarshal([]byte(rawOption), &opt)
	if err != nil {
		return model.Question{}, err
	}
	question.Options = opt
	return question, nil
}

func AddWord(db *sql.DB, req model.Word) error {
	var sqlStr = fmt.Sprintf("insert into parole (word, meaning) values ('%s', '%s')", req.Word, req.Meaning)
	_, err := db.Exec(sqlStr)
	return err
}

func QueryWordsFromDatabase(db *sql.DB, page, size int) ([]model.Word, error) {
	rows, err := db.Query("select word, meaning from parole")
	if err != nil {
		return []model.Word{}, err
	}

	var res []model.Word
	for rows.Next() {
		var w, m string
		rows.Scan(&w, &m)
		res = append(res, model.Word{Word: w, Meaning: m})
	}
	return res, nil
}