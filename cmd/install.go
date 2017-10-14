// Copyright © 2017 Daniel Jay Haskin <djhaskin987@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/djhaskin987/pask/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or update a package or packages",
	Long: `Read the spec file for a list of packages (zip archives). Download
and unzip them based on their locations in order of appearance in the spec file.
Any files in the "pask" folder inside the zip archives are treated specially,
and templating, pre-scripts and post-scripts are possible.
See the docs at pask.readthedocs.io.`,
	Run: func(cmd *cobra.Command, args []string) {
		if spec, err := pkg.ReadSpec(viper.Get("spec").(string)); err != nil {
			log.Fatalln("Error reading spec file: ", err)
		} else {
			spec.Install(viper.Get("base").(string))
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
