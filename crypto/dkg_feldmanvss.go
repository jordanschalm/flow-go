//go:build relic
// +build relic

package crypto

// #cgo CFLAGS: -g -Wall -std=c99
// #include "dkg_include.h"
import "C"

import (
	"errors"
	"fmt"
)

// Implements Feldman Verifiable Secret Sharing using
// the BLS set up on the BLS12-381 curve.

// The secret is a BLS private key generated by a single dealer/leader.
// (and hence this is a centralized generation).
// The leader generates key shares for a BLS-based
// threshold signature scheme and distributes the shares over the (n)
// partcipants including itself. The particpants validate their shares
// using a public verifiaction vector shared by the leader.

// Private keys are scalar in Zr, where r is the group order of G1/G2
// Public keys are in G2.

// feldman VSS protocol, implements DKGState
type feldmanVSSstate struct {
	// common DKG state
	*dkgCommon
	// participant leader index
	leaderIndex index
	// Polynomial P = a_0 + a_1*x + .. + a_t*x^t  in Zr[X], the vector size is (t+1)
	// a_0 is the group private key
	a []scalar
	// Public vector of the group, the vector size is (t+1)
	// A_0 is the group public key
	vA         []pointG2
	vAReceived bool
	// Private share of the current participant
	x         scalar
	xReceived bool
	// Public keys of the group participants, the vector size is (n)
	y []pointG2
	// true if the private share is valid
	validKey bool
}

// NewFeldmanVSS creates a new instance of Feldman VSS protocol.
//
// An instance is run by a single participant and is usable for only one protocol.
// In order to run the protocol again, a new instance needs to be created
func NewFeldmanVSS(size int, threshold int, myIndex int,
	processor DKGProcessor, leaderIndex int) (DKGState, error) {

	common, err := newDKGCommon(size, threshold, myIndex, processor, leaderIndex)
	if err != nil {
		return nil, err
	}

	fvss := &feldmanVSSstate{
		dkgCommon:   common,
		leaderIndex: index(leaderIndex),
	}
	fvss.init()
	return fvss, nil
}

func (s *feldmanVSSstate) init() {
	// set the bls context
	blsInstance.reInit()
	s.running = false
	s.y = nil
	s.xReceived = false
	s.vAReceived = false
	C.bn_new_wrapper((*C.bn_st)(&s.x))
}

// Start starts running the protocol in the current participant
// If the current participant is the leader, then the seed is used
// to generate the secret polynomial (including the group private key)
// if the current participant is not the leader, the seed is ignored.
func (s *feldmanVSSstate) Start(seed []byte) error {
	if s.running {
		return errors.New("dkg is already running")
	}

	s.running = true
	// Generate shares if necessary
	if s.leaderIndex == s.myIndex {
		return s.generateShares(seed)
	}
	return nil
}

// End ends the protocol in the current node.
// It returns the finalized public data and participants private key share.
// - the group public key corresponding to the group secret key
// - all the public key shares corresponding to the participants private
// key shares.
// - the finalized private key which is the current participant's own private key share
// - the returned erorr is :
//    - dkgFailureError if the private key and vector are inconsistent.
//    - other error if Start() was not called.
//    - nil otherwise.
func (s *feldmanVSSstate) End() (PrivateKey, PublicKey, []PublicKey, error) {
	if !s.running {
		return nil, nil, nil, errors.New("dkg is not running")
	}
	s.running = false
	if !s.validKey {
		return nil, nil, nil, dkgFailureErrorf("received private key is invalid")
	}
	// private key of the current participant
	x := newPrKeyBLSBLS12381(&s.x)

	// Group public key
	Y := newPubKeyBLSBLS12381(&s.vA[0])

	// The participants public keys
	y := make([]PublicKey, s.size)
	for i, p := range s.y {
		y[i] = newPubKeyBLSBLS12381(&p)
	}
	return x, Y, y, nil
}

const (
	shareSize = PrKeyLenBLSBLS12381
	// the actual verifVectorSize depends on the state and is:
	// PubKeyLenBLSBLS12381*(t+1)
	verifVectorSize = PubKeyLenBLSBLS12381
)

// HandleBroadcastMsg processes a new broadcasted message received by the current participant.
//
// orig is the message origin index.
func (s *feldmanVSSstate) HandleBroadcastMsg(orig int, msg []byte) error {
	if !s.running {
		return errors.New("dkg is not running")
	}
	if orig >= s.Size() || orig < 0 {
		return invalidInputsErrorf(
			"wrong origin input, should be less than %d, got %d",
			s.Size(),
			orig)
	}

	// In case a message is received by the origin participant,
	// the message is just ignored
	if s.myIndex == index(orig) {
		return nil
	}

	if len(msg) == 0 {
		s.processor.Disqualify(orig, "the received broadcast is empty")
		return nil
	}

	// msg = |tag| Data |
	if dkgMsgTag(msg[0]) == feldmanVSSVerifVec {
		s.receiveVerifVector(index(orig), msg[1:])
	} else {
		s.processor.Disqualify(orig,
			fmt.Sprintf("the broadcast header is invalid, got %d",
				dkgMsgTag(msg[0])))
	}
	return nil
}

// HandlePrivateMsg processes a new private message received by the current participant.
//
// orig is the message origin index.
func (s *feldmanVSSstate) HandlePrivateMsg(orig int, msg []byte) error {
	if !s.running {
		return errors.New("dkg is not running")
	}

	if orig >= s.Size() || orig < 0 {
		return invalidInputsErrorf(
			"wrong origin, should be positive less than %d, got %d",
			s.Size(),
			orig)
	}

	if len(msg) == 0 {
		// Ideally, the upper layer should stop sebsequent messages from sender
		s.processor.FlagMisbehavior(orig, "the private message is empty")
		return nil
	}

	// In case a private message is received by the origin participant,
	// the message is just ignored
	if s.myIndex == index(orig) {
		return nil
	}

	// msg = |tag| Data |
	if dkgMsgTag(msg[0]) == feldmanVSSShare {
		s.receiveShare(index(orig), msg[1:])
	} else {
		s.processor.FlagMisbehavior(orig,
			fmt.Sprintf("the private message header is invalid, got %d",
				dkgMsgTag(msg[0])))
	}
	return nil
}

// ForceDisqualify forces a participant to get disqualified
// for a reason outside of the DKG protocol
// The caller should make sure all honest participants call this function,
// otherwise, the protocol can be broken
func (s *feldmanVSSstate) ForceDisqualify(participant int) error {
	if !s.running {
		return errors.New("dkg is not running")
	}
	if participant >= s.Size() || participant < 0 {
		return invalidInputsErrorf(
			"wrong origin input, should be less than %d, got %d",
			s.Size(),
			participant)
	}
	if index(participant) == s.leaderIndex {
		s.validKey = false
	}
	return nil
}

// generates all private and public data by the leader
func (s *feldmanVSSstate) generateShares(seed []byte) error {
	err := seedRelic(seed)
	if err != nil {
		return fmt.Errorf("generating shares failed: %w", err)
	}

	// Generate a polyomial P in Zr[X] of degree t
	s.a = make([]scalar, s.threshold+1)
	s.vA = make([]pointG2, s.threshold+1)
	s.y = make([]pointG2, s.size)
	randZrStar(&s.a[0]) // non zero a[0]
	generatorScalarMultG2(&s.vA[0], &s.a[0])
	for i := 1; i < s.threshold+1; i++ {
		C.bn_new_wrapper((*C.bn_st)(&s.a[i]))
		randZr(&s.a[i])
		generatorScalarMultG2(&s.vA[i], &s.a[i])
	}

	// compute the shares
	for i := index(1); int(i) <= s.size; i++ {
		// the-leader-own share
		if i-1 == s.myIndex {
			xdata := make([]byte, shareSize)
			zrPolynomialImage(xdata, s.a, i, &s.y[i-1])
			C.bn_read_bin((*C.bn_st)(&s.x),
				(*C.uchar)(&xdata[0]),
				PrKeyLenBLSBLS12381,
			)
			continue
		}
		// the-other-participant shares
		data := make([]byte, shareSize+1)
		data[0] = byte(feldmanVSSShare)
		zrPolynomialImage(data[1:], s.a, i, &s.y[i-1])
		s.processor.PrivateSend(int(i-1), data)
	}
	// broadcast the vector
	vectorSize := verifVectorSize * (s.threshold + 1)
	data := make([]byte, vectorSize+1)
	data[0] = byte(feldmanVSSVerifVec)
	writeVerifVector(data[1:], s.vA)
	s.processor.Broadcast(data)

	s.vAReceived = true
	s.xReceived = true
	s.validKey = true
	return nil
}

// receives a private share from the leader
func (s *feldmanVSSstate) receiveShare(origin index, data []byte) {
	// only accept private shares from the leader.
	if origin != s.leaderIndex {
		return
	}

	if s.xReceived {
		s.processor.FlagMisbehavior(int(origin), "private share was already received")
		return
	}

	// at this point, tag the private message as received
	s.xReceived = true

	if (len(data)) != shareSize {
		s.validKey = false
		s.processor.FlagMisbehavior(int(origin),
			fmt.Sprintf("invalid share size, expects %d, got %d",
				shareSize, len(data)))
		return
	}

	// read the participant private share
	if C.bn_read_Zr_bin((*C.bn_st)(&s.x),
		(*C.uchar)(&data[0]),
		PrKeyLenBLSBLS12381,
	) != valid {
		s.validKey = false
		s.processor.FlagMisbehavior(int(origin),
			fmt.Sprintf("invalid share value %x", data))
		return
	}

	if s.vAReceived {
		s.validKey = s.verifyShare()
	}
}

// receives the public vector from the leader
func (s *feldmanVSSstate) receiveVerifVector(origin index, data []byte) {
	// only accept the verification vector from the leader.
	if origin != s.leaderIndex {
		return
	}

	if s.vAReceived {
		s.processor.FlagMisbehavior(int(origin),
			"verification vector was already received")
		return
	}

	if verifVectorSize*(s.threshold+1) != len(data) {
		s.vAReceived = true
		s.validKey = false
		s.processor.Disqualify(int(origin),
			fmt.Sprintf("invalid verification vector size, expects %d, got %d",
				verifVectorSize*(s.threshold+1), len(data)))
		return
	}
	// read the verification vector
	s.vA = make([]pointG2, s.threshold+1)
	err := readVerifVector(s.vA, data)
	if err != nil {
		if IsInvalidInputsError(err) { // case where vector format is invalid
			s.vAReceived = true
			s.validKey = false
			s.processor.Disqualify(int(origin),
				fmt.Sprintf("reading the verification vector failed: %s", err))
		}
		// verification vector should not be tagged as received if unexpected error
		return
	}

	s.y = make([]pointG2, s.size)
	s.computePublicKeys()

	s.vAReceived = true
	if s.xReceived {
		s.validKey = s.verifyShare()
	}
}

// zrPolynomialImage computes P(x) = a_0 + a_1*x + .. + a_n*x^n (mod r) in Z/Zr
// r being the order of G1
// P(x) is written in dest, while g2^P(x) is written in y
// x being a small integer
func zrPolynomialImage(dest []byte, a []scalar, x index, y *pointG2) {
	C.Zr_polynomialImage_export((*C.uchar)(&dest[0]),
		(*C.ep2_st)(y),
		(*C.bn_st)(&a[0]), (C.int)(len(a)),
		(C.uint8_t)(x),
	)
}

// writeVerifVector exports a vector A into an array of bytes
// assuming the array length matches the vector length
func writeVerifVector(dest []byte, A []pointG2) {
	C.ep2_vector_write_bin((*C.uchar)(&dest[0]),
		(*C.ep2_st)(&A[0]),
		(C.int)(len(A)),
	)
}

// readVerifVector imports A vector from an array of bytes,
// assuming the slice length matches the vector length
func readVerifVector(A []pointG2, src []byte) error {
	switch C.ep2_vector_read_bin((*C.ep2_st)(&A[0]),
		(*C.uchar)(&src[0]),
		(C.int)(len(A))) {
	case valid:
		return nil
	case invalid:
		return invalidInputsErrorf("the verifcation vector does not serialize G2 points")
	default:
		return errors.New("reading the verifcation vector failed")
	}
}

func (s *feldmanVSSstate) verifyShare() bool {
	// check y[current] == x.G2
	return C.verifyshare((*C.bn_st)(&s.x),
		(*C.ep2_st)(&s.y[s.myIndex])) == 1
}

// computePublicKeys extracts the participants public keys from the verification vector
// y[i] = Q(i+1) for all participants i, with:
//  Q(x) = A_0 + A_1*x + ... +  A_n*x^n  in G2
func (s *feldmanVSSstate) computePublicKeys() {
	C.G2_polynomialImages(
		(*C.ep2_st)(&s.y[0]), (C.int)(len(s.y)),
		(*C.ep2_st)(&s.vA[0]), (C.int)(len(s.vA)),
	)
}
