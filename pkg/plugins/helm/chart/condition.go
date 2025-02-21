package chart

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/updatecli/updatecli/pkg/core/scm"
)

// Condition check if a specific chart version exist
func (c *Chart) Condition(source string) (bool, error) {
	if c.Version != "" {
		logrus.Infof("Version %v, already defined from configuration file", c.Version)
	} else {
		c.Version = source
	}
	URL := fmt.Sprintf("%s/index.yaml", c.URL)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return false, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	index, err := loadIndex(body)
	if err != nil {
		return false, err
	}

	message := ""
	if c.Version != "" {
		message = fmt.Sprintf(" for version '%s'", c.Version)
	}

	if index.Has(c.Name, c.Version) {
		logrus.Infof("\u2714 Helm Chart '%s' is available on %s%s", c.Name, c.URL, message)
		return true, nil
	}

	logrus.Infof("\u2717 Helm Chart '%s' isn't available on %s%s", c.Name, c.URL, message)
	return false, nil
}

// ConditionFromSCM returns an error because it's not supported
func (c *Chart) ConditionFromSCM(source string, scm scm.Scm) (bool, error) {
	return false, fmt.Errorf("SCM configuration is not supported for Helm chart condition, aborting")
}
