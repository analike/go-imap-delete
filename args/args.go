/**
 * @package go-imap-delete (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Jul 19, 2026; 11:35 am
 */

package args

import (
	"fmt"
	"imap-delete/types"
	"time"

	"github.com/alexflint/go-arg"
)

var Args types.Args

func init() {
	Args = getArgs()
}

func checkDate(dt string) error {
	_, err := time.Parse(time.DateOnly, dt)
	return err
}

func getArgs() types.Args {
	a := types.Args{}
	p := arg.MustParse(&a)
	ap := &a
	if ap.HasBeforeDate() {
		e := checkDate(*ap.BeforeDate)
		if e != nil {
			p.Fail(fmt.Sprintf("invalid --before-date: %v", e))
		}
	}
	if ap.HasDate() {
		e := checkDate(*ap.Date)
		if e != nil {
			p.Fail(fmt.Sprintf("invalid --date: %v", e))
		}
	}
	if ap.HasAfterDate() {
		e := checkDate(*ap.SinceDate)
		if e != nil {
			p.Fail(fmt.Sprintf("invalid --after-date: %v", e))
		}
	}
	return a
}
