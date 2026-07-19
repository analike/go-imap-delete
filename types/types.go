/**
 * @package go-imap-delete (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Jul 19, 2026; 1:42 PM
 */

package types

import "strings"

type Args struct {
	Box        string  `arg:"positional, required" help:"MailBox to operate on as defined in the yaml config"`
	Config     *string `arg:"--config" help:"Path to yaml config file" default:"./imap.yaml"`
	Folder     *string `arg:"--folder" help:"Folder to work with" default:"INBOX"`
	Uid        *string `arg:"--uid" help:"UID of message to delete"`
	Answered   *string `arg:"--answered" help:"filter answered/unanswered messages (true/false/yes/no)"`
	Deleted    *string `arg:"--deleted" help:"filter deleted/undeleted messages (true/false/yes/no)"`
	Draft      *string `arg:"--draft" help:"filter draft/undraft messages (true/false/yes/no)"`
	Flagged    *string `arg:"--flagged" help:"filter flagged/unflagged messages (true/false/yes/no)"`
	Seen       *string `arg:"--seen" help:"filter seen/unseen messages (true/false/yes/no)"`
	Subject    *string `arg:"--subject" help:"filter messages with SUBJECT in the message title"`
	Body       *string `arg:"--body" help:"filter messages with BODY in the message body"`
	Text       *string `arg:"--text" help:"filter messages with TEXT in their message header or body"`
	From       *string `arg:"--from" help:"filter messages with FROM headers contains sender"`
	To         *string `arg:"--to" help:"filter messages with TO in to headers"`
	Cc         *string `arg:"--cc" help:"filter messages with CC in cc headers"`
	Bcc        *string `arg:"--bcc" help:"filter messages with BCC in bcc headers"`
	BeforeDate *string `arg:"--before" help:"filter messages received before BEFORE(date yyyy-mm-dd)"`
	SinceDate  *string `arg:"--since" help:"filter messages received on or after SINCE(date yyyy-mm-dd)"`
	Date       *string `arg:"--date" help:"filter messages received on DATE(date yyyy-mm-dd)"`
	Force      *bool   `arg:"--force" help:"force to run delete op when no filters are specified"` // cancel dry-run if no filter is provided and folder found
	DryRun     *bool   `arg:"--dry-run" help:"list found emails without deleting"`
}

func parseBool(val string) *bool {
	switch strings.ToLower(val) {
	case "true", "yes", "":
		return new(true)
	case "false", "no":
		return new(false)
	}
	return nil
}

func (*Args) Version() string {
	return "imap-delete 1.0.0"
}

func (a *Args) HasAnswered() bool {
	return a.Answered != nil && parseBool(*a.Answered) != nil
}

func (a *Args) GetAnswered() bool {
	return *parseBool(*a.Answered)
}

func (a *Args) HasDeleted() bool {
	return a.Deleted != nil && parseBool(*a.Deleted) != nil
}

func (a *Args) GetDeleted() bool {
	return *parseBool(*a.Deleted)
}

func (a *Args) HasDraft() bool {
	return a.Draft != nil && parseBool(*a.Draft) != nil
}

func (a *Args) GetDraft() bool {
	return *parseBool(*a.Draft)
}

func (a *Args) HasFlagged() bool {
	return a.Flagged != nil && parseBool(*a.Flagged) != nil
}

func (a *Args) GetFlagged() bool {
	return *parseBool(*a.Flagged)
}

func (a *Args) HasSeen() bool {
	return a.Seen != nil && parseBool(*a.Seen) != nil
}

func (a *Args) GetSeen() bool {
	return *parseBool(*a.Seen)
}

func (a *Args) HasBeforeDate() bool {
	return a.BeforeDate != nil
}

func (a *Args) HasAfterDate() bool {
	return a.SinceDate != nil
}

func (a *Args) HasDate() bool {
	return a.Date != nil
}
