package testMysql

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func TestMySQLConnection(t *testing.T) {
	db, err := gorm.Open("mysql", "root:root@(172.19.0.1:23306)/reddit?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Errorf("Error connecting to MySQL: %v", err)
		return // 连接失败时直接返回，避免在关闭未打开的数据库连接时出错
	}
	defer db.Close()

	if err := db.DB().Ping(); err != nil {
		t.Errorf("Failed to ping MySQL server: %v", err)
	}

	t.Logf("Connected to MySQL successfully")
	// Add more test cases here if needed
}
