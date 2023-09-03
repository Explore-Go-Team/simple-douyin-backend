package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"simple-douyin-backend/dal/db"
	"simple-douyin-backend/dal/db/dao"
)

var DSN = db.DSN

// Dynamic SQL
/*
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}
*/

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: filepath.Join(getProjectRoot(), "/dal/db/gorm_gen"),
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	var err error
	gormdb, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	g.UseDB(gormdb) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(dao.Message{}, dao.Relation{}, dao.UserDetail{})

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	//g.ApplyInterface(func(Querier) {}, model.User{})

	// Generate the code
	g.Execute()
}

func getProjectRoot() string {
	// 获取当前文件的路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("WARNING: fail to get project root, use work directory as default")
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("WARNING: fail to get work directory, use current directory as default")
			return "."
		}
		fmt.Println("Working directory: ", wd)
		return wd
	}
	// 获取当前文件的目录
	root := path.Dir(path.Dir(filename))
	fmt.Println("Project root: ", root)
	return root
}
