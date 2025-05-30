package verification

import (
	"testing"

	"github.com/onflow/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/consensus/hotstuff"
	"github.com/onflow/flow-go/consensus/hotstuff/mocks"
	"github.com/onflow/flow-go/consensus/hotstuff/model"
	"github.com/onflow/flow-go/consensus/hotstuff/safetyrules"
	"github.com/onflow/flow-go/consensus/hotstuff/signature"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/local"
	modulemock "github.com/onflow/flow-go/module/mock"
	msig "github.com/onflow/flow-go/module/signature"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/utils/unittest"
)

// Test that when beacon key is available for a view, a signed block can pass the validation
// the sig include both staking sig and random beacon sig.
func TestCombinedSignWithBeaconKey(t *testing.T) {
	identities := unittest.IdentityListFixture(4, unittest.WithRole(flow.RoleConsensus))

	// prepare data
	beaconKey := unittest.RandomBeaconPriv()
	pk := beaconKey.PublicKey()
	proposerView := uint64(20)

	proposerIdentity := identities[0]
	fblock := unittest.BlockFixture()
	fblock.Header.ProposerID = proposerIdentity.NodeID
	fblock.Header.View = proposerView
	fblock.Header.ParentView = proposerView - 1
	fblock.Header.LastViewTC = nil
	proposal := model.ProposalFromFlow(fblock.Header)
	signerID := fblock.Header.ProposerID

	beaconKeyStore := modulemock.NewRandomBeaconKeyStore(t)
	beaconKeyStore.On("ByView", proposerView).Return(beaconKey, nil)

	ourIdentity := unittest.IdentityFixture()
	stakingPriv := unittest.StakingPrivKeyFixture()
	ourIdentity.NodeID = signerID
	ourIdentity.StakingPubKey = stakingPriv.PublicKey()

	me, err := local.New(ourIdentity.IdentitySkeleton, stakingPriv)
	require.NoError(t, err)
	signer := NewCombinedSigner(me, beaconKeyStore)

	dkg := &mocks.DKG{}
	dkg.On("KeyShare", signerID).Return(pk, nil)

	committee := &mocks.DynamicCommittee{}
	committee.On("DKG", mock.Anything).Return(dkg, nil)
	committee.On("Self").Return(me.NodeID())
	committee.On("IdentityByBlock", fblock.Header.ID(), fblock.Header.ProposerID).Return(proposerIdentity, nil)
	committee.On("LeaderForView", proposerView).Return(signerID, nil).Maybe()

	packer := signature.NewConsensusSigDataPacker(committee)
	verifier := NewCombinedVerifier(committee, packer)

	persist := mocks.NewPersister(t)
	safetyData := &hotstuff.SafetyData{
		LockedOneChainView:      fblock.Header.ParentView,
		HighestAcknowledgedView: fblock.Header.ParentView,
	}
	persist.On("GetSafetyData", mock.Anything).Return(safetyData, nil).Once()
	persist.On("PutSafetyData", mock.Anything).Return(nil)
	safetyRules, err := safetyrules.New(signer, persist, committee)
	require.NoError(t, err)

	// check that the proposer's vote for their own block (i.e. the proposer signature in the header) passes verification
	vote, err := safetyRules.SignOwnProposal(proposal)
	require.NoError(t, err)

	err = verifier.VerifyVote(&ourIdentity.IdentitySkeleton, vote.SigData, proposal.Block.View, proposal.Block.BlockID)
	require.NoError(t, err)

	// check that a created proposal's signature is a combined staking sig and random beacon sig
	block := proposal.Block
	msg := MakeVoteMessage(block.View, block.BlockID)
	stakingSig, err := stakingPriv.Sign(msg, msig.NewBLSHasher(msig.ConsensusVoteTag))
	require.NoError(t, err)

	beaconSig, err := beaconKey.Sign(msg, msig.NewBLSHasher(msig.RandomBeaconTag))
	require.NoError(t, err)

	expectedSig := msig.EncodeDoubleSig(stakingSig, beaconSig)
	require.Equal(t, expectedSig, vote.SigData)

	// vote should be valid
	vote, err = signer.CreateVote(block)
	require.NoError(t, err)

	err = verifier.VerifyVote(&ourIdentity.IdentitySkeleton, vote.SigData, proposal.Block.View, proposal.Block.BlockID)
	require.NoError(t, err)

	// vote on different block should be invalid
	blockWrongID := *block
	blockWrongID.BlockID[0]++
	err = verifier.VerifyVote(&ourIdentity.IdentitySkeleton, vote.SigData, blockWrongID.View, blockWrongID.BlockID)
	require.ErrorIs(t, err, model.ErrInvalidSignature)

	// vote with a wrong proposerView should be invalid
	blockWrongView := *block
	blockWrongView.View++
	err = verifier.VerifyVote(&ourIdentity.IdentitySkeleton, vote.SigData, blockWrongID.View, blockWrongID.BlockID)
	require.ErrorIs(t, err, model.ErrInvalidSignature)

	// vote by different signer should be invalid
	wrongVoter := &identities[1].IdentitySkeleton
	wrongVoter.StakingPubKey = unittest.StakingPrivKeyFixture().PublicKey()
	err = verifier.VerifyVote(wrongVoter, vote.SigData, block.View, block.BlockID)
	require.ErrorIs(t, err, model.ErrInvalidSignature)

	// vote with changed signature should be invalid
	brokenSig := append([]byte{}, vote.SigData...) // copy
	brokenSig[4]++
	err = verifier.VerifyVote(&ourIdentity.IdentitySkeleton, brokenSig, block.View, block.BlockID)
	require.ErrorIs(t, err, model.ErrInvalidSignature)

	// Vote from a node that is _not_ part of the Random Beacon committee should be rejected.
	// Specifically, we expect that the verifier recognizes the `protocol.IdentityNotFoundError`
	// as a sign of an invalid vote and wraps it into a `model.InvalidSignerError`.
	*dkg = mocks.DKG{} // overwrite DKG mock with a new one
	dkg.On("KeyShare", signerID).Return(nil, protocol.IdentityNotFoundError{NodeID: signerID})
	err = verifier.VerifyVote(&ourIdentity.IdentitySkeleton, vote.SigData, proposal.Block.View, proposal.Block.BlockID)
	require.True(t, model.IsInvalidSignerError(err))
}

// Test that when beacon key is not available for a view, a signed block can pass the validation
// the sig only include staking sig
func TestCombinedSignWithNoBeaconKey(t *testing.T) {
	// prepare data
	beaconKey := unittest.RandomBeaconPriv()
	pk := beaconKey.PublicKey()
	proposerView := uint64(20)

	fblock := unittest.BlockFixture()
	fblock.Header.View = proposerView
	fblock.Header.ParentView = proposerView - 1
	fblock.Header.LastViewTC = nil
	proposal := model.ProposalFromFlow(fblock.Header)
	signerID := fblock.Header.ProposerID

	beaconKeyStore := modulemock.NewRandomBeaconKeyStore(t)
	beaconKeyStore.On("ByView", proposerView).Return(nil, module.ErrNoBeaconKeyForEpoch)

	ourIdentity := unittest.IdentityFixture()
	stakingPriv := unittest.StakingPrivKeyFixture()
	ourIdentity.NodeID = signerID
	ourIdentity.StakingPubKey = stakingPriv.PublicKey()

	me, err := local.New(ourIdentity.IdentitySkeleton, stakingPriv)
	require.NoError(t, err)
	signer := NewCombinedSigner(me, beaconKeyStore)

	dkg := &mocks.DKG{}
	dkg.On("KeyShare", signerID).Return(pk, nil)

	committee := &mocks.DynamicCommittee{}
	// even if the node failed DKG, and has no random beacon private key,
	// but other nodes, who completed and succeeded DKG, have a public key
	// for this failed node, which can be used to verify signature from
	// this failed node.
	committee.On("DKG", mock.Anything).Return(dkg, nil)
	committee.On("Self").Return(me.NodeID())
	committee.On("IdentityByBlock", fblock.Header.ID(), signerID).Return(ourIdentity, nil)
	committee.On("LeaderForView", mock.Anything).Return(signerID, nil).Maybe()

	packer := signature.NewConsensusSigDataPacker(committee)
	verifier := NewCombinedVerifier(committee, packer)

	persist := mocks.NewPersister(t)
	safetyData := &hotstuff.SafetyData{
		LockedOneChainView:      fblock.Header.ParentView,
		HighestAcknowledgedView: fblock.Header.ParentView,
	}
	persist.On("GetSafetyData", mock.Anything).Return(safetyData, nil).Once()
	persist.On("PutSafetyData", mock.Anything).Return(nil)
	safetyRules, err := safetyrules.New(signer, persist, committee)
	require.NoError(t, err)

	// check that the proposer's vote for their own block (i.e. the proposer signature in the header) passes verification
	vote, err := safetyRules.SignOwnProposal(proposal)
	require.NoError(t, err)

	err = verifier.VerifyVote(&ourIdentity.IdentitySkeleton, vote.SigData, proposal.Block.View, proposal.Block.BlockID)
	require.NoError(t, err)

	// As the proposer does not have a Random Beacon Key, it should sign solely with its staking key.
	// In this case, the SigData should be identical to the staking sig.
	expectedStakingSig, err := stakingPriv.Sign(
		MakeVoteMessage(proposal.Block.View, proposal.Block.BlockID),
		msig.NewBLSHasher(msig.ConsensusVoteTag),
	)
	require.NoError(t, err)
	require.Equal(t, expectedStakingSig, crypto.Signature(vote.SigData))
}

// Test_VerifyQC_EmptySigners checks that Verifier returns an `model.InsufficientSignaturesError`
// if `signers` input is empty or nil. This check should happen _before_ the Verifier calls into
// any sub-components, because some (e.g. `crypto.AggregateBLSPublicKeys`) don't provide sufficient
// sentinel errors to distinguish between internal problems and external byzantine inputs.
func Test_VerifyQC_EmptySigners(t *testing.T) {
	committee := &mocks.DynamicCommittee{}
	dkg := &mocks.DKG{}
	pk := &modulemock.PublicKey{}
	pk.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	dkg.On("GroupKey").Return(pk)
	committee.On("DKG", mock.Anything).Return(dkg, nil)

	packer := signature.NewConsensusSigDataPacker(committee)
	verifier := NewCombinedVerifier(committee, packer)

	header := unittest.BlockHeaderFixture()
	block := model.BlockFromFlow(header)

	// sigData with empty signers
	emptySignersInput := model.SignatureData{
		SigType:                      []byte{},
		AggregatedStakingSig:         unittest.SignatureFixture(),
		AggregatedRandomBeaconSig:    unittest.SignatureFixture(),
		ReconstructedRandomBeaconSig: unittest.SignatureFixture(),
	}
	encoder := new(model.SigDataPacker)
	sigData, err := encoder.Encode(&emptySignersInput)
	require.NoError(t, err)

	err = verifier.VerifyQC(flow.IdentitySkeletonList{}, sigData, block.View, block.BlockID)
	require.True(t, model.IsInsufficientSignaturesError(err))

	err = verifier.VerifyQC(nil, sigData, block.View, block.BlockID)
	require.True(t, model.IsInsufficientSignaturesError(err))
}

// TestCombinedSign_BeaconKeyStore_ViewForUnknownEpoch tests that if the beacon
// key store reports the view of the entity to sign has no known epoch, an
// exception should be raised.
func TestCombinedSign_BeaconKeyStore_ViewForUnknownEpoch(t *testing.T) {
	beaconKeyStore := modulemock.NewRandomBeaconKeyStore(t)
	beaconKeyStore.On("ByView", mock.Anything).Return(nil, model.ErrViewForUnknownEpoch)

	stakingPriv := unittest.StakingPrivKeyFixture()
	nodeID := unittest.IdentityFixture()
	nodeID.StakingPubKey = stakingPriv.PublicKey()

	me, err := local.New(nodeID.IdentitySkeleton, stakingPriv)
	require.NoError(t, err)
	signer := NewCombinedSigner(me, beaconKeyStore)

	fblock := unittest.BlockHeaderFixture()
	block := model.BlockFromFlow(fblock)

	vote, err := signer.CreateVote(block)
	require.Error(t, err)
	assert.Nil(t, vote)
}
