package nanoleaf

import "encoding/json"

// PanelUpdate contains the update for a given panel event.
// typeID must be set before deserializing the update as JSON
type PanelUpdate struct {
	TypeID int

	// State will be populated if typeID is set to 1
	State *PanelState

	// Layout will be populated if typeID is set to 2
	Layout *PanelLayout

	// Effect will be populated if typeID is set to 3
	Effect *PanelEffect

	// Gestures will be nonempty if typeID is set to 4
	Gestures []Gesture
}

// UnmarshalJSON allows us to decode the contents of the event
func (pu *PanelUpdate) UnmarshalJSON(b []byte) error {
	// For some reason type 4 is serialized differently...
	// Quickly handle it and then proceed if not the proper type
	if pu.TypeID == 4 {
		var tmp struct {
			Gestures []Gesture `json:"events"`
		}

		err := json.Unmarshal(b, &tmp)
		if err != nil {
			return err
		}

		pu.Gestures = tmp.Gestures
		return nil
	}

	var tmp struct {
		Events []avp `json:"events"`
	}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	switch pu.TypeID {
	case 1:
		state := &PanelState{}
		for _, event := range tmp.Events {
			switch event.Attribute {
			case 1:
				state.On = &BoolValue{}
				if err = json.Unmarshal(event.Value, &state.On.Value); err != nil {
					return err
				}
			case 2:
				state.Brightness = &IntRangeValue{}
				if err = json.Unmarshal(event.Value, &state.Brightness.Value); err != nil {
					return err
				}
			case 3:
				state.Hue = &IntRangeValue{}
				if err = json.Unmarshal(event.Value, &state.Hue.Value); err != nil {
					return err
				}
			case 4:
				state.Saturation = &IntRangeValue{}
				if err = json.Unmarshal(event.Value, &state.Saturation.Value); err != nil {
					return err
				}
			case 5:
				state.CT = &IntRangeValue{}
				if err = json.Unmarshal(event.Value, &state.CT.Value); err != nil {
					return err
				}
			case 6:
				if err = json.Unmarshal(event.Value, &state.ColorMode); err != nil {
					return err
				}
			}
		}
		pu.State = state
	case 2:
		layout := &PanelLayout{}
		for _, event := range tmp.Events {
			switch event.Attribute {
			case 1:
				if err = json.Unmarshal(event.Value, &layout.Panels); err != nil {
					return err
				}
			case 2:
				if err = json.Unmarshal(event.Value, &layout.Orientation); err != nil {
					return err
				}
			}
		}
		pu.Layout = layout
	case 3:
		effect := &PanelEffect{}
		for _, event := range tmp.Events {
			switch event.Attribute {
			case 1:
				if err = json.Unmarshal(event.Value, &effect.Current); err != nil {
					return err
				}
			}
		}
		pu.Effect = effect
	}

	return nil
}

// Gesture represents a detected touch event on supported panels.
type Gesture struct {
	GestureType int `json:"gesture"`
	// PanelID may be set to -1 if the specified gesture can't be targeted to a panel
	PanelID int `json:"panelId"`
}

type avp struct {
	Attribute int             `json:"attr"`
	Value     json.RawMessage `json:"Value"`
}
