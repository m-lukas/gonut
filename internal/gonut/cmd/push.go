// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"
	"github.com/gonvenience/text"
	"github.com/homeport/gonut/internal/gonut/assets"
	"github.com/homeport/gonut/internal/gonut/cf"
	"github.com/homeport/pina-golada/pkg/files"
)

// GonutAppPrefix is the prefeix for gonuts applications, it is also used by the
// cleanup command to decide whether an app is pushed by gonut or not.
var GonutAppPrefix = "gonut"

type sampleApp struct {
	caption       string
	buildpack     string
	command       string
	aliases       []string
	appNamePrefix string
	assetFunc     func() (files.Directory, error)
}

var (
	deleteSetting  string
	summarySetting string
	noPingSetting  bool
)

var sampleApps = []sampleApp{
	{
		caption:       "Golang",
		command:       "golang",
		buildpack:     "go_buildpack",
		aliases:       []string{"go"},
		appNamePrefix: fmt.Sprintf("%s-golang-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.GoSampleApp,
	},

	{
		caption:       "Python",
		command:       "python",
		buildpack:     "python_buildpack",
		aliases:       []string{},
		appNamePrefix: fmt.Sprintf("%s-python-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.PythonSampleApp,
	},

	{
		caption:       "PHP",
		command:       "php",
		buildpack:     "php_buildpack",
		aliases:       []string{},
		appNamePrefix: fmt.Sprintf("%s-php-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.PHPSampleApp,
	},

	{
		caption:       "Staticfile",
		command:       "staticfile",
		buildpack:     "staticfile_buildpack",
		aliases:       []string{"static"},
		appNamePrefix: fmt.Sprintf("%s-staticfile-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.StaticfileSampleApp,
	},

	{
		caption:       "Swift",
		command:       "swift",
		buildpack:     "swift_buildpack",
		aliases:       []string{},
		appNamePrefix: fmt.Sprintf("%s-swift-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.SwiftSampleApp,
	},

	{
		caption:       "NodeJS",
		command:       "nodejs",
		buildpack:     "nodejs_buildpack",
		aliases:       []string{"node"},
		appNamePrefix: fmt.Sprintf("%s-nodejs-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.NodeJSSampleApp,
	},

	{
		caption:       "Ruby",
		command:       "ruby",
		buildpack:     "ruby_buildpack",
		appNamePrefix: fmt.Sprintf("%s-ruby-sinatra-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.RubySampleApp,
	},

	{
		caption:       ".NET",
		command:       "dotnet",
		buildpack:     "dotnet-core",
		appNamePrefix: fmt.Sprintf("%s-dotnet-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.DotNetSampleApp,
	},

	{
		caption:       "Binary",
		command:       "binary",
		buildpack:     "binary_buildpack",
		appNamePrefix: fmt.Sprintf("%s-binary-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.BinarySampleApp,
	},

	{
		caption:       "Java",
		command:       "java",
		buildpack:     "java_buildpack",
		appNamePrefix: fmt.Sprintf("%s-java-app-", GonutAppPrefix),
		assetFunc:     assets.Provider.JavaSampleApp,
	},
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push a sample app to Cloud Foundry",
	Long:  `Use one of the sub-commands to select a sample app of a list of programming languages to be pushed to a Cloud Foundry instance.`,
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.PersistentFlags().StringVarP(&deleteSetting, "delete", "d", "always", "Delete application after push: always, never, on-success")
	pushCmd.PersistentFlags().StringVarP(&summarySetting, "summary", "s", "short", "Push summary detail level: quiet, short, full")
	pushCmd.PersistentFlags().BoolVarP(&noPingSetting, "no-ping", "p", false, "Do not ping application after push")

	for _, sampleApp := range sampleApps {
		pushCmd.AddCommand(&cobra.Command{
			Use:     sampleApp.command,
			Aliases: sampleApp.aliases,
			Short:   fmt.Sprintf("Push a %s sample app to Cloud Foundry", sampleApp.caption),
			Long:    fmt.Sprintf(`Push a %s sample app to Cloud Foundry. The application will be deleted after it was pushed successfully.`, sampleApp.caption),
			Run:     genericCommandFunc,
		})
	}

	pushCmd.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "Pushes all available sample apps to Cloud Foundry",
		Long:  `Pushes all available sample apps to Cloud Foundry. Each application will be deleted after it was pushed successfully.`,
		Run: func(cmd *cobra.Command, args []string) {
			for _, sampleApp := range sampleApps {
				if err := runSampleAppPush(sampleApp); err != nil {
					ExitGonut(err)
				}
			}
		},
	})
}

func lookUpSampleAppByName(name string) *sampleApp {
	for _, sampleApp := range sampleApps {
		if sampleApp.command == name {
			return &sampleApp
		}
	}

	return nil
}

func genericCommandFunc(cmd *cobra.Command, args []string) {
	sampleApp := lookUpSampleAppByName(cmd.Use)
	if sampleApp == nil {
		ExitGonut("failed to detect which sample app is to be tested")
	}

	if err := runSampleAppPush(*sampleApp); err != nil {
		ExitGonut(err)
	}
}

func runSampleAppPush(app sampleApp) error {
	hasBuildpack, err := cf.HasBuildpack(app.buildpack)
	if err != nil {
		return err
	}

	// Skip sample app push if desired buildpack is unavailable
	if !hasBuildpack {
		bunt.Printf("Skipping push of *%s* sample app, because there is no DarkSeaGreen{%s} installed.\n",
			app.caption,
			app.buildpack,
		)

		return nil
	}

	var cleanupSetting cf.AppCleanupSetting
	switch deleteSetting {
	case "always":
		cleanupSetting = cf.Always

	case "never":
		cleanupSetting = cf.Never

	case "on-success":
		cleanupSetting = cf.OnSuccess

	default:
		return fmt.Errorf("unsupported delete setting: %s", deleteSetting)
	}

	appName := text.RandomStringWithPrefix(app.appNamePrefix, 32)

	directory, err := app.assetFunc()
	if err != nil {
		return err
	}

	report, err := cf.PushApp(app.caption, appName, directory, cleanupSetting, noPingSetting)
	if err != nil {
		return err
	}

	switch strings.ToLower(summarySetting) {
	case "quiet":
		// Nothing to report

	case "short", "oneline":
		bunt.Printf("Successfully pushed *%s* sample app in CadetBlue{%s}.\n",
			app.caption,
			cf.HumanReadableDuration(report.ElapsedTime()),
		)

	case "json":
		out, err := neat.NewOutputProcessor(true, true, &neat.DefaultColorSchema).ToJSON(report.Export())
		if err != nil {
			return err
		}

		fmt.Println(out)

	case "yaml":
		out, err := neat.ToYAMLString(report.Export())
		if err != nil {
			return err
		}

		fmt.Println(out)

	case "full":
		headline := bunt.Sprintf("Successfully pushed *%s* sample app in CadetBlue{%s}",
			app.caption,
			cf.HumanReadableDuration(report.ElapsedTime()),
		)

		content, err := neat.Table(report.ExportTable(), neat.AlignRight(0))
		if err != nil {
			return err
		}

		neat.Box(os.Stdout, headline, strings.NewReader(content))
	}

	return nil
}
