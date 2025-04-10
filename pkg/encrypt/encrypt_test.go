package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	password := "admin1"
	salt := "saltsalt"

	// 手动计算预期结果
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	expected := hex.EncodeToString(hash.Sum(nil))

	// 调用被测试函数
	result := EncryptPassword(password, salt)

	// 打印结果
	t.Logf("EncryptPassword(%q, %q) = %q", password, salt, result)

	// 验证结果
	if result != expected {
		t.Errorf("EncryptPassword(%q, %q) = %q; want %q", password, salt, result, expected)
	}
}
