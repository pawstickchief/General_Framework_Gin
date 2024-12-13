package casbin

import (
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"log"

	"github.com/casbin/casbin/v2/model"
)

var Enforcer *casbin.Enforcer

// Init 初始化 Casbin 权限管理
func Init() {
	m, err := model.NewModelFromFile("casbin_model.conf")
	if err != nil {
		log.Fatalf("加载模型文件失败: %v", err)
	}

	a := fileadapter.NewAdapter("casbin_policy.csv")
	Enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("初始化 Casbin 失败: %v", err)
	}

	if err := Enforcer.LoadPolicy(); err != nil {
		log.Fatalf("加载策略文件失败: %v", err)
	}
}
