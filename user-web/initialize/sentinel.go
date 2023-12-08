package initialize

import (
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"

	sentinel "github.com/alibaba/sentinel-golang/api"
)

func InitializeSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "api-test",
			Threshold:              3,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       6000, // 6s <3
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
}
