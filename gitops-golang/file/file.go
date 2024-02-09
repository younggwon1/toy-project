package file

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadFromFile(filename string, out interface{}) error {
	// open file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// unmarshal data into cfg
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}

func ModifyFromYamlFile(helmRepo, yamlFile, valueToModify string) error {
	yamlFileName := "/tmp/" + helmRepo + "/" + yamlFile

	var parseYamlNode yaml.Node
	err := ReadFromFile(yamlFileName, &parseYamlNode)
	if err != nil {
		return err
	}

	var valueTemplate map[string]string
	err = yaml.Unmarshal([]byte(valueToModify), &valueTemplate)
	if err != nil {
		return err
	}

	var keys []string
	var values string
	for k, v := range valueTemplate {
		keys = strings.Split(k, ".")
		values = v
	}

	modifyNode(&parseYamlNode, keys, values)

	file, err := os.Create(yamlFileName)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	if err := encoder.Encode(&parseYamlNode); err != nil {
		return err
	}

	return nil
}

func modifyNode(node *yaml.Node, path []string, value string) {
	if len(path) == 0 {
		node.Value = value
		return
	}

	for i := 0; i < len(node.Content); i += 2 {
		if len(node.Content) == 1 {
			modifyNode(node.Content[i], path, value)
		}

		if node.Content[i].Value == path[0] {
			modifyNode(node.Content[i+1], path[1:], value)
			return
		}

		if strings.Contains(path[0], "cronjob") {
			split := strings.Split(path[0], "[")
			if node.Content[i].Value == split[0] {
				r := regexp.MustCompile(`\d+`)
				for _, str := range path {
					matches := r.FindStringSubmatch(str)
					if len(matches) > 0 {
						index, _ := strconv.Atoi(matches[0])
						modifyNode(node.Content[i+1].Content[index], path[1:], value)
						return
					}
				}
			}
		}
	}

	return
}

// func modifyNode(node *yaml.Node, path []string, value string) {
// 	// fmt.Println(node.Content[0].Content[10].Value)            // image
// 	// fmt.Println(node.Content[0].Content[11].Content[4].Value) // tag

// 	// fmt.Println(node.Content[0].Content[10].Value)                       // image
// 	// fmt.Println(node.Content[0].Content[11].Content[0].Value)            // frontend
// 	// fmt.Println(node.Content[0].Content[11].Content[1].Content[4].Value) // tag
// 	// fmt.Println(node.Content[0].Content[11].Content[2].Value)            // backend
// 	// fmt.Println(node.Content[0].Content[11].Content[3].Content[4].Value) // tag

// 	// fmt.Println(node.Content[0].Content[10].Value)                                  // cronjobs
// 	// fmt.Println(node.Content[0].Content[11].Content[0].Content[8].Value)            // cronjobs[0].image
// 	// fmt.Println(node.Content[0].Content[11].Content[0].Content[9].Content[2].Value) // cronjobs[0].image.tag
// 	// fmt.Println(node.Content[0].Content[11].Content[1].Content[8].Value)            // cronjobs[1].image
// 	// fmt.Println(node.Content[0].Content[11].Content[1].Content[9].Content[2].Value) // cronjobs[1].image.tag
// 	// fmt.Println(node.Content[0].Content[11].Content[2].Content[8].Value)            // cronjobs[2].image
// 	// fmt.Println(node.Content[0].Content[11].Content[2].Content[9].Content[2].Value) // cronjobs[2].image.tag

// 	if len(path) == 0 {
// 		node.Value = value
// 		return
// 	}

// 	for i := 0; i < len(node.Content); i += 2 {
// 		if len(node.Content) == 1 {
// 			modifyNode(node.Content[i], path, value)
// 		}
// 		if node.Content[i].Value == path[0] {
// 			modifyNode(node.Content[i+1], path[1:], value)
// 			return
// 		}
// 	}

// 	return
// }
