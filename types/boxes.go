/**
 * @package go-imap-delete (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Jul 19, 2026; 8:32 PM
 */

package types

type boxAuth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}
type Box struct {
	Host   string  `yaml:"host"`
	Port   int     `yaml:"port"`
	Secure bool    `yaml:"secure"`
	Auth   boxAuth `yaml:"auth"`
}

type timeout struct {
	Connect int `yaml:"connect"`
	Command int `yaml:"command"`
	Delay   int `yaml:"delay"`
}

type Boxes map[string]Box

type Config struct {
	Timeout   timeout `yaml:"timeout"`
	Mailboxes Boxes   `yaml:"mailboxes"`
}
