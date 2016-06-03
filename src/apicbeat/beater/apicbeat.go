package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	
	"apicbeat/apic"
)

type Apicbeat struct {
	period time.Duration
	Session *apic.Session
	
	AbConfig ConfigSettings
	events publisher.Client
	
	done chan struct{}
}

func New() *Apicbeat {
	return &Apicbeat{
	}
}

func (ab *Apicbeat) Config(b *beat.Beat) error {
	
	err := b.RawConfig.Unpack(&ab.AbConfig)
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}
	
	if ab.AbConfig.Apicbeat != nil && ab.AbConfig.Input != nil {
		return fmt.Errorf("'apicbeat' and 'input' are both set in config. Only " + "one can be enabled so use 'apicbeat'. ")
	}
	
	if ab.AbConfig.Input != nil {
		logp.Warn(" Use 'apicbeat' instead.")
		ab.AbConfig.Apicbeat = ab.AbConfig.Input
	}
	
	apicbeatConfig := ab.AbConfig.Apicbeat
	
	if apicbeatConfig.Period != nil {
		ab.period = time.Duration(*apicbeatConfig.Period) * time.Second
	} else {
		ab.period = 10 * time.Second
	}
	
	if apicbeatConfig.Addr != nil && apicbeatConfig.User != nil && apicbeatConfig.Passwd != nil {
		ab.Session = apic.NewSession(*apicbeatConfig.Addr, *apicbeatConfig.User, *apicbeatConfig.Passwd)
	}
	
	logp.Debug("apicbeat", "Init apicbeat")
	logp.Debug("apicbeat", "Period %v", ab.period)
	
	return nil
}

func (ab *Apicbeat) Setup(b *beat.Beat) error {
	ab.events = b.Publisher.Connect()
	ab.done = make(chan struct{})
	return nil
}

func (ab *Apicbeat) Run(b *beat.Beat) error {
	
	ticker := time.NewTicker(ab.period)
	defer ticker.Stop()
	
	for {
		select {
		case <-ab.done:
			return nil
		case <-ticker.C:
		}
		
		timerStart := time.Now()
		
		if *ab.AbConfig.Apicbeat.ForwardSet.Tennant_Endpoint == true {
			events, err := apic.GetFvTenantEndPoints(ab.Session); if err == nil { ab.events.PublishEvents(events) }
		}
		if *ab.AbConfig.Apicbeat.ForwardSet.Tennant_Endpoint_DN == true {
			events, err := apic.GetFvTenantEndPointDNs(ab.Session); if err == nil { ab.events.PublishEvents(events) }
		}
		if *ab.AbConfig.Apicbeat.ForwardSet.Tennant_Health == true {
			events, err := apic.GetFvTenantHealths(ab.Session); if err == nil { ab.events.PublishEvents(events) }
		}
		if *ab.AbConfig.Apicbeat.ForwardSet.Tennant_Health_Cur == true {
			events, err := apic.GetFvTenantHealthCurs(ab.Session); if err == nil { ab.events.PublishEvents(events) }
		}
		if *ab.AbConfig.Apicbeat.ForwardSet.Fault_Info == true {
			events, err := apic.GetFaultInfos(ab.Session); if err == nil { ab.events.PublishEvents(events) }
		}
		
		timerEnd := time.Now()
		duration := timerEnd.Sub(timerStart)
		if duration.Nanoseconds() > ab.period.Nanoseconds() {
			logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
		}
	}
	
	return nil
}

func (ab *Apicbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (ab *Apicbeat) Stop() {
	logp.Info("Send stop signal to apicbeat main loop")
	close(ab.done)
	ab.events.Close()
}
