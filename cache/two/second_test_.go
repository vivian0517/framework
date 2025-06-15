package two

import (
	"fmt"
	"framework/core/manager"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	manager.RegisterCase("描述sdfd", "p1", Test_ceshi01)
	manager.RegisterCase("描述sdf", "p0", Test_ceshi02)
}
func Test_ceshi01(t *testing.T) {

	assert.Equal(t, 1, 2)
}

func Test_ceshi02(t *testing.T) {

	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, 1)
	fmt.Println(time.Now())
	assert.Equal(t, 1, 2)
}
