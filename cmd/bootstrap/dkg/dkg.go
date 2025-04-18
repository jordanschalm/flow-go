package dkg

import (
	"fmt"

	"github.com/onflow/crypto"

	model "github.com/onflow/flow-go/model/dkg"
	"github.com/onflow/flow-go/module/signature"
)

// RandomBeaconKG is centralized BLS threshold signature key generation.
func RandomBeaconKG(n int, seed []byte) (model.ThresholdKeySet, error) {

	if n == 1 {
		sk, pk, pkGroup, err := thresholdSignKeyGenOneNode(seed)
		if err != nil {
			return model.ThresholdKeySet{}, fmt.Errorf("Beacon KeyGen failed: %w", err)
		}

		dkgData := model.ThresholdKeySet{
			PrivKeyShares: sk,
			PubGroupKey:   pkGroup,
			PubKeyShares:  pk,
		}
		return dkgData, nil
	}

	skShares, pkShares, pkGroup, err := crypto.BLSThresholdKeyGen(n,
		signature.RandomBeaconThreshold(n), seed)
	if err != nil {
		return model.ThresholdKeySet{}, fmt.Errorf("Beacon KeyGen failed: %w", err)
	}

	randomBeaconData := model.ThresholdKeySet{
		PrivKeyShares: skShares,
		PubGroupKey:   pkGroup,
		PubKeyShares:  pkShares,
	}

	return randomBeaconData, nil
}

// Beacon KG with one node
func thresholdSignKeyGenOneNode(seed []byte) ([]crypto.PrivateKey, []crypto.PublicKey, crypto.PublicKey, error) {
	sk, err := crypto.GeneratePrivateKey(crypto.BLSBLS12381, seed)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("KeyGen with one node failed: %w", err)
	}
	pk := sk.PublicKey()
	return []crypto.PrivateKey{sk},
		[]crypto.PublicKey{pk},
		pk,
		nil
}
