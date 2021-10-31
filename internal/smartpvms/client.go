package smartpvms

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	resty *resty.Client
}

func (c *Client) Login(u, sc string) (*LoginResult, string, error) {
	res, err := c.resty.NewRequest().
		SetBody(&LoginBody{Username: u, SystemCode: sc}).
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

func (c *Client) SetXSRFToken(t string) {
	c.resty.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		r.SetHeader("Xsrf-Token", t)
		return nil
	})
}

func NewClient(hu string) *Client {
	r := resty.New()
	r.SetHostURL(hu)

	return &Client{
		resty: r,
	}
}
