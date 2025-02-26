package ystruct

type Container struct {
	Name      string     `json:"name,omitempty" yaml:"name,omitempty"`
	Resources Resources  `json:"resources,omitempty" yaml:"resources,omitempty"`
	Attribute Attributes `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Cluster   string     `json:"cluster,omitempty" yaml:"cluster,omitempty"`
	Node      string     `json:"node,omitempty" yaml:"node,omitempty"`
}

type Resources struct {
	Requests ResourceDetails `json:"requests,omitempty" yaml:"requests,omitempty"`
	Limits   ResourceDetails `json:"limits,omitempty" yaml:"limits,omitempty"`
}

type Attributes struct {
	MaxReplicas            int `json:"maxReplicas,omitempty" yaml:"maxReplicas,omitempty"`
	TotalSize              int `json:"totalSize,omitempty" yaml:"totalSize,omitempty"`
	PredictedExecutionTime int `json:"predictedExecutionTime,omitempty" yaml:"predictedExecutionTime,omitempty"`
	Order                  int `json:"order,omitempty" yaml:"order,omitempty"`
}
