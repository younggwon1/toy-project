package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
티켓 입력 상황 정리
조건 : 여러 개의 티켓 입력 시 무조건 '_' 로 티켓이 연결되어야 한다.
1. jira ticket은 입력되어야한다.
2. CR-XX, SD-XX 의 기본 값이 jira ticket이 입력되는 경우는 배포를 막아야한다.
3. CR-1234 or CI-1234 와 같은 단일 티켓 입력
4. CR-1111_CR-2222_CR-3333_CR-4444_CR-5555 와 같은 복수 티켓 입력
5. CR-1111_CR-2222_CI-1234_CI-2345 와 같은 복수 티켓 입력
*/

func TestRegex(t *testing.T) {
	// set the test cases
	result := ValidateTicket("")
	expected := false
	assert.Equal(t, expected, result)

	result = ValidateTicket("N/A")
	expected = false
	assert.Equal(t, expected, result)

	result = ValidateTicket("CR-XX")
	expected = false
	assert.Equal(t, expected, result)

	result = ValidateTicket("SD-XX")
	expected = false
	assert.Equal(t, expected, result)

	result = ValidateTicket("CR-1234")
	expected = true
	assert.Equal(t, expected, result)

	result = ValidateTicket("CI-1234")
	expected = true
	assert.Equal(t, expected, result)

	result = ValidateTicket("CR-1111_CR-2222_CR-3333_CR-4444_CR-5555")
	expected = true
	assert.Equal(t, expected, result)

	result = ValidateTicket("CR-1111_CR-2222_CI-1234_CI-2345")
	expected = true
	assert.Equal(t, expected, result)

	result = ValidateTicket("CR-1111_CR-2222_CR-3333_CR-4444_CR-XX")
	expected = false
	assert.Equal(t, expected, result)

	result = ValidateTicket("CR-1111_CR-2222_CI-1234_CI-XX")
	expected = false
	assert.Equal(t, expected, result)
}
