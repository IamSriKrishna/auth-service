package output

type CompatibilityOutput struct {
	IsCompatible     bool   `json:"is_compatible"`
	BottleNeckFinish string `json:"bottle_neck_finish"`
	CapNeckFinish    string `json:"cap_neck_finish"`
	Message          string `json:"message"`
}
