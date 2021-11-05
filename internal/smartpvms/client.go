package smartpvms

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	resty *resty.Client
}

func (c *Client) Login(u, p string) (*LoginResult, string, error) {
	res, err := c.resty.NewRequest().
		SetBody(&LoginBody{Username: u, Password: p}).
		SetResult(&LoginResult{}).
		Post("/thirdData/login")

	if err != nil {
		return nil, "", err
	}

	return res.Result().(*LoginResult), res.Header().Get("Xsrf-Token"), nil
}

func (c *Client) Logout(t string) (*LogoutResult, error) {
	res, err := c.resty.NewRequest().
		SetBody(&LogoutBody{XSRFToken: t}).
		SetResult(&LogoutResult{}).
		Post("/thirdData/logout")

	if err != nil {
		return nil, err
	}

	return res.Result().(*LogoutResult), nil
}

func (c *Client) GetPlantList() (*GetPlantListResult, error) {
	res, err := c.resty.NewRequest().
		SetResult(&GetPlantListResult{}).
		Post("/thirdData/getStationList")

	if err != nil {
		return nil, err
	}

	return res.Result().(*GetPlantListResult), nil
}

func (c *Client) GetRealTimePlantData(cs []string) (*GetRealTimePlantDataResult, error) {
	res, err := c.resty.NewRequest().
		SetBody(&GetRealTimePlantDataBody{StationCodes: cs}).
		SetResult(&GetRealTimePlantDataResult{}).
		Post("/thirdData/getStationRealKpi")

	if err != nil {
		return nil, err
	}

	return res.Result().(*GetRealTimePlantDataResult), nil
}

func (c *Client) GetDeviceList(cs []string) (*GetDeviceListResult, error) {
	res, err := c.resty.NewRequest().
		SetBody(&GetDeviceListBody{StationCodes: cs}).
		SetResult(&GetDeviceListResult{}).
		Post("/thirdData/getDevList")

	if err != nil {
		return nil, err
	}

	return res.Result().(*GetDeviceListResult), nil
}

func (c *Client) GetDeviceData(ids []int64, t int) (*GetDeviceDataResult, error) {
	res, err := c.resty.NewRequest().
		SetBody(&GetDeviceDataBody{IDs: ids, Type: t}).
		SetResult(&GetDeviceDataResult{}).
		Post("/thirdData/getDevRealKpi")

	if err != nil {
		return nil, err
	}

	return res.Result().(*GetDeviceDataResult), nil
}

func (c *Client) SetXSRFToken(t string) {
	c.resty.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		r.SetHeader("Xsrf-Token", t)
		return nil
	})
}

func NewClient(bu string) *Client {
	r := resty.New()
	r.SetBaseURL(bu)

	return &Client{
		resty: r,
	}
}
