package smartpvms

import (
	"errors"
	"sync"

	"github.com/go-resty/resty/v2"
)

type XSRFTokenSource interface {
	XSRFToken() (*XSRFToken, error)
}

type Config struct {
	BaseURL  string
	Username string
	Password string
}

func (c *Config) Client() *resty.Client {
	return NewClient(c, c.XSRFTokenSource())
}

func (c *Config) XSRFTokenSource() XSRFTokenSource {
	return &xsrfReuseTokenSource{
		source: &xsrfTokenRefresher{
			config: c,
		},
	}
}

func NewClient(cfg *Config, src XSRFTokenSource) *resty.Client {
	r := resty.New().SetBaseURL(cfg.BaseURL)

	if src != nil {
		r.OnBeforeRequest(func(_ *resty.Client, req *resty.Request) error {
			t, err := src.XSRFToken()
			if err != nil {
				return err
			}

			req.SetHeader("Xsrf-Token", t.XSRFToken)

			return nil
		})
	}

	return r
}

type xsrfTokenRefresher struct {
	config *Config
}

func (r *xsrfTokenRefresher) XSRFToken() (*XSRFToken, error) {
	c := NewClient(r.config, nil)

	res, tkn, err := Login(c, r.config.Username, r.config.Password)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, errors.New("smartpvms: login failed")
	}

	return tkn, nil
}

type xsrfReuseTokenSource struct {
	source XSRFTokenSource

	mutex sync.Mutex
	token *XSRFToken
}

func (s *xsrfReuseTokenSource) XSRFToken() (*XSRFToken, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.token.IsValid() {
		return s.token, nil
	}

	t, err := s.source.XSRFToken()
	if err != nil {
		return nil, err
	}

	s.token = t

	return t, nil
}

func Login(c *resty.Client, u, p string) (*LoginResult, *XSRFToken, error) {
	res, err := c.NewRequest().
		SetBody(&LoginBody{Username: u, Password: p}).
		SetResult(&LoginResult{}).
		Post("/thirdData/login")

	if err != nil {
		return nil, nil, err
	}

	t := NewXSRFToken(
		res.Header().Get("Xsrf-Token"),
	)

	return res.Result().(*LoginResult), t, nil
}

func Logout(c *resty.Client, t string) (*LogoutResult, error) {
	res, err := c.NewRequest().
		SetBody(&LogoutBody{XSRFToken: t}).
		SetResult(&LogoutResult{}).
		Post("/thirdData/logout")

	if err != nil {
		return nil, err
	}

	return res.Result().(*LogoutResult), nil
}

func GetPlantList(c *resty.Client) (*GetPlantListResult, error) {
	res, err := c.NewRequest().
		SetResult(&GetPlantListResult{}).
		Post("/thirdData/getStationList")

	if err != nil {
		return nil, err
	}

	return res.Result().(*GetPlantListResult), nil
}

func GetRealtimePlantData(c *resty.Client, cs ...string) (*GetRealtimePlantDataResult, error) {
	res, err := c.NewRequest().
		SetBody(&GetRealtimePlantDataBody{StationCodes: cs}).
		SetResult(&GetRealtimePlantDataResult{}).
		Post("/thirdData/getStationRealKpi")

	if err != nil {
		return nil, err
	}

	return res.Result().(*GetRealtimePlantDataResult), nil
}

func GetDeviceList(c *resty.Client, cs ...string) (*GetDeviceListResult, error) {
	res, err := c.NewRequest().
		SetBody(&GetDeviceListBody{StationCodes: cs}).
		SetResult(&GetDeviceListResult{}).
		Post("/thirdData/getDevList")

	if err != nil {
		return nil, err
	}

	return res.Result().(*GetDeviceListResult), nil
}

func GetRealtimeDeviceData[T any](
	c *resty.Client,
	t DeviceType,
	ids ...int64,
) (*GetRealtimeDeviceDataResult[T], error) {
	res, err := c.NewRequest().
		SetBody(&GetRealtimeDeviceDataBody{Type: t, IDs: ids}).
		SetResult(&GetRealtimeDeviceDataResult[T]{}).
		Post("/thirdData/getDevRealKpi")

	if err != nil {
		return nil, err
	}

	return res.Result().(*GetRealtimeDeviceDataResult[T]), nil
}
