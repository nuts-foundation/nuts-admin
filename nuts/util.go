package nuts

import (
	"errors"
	"net"
)

var ErrNutsNodeUnreachable = errors.New("nuts node unreachable")

func UnwrapAPIError(err error) error {
	if _, ok := err.(net.Error); ok {
		return errors.Join(ErrNutsNodeUnreachable, err)
	}
	return err
}
