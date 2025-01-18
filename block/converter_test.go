package block

import (
	"encoding/base64"
	"github.com/streamingfast/firehose-aelf/pb/aelf"
	"github.com/test-go/testify/assert"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestConvertBlock(t *testing.T) {
	blockSample := "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eW1ib2wQAVjbogFgAWrjEAq2BAo6cEdhNGU1aE5Hc2drZmpFR203MlRFdmJGN2FSRHFLQmQ0THVYdGFiNHVjTWJYTGNnSi9Sb3VuZHMvMRL3AwgBEu4DCoIBMDQ0ZjM1NWJkY2I3Y2MwYWY3MjhlZjNjY2ViOTYxNWQ5MDY4NGJiNWIyY2E1Zjg1OWFiMGYwYjcwNDA3NTg3MWFhMzg1YjZiMWI4ZWFkODA5Y2E2NzQ1NGQ5NjgzZmNmMmJhMDM0NTZkNmZlMmM0YWJlMmIwN2YwZmJkYmIyZjFjMRLmAggBEAEiIgog3QrTQ24iKMMmOBm57b5IyP79Sd+Souk/YOES5+gEk7QqIgogJIhGb4ACb4d7F82U/co9wkogOot1Gcg6gF74aGatXP4yAggEOAdKggEwNDRmMzU1YmRjYjdjYzBhZjcyOGVmM2NjZWI5NjE1ZDkwNjg0YmI1YjJjYTVmODU5YWIwZjBiNzA0MDc1ODcxYWEzODViNmIxYjhlYWQ4MDljYTY3NDU0ZDk2ODNmY2YyYmEwMzQ1NmQ2ZmUyYzRhYmUyYjA3ZjBmYmRiYjJmMWMxUiIKIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWAFgAWoMCMDvrrwGELj/8fUCagsIwe+uvAYQgM2AemoMCMHvrrwGEIjTsckBagwIwe+uvAYQrIGGnwJqDAjB7668BhD02ozIAmoMCMHvrrwGEID73fQCagwIwe+uvAYQ5ObFqAOAAQeIAQIwAUgBCtsBCk9wR2E0ZTVoTkdzZ2tmakVHbTcyVEV2YkY3YVJEcUtCZDRMdVh0YWI0dWNNYlhMY2dKL0xhdGVzdFB1YmtleVRvVGlueUJsb2Nrc0NvdW50EocBCoIBMDQ0ZjM1NWJkY2I3Y2MwYWY3MjhlZjNjY2ViOTYxNWQ5MDY4NGJiNWIyY2E1Zjg1OWFiMGYwYjcwNDA3NTg3MWFhMzg1YjZiMWI4ZWFkODA5Y2E2NzQ1NGQ5NjgzZmNmMmJhMDM0NTZkNmZlMmM0YWJlMmIwN2YwZmJkYmIyZjFjMRABCsgECkxwR2E0ZTVoTkdzZ2tmakVHbTcyVEV2YkY3YVJEcUtCZDRMdVh0YWI0dWNNYlhMY2dKL1JvdW5kQmVmb3JlTGF0ZXN0RXhlY3V0aW9uEvcDCAES7gMKggEwNDRmMzU1YmRjYjdjYzBhZjcyOGVmM2NjZWI5NjE1ZDkwNjg0YmI1YjJjYTVmODU5YWIwZjBiNzA0MDc1ODcxYWEzODViNmIxYjhlYWQ4MDljYTY3NDU0ZDk2ODNmY2YyYmEwMzQ1NmQ2ZmUyYzRhYmUyYjA3ZjBmYmRiYjJmMWMxEuYCCAEQASIiCiDdCtNDbiIowyY4GbntvkjI/v1J35Ki6T9g4RLn6ASTtCoiCiAkiEZvgAJvh3sXzZT9yj3CSiA6i3UZyDqAXvhoZq1c/jICCAQ4B0qCATA0NGYzNTViZGNiN2NjMGFmNzI4ZWYzY2NlYjk2MTVkOTA2ODRiYjViMmNhNWY4NTlhYjBmMGI3MDQwNzU4NzFhYTM4NWI2YjFiOGVhZDgwOWNhNjc0NTRkOTY4M2ZjZjJiYTAzNDU2ZDZmZTJjNGFiZTJiMDdmMGZiZGJiMmYxYzFSIgogAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABYAWABagwIwO+uvAYQuP/x9QJqCwjB7668BhCAzYB6agwIwe+uvAYQiNOxyQFqDAjB7668BhCsgYafAmoMCMHvrrwGEPTajMgCagwIwe+uvAYQgPvd9AJqDAjB7668BhDk5sWoA4ABB4gBAjABSAEKZgpAcEdhNGU1aE5Hc2drZmpFR203MlRFdmJGN2FSRHFLQmQ0THVYdGFiNHVjTWJYTGNnSi9SYW5kb21IYXNoZXMvOBIiCiDtZaBXo1+kSrpxLwECHZ9RuhcxMCI9zUvyMnkVocQ+SwpLCkZwR2E0ZTVoTkdzZ2tmakVHbTcyVEV2YkY3YVJEcUtCZDRMdVh0YWI0dWNNYlhMY2dKL0xhdGVzdEV4ZWN1dGVkSGVpZ2h0EgEQEkgKRHBHYTRlNWhOR3Nna2ZqRUdtNzJURXZiRjdhUkRxS0JkNEx1WHRhYjR1Y01iWExjZ0ovQ3VycmVudFJvdW5kTnVtYmVyEAESPgo6cEdhNGU1aE5Hc2drZmpFR203MlRFdmJGN2FSRHFLQmQ0THVYdGFiNHVjTWJYTGNnSi9Sb3VuZHMvMRABEkEKPXBHYTRlNWhOR3Nna2ZqRUdtNzJURXZiRjdhUkRxS0JkNEx1WHRhYjR1Y01iWExjZ0ovSXNNYWluQ2hhaW4QARJTCk9wR2E0ZTVoTkdzZ2tmakVHbTcyVEV2YkY3YVJEcUtCZDRMdVh0YWI0dWNNYlhMY2dKL0xhdGVzdFB1YmtleVRvVGlueUJsb2Nrc0NvdW50EAESUApMcEdhNGU1aE5Hc2drZmpFR203MlRFdmJGN2FSRHFLQmQ0THVYdGFiNHVjTWJYTGNnSi9Sb3VuZEJlZm9yZUxhdGVzdEV4ZWN1dGlvbhABEkQKQHBHYTRlNWhOR3Nna2ZqRUdtNzJURXZiRjdhUkRxS0JkNEx1WHRhYjR1Y01iWExjZ0ovUmFuZG9tSGFzaGVzLzcQARJECkBwR2E0ZTVoTkdzZ2tmakVHbTcyVEV2YkY3YVJEcUtCZDRMdVh0YWI0dWNNYlhMY2dKL1JhbmRvbUhhc2hlcy84EAESSgpGcEdhNGU1aE5Hc2drZmpFR203MlRFdmJGN2FSRHFLQmQ0THVYdGFiNHVjTWJYTGNnSi9MYXRlc3RFeGVjdXRlZEhlaWdodBABGusGCiIKIJQSvf6ha65o1SJuVxZv3z1N0ZtJ3joScrcOZHSOykYnIqEBCiIKIEHhX/4StnOVIjpvGutfDK0HGiL26HpvAaHpycr2iY07EiIKICeR6ZKlfyjnWhHxOvLArsiw6zXS8EjULrqJAckuA3jcGAcqFUNoYXJnZVRyYW5zYWN0aW9uRmVlczI+ChNEb25hdGVSZXNvdXJjZVRva2VuEiIKICeR6ZKlfyjnWhHxOvLArsiw6zXS8EjULrqJAckuA3jcGKDkwgwqfgoiCiAltq3/QBU35Yj5PKR2ZoXjX3pupniJel4M1uQpInZ+5RICCAFYm31gAWpPEk0KSUpSbUJkdWg0blhXaTFhWGdkVXNqNWdKcnplWmIyTHhtckFiZjdXOTlmYVpTdm9BYUUvQ2hhaW5QcmltYXJ5VG9rZW5TeW1ib2wQAViSpAFgAWqaBAp5ClNKUm1CZHVoNG5YV2kxYVhnZFVzajVnSnJ6ZVpiMkx4bXJBYmY3Vzk5ZmFaU3ZvQWFFL0xhdGVzdFRvdGFsUmVzb3VyY2VUb2tlbnNNYXBzSGFzaBIiCiBcOi1ZaO4gPsHjgGSMcSRB609M3ip9cGpJZdGktp43DgpXClJKUm1CZHVoNG5YV2kxYVhnZFVzajVnSnJ6ZVpiMkx4bXJBYmY3Vzk5ZmFaU3ZvQWFFL0RvbmF0ZVJlc291cmNlVG9rZW5FeGVjdXRlSGVpZ2h0EgESEkcKQ0pSbUJkdWg0blhXaTFhWGdkVXNqNWdKcnplWmIyTHhtckFiZjdXOTlmYVpTdm9BYUUvQ29uc2Vuc3VzQ29udHJhY3QQARJKCkZKUm1CZHVoNG5YV2kxYVhnZFVzajVnSnJ6ZVpiMkx4bXJBYmY3Vzk5ZmFaU3ZvQWFFL0RpdmlkZW5kUG9vbENvbnRyYWN0EAESVwpTSlJtQmR1aDRuWFdpMWFYZ2RVc2o1Z0pyemVaYjJMeG1yQWJmN1c5OWZhWlN2b0FhRS9MYXRlc3RUb3RhbFJlc291cmNlVG9rZW5zTWFwc0hhc2gQARJWClJKUm1CZHVoNG5YV2kxYVhnZFVzajVnSnJ6ZVpiMkx4bXJBYmY3Vzk5ZmFaU3ZvQWFFL0RvbmF0ZVJlc291cmNlVG9rZW5FeGVjdXRlSGVpZ2h0EAEioA4KUApJSlJtQmR1aDRuWFdpMWFYZ2RVc2o1Z0pyemVaYjJMeG1yQWJmN1c5OWZhWlN2b0FhRS9DaGFpblByaW1hcnlUb2tlblN5bWJvbBIDRUxGCkkKRHBHYTRlNWhOR3Nna2ZqRUdtNzJURXZiRjdhUkRxS0JkNEx1WHRhYjR1Y01iWExjZ0ovQ3VycmVudFJvdW5kTnVtYmVyEgECCkIKPXBHYTRlNWhOR3Nna2ZqRUdtNzJURXZiRjdhUkRxS0JkNEx1WHRhYjR1Y01iWExjZ0ovSXNNYWluQ2hhaW4SAQEKSwpGcEdhNGU1aE5Hc2drZmpFR203MlRFdmJGN2FSRHFLQmQ0THVYdGFiNHVjTWJYTGNnSi9MYXRlc3RFeGVjdXRlZEhlaWdodBIBDgrbAQpPcEdhNGU1aE5Hc2drZmpFR203MlRFdmJGN2FSRHFLQmQ0THVYdGFiNHVjTWJYTGNnSi9MYXRlc3RQdWJrZXlUb1RpbnlCbG9ja3NDb3VudBKHAQqCATA0NGYzNTViZGNiN2NjMGFmNzI4ZWYzY2NlYjk2MTVkOTA2ODRiYjViMmNhNWY4NTlhYjBmMGI3MDQwNzU4NzFhYTM4NWI2YjFiOGVhZDgwOWNhNjc0NTRkOTY4M2ZjZjJiYTAzNDU2ZDZmZTJjNGFiZTJiMDdmMGZiZGJiMmYxYzEQAgpmCkBwR2E0ZTVoTkdzZ2tmakVHbTcyVEV2YkY3YVJEcUtCZDRMdVh0YWI0dWNNYlhMY2dKL1JhbmRvbUhhc2hlcy83EiIKICuu7iuDQt4NeB3jUep8LaROEkRu3QidA0h4BAxSNcALCkIKQHBHYTRlNWhOR3Nna2ZqRUdtNzJURXZiRjdhUkRxS0JkNEx1WHRhYjR1Y01iWExjZ0ovUmFuZG9tSGFzaGVzLzgKugQKTHBHYTRlNWhOR3Nna2ZqRUdtNzJURXZiRjdhUkRxS0JkNEx1WHRhYjR1Y01iWExjZ0ovUm91bmRCZWZvcmVMYXRlc3RFeGVjdXRpb24S6QMIARLgAwqCATA0NGYzNTViZGNiN2NjMGFmNzI4ZWYzY2NlYjk2MTVkOTA2ODRiYjViMmNhNWY4NTlhYjBmMGI3MDQwNzU4NzFhYTM4NWI2YjFiOGVhZDgwOWNhNjc0NTRkOTY4M2ZjZjJiYTAzNDU2ZDZmZTJjNGFiZTJiMDdmMGZiZGJiMmYxYzES2AIIARABIiIKIN0K00NuIijDJjgZue2+SMj+/UnfkqLpP2DhEufoBJO0KiIKICSIRm+AAm+HexfNlP3KPcJKIDqLdRnIOoBe+GhmrVz+MgIIBDgGSoIBMDQ0ZjM1NWJkY2I3Y2MwYWY3MjhlZjNjY2ViOTYxNWQ5MDY4NGJiNWIyY2E1Zjg1OWFiMGYwYjcwNDA3NTg3MWFhMzg1YjZiMWI4ZWFkODA5Y2E2NzQ1NGQ5NjgzZmNmMmJhMDM0NTZkNmZlMmM0YWJlMmIwN2YwZmJkYmIyZjFjMVIiCiAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFgBYAFqDAjA7668BhC4//H1AmoLCMHvrrwGEIDNgHpqDAjB7668BhCI07HJAWoMCMHvrrwGEKyBhp8CagwIwe+uvAYQ9NqMyAJqDAjB7668BhCA+930AoABBogBAjABSAEKqAQKOnBHYTRlNWhOR3Nna2ZqRUdtNzJURXZiRjdhUkRxS0JkNEx1WHRhYjR1Y01iWExjZ0ovUm91bmRzLzES6QMIARLgAwqCATA0NGYzNTViZGNiN2NjMGFmNzI4ZWYzY2NlYjk2MTVkOTA2ODRiYjViMmNhNWY4NTlhYjBmMGI3MDQwNzU4NzFhYTM4NWI2YjFiOGVhZDgwOWNhNjc0NTRkOTY4M2ZjZjJiYTAzNDU2ZDZmZTJjNGFiZTJiMDdmMGZiZGJiMmYxYzES2AIIARABIiIKIN0K00NuIijDJjgZue2+SMj+/UnfkqLpP2DhEufoBJO0KiIKICSIRm+AAm+HexfNlP3KPcJKIDqLdRnIOoBe+GhmrVz+MgIIBDgGSoIBMDQ0ZjM1NWJkY2I3Y2MwYWY3MjhlZjNjY2ViOTYxNWQ5MDY4NGJiNWIyY2E1Zjg1OWFiMGYwYjcwNDA3NTg3MWFhMzg1YjZiMWI4ZWFkODA5Y2E2NzQ1NGQ5NjgzZmNmMmJhMDM0NTZkNmZlMmM0YWJlMmIwN2YwZmJkYmIyZjFjMVIiCiAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFgBYAFqDAjA7668BhC4//H1AmoLCMHvrrwGEIDNgHpqDAjB7668BhCI07HJAWoMCMHvrrwGEKyBhp8CagwIwe+uvAYQ9NqMyAJqDAjB7668BhCA+930AoABBogBAjABSAEigwkKUApJSlJtQmR1aDRuWFdpMWFYZ2RVc2o1Z0pyemVaYjJMeG1yQWJmN1c5OWZhWlN2b0FhRS9DaGFpblByaW1hcnlUb2tlblN5bWJvbBIDRUxGCmkKQ0pSbUJkdWg0blhXaTFhWGdkVXNqNWdKcnplWmIyTHhtckFiZjdXOTlmYVpTdm9BYUUvQ29uc2Vuc3VzQ29udHJhY3QSIgoga1Qr/CdR220Mwa70uLBxsTSAyUfcmyY7gZH7gQRQZD0KbApGSlJtQmR1aDRuWFdpMWFYZ2RVc2o1Z0pyemVaYjJMeG1yQWJmN1c5OWZhWlN2b0FhRS9EaXZpZGVuZFBvb2xDb250cmFjdBIiCiApt8V4G15ozIdPMRxVf8x/KUnRLPB0vBy54ZazL8b95ApXClJKUm1CZHVoNG5YV2kxYVhnZFVzajVnSnJ6ZVpiMkx4bXJBYmY3Vzk5ZmFaU3ZvQWFFL0RvbmF0ZVJlc291cmNlVG9rZW5FeGVjdXRlSGVpZ2h0EgEQCnkKU0pSbUJkdWg0blhXaTFhWGdkVXNqNWdKcnplWmIyTHhtckFiZjdXOTlmYVpTdm9BYUUvTGF0ZXN0VG90YWxSZXNvdXJjZVRva2Vuc01hcHNIYXNoEiIKIKBwuXw0/Cr3gPJ3g5jn6ucLheMcmAUQk2apXKiAxvs/CkkKRHBHYTRlNWhOR3Nna2ZqRUdtNzJURXZiRjdhUkRxS0JkNEx1WHRhYjR1Y01iWExjZ0ovQ3VycmVudFJvdW5kTnVtYmVyEgECCrYECjpwR2E0ZTVoTkdzZ2tmakVHbTcyVEV2YkY3YVJEcUtCZDRMdVh0YWI0dWNNYlhMY2dKL1JvdW5kcy8xEvcDCAES7gMKggEwNDRmMzU1YmRjYjdjYzBhZjcyOGVmM2NjZWI5NjE1ZDkwNjg0YmI1YjJjYTVmODU5YWIwZjBiNzA0MDc1ODcxYWEzODViNmIxYjhlYWQ4MDljYTY3NDU0ZDk2ODNmY2YyYmEwMzQ1NmQ2ZmUyYzRhYmUyYjA3ZjBmYmRiYjJmMWMxEuYCCAEQASIiCiDdCtNDbiIowyY4GbntvkjI/v1J35Ki6T9g4RLn6ASTtCoiCiAkiEZvgAJvh3sXzZT9yj3CSiA6i3UZyDqAXvhoZq1c/jICCAQ4B0qCATA0NGYzNTViZGNiN2NjMGFmNzI4ZWYzY2NlYjk2MTVkOTA2ODRiYjViMmNhNWY4NTlhYjBmMGI3MDQwNzU4NzFhYTM4NWI2YjFiOGVhZDgwOWNhNjc0NTRkOTY4M2ZjZjJiYTAzNDU2ZDZmZTJjNGFiZTJiMDdmMGZiZGJiMmYxYzFSIgogAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABYAWABagwIwO+uvAYQuP/x9QJqCwjB7668BhCAzYB6agwIwe+uvAYQiNOxyQFqDAjB7668BhCsgYafAmoMCMHvrrwGEPTajMgCagwIwe+uvAYQgPvd9AJqDAjB7668BhDk5sWoA4ABB4gBAjABSAE="
	data, err := base64.StdEncoding.DecodeString(blockSample)
	assert.NoError(t, err)
	var blk aelf.Block
	err = proto.Unmarshal(data, &blk)
	assert.NoError(t, err)
	newBlck := ConvertBlock("0119d9910a9d3e3b114679b262974cc4f4a15f82990552b2585f5b56d424c47d", &blk)
	assert.Equal(t, int64(8), newBlck.Height)
	assert.Equal(t, "0119d9910a9d3e3b114679b262974cc4f4a15f82990552b2585f5b56d424c47d", newBlck.BlockHash)
}
