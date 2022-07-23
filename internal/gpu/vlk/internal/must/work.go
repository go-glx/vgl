package must

import (
	"fmt"
	"log"
	"runtime"

	"github.com/vulkan-go/vulkan"
)

// Work will panic when vkResult is not success
func Work(vkResult vulkan.Result) {
	if vkResult == vulkan.Success {
		return
	}

	panic(asGoError(vkResult, false))
}

// NotCare will do nothing when vkResult is success
// and log error, when is not.
// also return true when vkResult is success
func NotCare(vkResult vulkan.Result) bool {
	if vkResult == vulkan.Success {
		return true
	}

	log.Println(asGoError(vkResult, true).Error())
	return false
}

func asGoError(vkResult vulkan.Result, short bool) error {
	errorName := "unknown"
	description := ""

	if ref, ok := codeNameReferences[vkResult]; ok {
		errorName, description = ref[0], ref[1]
	}

	where := "unknown"
	if _, file, line, ok := runtime.Caller(2); ok {
		where = fmt.Sprintf("%s:%d", file, line)
	}

	if short {
		return fmt.Errorf("vk: Err: %d (%s), at %s",
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
