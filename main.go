package main

import (
	"fmt"
	"os"
	"todo/cli"
	h "todo/handlers"
	"todo/todo"
)

var todoTableColsNameType = map[string]string{
	"id":          "INTEGER PRIMARY KEY AUTOINCREMENT",
	"description": "TEXT",
	"completed":   "BOOLEAN",
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	todoFields := []string{}
	for k := range todoTableColsNameType {
		todoFields = append(todoFields, k)
	}
	tc := todo.NewTodoConfig()
	tc.TodoDir = fmt.Sprintf("%s/.todo", home)
	tc.TodoDb = "todos.db"
	tc.TodoTable = "todos"
	tc.TodoFields = todoFields
	tc.Init()
	db := todo.NewDB(tc.Source())
	defer db.Close()
	db.CreateTable("todos", todoTableColsNameType)

	err = RunCli()
	if err != nil {
		fmt.Println(err)
	}
}

func setCommands(cs *cli.Commands) {
	list := cli.Command{
		Name:    "list",
		Handler: h.HandleList,
		FlagTypes: []cli.FlagArg{
			{Name: "completed", ShortName: "c", Type: "bool", Required: false, Help: "Also show completed todos."},
		},
		Help: "Show completed todos.",
	}
	create := cli.Command{
		Name:    "create",
		Handler: h.HandleCreate,
		ArgTypes: []cli.PositonalArg{
			{Name: "description", Type: "string", Required: true, Help: "The description of the todo."},
		},
		Help: "Create a new todo.",
	}
	complete := cli.Command{
		Name:    "complete",
		Handler: h.HandleComplete,
		ArgTypes: []cli.PositonalArg{
			{Name: "id", Type: "int", Required: true, Help: "The id of the todo to complete."},
		},
		Help: "Complete a todo.",
	}
	delete := cli.Command{
		Name:    "delete",
		Handler: h.HandleDelete,
		ArgTypes: []cli.PositonalArg{
			{Name: "id", Type: "int", Required: true, Help: "The id of the todo to delete."},
		},
		Help: "Delete a todo.",
	}
	nukeDb := cli.Command{
		Name:    "nuke",
		Handler: h.HandleNuke,
		Help:    "Delete database.",
		FlagTypes: []cli.FlagArg{
			{Name: "force", ShortName: "f", Type: "bool", Required: false, Help: "Force delete the database."},
		},
	}
	cs.AddCommand(&list)
	cs.AddCommand(&create)
	cs.AddCommand(&complete)
	cs.AddCommand(&delete)
	cs.AddCommand(&nukeDb)
}

func RunCli() error {
	cs := cli.NewCommandSet("todo")
	setCommands(cs)
	cmdName, args, err := cs.ParseArgs()
	if err != nil {
		return err
	}
	return cs.Run(cmdName, args)

}
