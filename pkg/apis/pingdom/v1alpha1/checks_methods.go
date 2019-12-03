package v1alpha1

import (
	"github.com/russellcardullo/go-pingdom/pingdom"
)

// SetupHTTP setups HTTP check
func (checks *Checks) SetupHTTP() error {
	client, err := checks.makeClient()
	if err != nil {
		return err
	}

	if checks.Status.ID != 0 {
		return checks.update(client)
	}

	return checks.create(client)
}

// DeleteHTTP delete HTTP check
func (checks *Checks) DeleteHTTP() error {
	client, err := checks.makeClient()
	if err != nil {
		return err
	}

	if checks.Status.ID == 0 {
		return nil
	}

	return checks.delete(client)
}

// -- helpers --

func (checks *Checks) makeClient() (*pingdom.Client, error) {
	return pingdom.NewClientWithConfig(pingdom.ClientConfig{
		User:     checks.Spec.User,
		Password: checks.Spec.Password,
		APIKey:   checks.Spec.Key,
	})
}

func (checks *Checks) create(client *pingdom.Client) error {
	details, err := client.Checks.Create(convert(checks.Spec.HTTP))
	if err != nil {
		return err
	}

	checks.Status.ID = details.ID

	return nil
}

func (checks *Checks) update(client *pingdom.Client) error {
	_, err := client.Checks.Update(
		checks.Status.ID,
		convert(checks.Spec.HTTP),
	)
	if err != nil {
		return err
	}

	return nil
}

func (checks *Checks) delete(client *pingdom.Client) error {
	_, err := client.Checks.Delete(checks.Status.ID)
	if err != nil {
		return err
	}

	checks.Status.ID = 0

	return nil
}

func convert(http *HTTPCheck) *pingdom.HttpCheck {
	tmp := pingdom.HttpCheck(*http)
	return &tmp
}
