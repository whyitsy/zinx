package znet

import (
	"zinx/zInterface"
)

/*
BaseRouter实现IRouter接口的三个方法, 但是什么都不做的目的：为使用者提供基类, 可按需重写, 不需要实现接口的所有方法.
因为一般就只需要使用Handle方法, PreHandle和PostHandle是按需使用的.
*/
type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(request zInterface.IRequest) {
}

func (b *BaseRouter) Handle(request zInterface.IRequest) {
}

func (b *BaseRouter) PostHandle(request zInterface.IRequest) {

}
