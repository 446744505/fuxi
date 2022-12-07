package core

type Net interface {
	Controler
	AddService(service Service)
}

type CoreNet struct {
	services []Service
}

func (self *CoreNet) Start() {
	for _, service := range self.services {
		service.Start()
	}
}

func (self *CoreNet) Stop() {
	for _, service := range self.services {
		service.Stop()
	}
}

func (self *CoreNet) AddService(service Service) {
	self.services = append(self.services, service)
}