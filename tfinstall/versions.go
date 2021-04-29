package tfinstall

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-version"
)

type versionIndex struct {
	Versions map[string]interface{}
}

// ListVersions will return a sorted list of available Terraform versions.
// https://releases.hashicorp.com/terraform/index.json
func ListVersions(ctx context.Context) (version.Collection, error) {
	c := retryablehttp.NewClient()
	url := fmt.Sprintf("%s/%s", baseURL, "index.json")
	r, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(r)
	v := &versionIndex{}
	if err := dec.Decode(v); err != nil {
		return nil, err
	}

	versions := make(version.Collection, 0, len(v.Versions))
	for vx, _ := range v.Versions {
		sv, err := version.NewSemver(vx)
		if err != nil {
			return nil, err
		}

		versions = append(versions, sv)
	}

	sort.Sort(versions)

	return versions, nil
}