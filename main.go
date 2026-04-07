package main

import (
	"github.com/tss182/migration/cmd"
	_ "github.com/tss182/migration/database/migrations"
	_ "github.com/tss182/migration/database/seeders"
)

func main() {
	cmd.Execute()
}
