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
	Versions map[string]struct{}
}

// ListVersions will return a sorted list of available Terraform versions.
// versions are sourced from
// https://releases.hashicorp.com/terraform/index.json
func ListVersions(ctx context.Context) (version.Collection, error) {
	c := retryablehttp.NewClient()
	url := fmt.Sprintf("%s/%s", baseURL, "index.json")

	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	r, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(r.Body)
	v := &versionIndex{}
	if err := dec.Decode(v); err != nil {
		return nil, err
	}

	versions := make(version.Collection, 0, len(v.Versions))
	for vx := range v.Versions {
		sv, err := version.NewSemver(vx)
		if err != nil {
			return nil, err
		}
		versions = append(versions, sv)
	}

	sort.Sort(versions)

	return versions, nil
}
