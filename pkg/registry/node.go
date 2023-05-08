/*
 * Copyright 2023 CoreLayer BV
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package registry

import "github.com/corelayer/netscaleradc-nitro-go/pkg/nitro"

type Node struct {
	Name    string `json:"name" yaml:"name"`
	Address string `json:"address" yaml:"address"`
}

func (n Node) GetNodeUrl(s Settings) string {
	return s.GetUrlScheme() + n.Address
}

func (n Node) GetNitroSettings(s Settings, c Credentials) nitro.Settings {
	return nitro.Settings{
		Name:                      n.Name,
		BaseUrl:                   n.GetNodeUrl(s),
		Username:                  c.Username,
		Password:                  c.Password,
		AutoLogin:                 s.AutoLogin,
		Timeout:                   s.Timeout,
		ValidateServerCertificate: s.ValidateServerCertificate,
		LogTlsSecrets:             s.LogTlsSecrets,
		LogTlsSecretsDestination:  s.LogTlsSecretsDestination,
	}
}
