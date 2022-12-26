package main

import (
	"fmswift.com.cn/watchmen/core/component"
)

type ComponentManager struct {
	root string // root directory

}

func (cm *ComponentManager) ListAll() ([]component.Description, error) {

	return nil, nil
}

// scanRoot - 扫描组件根目录
func (cm *ComponentManager) scanRoot() ([]component.Description, error) {

	return nil, nil
}
