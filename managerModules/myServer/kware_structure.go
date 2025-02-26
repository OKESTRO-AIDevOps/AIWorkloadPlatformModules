package ystruct

type ReqResource struct {
	Version string  `json:"version,omitempty" yaml:"version,omitempty"`
	Request Request `json:"request,omitempty" yaml:"request,omitempty"`
}

type Request struct {
	Name       string      `json:"name,omitempty" yaml:"name,omitempty"`
	ID         string      `json:"id,omitempty" yaml:"id,omitempty"`
	Date       string      `json:"date,omitempty" yaml:"date,omitempty"`
	Containers []Container `json:"containers,omitempty" yaml:"containers,omitempty"`
	Attribute  Attribute   `json:"attribute,omitempty" yaml:"attribute,omitempty"`
}

type Attribute struct {
	WorkloadType     string      `json:"workloadType,omitempty" yaml:"workloadType,omitempty"`
	IsCronJob        bool        `json:"isCronJob,omitempty" yaml:"isCronJob,omitempty"`
	DevOpsType       string      `json:"devOpsType,omitempty" yaml:"devOpsType,omitempty"`
	CudaVersion      float64     `json:"cudaVersion,omitempty" yaml:"cudaVersion,omitempty"`
	GPUDriverVersion float64     `json:"gpuDriverVersion,omitempty" yaml:"gpuDriverVersion,omitempty"`
	WorkloadFeature  int      	 `json:"workloadFeature,omitempty" yaml:"workloadFeature,omitempty"`
	UserID           string      `json:"userId,omitempty" yaml:"userId,omitempty"`
	Checkpoint       bool		 `json:"checkpoint" yaml:"checkpoint"`
	Yaml             string      `json:"yaml,omitempty" yaml:"yaml,omitempty"`
	Dag              interface{} `json:"dag,omitempty" yaml:"dag,omitempty"`
}
