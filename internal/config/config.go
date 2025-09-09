/*
   Copyright 2025 Tom Craggs

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

package config

type DevnoteConfig struct {
	// note settings
	NotesPath     string `yaml:"notes_path"`     //devnotes
	FileExtension string `yaml:"file_extension"` //.md

	// editor settings
	Editor     string   `yaml:"editor"`      // auto
	EditorArgs []string `yaml:"editor_args"` // []

	// behaviour
	OpenAfterNew bool `yaml:"open_after_new"` // true

	// git integration
	AutoCommit   bool    `yaml:"auto_commit"`    // false
	AutoAdd      bool    `yaml:"auto_add"`       // false
	GitCommitMsg *string `yaml:"git_commit_msg"` // null

}
