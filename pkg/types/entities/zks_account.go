package entities

const ZKSCurveBN254 = "bn254"
const ZKSAlgorithmEDDSA = "eddsa"

type ZKSAccount struct {
	Curve            string `json:"curve"`
	SigningAlgorithm string `json:"signingAlgorithm"`
	PublicKey        string `json:"publicKey"`
	Namespace        string `json:"namespace,omitempty"`
}
