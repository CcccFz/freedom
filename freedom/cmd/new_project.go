package cmd

import (
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/8treenet/freedom/freedom/template/project"

	"github.com/spf13/cobra"
)

var (
	NewProjectCmd = &cobra.Command{
		Use:   "new-project [project_name]",
		Short: "New a microservice project based on freedom",
		Long:  `New project from freedom project template. `,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			sysPath, err := filepath.Abs(args[0])
			if err != nil {
				return
			}

			projectPath := strings.Replace(sysPath, build.Default.GOPATH+"/src/", "", 1)
			projectName := args[0]
			pdata := map[string]interface{}{
				"PackagePath": projectPath,
				"PackageName": projectName,
			}

			mkdirAll(sysPath)
			m := project.FileContent()
			for filepath, content := range m {
				var pf *os.File
				pf, err = os.Create(sysPath + filepath)
				if err != nil {
					return err
				}
				tmpl, err := template.New(projectName).Parse(content)
				if err = tmpl.Execute(pf, pdata); err != nil {
					return err
				}
			}
			exec.Command("gofmt", "-w", sysPath).Output()
			return nil
		},
	}
)

func init() {
	AddCommand(NewProjectCmd)
}

func mkdirAll(projectPath string) {
	os.MkdirAll(projectPath+"/cmd", os.ModePerm)
	os.MkdirAll(projectPath+"/cmd/conf", os.ModePerm)
	os.MkdirAll(projectPath+"/business/controllers", os.ModePerm)
	os.MkdirAll(projectPath+"/models", os.ModePerm)
	os.MkdirAll(projectPath+"/components/config", os.ModePerm)
	os.MkdirAll(projectPath+"/business/repositorys", os.ModePerm)
	os.MkdirAll(projectPath+"/business/services", os.ModePerm)
}