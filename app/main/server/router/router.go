package router

type Router interface {
	Group(prefix string, group func(group RouterGroup)) RouterGroup
}

type RouterGroup interface {
	Router
}
