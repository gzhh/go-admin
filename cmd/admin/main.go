package main

import (
	"go-admin/cmd/admin/app"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app.Execute()
}
