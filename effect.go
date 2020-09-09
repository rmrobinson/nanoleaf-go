package nanoleaf

import "context"

const (
	// WheelPluginUUID is the UUID of the wheel plugin
	WheelPluginUUID = "6970681a-20b5-4c5e-8813-bdaebc4ee4fa"
	// FlowPluginUUID is the UUID of the flow plugin
	FlowPluginUUID = "027842e4-e1d6-4a4c-a731-be74a1ebd4cf"
	// ExplodePluginUUID is the UUID of the explode plugin
	ExplodePluginUUID = "713518c1-d560-47db-8991-de780af71d1e"
	// FadePluginUUID is the UUID of the fade plugin
	FadePluginUUID = "b3fd723a-aae8-4c99-bf2b-087159e0ef53"
	// RandomPluginUUID is the UUID of the random plugin
	RandomPluginUUID = "ba632d3e-9c2b-4413-a965-510c839b3f71"
	// HighlightPluginUUID is the UUID of the highlight plugin
	HighlightPluginUUID = "70b7c636-6bf8-491f-89c1-f4103508d642"
)

// GetEffect retrieves the specified effect
func (c *Client) GetEffect(ctx context.Context, effectName string) (*Effect, error) {
	var req struct {
		Body struct {
			Command       string `json:"command"`
			AnimationName string `json:"animName"`
		} `json:"write"`
	}

	req.Body.Command = "request"
	req.Body.AnimationName = effectName

	resp := &Effect{}
	err := c.put(ctx, "effects", req, resp)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetEffects retrieves all the configured effects
func (c *Client) GetEffects(ctx context.Context) ([]Effect, error) {
	var req struct {
		Body struct {
			Command string `json:"command"`
		} `json:"write"`
	}

	req.Body.Command = "requestAll"

	var resp struct {
		Effects []Effect `json:"animations"`
	}
	err := c.put(ctx, "effects", req, &resp)

	if err != nil {
		return nil, err
	}
	return resp.Effects, nil
}

// Effect represents a single effect in the panel
type Effect struct {
	Name            string `json:"animName"`
	Version         string `json:"version"`
	PluginType      string `json:"pluginType"`
	PluginUUID      string `json:"pluginUuid"`
	Palette         []HSB  `json:"palette"`
	BrightnessRange MaxMin `json:"brightnessRange"`
	TransitionTime  MaxMin `json:"transTime"`
	DelayTime       MaxMin `json:"delayTime"`
	ColorType       string `json:"colorType"`
	AnimationType   string `json:"animType"`
	FlowFactor      int    `json:"flowFactor"`
	ExplodeFactor   int    `json:"explodeFactor"`
	WindowSize      int    `json:"windowSize"`
	Direction       string `json:"direction"`
	Loop            bool   `json:"loop"`
}

// HSB represents a hue/saturation/brightness entry
type HSB struct {
	// Hue is a 0-359 value
	Hue int `json:"hue"`
	// Saturation is a 0-100 value
	Saturation int `json:"saturation"`
	// Brightness is a 0-100 value
	Brightness int `json:"brightness"`
	// Probability reflects the chance the above HSB value will apply
	Probability float64 `json:"probability"`
}

// MaxMin represents a pair of values for max and min
type MaxMin struct {
	Maximum int `json:"maxValue"`
	Minimum int `json:"minValue"`
}
