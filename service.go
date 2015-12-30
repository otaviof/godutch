package godutch

type Service struct {
	Name string
	g    *GoDutch
	Srvc *NrpeService
	cfg  *ServiceConfig
}

func NewService(cfg *ServiceConfig, g *GoDutch) *Service {
	var s *Service = &Service{Name: cfg.Name, g: g}
	s.Srvc = NewNrpeService(cfg, g)
	return s
}

func (s *Service) Shutdown() error {
	return nil
}

func (s *Service) Bootstrap() error {
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
	var component Component
	component = Component{
		Name:     s.Name,
		Checks:   []string{},
		Type:     "service",
		Instance: s.Srvc,
	}
	return &component
}

/* EOF */
