/*
Copyright 2021 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/inspect"
	buildEnv "github.com/GoogleContainerTools/skaffold/pkg/skaffold/inspect/buildEnv"
)

var buildEnvFlags = struct {
	profile     string
	projectID   string
	diskSizeGb  int64
	machineType string
	timeout     string
	concurrency int
}{}

func cmdBuildEnv() *cobra.Command {
	return NewCmd("build-env").
		WithDescription("Interact with skaffold build environment definitions.").
		WithPersistentFlagAdder(cmdBuildEnvFlags).
		WithCommands(cmdBuildEnvList(), cmdBuildEnvAdd())
}

func cmdBuildEnvList() *cobra.Command {
	return NewCmd("list").
		WithExample("Get list of target build environments with activated profiles p1 and p2", "inspect build-env list -p p1,p2 --format json").
		WithDescription("Print the list of active build environments.").
		WithFlagAdder(cmdBuildEnvListFlags).
		NoArgs(listBuildEnv)
}

func cmdBuildEnvAdd() *cobra.Command {
	return NewCmd("add").
		WithDescription("Add a new build environment to the default pipeline or to a new or existing profile.").
		WithPersistentFlagAdder(cmdBuildEnvAddFlags).
		WithCommands(cmdBuildEnvAddGcb())
}

func cmdBuildEnvAddGcb() *cobra.Command {
	return NewCmd("googleCloudBuild").
		WithDescription("Add a new GoogleCloudBuild build environment definition").
		WithLongDescription(`Add a new GoogleCloudBuild build environment definition.
Without the '--profile' flag the new environment definition is added to the default pipeline. With the '--profile' flag it will create a new profile with this build env definition. 
In these respective scenarios, it will fail if the build env definition for the default pipeline or the named profile already exists. To override an existing definition use 'skaffold inspect build-env modify' command instead. 
Use the '--module' filter to specify the individual module to target. Otherwise, it'll be applied to all modules defined in the target file. Also, with the '--profile' flag if the target config imports other configs as dependencies, then the new profile will be recursively created in all the imported configs also.`).
		WithExample("Add a new profile named 'gcb' targeting the builder 'googleCloudBuild' against the GCP project ID '1234'.", "inspect build-env add googleCloudBuild --profile gcb --projectID 1234 -f skaffold.yaml").
		WithFlagAdder(cmdBuildEnvAddGcbFlags).
		NoArgs(addGcbBuildEnv)
}

func listBuildEnv(ctx context.Context, out io.Writer) error {
	return buildEnv.PrintBuildEnvsList(ctx, out, printBuildEnvsListOptions())
}

func addGcbBuildEnv(ctx context.Context, out io.Writer) error {
	return buildEnv.AddGcbBuildEnv(ctx, out, addGcbBuildEnvOptions())
}

func cmdBuildEnvAddFlags(f *pflag.FlagSet) {
	f.StringVarP(&buildEnvFlags.profile, "profile", "p", "", `Profile name to add the new build env definition in. If the profile name doesn't exist then the profile will be created in all the target configs. If this flag is not specified then the build env is added to the default pipeline of the target configs.`)
}

func cmdBuildEnvAddGcbFlags(f *pflag.FlagSet) {
	f.StringVar(&buildEnvFlags.projectID, "projectId", "", `ID of the Cloud Platform Project.`)
	f.Int64Var(&buildEnvFlags.diskSizeGb, "diskSizeGb", 0, `Disk size of the VM that runs the build`)
	f.StringVar(&buildEnvFlags.machineType, "machineType", "", `Type of VM that runs the build`)
	f.StringVar(&buildEnvFlags.timeout, "timeout", "", `Build timeout (in seconds)`)
	f.IntVar(&buildEnvFlags.concurrency, "concurrency", -1, `number of artifacts to build concurrently. 0 means "no-limit"`)
}

func cmdBuildEnvFlags(f *pflag.FlagSet) {
	f.StringSliceVarP(&inspectFlags.modules, "module", "m", nil, "Names of modules to filter target action by.")
}

func cmdBuildEnvListFlags(f *pflag.FlagSet) {
	f.StringSliceVarP(&inspectFlags.profiles, "profile", "p", nil, `Profile names to activate`)
}

func printBuildEnvsListOptions() inspect.Options {
	return inspect.Options{
		Filename:  inspectFlags.fileName,
		OutFormat: inspectFlags.outFormat,
		Modules:   inspectFlags.modules,
		BuildEnvOptions: inspect.BuildEnvOptions{
			Profiles: inspectFlags.profiles,
		},
	}
}
func addGcbBuildEnvOptions() inspect.Options {
	return inspect.Options{
		Filename:  inspectFlags.fileName,
		OutFormat: inspectFlags.outFormat,
		Modules:   inspectFlags.modules,
		BuildEnvOptions: inspect.BuildEnvOptions{
			Profile:     buildEnvFlags.profile,
			ProjectID:   buildEnvFlags.projectID,
			DiskSizeGb:  buildEnvFlags.diskSizeGb,
			MachineType: buildEnvFlags.machineType,
			Timeout:     buildEnvFlags.timeout,
			Concurrency: buildEnvFlags.concurrency,
		},
	}
}