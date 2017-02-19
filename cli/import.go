// Copyright © 2016 Abcum Ltd
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

package cli

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/abcum/surreal/log"
)

var importCmd = &cobra.Command{
	Use:     "import [flags] <file>",
	Short:   "Import data into an existing database",
	Example: "  surreal import --auth root:root backup.db",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		var fle *os.File
		var req *http.Request
		var res *http.Response

		// Ensure that the command has a filepath
		// as the output file argument. If no filepath
		// has been provided then return an error.

		if len(args) != 1 {
			log.Fatalln("No filepath provided.")
			return
		}

		// Attempt to open or create the specified file
		// in write-only mode, and if there is a problem
		// creating the file, then return an error.

		if fle, err = os.OpenFile(args[0], os.O_RDONLY, 0644); err != nil {
			log.Fatalln("Import failed - please check the filepath and try again.")
			return
		}

		// Ensure that we properly close the file handle
		// when we have finished with the file so that
		// the file descriptor is released.

		defer fle.Close()

		// Configure the export connection endpoint url
		// and specify the authentication header using
		// basic auth for root login.

		url := fmt.Sprintf("http://%s@%s:%s/import", opts.Auth.Auth, opts.DB.Host, opts.DB.Port)

		// Create a new http request object that we
		// can use to connect to the import endpoint
		// using a POST http request type.

		if req, err = http.NewRequest("POST", url, fle); err != nil {
			log.Fatalln("Connection failed - check the connection details and try again.")
			return
		}

		// Specify that the request is an octet stream
		// so that we can stream the file contents to
		// the server without reading the whole file.

		req.Header.Set("Content-Type", "application/octet-stream")

		// Attempt to dial the import endpoint and
		// if there is an error then stop execution
		// and return the connection error.

		if res, err = http.DefaultClient.Do(req); err != nil {
			log.Fatalln("Connection failed - check the connection details and try again.")
			return
		}

		// Ensure that we received a http 200 status
		// code back from the server, otherwise there
		// was a problem with our authentication.

		if res.StatusCode != 200 {
			log.Fatalln("Connection failed - check the connection details and try again.")
			return
		}

		return

	},
}

func init() {

	importCmd.PersistentFlags().StringVar(&opts.Auth.Auth, "auth", "root:root", "Master authentication details to use when connecting.")
	importCmd.PersistentFlags().StringVar(&opts.DB.Host, "host", "127.0.0.1", "Database server host to connect to.")
	importCmd.PersistentFlags().StringVar(&opts.DB.Port, "port", "8000", "Database server port to connect to.")

}
