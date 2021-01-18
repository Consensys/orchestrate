package entities

const ZKSCurveBN256 = "bn256"
const ZKSAlgorithmEDDSA = "eddsa"

type ZKSAccount struct {
	Curve            string `json:"curve"`
	SigningAlgorithm string `json:"signingAlgorithm"`
	PublicKey        string `json:"publicKey"`
	Namespace        string `json:"namespace,omitempty"`
}
