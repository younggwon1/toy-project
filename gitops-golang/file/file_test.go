package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Case of image value to modify.
--values "{\"image.tag\":\"dev-c7630b\"}"
--values "{\"image.frontend.tag\":\"dev-c7630b\"}"
--values "{\"image.backend.tag\":\"dev-c7630b\"}"
--values "{\"image.grpc.tag\":\"dev-c7630b\"}"
--values "{\"images.grpc.tag\":\"dev-c7630b\"}"
--values "{\"cronjobs[0].image.tag\":\"dev-14b5fa\"}"
--values "{\"cronjobs[1].image.tag\":\"dev-14b5fa\"}"
*/

/*
The Case I Think of image value to modify.
--values "{\"image.cronjobs[0].tag\":\"dev-14b5fa\"}"
--values "{\"image.cronjobs[1].tag\":\"dev-14b5fa\"}"
*/

/*
cronjobs[0], cronjobs[1] 이러한 list 형식의 input 이 들어올 때, 이러한 list 형식의 input 을 어떻게 처리할 것인가?
이 부분에 대한 로직 구현이 필요해보인다.
*/

func TestModify(t *testing.T) {
	// set the target value to be modified
	targetValues := `image:
  name: "default"
  tag: tmp
  frontend:
    name: "frontend"
    tag: tmp
  backend:
    name: "backend"
    tag: tmp
  grpc:
    name: "grpc"
    tag: tmp
cronjobs:
  - name: "cronjobOne"
    image:
      tag: tmp
  - name: "cronjobTwo"
    image:
      tag: tmp
  - name: "cronjobThree"
    image:
      tag: tmp
arrays:
  - name: "life"
    definition: "wonderful"
multi:
  frontend:
    tag: tmp
  backend:
    tag: tmp
`

	// set the values to modify
	valuesToModify := [...]string{
		`{"image.tag":"v3.14"}`,
		`{"image.frontend.tag":"v3.14"}`,
		`{"image.backend.tag":"v3.14"}`,
		`{"image.grpc.tag":"v3.14"}`,
		`{"arrays[0].definition":"beautiful"}`,
		`{"cronjobs[0].image.tag":"v3.14", "cronjobs[1].image.tag":"v3.14", "cronjobs[2].image.tag":"v3.14"}`,
		`{"multi.frontend.tag":"v3.14", "multi.backend.tag":"v3.14"}`,
	}

	// set expected values
	expectedValues := `image:
  name: "default"
  tag: v3.14
  frontend:
    name: "frontend"
    tag: v3.14
  backend:
    name: "backend"
    tag: v3.14
  grpc:
    name: "grpc"
    tag: v3.14
cronjobs:
  - name: "cronjobOne"
    image:
      tag: v3.14
  - name: "cronjobTwo"
    image:
      tag: v3.14
  - name: "cronjobThree"
    image:
      tag: v3.14
arrays:
  - name: "life"
    definition: "beautiful"
multi:
  frontend:
    tag: v3.14
  backend:
    tag: v3.14
`

	// write target value to yaml file
	err := os.WriteFile("/tmp/value.yaml", []byte(targetValues), 0644)
	assert.NoError(t, err)
	defer os.Remove("/tmp/value.yaml")

	// mofidy values
	for _, value := range valuesToModify {
		err = ModifyYamlFile(
			"/tmp/value.yaml",
			value,
		)
		assert.NoError(t, err)
	}

	// read changed yaml file
	data, err := os.ReadFile("/tmp/value.yaml")
	assert.NoError(t, err)

	// compare expected and changed values
	assert.Equal(t, expectedValues, string(data))
}
