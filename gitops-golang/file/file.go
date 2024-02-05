package file

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ModifyImageTagInYAMLFile(helmRepo string, yamlFile string, newImageValue string) error {
	yamlFileName := "/tmp/" + helmRepo + "/" + yamlFile
	// 파일에서 yaml을 읽습니다.
	data, err := os.ReadFile(yamlFileName)
	if err != nil {
		return err
	}

	// yaml을 노드 트리로 파싱합니다.
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return err
	}

	var imageValue map[string]string
	if err := yaml.Unmarshal([]byte(newImageValue), &imageValue); err != nil {
		return err
	}

	var keypart []string
	var keyvalue string
	for key, value := range imageValue {
		keypart = strings.Split(key, ".")
		keyvalue = value
	}
	modifyNode(&root, keypart, keyvalue)

	file, err := os.Create(yamlFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	if err := encoder.Encode(&root); err != nil {
		return err
	}

	return nil
}

func modifyNode(node *yaml.Node, pathParts []string, newValue string) {
	if len(pathParts) == 0 {
		node.Value = newValue
		return
	}

	if len(pathParts) == 1 {
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			if keyNode.Value == pathParts[0] {
				valueNode.Value = newValue
				return
			}
		}
	} else if len(pathParts) > 1 {
		for i := 0; i < len(node.Content[0].Content); i += 2 {
			keyNode := node.Content[0].Content[i]
			valueNode := node.Content[0].Content[i+1]
			if keyNode.Value == pathParts[0] {
				modifyNode(valueNode, pathParts[1:], newValue)
			}
		}
	}
	return
}
