//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package editor

import (
	"github.com/cloustone/pandas/apimachinery/models"
	"github.com/cloustone/pandas/pkg/auth"
)

type Manager interface {
	// Widget
	// GetWigets return model's widget used by web console to construct
	// compound device model. There are no special data store to manage widget,
	// Widget is created dynamically using device model and categorized by
	// domain. Preset models are system widgets and user's modes are
	// customization widgets. (in futer, all widgets will be stored into db)
	GetWidgets(principal auth.Principal) ([]models.Widget, error)
}
