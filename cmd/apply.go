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

	"github.com/spf13/cobra"
	"golang.org/x/lint"
)

var applyFile, projectType string
var lintFiles bool

func init() {

	applyCmd.Flags().StringVarP(&applyFile, "filename", "f", "default", "config file to generate the project")
	applyCmd.Flags().StringVarP(&projectType, "type", "t", "http_mongo", "type of the project http_mongo | http_sql_pg | grpc_mogo | grpc_sql | cloudevent")
	applyCmd.Flags().BoolVarP(&lintFiles, "lint", "l", false, "lints all generated files and gives warnings and errors")

	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies a configuration file",
	Long:  `apply applies a configuration file for the project generation.User must supply a configuration file`,
	Run: func(cmd *cobra.Command, args []string) {

		ops := boxops.New("../box")
		mapping, err := mapping.New(ops, "configs/mappings.json", projectType)
		if err != nil {
			log.Fatal(err)
		}

		if applyFile == "default" {
			log.Fatal("apply must supply corrosponding configuration file")
		}

		tg, err := generate.New(&applyFile, mapping)
		if err != nil {
			log.Fatal(err)
		}
		err = tg.CreateAll()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(lintFiles)
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
				log.Println(err)
			}

		}
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
