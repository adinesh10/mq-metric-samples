package main

/*
  Copyright (c) IBM Corporation 2016

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific

   Contributors:
     Mark Taylor - Initial Contribution
*/

import (
	"flag"
	"os"

	"github.com/ibm-messaging/mq-golang/mqmetric"

	log "github.com/sirupsen/logrus"
)

type mqTTYConfig struct {
	qMgrName            string
	replyQ              string
	monitoredQueues     string
	monitoredQueuesFile string
	hostname            string
	hostlabel           string // Used in the output string

	cc mqmetric.ConnectionConfig

	interval string

	logLevel string
}

var config mqTTYConfig

/*
initConfig parses the command line parameters.
*/
func initConfig() {

	flag.StringVar(&config.qMgrName, "ibmmq.queueManager", "", "Queue Manager name")
	flag.StringVar(&config.replyQ, "ibmmq.replyQueue", "SYSTEM.DEFAULT.MODEL.QUEUE", "Reply Queue to collect data")
	flag.StringVar(&config.monitoredQueues, "ibmmq.monitoredQueues", "", "Patterns of queues to monitor")
	flag.StringVar(&config.monitoredQueuesFile, "ibmmq.monitoredQueuesFile", "", "File with patterns of queues to monitor")

	flag.StringVar(&config.interval, "ibmmq.interval", "10", "How many seconds between each collection")
	flag.StringVar(&config.hostname, "ibmmq.hostname", "localhost", "Host to connect to")

	flag.BoolVar(&config.cc.ClientMode, "ibmmq.client", false, "Connect as MQ client")

	flag.StringVar(&config.logLevel, "log.level", "error", "Log level - debug, info, error")

	flag.Parse()

	if config.monitoredQueuesFile != "" {
		var err error
		config.monitoredQueues, err = mqmetric.ReadPatterns(config.monitoredQueuesFile)
		if err != nil {
			log.Errorf("Failed to parse monitored queues file - %v", err)
		}
	}

	// Don't want to use "localhost" as the tag in the metric printing
	if config.hostname == "localhost" {
		config.hostlabel, _ = os.Hostname()
	} else {
		config.hostlabel = config.hostname
	}
}
