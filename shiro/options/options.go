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
package options

import (
	broadcast_options "github.com/cloustone/pandas/pkg/broadcast"
	genericoptions "github.com/cloustone/pandas/pkg/server/options"
	"github.com/cloustone/pandas/pkg/sms"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/xid"
	"github.com/spf13/pflag"
)

type ServingOptions struct {
	SecureServing    *genericoptions.SecureServingOptions
	BroadcastServing *broadcast_options.ServingOptions
	ServiceID        string
	InitFile         string
	RealmConfigFile  string
	MFA              string
	SmsOptions       *sms.ServingOptions
	RolesFile        string
	BackstorePath    string
}

func NewServingOptions() *ServingOptions {
	s := ServingOptions{
		SecureServing:    genericoptions.NewSecureServingOptions("dmms"),
		BroadcastServing: broadcast_options.NewServingOptions(),
		SmsOptions:       sms.NewServingOptions(),
		ServiceID:        xid.New().String(),
		RealmConfigFile:  "./shiro/shiro-realms.json",
		RolesFile:        "./shiro/shiro-roles.json",
		//BackstorePath:    "localhost:27017",
		BackstorePath: "sqlite3",
	}
	return &s
}

func (s *ServingOptions) AddFlags(fs *pflag.FlagSet) {
	s.BroadcastServing.AddFlags(fs)
	s.SmsOptions.AddFlags(fs)
	fs.StringVar(&s.ServiceID, "service-id", s.ServiceID, "shiro service ID")
	fs.StringVar(&s.InitFile, "init-file", s.InitFile, "initial shiro config file")
	fs.StringVar(&s.RealmConfigFile, "realm-config-file", s.RealmConfigFile, "realm config file")
	fs.StringVar(&s.MFA, "mfa", s.MFA, "multiple factor authentication")
	fs.StringVar(&s.RolesFile, "roles-file", s.RolesFile, "builtin roles definitions")
	fs.StringVar(&s.BackstorePath, "backstore-path", s.BackstorePath, "backstore path for postgres")
}
