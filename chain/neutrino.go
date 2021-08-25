package chain

import (
	"github.com/btcsuite/btcd/txscript"
)

// buildFilterBlocksWatchList constructs a watchlist used for matching against a
// cfilter from a FilterBlocksRequest. The watchlist will be populated with all
// external addresses, internal addresses, and outpoints contained in the
// request.
func buildFilterBlocksWatchList(req *FilterBlocksRequest) ([][]byte, error) {
	// Construct a watch list containing the script addresses of all
	// internal and external addresses that were requested, in addition to
	// the set of outpoints currently being watched.
	watchListSize := len(req.ExternalAddrs) +
		len(req.InternalAddrs) +
		len(req.WatchedOutPoints)

	watchList := make([][]byte, 0, watchListSize)

	for _, addr := range req.ExternalAddrs {
		p2shAddr, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, err
		}

		watchList = append(watchList, p2shAddr)
	}

	for _, addr := range req.InternalAddrs {
		p2shAddr, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, err
		}

		watchList = append(watchList, p2shAddr)
	}

	for _, addr := range req.WatchedOutPoints {
		addr, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, err
		}

		watchList = append(watchList, addr)
	}

	return watchList, nil
}
