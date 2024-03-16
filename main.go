package main

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/tpyle/testamint/lib/types"
	"sigs.k8s.io/yaml"
)

func main() {
	fileName := "./examples/test-config.yaml"

	body, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%s\n", body)

	var config types.TestConfig
	// Load the test configuration from the file
	// j, err := yaml.YAMLToJSON(body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", j)
	err = yaml.Unmarshal(body, &config)
	if err != nil {
		panic(err)
	}

	fileParent := path.Dir(fileName)
	// fmt.Printf("File parent: %s\n", fileParent)
	err = os.Chdir(fileParent)
	if err != nil {
		panic(err)
	}

	for _, test := range config.Tests {

		fmt.Printf("%+v\n", test)
		fmt.Printf("%+v\n", test.Setup)
		for _, readyCheck := range test.ReadyChecks {
			fmt.Printf("%+v\n", readyCheck)
		}
		// fmt.Printf("%s\n", test.Runner)

		err = test.Setup.Validate()
		if err != nil {
			logrus.WithError(err).Fatal("Invalid Setup")
		}

		// setupContext, err := test.Setup.Setup(nil, os.Stdout)
		// if err != nil {
		// 	logrus.WithError(err).Fatal("Failed to setup environment")
		// }

		// fmt.Printf("%+v\n", setupContext)

		// fmt.Printf("%+v\n", test.ReadyChecks[0].Validate(setupContext, os.Stdout))

		// err = test.Setup.Teardown(os.Stdout)
		// if err != nil {
		// 	logrus.WithError(err).Fatal("Failed to teardown setup")
		// }

		// fmt.Printf("%s\n", test.Setup)

	}
}
