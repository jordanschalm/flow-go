package debug

import (
	"encoding/json"
	"math/big"
	"testing"

	gethCommon "github.com/onflow/go-ethereum/common"
	gethTypes "github.com/onflow/go-ethereum/core/types"
	"github.com/onflow/go-ethereum/core/vm"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/fvm/evm/testutils"
	"github.com/onflow/flow-go/model/flow"
)

func Test_CallTracer(t *testing.T) {
	t.Run("collect traces and upload them", func(t *testing.T) {
		txID := gethCommon.Hash{0x05}
		blockID := flow.Identifier{0x01}
		var res json.RawMessage

		mockUpload := &testutils.MockUploader{
			UploadFunc: func(id string, message json.RawMessage) error {
				require.Equal(t, TraceID(txID, blockID), id)
				require.Equal(t, res, message)
				return nil
			},
		}

		tracer, err := NewEVMCallTracer(mockUpload, zerolog.Nop())
		require.NoError(t, err)
		tracer.WithBlockID(blockID)

		from := gethCommon.HexToAddress("0x01")
		to := gethCommon.HexToAddress("0x02")
		nonce := uint64(10)
		data := []byte{0x02, 0x04}
		amount := big.NewInt(1)

		tr := tracer.TxTracer()
		require.NotNil(t, tr)

		tr.OnEnter(0, 0, from, to, []byte{0x01, 0x02}, 10, big.NewInt(1))
		tx := gethTypes.NewTransaction(nonce, to, amount, 100, big.NewInt(10), data)
		tr.OnTxStart(nil, tx, from)
		tr.OnEnter(1, byte(vm.ADD), from, to, data, 20, big.NewInt(2))
		tr.OnTxEnd(&gethTypes.Receipt{}, nil)
		tr.OnExit(0, []byte{0x02}, 200, nil, false)

		res, err = tr.GetResult()
		require.NoError(t, err)
		require.NotNil(t, res)

		tracer.Collect(txID)
	})

	t.Run("collector panic recovery", func(t *testing.T) {
		txID := gethCommon.Hash{0x05}
		blockID := flow.Identifier{0x01}

		mockUpload := &testutils.MockUploader{
			UploadFunc: func(id string, message json.RawMessage) error {
				panic("nooooo")
			},
		}

		tracer, err := NewEVMCallTracer(mockUpload, zerolog.Nop())
		require.NoError(t, err)
		tracer.WithBlockID(blockID)

		require.NotPanics(t, func() {
			tracer.Collect(txID)
		})
	})

	t.Run("nop tracer", func(t *testing.T) {
		tracer := nopTracer{}
		require.Nil(t, tracer.TxTracer())
	})
}
