// Code generated by "stringer -linecomment -type UpsertOperation"; DO NOT EDIT.

package common

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UpsertOperationInsert-1]
	_ = x[UpsertOperationUpdate-2]
}

const _UpsertOperation_name = "UpsertOperationInsertUpsertOperationUpdate"

var _UpsertOperation_index = [...]uint8{0, 21, 42}

func (i UpsertOperation) String() string {
	i -= 1
	if i >= UpsertOperation(len(_UpsertOperation_index)-1) {
		return "UpsertOperation(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _UpsertOperation_name[_UpsertOperation_index[i]:_UpsertOperation_index[i+1]]
}
