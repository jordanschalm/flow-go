// Code generated by "stringer -type=ContainerKind"; DO NOT EDIT.

package sema

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ContainerKindUnknown-0]
	_ = x[ContainerKindInterface-1]
	_ = x[ContainerKindComposite-2]
}

const _ContainerKind_name = "ContainerKindUnknownContainerKindInterfaceContainerKindComposite"

var _ContainerKind_index = [...]uint8{0, 20, 42, 64}

func (i ContainerKind) String() string {
	if i >= ContainerKind(len(_ContainerKind_index)-1) {
		return "ContainerKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ContainerKind_name[_ContainerKind_index[i]:_ContainerKind_index[i+1]]
}