package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ContribType int

const (
	ACTION ContribType = 1 + iota
	TRIGGER
	ACTIVITY
	FUNCTION
)

var dir = flag.String("dir", "", "Give dir to geneate all legacy metadata file under this dir")

func main() {
	flag.Parse()

	if FileExist(*dir) {
		err := GenerateLegacyMetdata(*dir)
		if err != nil {
			fmt.Println("generate metadata file failed, due to ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Generate Legacy Metadata successfully")
	} else {
		fmt.Println("Dir ", *dir, " not exist!")
		os.Exit(1)
	}
}

// ListDependencies lists all installed dependencies
func GenerateLegacyMetdata(dir string) error {
	err := filepath.Walk(dir, func(filePath string, info os.FileInfo, _ error) error {

		if !info.IsDir() {

			switch info.Name() {
			case "trigger.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0 : len(filePath)-12]
				if hasModelFile(dir) {
					//new model
					return nil
				}
				pkg := filepath.Base(dir)
				if _, err := os.Stat(fmt.Sprintf("%s/../trigger.json", dir)); err == nil {
					//old trigger.json, ignore
					return nil
				}

				desc, err := readDescriptor(filePath, info)
				if err == nil && desc.Type == "flogo:trigger" {
					raw, err := ioutil.ReadFile(filePath)
					if err != nil {
						return err
					}
					err = generateMetdataFile(pkg, string(raw), filepath.Join(dir, "trigger_metadata.go"), tplTriggerMetadataGoFile)
					if err != nil {
						return err
					}
				}
			case "activity.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0 : len(filePath)-13]
				if hasModelFile(dir) {
					//new model
					return nil
				}
				if _, err := os.Stat(fmt.Sprintf("%s/../activity.json", dir)); err == nil {
					//old activity.json, ignore
					return nil
				}
				pkg := filepath.Base(dir)
				desc, err := readDescriptor(filePath, info)
				if err == nil && desc.Type == "flogo:activity" {
					raw, err := ioutil.ReadFile(filePath)
					if err != nil {
						return err
					}
					err = generateMetdataFile(pkg, string(raw), filepath.Join(dir, "activity_metadata.go"), tplActivityMetadataGoFile)
					if err != nil {
						return err
					}
				}
			}

		}

		return nil
	})

	return err
}

func hasModelFile(dir string) bool {
	return FileExist(filepath.Join(dir, "model_1.1.0.txt"))
}

func generateMetdataFile(pkg, jsonContent, metdataGoFilePath, tplMetdata string) error {
	info := &struct {
		Package      string
		MetadataJSON string
	}{
		Package:      pkg,
		MetadataJSON: jsonContent,
	}

	f, err := os.Create(metdataGoFilePath)
	if err != nil {
		return err
	}
	RenderMetdataTemplate(f, tplMetdata, info)
	return nil
}

func readDescriptor(filepath string, info os.FileInfo) (*AppDescriptor, error) {

	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return nil, err
	}

	return ParseDescriptor(string(raw))
}

var tplActivityMetadataGoFile = `package {{.Package}}

import (
	"github.com/project-flogo/legacybridge"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	legacybridge.RegisterLegacyActivity(NewActivity(md))
}
`

var tplTriggerMetadataGoFile = `package {{.Package}}

import (
	"github.com/project-flogo/legacybridge"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register trigger factory
func init() {
	md := trigger.NewMetadata(jsonMetadata)
	legacybridge.RegisterLegacyTriggerFactory(md.ID, NewFactory(md))
}
`

//RenderTemplate renders the specified template
func RenderMetdataTemplate(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

type Descriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type AppDescriptor struct {
	Descriptor
	Ref     string `json:"ref"`
	Display struct {
		Category string `json:"category"`
	} `json:"display"`
}

// ParseDescriptor parse a descriptor
func ParseDescriptor(descJson string) (*AppDescriptor, error) {
	descriptor := &AppDescriptor{}

	err := json.Unmarshal([]byte(descJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

type WIRefs struct {
	Type         ContribType
	Ref          string
	Path         string
	CategoryName string
}

func (wiref *WIRefs) String() string {
	return "[" + wiref.CategoryName + "]" + "===>" + wiref.Ref
}

// ListDependencies lists all installed dependencies
func ListDependencies(dir string, cType ContribType) ([]*WIRefs, error) {
	var deps []*WIRefs

	err := filepath.Walk(dir, func(filePath string, info os.FileInfo, _ error) error {

		if !info.IsDir() {

			switch info.Name() {
			case "action.json":
				if cType == 0 || cType == ACTION {
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:action" {
						deps = append(deps, &WIRefs{Type: ACTION, Ref: desc.Ref, CategoryName: desc.Display.Category, Path: filePath})
					}
				}
			case "trigger.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0 : len(filePath)-12]
				if _, err := os.Stat(fmt.Sprintf("%s/../trigger.json", dir)); err == nil {
					//old trigger.json, ignore
					return nil
				}
				if cType == 0 || cType == TRIGGER {
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:trigger" {
						deps = append(deps, &WIRefs{Type: TRIGGER, Ref: desc.Ref, CategoryName: desc.Display.Category, Path: filePath})
					}
				}
			case "activity.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0 : len(filePath)-13]
				if _, err := os.Stat(fmt.Sprintf("%s/../activity.json", dir)); err == nil {
					//old activity.json, ignore
					return nil
				}
				if cType == 0 || cType == ACTIVITY {
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:activity" {
						deps = append(deps, &WIRefs{Type: ACTIVITY, Ref: desc.Ref, CategoryName: desc.Display.Category, Path: filePath})
					}
				}

			case "descriptor.json":
				//Maybe function/trigger/activity etc...
				desc, err := readDescriptor(filePath, info)
				if err != nil {
					return err
				}
				//So descriptor.json for function only
				if desc.Type == "flogo:activity" {

				} else if desc.Type == "flogo:function" {
					ref, err := GetRefFromModFile(filepath.Dir(filePath))
					if err != nil {
						return err
					}
					deps = append(deps, &WIRefs{Type: FUNCTION, Ref: ref, CategoryName: desc.Name, Path: filePath})

				} else if desc.Type == "flogo:trigger" {

				}
			}

		}

		return nil
	})

	return deps, err
}

func GetRefFromModFile(moddir string) (string, error) {
	if FileExist(filepath.Join(moddir, "go.mod")) {
		inFile, _ := os.Open(filepath.Join(moddir, "go.mod"))
		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "module") {
				return strings.TrimSpace(line[6:]), nil
			}
		}
	}
	return "", fmt.Errorf("No go.mod file found from path [%s], go.mod file must present for function contribution", moddir)
}

func FileExist(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
}
