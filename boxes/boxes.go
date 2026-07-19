/**
 * @package go-imap-delete (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Jul 19, 2026; 12:53 PM
 */

package boxes

import (
	"fmt"
	"imap-delete/args"
	"imap-delete/types"
	"imap-delete/utils"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/BrianLeishman/go-imap"
	"github.com/goccy/go-yaml"
	"go.analike.dev/logger"
)

type Client struct {
	client  *imap.Dialer
	folders []string
}

var (
	config = &types.Config{}
	boxes  *types.Boxes
)

func init() {
	parseConfig()
}

func Run() {
	var (
		ag  = &args.Args
		box = findBox(ag.Box)
	)
	if box == nil {
		logger.ErrorMsgf("mailbox `%s` not found in config file (%s)", ag.Box, *ag.Config)
		return
	}
	conn := login(box)
	if conn != nil {
		defer conn.client.Close()
		var (
			client    = conn.client
			boxPrefix = fmt.Sprintf("[%s/%s]", box.Host, box.Auth.User)
			selErr    = client.SelectFolder(*ag.Folder)
		)
		if selErr != nil {
			logger.ErrorMsgf("%s could not select folder `%s`: %v", boxPrefix, *ag.Folder, selErr)
			return
		}

		var (
			se, filCount = getSearchConfig(ag)
			agDry        = ag.DryRun != nil && *ag.DryRun == true
			seDry        = filCount < 1
			agForce      = ag.Force != nil && *ag.Force == true
			dryRun       = agDry
		)

		if !agDry && seDry && !agForce {
			logger.Infof("%s forcing dry-run as no filters supplied", boxPrefix)
			dryRun = true
		}

		ids, sEr := client.SearchUIDs(se)
		if sEr != nil {
			logger.ErrorMsgf("%s Could not search mailbox: %v", boxPrefix, sEr)
			return
		}

		if len(ids) < 1 {
			if filCount > 0 {
				logger.Infof("%s No messages match given filters", boxPrefix)
			} else {
				logger.Infof("%s No messages found in %s", boxPrefix, *ag.Folder)
			}
			return
		}

		logger.Infof("%s Found %d messages", boxPrefix, len(ids))

		overviews, ovErr := client.GetOverviews(ids...)
		if ovErr != nil {
			logger.ErrorMsgf("%s Could not get message info for found items: %v", boxPrefix, ovErr)
			return
		}
		var (
			total = len(overviews)
			i     = 1
			toDel = 0
			wait  = config.Timeout.Delay
		)
		for id, mail := range overviews {
			var (
				pos = utils.FormatIndex(i, total)
				sdr = utils.GetFirstEmail(mail.From)
				sub = mail.Subject
				dt  = mail.Received.Format(time.RFC1123)
			)
			iPrefix := fmt.Sprintf("%s %s: %s <%s>", pos, sdr, sub, dt)
			if dryRun {
				logger.Infof("%s %s", boxPrefix, iPrefix)
			} else {
				delEr := client.DeleteEmail(id)
				if delEr == nil {
					logger.Successf("%s %s", boxPrefix, iPrefix)
					toDel++
				} else {
					logger.ErrorMsgf("%s %s [%v]", boxPrefix, iPrefix, delEr)
				}
			}
			i++
			if wait > 0 {
				time.Sleep(time.Duration(wait) * time.Millisecond)
			}
		}
		if toDel > 0 {
			expEr := client.Expunge()
			if expEr != nil {
				logger.ErrorMsgf("%s Could not expunge deleted items: %v", boxPrefix, expEr)
			} else {
				logger.Successf("%s Expunged %d messages", boxPrefix, toDel)
			}
		}
	}
}

func login(box *types.Box) *Client {
	tOut := &config.Timeout
	imap.DialTimeout = time.Duration(tOut.Connect) * time.Second
	imap.CommandTimeout = time.Duration(tOut.Command) * time.Second
	imap.TLSSkipVerify = !box.Secure

	d, err := imap.New(box.Auth.User, box.Auth.Pass, box.Host, box.Port)
	prefix := fmt.Sprintf("[%s:%d/%s]", box.Host, box.Port, box.Auth.User)
	if err != nil {
		logger.Fatalf("%s Failed to login to: %v", prefix, err)
	}
	var folders []string
	folders, err = d.GetFolders()
	if err != nil {
		logger.Fatalf("%s Failed to get folders: %v", prefix, err)
	}
	logger.Infof("%s connected; found %d folder(s) [%s]", prefix, len(folders), strings.Join(folders, ", "))
	workDir := *args.Args.Folder
	if strings.ToLower(workDir) != "inbox" && !slices.Contains(folders, workDir) {
		logger.Fatalf("%s folder `%s` not found", prefix, workDir)
	}
	return &Client{
		client:  d,
		folders: folders,
	}
}

func getSearchConfig(arg *types.Args) (*imap.SearchBuilder, int) {
	build := imap.Search()
	count := 0
	if arg.Uid != nil {
		build.UID(*arg.Uid)
		count++
	}
	if arg.HasAnswered() {
		if arg.GetAnswered() {
			build.Answered()
		} else {
			build.Unanswered()
		}
		count++
	}
	if arg.HasDeleted() {
		if arg.GetDeleted() {
			build.Deleted()
		} else {
			build.Undeleted()
		}
		count++
	}
	if arg.HasDraft() {
		if arg.GetDraft() {
			build.Draft()
		} else {
			build.Undraft()
		}
		count++
	}
	if arg.HasFlagged() {
		if arg.GetFlagged() {
			build.Flagged()
		} else {
			build.Unflagged()
		}
		count++
	}
	if arg.HasSeen() {
		if arg.GetSeen() {
			build.Seen()
		} else {
			build.Unseen()
		}
		count++
	}
	if arg.Subject != nil {
		build.Subject(*arg.Subject)
		count++
	}
	if arg.Body != nil {
		build.Body(*arg.Body)
		count++
	}
	if arg.Text != nil {
		build.Text(*arg.Text)
		count++
	}
	if arg.From != nil {
		build.From(*arg.From)
		count++
	}
	if arg.To != nil {
		build.To(*arg.To)
		count++
	}
	if arg.Cc != nil {
		build.CC(*arg.Cc)
		count++
	}
	if arg.Bcc != nil {
		build.BCC(*arg.Bcc)
		count++
	}
	if arg.BeforeDate != nil {
		dt, _ := time.Parse(time.DateOnly, *arg.BeforeDate)
		build.Before(dt)
		count++
	}
	if arg.SinceDate != nil {
		dt, _ := time.Parse(time.DateOnly, *arg.SinceDate)
		build.Since(dt)
		count++
	}
	if arg.Date != nil {
		dt, _ := time.Parse(time.DateOnly, *arg.Date)
		build.On(dt)
		count++
	}

	return build, count
}

func findBox(name string) *types.Box {
	for key, box := range *boxes {
		if key == name {
			return &box
		}
	}
	return nil
}

func parseConfig() {
	var (
		confFile = *args.Args.Config
	)

	if !utils.FileExists(confFile) {
		logger.Fatalf("Config file does not exists: %s", confFile)
	}
	confBytes, err := os.ReadFile(confFile)
	if err != nil {
		logger.Fatalf("Could not read from config file %v", err)
	}
	err = yaml.Unmarshal(confBytes, config)
	if err != nil {
		logger.Fatalf("Could not parse config file %v", err)
	}
	if len(config.Mailboxes) < 1 {
		logger.Fatal("No mailbox found in config file")
	}
	boxes = &config.Mailboxes
}
