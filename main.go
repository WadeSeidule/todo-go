package main

import (
	"flag"
	"fmt"
	"os"
	"todo/cli"
	h "todo/handlers"
	"todo/todo"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	tc := todo.NewTodoConfig()
	tc.TodoDir = fmt.Sprintf("%s/.todo", home)
	tc.TodoDb = "todos.db"
	tc.TodoTable = "todos"
	tc.TodoFields = []string{"id", "description", "completed"}
	tc.Init()
	db := todo.NewDB(tc.Source())
	defer db.Close()
	db.CreateTable("todos", map[string]string{
		"id":          "INTEGER PRIMARY KEY AUTOINCREMENT",
		"description": "TEXT",
		"completed":   "BOOLEAN",
	})

	err = RunCli()
	if err != nil {
		fmt.Println(err)
	}
}


func getCommands() cli.Commands {
	list := cli.Command{
		Name:    "list",
		Handler: h.HandleList,
		FlagTypes: []cli.FlagArg{
			{Name: "completed", ShortName: "c", Type: "bool", Required: false, Help: "Also show completed todos."},
		},
		Help:    "Show completed todos. 'todos list {--completed, -c}'",
		FlagSet: flag.NewFlagSet("list", flag.ExitOnError),
	}
	create := cli.Command{
		Name:    "create",
		Handler: h.HandleCreate,
		ArgTypes: []cli.PositonalArg{
			{Name: "description", Type: "string", Required: true, Help: "The description of the todo."},
		},
		Help:    "Create a new todo. 'todos create [description]'",
		FlagSet: flag.NewFlagSet("create", flag.ExitOnError),
	}
	complete := cli.Command{
		Name:    "complete",
		Handler: h.HandleComplete,
		ArgTypes: []cli.PositonalArg{
			{Name: "id", Type: "int", Required: true, Help: "The id of the todo to complete."},
		},
		Help:    "Complete a todo. 'todos complete [id]'",
		FlagSet: flag.NewFlagSet("complete", flag.ExitOnError),
	}
	delete := cli.Command{
		Name:    "delete",
		Handler: h.HandleDelete,
		ArgTypes: []cli.PositonalArg{
			{Name: "id", Type: "int", Required: true, Help: "The id of the todo to delete."},
		},
		Help:    "Delete a todo. 'todos delete [id]'",
		FlagSet: flag.NewFlagSet("delete", flag.ExitOnError),
	}
	nukeDb := cli.Command{
		Name:    "nuke",
		Handler: h.HandleNuke,
		Help:    "Delete database. 'todos nuke'",
		FlagSet: flag.NewFlagSet("nuke", flag.ExitOnError),
		FlagTypes: []cli.FlagArg{
			{Name: "force", ShortName: "f", Type: "bool", Required: false, Help: "Force delete the database."},
		},
	}
	cs := cli.Commands{}
	cs.AddCommand(&list)
	cs.AddCommand(&create)
	cs.AddCommand(&complete)
	cs.AddCommand(&delete)
	cs.AddCommand(&nukeDb)
	return cs
}

func RunCli() error {
	cs := getCommands()
	cmdName, args, err := cs.ParseArgs()
	if err != nil {
		return err
	}
	return cs.Run(cmdName, args)

}
