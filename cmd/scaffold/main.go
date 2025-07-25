// go build -o ./scaffold.exe ./cmd/scaffold
// go run ./cmd/scaffold

package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func main() {
	var scaffoldType string
	survey.AskOne(&survey.Select{
		Message: "작업할 내용을 선택하세요.",
		Options: []string{
			"프로젝트 생성",
			"프로젝트 삭제",
			"프로젝트 배포",
			"프로젝트 배포 회수"},
	}, &scaffoldType)

	switch scaffoldType {
	case "프로젝트 생성":
		var projectName string
		var port string

		survey.AskOne(&survey.Input{
			Message: "프로젝트 명: ",
		}, &projectName, survey.WithValidator(survey.Required))

		survey.AskOne(&survey.Input{
			Message: "포트 번호: ",
		}, &port, survey.WithValidator(survey.Required))

		fmt.Println("프로젝트 생성")
		if _, err := os.Stat(fmt.Sprintf("projects/%s", projectName)); os.IsNotExist(err) {
			if err := os.Mkdir(fmt.Sprintf("projects/%s", projectName), 0755); err != nil {
				panic(err)
			}
		}

		fmt.Printf("project/%s 폴더 생성\n", projectName)

		if _, err := os.Stat(fmt.Sprintf("cmd/%s/main.go", projectName)); os.IsNotExist(err) {
			if err := os.Mkdir(fmt.Sprintf("cmd/%s", projectName), 0755); err != nil {
				panic(err)
			}
			f, err := os.Create(fmt.Sprintf("cmd/%s/main.go", projectName))
			if err != nil {
				panic(err)
			}
			defer f.Close()

			_, err = f.WriteString(`
package main

import ( 
	"fmt"
) 

func main() {
	fmt.Printf("Hello, World! from %s \n", "` + projectName + `")
}
`)
			if err != nil {
				panic(err)
			}
		}

		fmt.Printf("cmd/%s/main.go 파일 생성\n", projectName)

	case "프로젝트 삭제":

		fmt.Println("프로젝트 삭제")
	case "프로젝트 배포":

		fmt.Println("프로젝트 배포")
	case "프로젝트 배포 회수":

		fmt.Println("프로젝트 배포 회수")
	}
}
