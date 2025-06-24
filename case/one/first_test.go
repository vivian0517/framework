package one

import (
	"fmt"

	"framework/core/manager"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	manager.RegisterCase("描述1", "p0", Test_ceshi01)
	manager.RegisterCase("描述2", "p1", Test_ceshi02)
	manager.RegisterCase("描述3", "p0", Test_ceshi03)
}

func Test_ceshi01(t *testing.T) {
	fmt.Println("一个失败的case")
	//业务逻辑
	assert.Equal(t, 1, 2)
}

func Test_ceshi02(t *testing.T) {
	fmt.Println("一个成功的case")
	assert.Equal(t, 1, 1)
}

func Test_ceshi03(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}

// Add 返回两个整数的和
func Add(a, b int) int {
	return a + b
}
