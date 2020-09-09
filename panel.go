package nanoleaf

import "context"

// GetPanel retrieves the panel details
func (c *Client) GetPanel(ctx context.Context) (*LightPanel, error) {
	panel := &LightPanel{}
	err := c.get(ctx, "", panel)
	if err != nil {
		return nil, err
	}

	return panel, nil
}

// SetOn sets the panel to either be on or off
func (c *Client) SetOn(ctx context.Context, on bool) error {
	var req struct {
		On struct {
			Value bool `json:"value"`
		} `json:"on"`
	}

	req.On.Value = on
	return c.put(ctx, "state", req, nil)
}

// SetScene selects the specified scene name on the panel and applies it
func (c *Client) SetScene(ctx context.Context, sceneName string) error {
	var req struct {
		SceneName string `json:"select"`
	}

	req.SceneName = sceneName
	return c.put(ctx, "effects", req, nil)
}

// SetOrientation updates the orientation of the panels (orientation is 360 degrees - from 0 to 359)
func (c *Client) SetOrientation(ctx context.Context, orientation int) error {
	var req struct {
		Orientation struct {
			Value int `json:"value"`
		} `json:"globalOrientation"`
	}

	req.Orientation.Value = orientation
	return c.put(ctx, "panelLayout", req, nil)
}

// SetBrightness will set the brightness level and (optionally) duration in seconds
func (c *Client) SetBrightness(ctx context.Context, level int, duration int) error {
	var req struct {
		Brightness struct {
			Value    int `json:"value"`
			Duration int `json:"duration,omitempty"`
		} `json:"brightness"`
	}

	req.Brightness.Value = level
	req.Brightness.Duration = duration

	return c.put(ctx, "state", req, nil)
}

// IncrementBrightness will increment the brightness level. Both positive and negative values are supported.
func (c *Client) IncrementBrightness(ctx context.Context, amount int) error {
	var req struct {
		Brightness struct {
			Increment int `json:"increment"`
		} `json:"brightness"`
	}

	req.Brightness.Increment = amount

	return c.put(ctx, "state", req, nil)
}

// SetHue will set the hue of the light
func (c *Client) SetHue(ctx context.Context, hue int) error {
	var req struct {
		Hue struct {
			Value int `json:"value"`
		} `json:"hue"`
	}

	req.Hue.Value = hue

	return c.put(ctx, "state", req, nil)
}

// IncrementHue will increment the hue of the light. Both positive and negative values are supported.
func (c *Client) IncrementHue(ctx context.Context, amount int) error {
	var req struct {
		Hue struct {
			Increment int `json:"increment"`
		} `json:"hue"`
	}

	req.Hue.Increment = amount

	return c.put(ctx, "state", req, nil)
}

// SetSaturation will set the saturation of the light
func (c *Client) SetSaturation(ctx context.Context, sat int) error {
	var req struct {
		Saturation struct {
			Value int `json:"value"`
		} `json:"sat"`
	}

	req.Saturation.Value = sat

	return c.put(ctx, "state", req, nil)
}

// IncrementSaturation will increment the saturation of the light. Both positive and negative values are supported.
func (c *Client) IncrementSaturation(ctx context.Context, amount int) error {
	var req struct {
		Saturation struct {
			Increment int `json:"increment"`
		} `json:"sat"`
	}

	req.Saturation.Increment = amount

	return c.put(ctx, "state", req, nil)
}

// SetCT will set the colour temperature of the light
func (c *Client) SetCT(ctx context.Context, ct int) error {
	var req struct {
		ColorTemperature struct {
			Value int `json:"value"`
		} `json:"ct"`
	}

	req.ColorTemperature.Value = ct

	return c.put(ctx, "state", req, nil)
}

// IncrementCT will increment the colour temperature of the light. Both positive and negative values are supported.
func (c *Client) IncrementCT(ctx context.Context, amount int) error {
	var req struct {
		ColorTemperature struct {
			Increment int `json:"increment"`
		} `json:"ct"`
	}

	req.ColorTemperature.Increment = amount

	return c.put(ctx, "state", req, nil)
}

// LightPanel represents the current state of a Nanoleaf Light Panel
type LightPanel struct {
	Name            string `json:"name"`
	SerialNumber    string `json:"serialNo"`
	Manufacturer    string `json:"manufacturer"`
	FirmwareVersion string `json:"firmwareVersion"`
	ModelNumber     string `json:"model"`

	State  PanelState  `json:"state"`
	Effect PanelEffect `json:"effects"`
	Layout struct {
		Orientation IntRangeValue `json:"globalOrientation"`
		Panels      PanelLayout   `json:"layout"`
	} `json:"panelLayout"`
	Rhythm Rhythm `json:"rhythm"`
}

// PanelState contains the current state of the light panel
type PanelState struct {
	On         BoolValue     `json:"on"`
	Brightness IntRangeValue `json:"brightness"`
	Hue        IntRangeValue `json:"hue"`
	Saturation IntRangeValue `json:"sat"`
	CT         IntRangeValue `json:"ct"`
	ColorMode  string        `json:"colorMode"`
}

// PanelEffect represents the current and possible set of effects on this light panel
type PanelEffect struct {
	Current string   `json:"select"`
	Options []string `json:"effectsList"`
}

// BoolValue contains a serialized boolean
type BoolValue struct {
	Value bool `json:"value"`
}

// IntRangeValue contains a serialized integer value with its respective max/min values
type IntRangeValue struct {
	Value int `json:"value"`
	Max   int `json:"max"`
	Min   int `json:"min"`
}

// PanelLayout represents the layout of all the panels making up this light
type PanelLayout struct {
	PanelCount int             `json:"numPanels"`
	SideLength int             `json:"sideLength"`
	Panels     []PanelPosition `json:"positionData"`
}

// PanelPosition contains information about the relative layout of a single panel
type PanelPosition struct {
	PanelID     int `json:"panelId"`
	X           int `json:"x"`
	Y           int `json:"y"`
	Orientation int `json:"o"`
	Type        int `json:"shapeType"`
}

// Rhythm contains information about the Rhythm module
type Rhythm struct {
	Connected       bool   `json:"rhythmConnected"`
	Active          bool   `json:"rhythmActive"`
	ID              int    `json:"rhythmId"`
	HardwareVersion string `json:"hardwareVersion"`
	FirmwareVersion string `json:"firmwareVersion"`
	AuxAvailable    bool   `json:"auxAvailable"`
	Mode            int    `json:"rhythmMode"`
	// Position will only have the x, y and o fields filled in
	Position PanelPosition `json:"rhythmPos"`
}
