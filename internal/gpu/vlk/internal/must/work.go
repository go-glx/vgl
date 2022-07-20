package must

import (
	"fmt"
	"runtime"

	"github.com/vulkan-go/vulkan"
)

func Work(resultCode vulkan.Result) {
	if resultCode == vulkan.Success {
		return
	}

	errorName := "unknown"
	description := ""

	if ref, ok := codeNameReferences[resultCode]; ok {
		errorName, description = ref[0], ref[1]
	}

	where := "unknown"
	if _, file, line, ok := runtime.Caller(1); ok {
		where = fmt.Sprintf("%s:%d", file, line)
	}

	panic(fmt.Errorf(
		"\n----------- VK Error -----------\n"+
			"  code: %d\n"+
			"   err: %s\n"+
			"    at: %s\n"+
			"reason: %s\n"+
			"  more: %s\n"+
			"----------- ^ -----------\n\n",
		resultCode,
		errorName,
		where,
		description,
		errorCodesUrl,
	))
}
