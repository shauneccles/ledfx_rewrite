package virtual

import (
	"errors"
	"fmt"
	"ledfx/color"
	"ledfx/config"
	"ledfx/device"
	"ledfx/effect"
	log "ledfx/logger"
)

// TODO: this should belong to the virtual instance
var done chan bool

func PlayVirtual(virtualID string, playState bool, clr string, effectType string) (err error) {
	fmt.Println("Set PlayState of ", virtualID, " to ", playState)

	if virtualID == "" {
		return errors.New("virtual id is empty. Please provide Id to add virtual to config")
	}

	var virtualExists bool

	newColor, err := color.NewColor(clr)
	if err != nil {
		return fmt.Errorf("error generating new color: %w", err)
	}

	for i, d := range config.GlobalConfig.Virtuals {
		if d.Id == virtualID {
			if effectType != "" {
				d.Effect.Type = effectType
			}
			if clr != "" {
				d.Effect.Config.Color = clr
			}
			virtualExists = true
			config.GlobalConfig.Virtuals[i].Active = playState
			if config.GlobalConfig.Virtuals[i].IsDevice != "" {
				for in, de := range config.GlobalConfig.Devices {
					if de.Id == config.GlobalConfig.Virtuals[i].IsDevice {
						var dev = &device.UDPDevice{
							Name:     config.GlobalConfig.Devices[in].Config.Name,
							Port:     config.GlobalConfig.Devices[in].Config.Port,
							Protocol: device.UDPProtocols[config.GlobalConfig.Devices[in].Config.UdpPacketType],
							Config:   config.GlobalConfig.Devices[in].Config,
						}

						data := make([]color.Color, de.Config.PixelCount)
						for i2 := range data {
							data[i2] = newColor
						}

						if err := dev.Init(); err != nil {
							return fmt.Errorf("error initializing dev: %w", err)
						}

						var timeout byte
						if playState {
							timeout = 0xff
						} else {
							timeout = 0x00
						}

						if err := dev.SendData(data, timeout); err != nil {
							return fmt.Errorf("error sending data to WLED: %w", err)
						}
					}
				}
			}
		}
	}

	if !virtualExists {
		return fmt.Errorf("virtual with ID %q does not exist", virtualID)
	}

	config.GlobalViper.Set("virtuals", config.GlobalConfig.Virtuals)
	return config.GlobalViper.WriteConfig()
}

func StopVirtual(virtualid string) (err error) {
	log.Logger.WithField("category", "Virtual Stopper").Infof("CLEAR EFFECT OF: %s", virtualid)

	if virtualid == "" {
		return errors.New("virtual id is empty, please provide id to add virtual to config")
	}

	for i, d := range config.GlobalConfig.Virtuals {
		if d.Id == virtualid {
			if config.GlobalConfig.Virtuals[i].IsDevice != "" {
				log.Logger.WithField("category", "Virtual Stopper").Warnf("WTC Clear Effect of %s", config.GlobalConfig.Virtuals[i].Effect.Name)
				for in, de := range config.GlobalConfig.Devices {
					if de.Id == config.GlobalConfig.Virtuals[i].IsDevice {
						go func() {
							var currentEffect effect.Effect = &effect.PulsingEffect{}
							if err := effect.StopEffect(config.GlobalConfig.Devices[in].Config, currentEffect, "#000000", 60, done); err != nil {
								log.Logger.WithField("category", "Virtual Stopper").Warnf("Error stopping current effect: %v", err)
							}
						}()
					}
				}
			}
		}
	}
	return
}

// LoadVirtuals loads the virtuals from the config file and plays any effects that are active on them
func LoadVirtuals() (err error) {
	// TODO: load all virtuals from config

	// c := &config.GlobalConfig
	// v := config.GlobalViper

	// for i, virtualConfig := range config.GlobalConfig.Virtuals {
	// 	if config.GlobalConfig.Virtuals[i].Active == true {
	// 		if config.GlobalConfig.Virtuals[i].IsDevice != "" {
	//       // TODO: instantiate a virtual
	// 			// PlayVirtual(virtual)
	// 		}
	// 	}
	// }

	return nil
}
