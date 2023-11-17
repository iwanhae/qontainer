// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dns

import (
	"time"
)

var (
	defaultNS = []string{"127.0.0.1", "[::1]"}
)

type DnsConfig struct {
	Servers       []string      // server addresses (in host:port form) to use
	Search        []string      // rooted suffixes to append to local name
	ndots         int           // number of dots in name to trigger absolute lookup
	timeout       time.Duration // wait before giving up on a query, including retries
	attempts      int           // lost packets before giving up on server
	rotate        bool          // round robin among servers
	unknownOpt    bool          // anything unknown was encountered
	lookup        []string      // OpenBSD top-level database "lookup" order
	err           error         // any error that occurs during open of resolv.conf
	mtime         time.Time     // time of resolv.conf modification
	singleRequest bool          // use sequential A and AAAA queries instead of parallel queries
	useTCP        bool          // force usage of TCP for DNS resolutions
	trustAD       bool          // add AD flag to queries
	noReload      bool          // do not check for config file updates
}
