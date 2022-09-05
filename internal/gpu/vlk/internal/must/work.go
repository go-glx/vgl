package must

import (
	"fmt"
	"runtime"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/vgl/shared/vlkext"
)

// Work will panic when vkResult is not success
func Work(vkResult vulkan.Result) {
	if vkResult == vulkan.Success {
		return
	}

	panic(asGoError(vkResult, false, 2))
}

// NotCare will do nothing when vkResult is success
// and log error, when is not.
// also return true when vkResult is success
// NotCare MUST be called from some util function, not directly
// Code -> utilFn -> NotCare -> goError (3)
// This will show n-3 stack trace line, where code is halted
func NotCare(logger vlkext.Logger, vkResult vulkan.Result) bool {
	if vkResult == vulkan.Success {
		return true
	}

	logger.Notice(asGoError(vkResult, true, 3).Error())
	return false
}

func asGoError(vkResult vulkan.Result, short bool, stackLevel int) error {
	errorName := "unknown"
	description := ""

	if ref, ok := codeNameReferences[vkResult]; ok {
		errorName, description = ref[0], ref[1]
	}

	where := "unknown"
	if _, file, line, ok := runtime.Caller(stackLevel); ok {
		where = fmt.Sprintf("%s:%d", file, line)
	}

	if short {
		return fmt.Errorf("bad result, code=%d (%s), at %s",
			vkResult,
			errorName,
			where,
		)
	}

	return fmt.Errorf(
		"\n----------- VK Error -----------\n"+
			"  code: %d\n"+
			"   err: %s\n"+
			"    at: %s\n"+
			"reason: %s\n"+
			"  more: %s\n"+
			"----------- ^ -----------\n\n",
		vkResult,
		errorName,
		where,
		description,
		errorCodesUrl,
	)
}
