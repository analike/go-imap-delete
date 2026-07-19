/**
 * @package go-imap-delete (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Jul 19, 2026; 12:56 PM
 */

package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BrianLeishman/go-imap"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func FormatIndex(cur, total int) string {
	totalStr := strconv.Itoa(total)
	width := len(totalStr)
	curStr := fmt.Sprintf("%0*d", width, cur)

	return curStr + "/" + totalStr
}

func GetFirstEmail(ad imap.EmailAddresses) string {
	if len(ad) > 0 {
		for k := range ad {
			return k
		}
	}
	return ""
}
