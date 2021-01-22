package api

type SubServer interface {
	AddSubServerProxy(interface{})
	AddSubServerProxyGroup(interface{})
	ToFile() error
}
