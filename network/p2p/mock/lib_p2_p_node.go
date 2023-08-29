// Code generated by mockery v2.21.4. DO NOT EDIT.

package mockp2p

import (
	component "github.com/onflow/flow-go/module/component"
	channels "github.com/onflow/flow-go/network/channels"

	context "context"

	corenetwork "github.com/libp2p/go-libp2p/core/network"

	flow "github.com/onflow/flow-go/model/flow"

	host "github.com/libp2p/go-libp2p/core/host"

	irrecoverable "github.com/onflow/flow-go/module/irrecoverable"

	kbucket "github.com/libp2p/go-libp2p-kbucket"

	mock "github.com/stretchr/testify/mock"

	network "github.com/onflow/flow-go/network"

	p2p "github.com/onflow/flow-go/network/p2p"

	peer "github.com/libp2p/go-libp2p/core/peer"

	protocol "github.com/libp2p/go-libp2p/core/protocol"

	protocols "github.com/onflow/flow-go/network/p2p/unicast/protocols"

	routing "github.com/libp2p/go-libp2p/core/routing"
)

// LibP2PNode is an autogenerated mock type for the LibP2PNode type
type LibP2PNode struct {
	mock.Mock
}

// ActiveClustersChanged provides a mock function with given fields: _a0
func (_m *LibP2PNode) ActiveClustersChanged(_a0 flow.ChainIDList) {
	_m.Called(_a0)
}

// ConnectToPeerAddrInfo provides a mock function with given fields: ctx, peerInfo
func (_m *LibP2PNode) ConnectToPeerAddrInfo(ctx context.Context, peerInfo peer.AddrInfo) error {
	ret := _m.Called(ctx, peerInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, peer.AddrInfo) error); ok {
		r0 = rf(ctx, peerInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Done provides a mock function with given fields:
func (_m *LibP2PNode) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// GetIPPort provides a mock function with given fields:
func (_m *LibP2PNode) GetIPPort() (string, string, error) {
	ret := _m.Called()

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func() (string, string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() string); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetPeersForProtocol provides a mock function with given fields: pid
func (_m *LibP2PNode) GetPeersForProtocol(pid protocol.ID) peer.IDSlice {
	ret := _m.Called(pid)

	var r0 peer.IDSlice
	if rf, ok := ret.Get(0).(func(protocol.ID) peer.IDSlice); ok {
		r0 = rf(pid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(peer.IDSlice)
		}
	}

	return r0
}

// HasSubscription provides a mock function with given fields: topic
func (_m *LibP2PNode) HasSubscription(topic channels.Topic) bool {
	ret := _m.Called(topic)

	var r0 bool
	if rf, ok := ret.Get(0).(func(channels.Topic) bool); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Host provides a mock function with given fields:
func (_m *LibP2PNode) Host() host.Host {
	ret := _m.Called()

	var r0 host.Host
	if rf, ok := ret.Get(0).(func() host.Host); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(host.Host)
		}
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *LibP2PNode) ID() peer.ID {
	ret := _m.Called()

	var r0 peer.ID
	if rf, ok := ret.Get(0).(func() peer.ID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(peer.ID)
	}

	return r0
}

// IsConnected provides a mock function with given fields: peerID
func (_m *LibP2PNode) IsConnected(peerID peer.ID) (bool, error) {
	ret := _m.Called(peerID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(peer.ID) (bool, error)); ok {
		return rf(peerID)
	}
	if rf, ok := ret.Get(0).(func(peer.ID) bool); ok {
		r0 = rf(peerID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(peer.ID) error); ok {
		r1 = rf(peerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsDisallowListed provides a mock function with given fields: peerId
func (_m *LibP2PNode) IsDisallowListed(peerId peer.ID) ([]network.DisallowListedCause, bool) {
	ret := _m.Called(peerId)

	var r0 []network.DisallowListedCause
	var r1 bool
	if rf, ok := ret.Get(0).(func(peer.ID) ([]network.DisallowListedCause, bool)); ok {
		return rf(peerId)
	}
	if rf, ok := ret.Get(0).(func(peer.ID) []network.DisallowListedCause); ok {
		r0 = rf(peerId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]network.DisallowListedCause)
		}
	}

	if rf, ok := ret.Get(1).(func(peer.ID) bool); ok {
		r1 = rf(peerId)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// ListPeers provides a mock function with given fields: topic
func (_m *LibP2PNode) ListPeers(topic string) []peer.ID {
	ret := _m.Called(topic)

	var r0 []peer.ID
	if rf, ok := ret.Get(0).(func(string) []peer.ID); ok {
		r0 = rf(topic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]peer.ID)
		}
	}

	return r0
}

// OnAllowListNotification provides a mock function with given fields: id, cause
func (_m *LibP2PNode) OnAllowListNotification(id peer.ID, cause network.DisallowListedCause) {
	_m.Called(id, cause)
}

// OnDisallowListNotification provides a mock function with given fields: id, cause
func (_m *LibP2PNode) OnDisallowListNotification(id peer.ID, cause network.DisallowListedCause) {
	_m.Called(id, cause)
}

// OpenProtectedStream provides a mock function with given fields: ctx, peerID, protectionTag, writingLogic
func (_m *LibP2PNode) OpenProtectedStream(ctx context.Context, peerID peer.ID, protectionTag string, writingLogic func(corenetwork.Stream) error) error {
	ret := _m.Called(ctx, peerID, protectionTag, writingLogic)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, peer.ID, string, func(corenetwork.Stream) error) error); ok {
		r0 = rf(ctx, peerID, protectionTag, writingLogic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PeerManagerComponent provides a mock function with given fields:
func (_m *LibP2PNode) PeerManagerComponent() component.Component {
	ret := _m.Called()

	var r0 component.Component
	if rf, ok := ret.Get(0).(func() component.Component); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(component.Component)
		}
	}

	return r0
}

// PeerScoreExposer provides a mock function with given fields:
func (_m *LibP2PNode) PeerScoreExposer() p2p.PeerScoreExposer {
	ret := _m.Called()

	var r0 p2p.PeerScoreExposer
	if rf, ok := ret.Get(0).(func() p2p.PeerScoreExposer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.PeerScoreExposer)
		}
	}

	return r0
}

// Publish provides a mock function with given fields: ctx, messageScope
func (_m *LibP2PNode) Publish(ctx context.Context, messageScope network.OutgoingMessageScope) error {
	ret := _m.Called(ctx, messageScope)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, network.OutgoingMessageScope) error); ok {
		r0 = rf(ctx, messageScope)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *LibP2PNode) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// RemovePeer provides a mock function with given fields: peerID
func (_m *LibP2PNode) RemovePeer(peerID peer.ID) error {
	ret := _m.Called(peerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(peer.ID) error); ok {
		r0 = rf(peerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RequestPeerUpdate provides a mock function with given fields:
func (_m *LibP2PNode) RequestPeerUpdate() {
	_m.Called()
}

// Routing provides a mock function with given fields:
func (_m *LibP2PNode) Routing() routing.Routing {
	ret := _m.Called()

	var r0 routing.Routing
	if rf, ok := ret.Get(0).(func() routing.Routing); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(routing.Routing)
		}
	}

	return r0
}

// RoutingTable provides a mock function with given fields:
func (_m *LibP2PNode) RoutingTable() *kbucket.RoutingTable {
	ret := _m.Called()

	var r0 *kbucket.RoutingTable
	if rf, ok := ret.Get(0).(func() *kbucket.RoutingTable); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kbucket.RoutingTable)
		}
	}

	return r0
}

// SetComponentManager provides a mock function with given fields: cm
func (_m *LibP2PNode) SetComponentManager(cm *component.ComponentManager) {
	_m.Called(cm)
}

// SetPubSub provides a mock function with given fields: ps
func (_m *LibP2PNode) SetPubSub(ps p2p.PubSubAdapter) {
	_m.Called(ps)
}

// SetRouting provides a mock function with given fields: r
func (_m *LibP2PNode) SetRouting(r routing.Routing) {
	_m.Called(r)
}

// SetUnicastManager provides a mock function with given fields: uniMgr
func (_m *LibP2PNode) SetUnicastManager(uniMgr p2p.UnicastManager) {
	_m.Called(uniMgr)
}

// Start provides a mock function with given fields: ctx
func (_m *LibP2PNode) Start(ctx irrecoverable.SignalerContext) {
	_m.Called(ctx)
}

// Stop provides a mock function with given fields:
func (_m *LibP2PNode) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: topic, topicValidator
func (_m *LibP2PNode) Subscribe(topic channels.Topic, topicValidator p2p.TopicValidatorFunc) (p2p.Subscription, error) {
	ret := _m.Called(topic, topicValidator)

	var r0 p2p.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(channels.Topic, p2p.TopicValidatorFunc) (p2p.Subscription, error)); ok {
		return rf(topic, topicValidator)
	}
	if rf, ok := ret.Get(0).(func(channels.Topic, p2p.TopicValidatorFunc) p2p.Subscription); ok {
		r0 = rf(topic, topicValidator)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(channels.Topic, p2p.TopicValidatorFunc) error); ok {
		r1 = rf(topic, topicValidator)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Unsubscribe provides a mock function with given fields: topic
func (_m *LibP2PNode) Unsubscribe(topic channels.Topic) error {
	ret := _m.Called(topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(channels.Topic) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithDefaultUnicastProtocol provides a mock function with given fields: defaultHandler, preferred
func (_m *LibP2PNode) WithDefaultUnicastProtocol(defaultHandler corenetwork.StreamHandler, preferred []protocols.ProtocolName) error {
	ret := _m.Called(defaultHandler, preferred)

	var r0 error
	if rf, ok := ret.Get(0).(func(corenetwork.StreamHandler, []protocols.ProtocolName) error); ok {
		r0 = rf(defaultHandler, preferred)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithPeersProvider provides a mock function with given fields: peersProvider
func (_m *LibP2PNode) WithPeersProvider(peersProvider p2p.PeersProvider) {
	_m.Called(peersProvider)
}

type mockConstructorTestingTNewLibP2PNode interface {
	mock.TestingT
	Cleanup(func())
}

// NewLibP2PNode creates a new instance of LibP2PNode. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLibP2PNode(t mockConstructorTestingTNewLibP2PNode) *LibP2PNode {
	mock := &LibP2PNode{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
