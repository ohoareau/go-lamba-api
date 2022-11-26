package gola

import (
	"fmt"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/ohoareau/gola/common"
	"github.com/ohoareau/gola/services"
	"github.com/ohoareau/gola/utils"
	"net/http"
	"os"
)

func Main(options common.Options) {
	if !utils.HasEnvVar("AWS_LAMBDA_FUNCTION_NAME") {
		port := utils.GetEnvVar("PORT", "5000")
		fmt.Println("🚀 Server ready at " + "http://localhost:" + port)
		err := http.ListenAndServe(":"+port, services.CreateHttpRouter(&options))
		if nil != err {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	h := services.CreateHandler(&options)

	if nil != options.HandlerWrapper {
		h = options.HandlerWrapper(h)
	}

	runtime.StartHandler(h)
}
