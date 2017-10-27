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

// runCmd represents the run command
var runCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "run TASK",
	Short: "Run a packaged task",
	Long: `Calls any executable file called TASK found in the folder
"<root-path>/pask/<pkg>/<vers>/tasks". Treats each package folder in order
of appearance of the packages under the "packages" key in the spec file.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec := viper.Get("spec").(string)
		base := viper.Get("base").(string)
		log.Println("Running packaged tasks...")
		log.Printf("Using spec file `%s`\n", spec)
		log.Println("Using project base `%s`\n", base)
		if spec, err := pkg.ReadSpec(spec); err != nil {
			log.Fatalln("Error reading spec file:", err)
		} else {
			log.Printf("Using base directory `%s`\n", base)
			for _, task := range args {
				if err := spec.Run(base, task); err != nil {
					log.Fatalf("Problem running task `%s`: `%s`\n", task, err)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
