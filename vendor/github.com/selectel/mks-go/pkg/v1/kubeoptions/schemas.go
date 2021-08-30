package kubeoptions

// View represents list of feature-gates/admission-controllers by kubernetes version.
type View struct {
	// KubeVersion represents the Kubernetes minor version in format: "X.Y".
	KubeVersion string `json:"KubeVersionMinor"`

	// Names represents list of feature-gate names.
	Names []string `json:"Names"`
}
