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

func ModifyYamlFile(targetPath, inputValue string) error {
	var parseYamlNode yaml.Node
	err := ReadFromFile(targetPath, &parseYamlNode)

	if err != nil {
		return err
	}

	var valueTemplate map[string]string
	err = yaml.Unmarshal([]byte(inputValue), &valueTemplate)
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

	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := yaml.NewEncoder(f)
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
		// case for single node
		if len(node.Content) == 1 {
			modifyNode(node.Content[i], path, value)
		}

		// case for map
		if node.Content[i].Value == path[0] {
			modifyNode(node.Content[i+1], path[1:], value)
			return
		}

		// case for array
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
