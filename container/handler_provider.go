// Package container 提供 HTTP 处理器的初始化和管理
package container

// HandlerProvider HTTP 处理器提供者
// 可以在此管理各种 HTTP Handler 的实例化逻辑
type HandlerProvider struct {
	container *AppContainer
}

// NewHandlerProvider 创建新的处理器提供者
func NewHandlerProvider(container *AppContainer) *HandlerProvider {
	return &HandlerProvider{
		container: container,
	}
}

// 根据实际需求，可以在这里添加各种 Handler 的工厂方法
// 例如:
// func (p *HandlerProvider) GetUserHandler() *handler.UserHandler {
//     return handler.NewUserHandler(p.container.GetDB())
// }
//
// func (p *HandlerProvider) GetAdminHandler() *handler.AdminHandler {
//     return handler.NewAdminHandler(p.container.GetDB())
// }
