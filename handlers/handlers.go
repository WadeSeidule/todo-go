package handlers

import (
	"fmt"
	"os"
	"todo/cli"
	"todo/todo"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Table struct {
	Fields        []string
	Rows          []table.Row
	ExcludeFields []string
	writer        table.Writer
}

func (t *Table) create() {
	t.writer = table.NewWriter()
	t.writer.SetOutputMirror(os.Stdout)
	t.setHeader()
	for _, row := range t.Rows {
		t.writer.AppendRow(row)
	}
}

func (t *Table) setHeader() {
	header := table.Row{}
	for _, field := range t.Fields {
		if containsString(t.ExcludeFields, field) {
			continue
		}
		header = append(header, field)
	}
	t.writer.AppendHeader(header)
}

func (t *Table) Render() {
	t.create()
	t.writer.Render()
}

func HandleList(c *cli.Command) error {
	return list(c.GetBoolFlag("completed"))
}

func list(showCompleted bool) error {
	todos, err := todo.GetTodos()
	if err != nil {
		fmt.Println("Error getting todos:", err)
		return err
	}

	exclude := []string{}
	if !showCompleted {
		exclude = append(exclude, "completed")
	}

	rows := []table.Row{}
	for _, todo := range todos {
		if !todo.Completed && !showCompleted {
			rows = append(rows, table.Row{todo.Id, todo.Description})
		} else if showCompleted {
			rows = append(rows, table.Row{todo.Id, todo.Description, todo.Completed})
		}
	}
	table := Table{
		Fields:        todo.GetTodoFields(),
		Rows:          rows,
		ExcludeFields: exclude,
	}
	table.Render()
	return nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func HandleCreate(c *cli.Command) error {
	description := c.GetStringArg("description")
	t := todo.Todo{Description: description}
	err := t.Save()
	if err != nil {
		return err
	}
	return list(false)
}

func HandleComplete(c *cli.Command) error {
	id := c.GetIntArg("id")
	t := todo.Todo{Id: id}
	err := t.Get()
	if err != nil {
		return err
	}
	err = t.Complete()
	if err != nil {
		return err
	}
	return list(false)
}

func HandleDelete(c *cli.Command) error {
	id := c.GetIntArg("id")
	t := todo.Todo{Id: id}
	err := t.Delete()
	if err != nil {
		return err
	}
	return list(false)
}

func HandleNuke(c *cli.Command) error{
	force := c.GetBoolFlag("force")
	todo.DeleteDB(force)
	return nil

}