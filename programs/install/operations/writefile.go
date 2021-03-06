/*
 Copyright 2016 Padduck, LLC

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

package operations

import (
	"github.com/pufferpanel/pufferd/environments"
	"io/ioutil"
	"github.com/pufferpanel/pufferd/utils"
)

type WriteFile struct {
	TargetFile  string
	Environment environments.Environment
	Text string
}

func (c *WriteFile) Run() error {
	target := utils.JoinPath(c.Environment.GetRootDirectory(), c.TargetFile)
	ioutil.WriteFile(target, []byte(c.Text), 0644)
	return nil
}
