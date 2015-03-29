package netutil

import (
	"io"
)

func WriteAllTo(w io.Writer, wt ...io.WriterTo) (int64, error) {
	var n int64

	for _, wrt := range wt {
		nn, err := wrt.WriteTo(w)
		if n += nn; err != nil {
			return n, err
		}
	}

	return n, nil
}
