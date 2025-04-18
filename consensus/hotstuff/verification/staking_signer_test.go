package verification

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/consensus/hotstuff/helper"
	"github.com/onflow/flow-go/consensus/hotstuff/model"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/local"
	modulemock "github.com/onflow/flow-go/module/mock"
	"github.com/onflow/flow-go/utils/unittest"
)

// TestStakingSigner_CreateVote verifies that StakingSigner can produce correctly signed vote
// that can be verified later using StakingVerifier.
// Additionally, we check cases where errors during signing are happening.
func TestStakingSigner_CreateVote(t *testing.T) {
	stakingPriv := unittest.StakingPrivKeyFixture()
	signer := unittest.IdentityFixture()
	signer.StakingPubKey = stakingPriv.PublicKey()
	signerID := signer.NodeID

	t.Run("could-not-sign", func(t *testing.T) {
		signException := errors.New("sign-exception")
		me := &modulemock.Local{}
		me.On("NodeID").Return(signerID)
		me.On("Sign", mock.Anything, mock.Anything).Return(nil, signException).Once()
		signer := NewStakingSigner(me)

		block := helper.MakeBlock()
		proposal, err := signer.CreateVote(block)
		require.ErrorAs(t, err, &signException)
		require.Nil(t, proposal)
	})
	t.Run("created-vote", func(t *testing.T) {
		me, err := local.New(signer.IdentitySkeleton, stakingPriv)
		require.NoError(t, err)

		signerIdentity := &unittest.IdentityFixture(unittest.WithNodeID(signerID),
			unittest.WithStakingPubKey(stakingPriv.PublicKey())).IdentitySkeleton

		signer := NewStakingSigner(me)

		block := helper.MakeBlock(helper.WithBlockProposer(signerID))
		vote, err := signer.CreateVote(block)
		require.NoError(t, err)
		require.NotNil(t, vote)

		verifier := NewStakingVerifier()
		err = verifier.VerifyVote(signerIdentity, vote.SigData, block.View, block.BlockID)
		require.NoError(t, err)
	})
}

// TestStakingSigner_VerifyQC checks that a QC without any signers is rejected right away without calling into any sub-components
func TestStakingSigner_VerifyQC(t *testing.T) {
	header := unittest.BlockHeaderFixture()
	block := model.BlockFromFlow(header)
	sigData := unittest.RandomBytes(127)

	verifier := NewStakingVerifier()
	err := verifier.VerifyQC(flow.IdentitySkeletonList{}, sigData, block.View, block.BlockID)
	require.True(t, model.IsInsufficientSignaturesError(err))

	err = verifier.VerifyQC(nil, sigData, block.View, block.BlockID)
	require.True(t, model.IsInsufficientSignaturesError(err))
}
