package godutch

//
// Type Service must implement Composer interface and keep track of the running
// service instance (Suture interface compatible) running underneath. Here the
// actual service boostrap should take place.
//
type Service struct {
	Name string
	g    *GoDutch
	Srvc *NrpeService
	cfg  *ServiceConfig
}

// Creates a new Service instance spawning a Nrpe listener.
func NewService(cfg *ServiceConfig, g *GoDutch) *Service {
	var s *Service = &Service{Name: cfg.Name, g: g}
	// TODO
	//  * How to identify and only load the informed service? Here we have a
	//    hardcoded NrpeService being managed;
	s.Srvc = NewNrpeService(cfg, g)
	return s
}

// Dummy method on Service, there's no bootstrap here.
func (s *Service) Bootstrap() error {
	return nil
}

// Shutdown will cover calling for running service stop.
func (s *Service) Shutdown() error {
	s.Srvc.Stop()
	return nil
}

// Execute method is not implememented on this Object, atlhough, still required
// to be part of "component" interface.
func (n *Service) Execute(req []byte) (*Response, error) {
	var err error
	var resp *Response = &Response{}
	return resp, err
}

// Displays the name for this service (component).
func (s *Service) ComponentInfo() *Component {
	var component *Component
	component = &Component{
		Name:     s.Name,
		Checks:   []string{},
		Type:     "service",
		Instance: s.Srvc,
	}
	return component
}

/* EOF */
