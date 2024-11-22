package block

import (
	"encoding/hex"
	"fmt"
	"github.com/streamingfast/firehose-aelf/pb/aelf"
	pbaelf "github.com/streamingfast/firehose-aelf/pb/sf/aelf/type/v1"
	"google.golang.org/protobuf/proto"
	"log"
)

func ConvertBlock(blockHash string, block *aelf.Block) *pbaelf.Block {
	return &pbaelf.Block{
		Version:           1,
		BlockHash:         blockHash,
		Height:            block.Header.Height,
		Header:            convertBlockHeader(block.Header),
		TransactionTraces: prepareTransactionTraces(block),
	}
}

func convertBlockHeader(left *aelf.BlockHeader) *pbaelf.BlockHeader {
	return &pbaelf.BlockHeader{
		Version:                           left.Version,
		ChainId:                           left.ChainId,
		PreviousBlockHash:                 left.PreviousBlockHash.ToHex(),
		MerkleTreeRootOfTransactions:      left.MerkleTreeRootOfTransactions.ToHex(),
		MerkleTreeRootOfWorldState:        left.MerkleTreeRootOfWorldState.ToHex(),
		Bloom:                             left.Bloom,
		Height:                            left.Height,
		ExtraData:                         left.ExtraData,
		Time:                              left.Time,
		MerkleTreeRootOfTransactionStatus: left.MerkleTreeRootOfTransactionStatus.ToHex(),
		SignerPubkey:                      left.SignerPubkey,
		Signature:                         left.Signature,
	}
}

type TrackedTransactionTrace struct {
	*aelf.TransactionTrace
	PreTrackedTraces    []*TrackedTransactionTrace
	InlineTrackedTraces []*TrackedTransactionTrace
	PostTrackedTraces   []*TrackedTransactionTrace
	IsReverted          bool
}

func (t *TrackedTransactionTrace) SetReverted() {
	t.IsReverted = true
	for _, inlineTrace := range t.InlineTrackedTraces {
		inlineTrace.SetReverted()
	}
}

func prepareTransactionTraces(block *aelf.Block) []*pbaelf.TransactionTrace {
	var pbTraces []*pbaelf.TransactionTrace
	for i, txIdInHash := range block.Body.TransactionIds {
		txId := txIdInHash.ToHex()
		tx := block.FirehoseBody.Transactions[i]

		trace := block.FirehoseBody.TransactionTraces[i]
		calls, mainCallIndex := extractCalls(tx, convertTraceToTracked(trace), txId, "", 0)

		pbTrace := &pbaelf.TransactionTrace{
			TransactionId:  txId,
			RawTransaction: serializeTransaction(tx), // TODO: Check if this is reliable
			Signature:      tx.Signature,
			Calls:          calls,
			MainCallIndex:  mainCallIndex,
		}
		pbTraces = append(pbTraces, pbTrace)
	}
	return pbTraces
}

func convertTraceToTracked(trace *aelf.TransactionTrace) *TrackedTransactionTrace {
	var (
		pre    []*TrackedTransactionTrace
		inline []*TrackedTransactionTrace
		post   []*TrackedTransactionTrace
	)
	reverted := false
	for _, preTrace := range trace.PreTraces {
		convertedPre := convertTraceToTracked(preTrace)
		pre = append(pre, convertedPre)
		reverted = reverted || convertedPre.IsReverted
	}
	for _, inlineTrace := range trace.InlineTraces {
		convertedInline := convertTraceToTracked(inlineTrace)
		inline = append(inline, convertedInline)
		reverted = reverted || convertedInline.IsReverted
	}
	for _, postTrace := range trace.PostTraces {
		convertedPost := convertTraceToTracked(postTrace)
		post = append(post, convertedPost)
		reverted = reverted || convertedPost.IsReverted
	}
	reverted = reverted || trace.ExecutionStatus != aelf.ExecutionStatus_EXECUTED

	newTrace := &TrackedTransactionTrace{
		TransactionTrace:    trace,
		PreTrackedTraces:    pre,
		InlineTrackedTraces: inline,
		PostTrackedTraces:   post,
		IsReverted:          reverted,
	}

	if reverted {
		// need to propagate down the reversion status
		newTrace.SetReverted()
	}

	return newTrace
}

func serializeTransaction(tx *aelf.Transaction) []byte {
	data, err := proto.Marshal(tx)
	if err != nil {
		log.Fatalf("Failed to marshal transaction: %v", err)
		return nil
	}
	return data
}

func extractCalls(tx *aelf.Transaction, trace *TrackedTransactionTrace, txId string, callPathPrefix string, index int) ([]*pbaelf.Call, int32) {

	log.Println(fmt.Sprintf("extract %s %s", txId, callPathPrefix))
	var flattenedCalls []*pbaelf.Call
	thisCallPath := fmt.Sprintf("%s:%d", callPathPrefix, index)
	for i, preTrace := range trace.PreTrackedTraces {
		preCallPathPrefix := fmt.Sprintf("%s:pre", thisCallPath)
		preTx := trace.PreTransactions[i]
		childrenCalls, _ := extractCalls(preTx, preTrace, txId, preCallPathPrefix, i)
		for _, call := range childrenCalls {
			flattenedCalls = append(flattenedCalls, call)
		}
	}

	mainCallIndex := len(flattenedCalls)

	mainCall := &pbaelf.Call{
		TransactionId:   txId,
		CallPath:        thisCallPath,
		RefBlockNumber:  tx.RefBlockNumber,
		RefBlockPrefix:  hex.EncodeToString(tx.RefBlockPrefix),
		From:            tx.From.ToBase58(),
		To:              tx.To.ToBase58(),
		MethodName:      tx.MethodName,
		Params:          tx.Params,
		ExecutionStatus: pbaelf.ExecutionStatus(trace.ExecutionStatus),
		ReturnValue:     trace.ReturnValue,
		Error:           trace.Error,
		StateSet: &pbaelf.TransactionExecutingStateSet{
			Writes:  trace.StateSet.Writes,
			Reads:   trace.StateSet.Reads,
			Deletes: trace.StateSet.Deletes,
		},
		Logs:       convertLogs(trace.Logs),
		IsReverted: trace.IsReverted,
	}
	flattenedCalls = append(flattenedCalls, mainCall)

	for i, inlineTrace := range trace.InlineTrackedTraces {
		inlineCallPathPrefix := thisCallPath
		inlineTx := trace.InlineTransactions[i]
		childrenCalls, _ := extractCalls(inlineTx, inlineTrace, txId, inlineCallPathPrefix, i)
		for _, call := range childrenCalls {
			flattenedCalls = append(flattenedCalls, call)
		}
	}
	for i, postTrace := range trace.PostTrackedTraces {
		postCallPathPrefix := fmt.Sprintf("%s:post", thisCallPath)
		postTx := trace.PostTransactions[i]
		childrenCalls, _ := extractCalls(postTx, postTrace, txId, postCallPathPrefix, i)
		for _, call := range childrenCalls {
			flattenedCalls = append(flattenedCalls, call)
		}
	}
	return flattenedCalls, int32(mainCallIndex)
}

func convertLogs(original []*aelf.LogEvent) []*pbaelf.LogEvent {
	var output []*pbaelf.LogEvent
	for _, log := range original {
		newLog := &pbaelf.LogEvent{
			Address:    log.Address.ToBase58(),
			Name:       log.Name,
			Indexed:    log.Indexed,
			NonIndexed: log.NonIndexed,
		}
		output = append(output, newLog)
	}
	return output
}
