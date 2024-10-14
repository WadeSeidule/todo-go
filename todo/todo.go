package todo

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var tc *TodoConfig

type TodoConfig struct {
	TodoDir    string
	TodoDb     string
	TodoTable  string
	TodoFields []string
}

func (t *TodoConfig) createTodoDir() {
	_, err := os.Stat(t.TodoDir)
	exists := os.IsExist(err)
	if !exists {
		err := os.MkdirAll(t.TodoDir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func (t *TodoConfig) Init() {
	t.createTodoDir()
}

func (t *TodoConfig) Source() string {
	return fmt.Sprintf("%s/%s", t.TodoDir, t.TodoDb)
}

func NewTodoConfig() *TodoConfig {
	tc = &TodoConfig{}
	return tc
}

func GetTodoFields() []string {
	return tc.TodoFields
}

func NewTodo(desc string) *Todo {
	return &Todo{Id: 0, Description: desc, Completed: false}
}

type Todo struct {
	Id          int    `db:"id"`
	Description string `db:"description"`
	Completed   bool   `db:"completed"`
}

func (t *Todo) doesTodoExist() bool {
	if t.Id == 0 {
		return false
	}
	var count int
	db.engine.Get(&count, "SELECT COUNT(*) FROM %s WHERE id=?", tc.TodoTable, t.Id)
	return count > 0
}

func (t *Todo) Insert() error {
	insert_q := "INSERT INTO todos(description, completed) VALUES (:description, :completed)"
	_, err := db.engine.NamedExec(insert_q, t)
	return err
}

func (t *Todo) Update() error {
	update_q := "UPDATE todos SET description=:description, completed=:completed WHERE id=:id"
	_, err := db.engine.NamedExec(update_q, t)
	return err
}

func (t *Todo) Save() error {
	if t.doesTodoExist() {
		return t.Update()
	}
	return t.Insert()
}

func (t *Todo) Delete() error {
	q := "DELETE FROM todos WHERE id=?"
	_, err := db.engine.Exec(q, t.Id)
	return err
}

func (t *Todo) Complete() error {
	t.Completed = true
	return t.Save()
}

func (t *Todo) Get() error {
	q := "SELECT * FROM todos WHERE id=?"
	return db.engine.Get(t, q, t.Id)
}

func GetTodos() ([]Todo, error) {
	q := "SELECT * FROM todos ORDER BY id"
	todos := []Todo{}
	err := db.engine.Select(&todos, q)
	return todos, err

}
