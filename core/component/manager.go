package component

import "errors"

var (
	gErrNotImplemented = errors.New("not implemented")
)

// Manager - 基于本机的组件管理
type Manager struct {
	root string // root directory
}

// ListAll - 列出所有组件
func (m Manager) ListAll() ([]*Description, error) {
	return nil, gErrNotImplemented
}

// Install - 将指定路径下的组件安装到 root 中
func (m Manager) Install(src string, start bool) error {
	return gErrNotImplemented
}

// Uninstall - 移除指定名称的组件
func (m Manager) Uninstal(name string) error {
	return gErrNotImplemented
}

// Start - start component with args
func (m Manager) Start(name string, args ...string) error {
	return gErrNotImplemented
}

// Stop - stop component
func (m Manager) Stop(name string) error {
	return gErrNotImplemented
}
