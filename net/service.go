package net

type Servicer interface {
	Start()
	Stop()
}

type Service struct {
	Conf  *ServiceConf
	Ports []Porter
}


type ServiceConf struct {
	Name      string
	Factory   ServiceFactory
	PortConfs []*PortConf
	Ops       *NetOptions
}

func NewServiceConf(name string) *ServiceConf {
	return &ServiceConf{Name: name}
}

func NewService(conf *ServiceConf) *Service {
	service := &Service{Conf: conf}
	service.Init()
	return service
}

func (self *Service) Init() {
	factory := self.Conf.Ops.factory
	for _, portConf := range self.Conf.PortConfs {
		port := factory.NewPort(portConf)
		self.Ports = append(self.Ports, port)
	}
}

func (self *Service) Start() {
	for _, port := range self.Ports {
		port.Start()
	}
}

func (self *Service) Stop() {
	for _, port := range self.Ports {
		port.Stop()
	}
}

func (self *ServiceConf) NewPort(typ PortType) *PortConf {
	newConf := NewPortConf(typ)
	newConf.ServiceConf = self
	self.PortConfs = append(self.PortConfs, newConf)
	return newConf
}

func (self *ServiceConf) BuildFactory(factory ServiceFactory) *ServiceConf {
	self.Factory = factory
	return self
}

