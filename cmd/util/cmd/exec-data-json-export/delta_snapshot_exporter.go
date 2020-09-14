package jsonexporter

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dapperlabs/flow-go/cmd/util/cmd/common"
	"github.com/dapperlabs/flow-go/engine/execution/state/delta"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/module/metrics"
	"github.com/dapperlabs/flow-go/storage/badger"
	"github.com/dapperlabs/flow-go/storage/badger/operation"
)

type dSnapshot struct {
	DeltaJSONStr string `json:"delta_json_str"`
	SpockSecret  string `json:"spock_secret_data"`
}

// ExportDeltaSnapshots exports all the delta snapshots
func ExportDeltaSnapshots(blockID flow.Identifier, dbPath string, outputPath string) error {

	// traverse backward from the given block (parent block) and fetch by blockHash
	db := common.InitStorage(dbPath)
	defer db.Close()

	cacheMetrics := &metrics.NoopCollector{}
	headers := badger.NewHeaders(cacheMetrics, db)

	activeBlockID := blockID
	outputFile := filepath.Join(outputPath, "delta.jsonl")

	fi, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("could not create delta snapshot output file %w", err)
	}
	defer fi.Close()

	writer := bufio.NewWriter(fi)
	defer writer.Flush()

	for {
		header, err := headers.ByBlockID(activeBlockID)
		if err != nil {
			// no more header is available
			return nil
		}

		var snap []*delta.Snapshot
		err = db.View(operation.RetrieveExecutionStateInteractions(activeBlockID, &snap))
		if err != nil {
			return fmt.Errorf("could not load delta snapshot: %w", err)
		}

		if len(snap) < 1 {
			// end of snapshots
			return nil
		}
		m, err := snap[0].Delta.MarshalJSON()
		if err != nil {
			return fmt.Errorf("could not load delta snapshot: %w", err)
		}

		data := dSnapshot{
			DeltaJSONStr: string(m),
			SpockSecret:  hex.EncodeToString(snap[0].SpockSecret),
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("could not create a json obj for a delta snapshot: %w", err)
		}
		_, err = writer.WriteString(string(jsonData) + "\n")
		if err != nil {
			return fmt.Errorf("could not write delta snapshot json to the file: %w", err)
		}
		writer.Flush()

		activeBlockID = header.ParentID
	}
}
