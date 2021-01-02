package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"readyGo/boxops"
	"readyGo/generate"
	"readyGo/mapping"
	"readyGo/scaler"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/lint"
)

var applyFile, projectType string
var lintFiles bool

func init() {

	applyCmd.Flags().StringVarP(&applyFile, "filename", "f", "", "user has to privide the file.There is no default file.")
	applyCmd.Flags().StringVarP(&projectType, "type", "t", "http_mongo", "type of the project can be http_mongo | http_sql_pg | grpc_mongo | grpc_sql_pg")
	applyCmd.Flags().BoolVarP(&lintFiles, "lint", "l", false, "lints all generated files and gives warnings and errors")

	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies a configuration file",
	Long:  `apply applies a configuration file for the project generation.User must supply a configuration file`,
	Run: func(cmd *cobra.Command, args []string) {

		ops := boxops.New("../box")
		mapping, err := mapping.New(ops, "configs/mappings/"+projectType+".json", projectType)
		if err != nil {
			log.Fatal(Fata(err))
		}

		scaler, err := scaler.New(ops, "configs/scalers.json")

		if err != nil {
			log.Fatal(Fata(err))
		}

		if applyFile == "default" {
			log.Fatal(Fata("apply must supply corrosponding configuration file"))
		}

		tg, err := generate.New(&applyFile, mapping, scaler)
		if err != nil {
			log.Fatal(Fata(err))
		}
		err = tg.CreateAll()
		if err != nil {
			tg.RmDir()
			log.Fatal(Fata(err))
		}
		if lintFiles {
			err := filepath.Walk("./"+tg.Project,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					extension := filepath.Ext(path)
					if !info.IsDir() && extension == ".go" {
						lintFile(path, 0.2)
					}
					return nil
				})

			if err != nil {
				log.Println(Warn(err))
			}

		}
		DisplayInfo(projectType)
	},
}

func lintFile(filename string, minConfidence float64) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	l := new(lint.Linter)
	ps, err := l.Lint(filename, src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v:%v\n", filename, err)
		return
	}
	for _, p := range ps {
		if p.Confidence >= minConfidence {
			fmt.Printf("%s:%v: %s\n", filename, p.Position, p.Text)
		}
	}
}

// DisplayInfo is to display information regarding project
func DisplayInfo(projectType string) {
	fmt.Println(Info("Attention please"))
	fmt.Println()
	if strings.Contains(projectType, "http") || strings.Contains(projectType, "grpc") {
		fmt.Println(Warn("port information"))
		fmt.Println(Details("readyGo does not know whether the port is available or not."))
		fmt.Println(Details("User has to make sure that the port is available and not behind firewall"))
	}
	fmt.Println()
	if strings.Contains(projectType, "grpc") {
		fmt.Println(Warn("grpc protocol buffer information"))
		fmt.Println(Details("readyGo does not generate proto buffer go files for you."))
		fmt.Println(Details("User has to make sure that protoc , proto_gen_go and protoc_gen_go_grpc tools are installed w.r.t the OS"))
	}
	fmt.Println()
	if strings.Contains(projectType, "mongo") {
		fmt.Println(Warn("mongo database information"))
		fmt.Println(Details("readyGo does not start the database."))
		fmt.Println(Details("Make sure your mongodb database is started , up and running"))
	}
	fmt.Println()
	if strings.Contains(projectType, "sql") {
		fmt.Println(Warn("sql database information"))
		fmt.Println(Details("readyGo does not start the database."))
		fmt.Println(Details("Make sure your sql database is started , up and running"))
	}
	fmt.Println()
}
