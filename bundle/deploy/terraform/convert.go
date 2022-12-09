package terraform

import (
	"encoding/json"

	"github.com/databricks/bricks/bundle/config"
	"github.com/databricks/bricks/bundle/internal/tf/schema"
)

func conv(from any, to any) {
	buf, _ := json.Marshal(from)
	json.Unmarshal(buf, &to)
}

// BundleToTerraform converts resources in a bundle configuration
// to the equivalent Terraform JSON representation.
//
// NOTE: THIS IS CURRENTLY A HACK. WE NEED A BETTER WAY TO
// CONVERT TO/FROM TERRAFORM COMPATIBLE FORMAT.
func BundleToTerraform(config *config.Root) *schema.Root {
	tfroot := schema.NewRoot()
	tfroot.Provider.Databricks.Profile = config.Workspace.Profile

	for k, src := range config.Resources.Jobs {
		var dst schema.ResourceJob
		conv(src, &dst)

		for _, v := range src.Tasks {
			var t schema.ResourceJobTask
			conv(v, &t)
			dst.Task = append(dst.Task, t)
		}

		for _, v := range src.JobClusters {
			var t schema.ResourceJobJobCluster
			conv(v, &t)
			dst.JobCluster = append(dst.JobCluster, t)
		}

		// Unblock downstream work. To be addressed more generally later.
		if git := src.GitSource; git != nil {
			dst.GitSource = &schema.ResourceJobGitSource{
				Url:      git.GitUrl,
				Branch:   git.GitBranch,
				Commit:   git.GitCommit,
				Provider: string(git.GitProvider),
				Tag:      git.GitTag,
			}
		}

		tfroot.Resource.Job[k] = &dst
	}

	for k, src := range config.Resources.Pipelines {
		var dst schema.ResourcePipeline
		conv(src, &dst)

		for _, v := range src.Libraries {
			var l schema.ResourcePipelineLibrary
			conv(v, &l)
			dst.Library = append(dst.Library, l)
		}

		for _, v := range src.Clusters {
			var l schema.ResourcePipelineCluster
			conv(v, &l)
			dst.Cluster = append(dst.Cluster, l)
		}

		tfroot.Resource.Pipeline[k] = &dst
	}

	// Clear data sources because we don't have any.
	tfroot.Data = nil

	return tfroot
}
