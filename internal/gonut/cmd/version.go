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
	"image/color"

	"github.com/gonvenience/bunt"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/spf13/cobra"
)

var version string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version",
	Long:  `Shows the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(GetVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// GetVersion returns the version of the tool
func GetVersion() string {
	if len(version) == 0 {
		version = "development"
	}

	lightblue, _ := colorful.MakeColor(color.RGBA{77, 173, 233, 255})
	otherblue, _ := colorful.MakeColor(color.RGBA{63, 143, 231, 255})

	return fmt.Sprintf("%s%s version %s\n",
		bunt.Style("go", bunt.Foreground(lightblue)),
		bunt.Style("nut", bunt.Foreground(otherblue), bunt.Bold()),
		bunt.Style(version, bunt.Foreground(bunt.DimGray)),
	)
}
