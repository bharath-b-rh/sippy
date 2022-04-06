package testgridconversion

import (
	"strings"
	"testing"

	testgridv1 "github.com/openshift/sippy/pkg/apis/testgrid/v1"
	"github.com/openshift/sippy/pkg/testgridanalysis/testgridanalysisapi"
	"github.com/stretchr/testify/assert"
)

func TestProcessJobDetails(t *testing.T) {
	// need to generate JobDetails with a mix of expected random and non random test names
	// validate the RawJobResult has the corrected names

	testNames := []string{
		"Operator results test operator install install_operatorname",
		testgridanalysisapi.OperatorUpgradePrefix + "upgrade_operatorname",
		"\"Installing \"Red Hat Integration - 3scale\" operator in test-nbqyx.Installing \"Red Hat Integration - 3scale\" operator in test-nbqyx Installs Red Hat Integration - 3scale operator in test-nbqyx and creates 3scale Backend Schema operand instance\"",
		"This test name should not be modified",
	}

	validationStrings := []string{
		"Operator results.operator conditions  install_operatorname",
		"Operator results.operator conditions  upgrade_operatorname",
		"\"Installing \"Red Hat Integration - 3scale\" operator in test-random.Installing \"Red Hat Integration - 3scale\" operator in test-random Installs Red Hat Integration - 3scale operator in test-random and creates 3scale Backend Schema operand instance\"",
		"This test name should not be modified",
	}

	result := processJobDetails(buildFakeJobDetails(testNames), 0, 1)

	assert.NotNil(t, result, "Nil response from processJobDetails")

	// check the keys of the map and validate they match our expectations
	assert.Equal(t, len(result.TestResults), len(testNames), "Unexpected test resulsts size %d", len(result.TestResults))

	for _, s := range validationStrings {
		assert.NotNil(t, result.TestResults[s], "Expected non nil test result for %s", s)
	}

}

func TestCleanTestName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		vercount int
		output   string
	}{
		{
			name:     "generic 1",
			input:    "\"Installing \"Red Hat Integration - 3scale\" operator in test-nbqyx.Installing \"Red Hat Integration - 3scale\" operator in test-nbqyx Installs Red Hat Integration - 3scale operator in test-nbqyx and creates 3scale Backend Schema operand instance\"",
			vercount: 3,
		},
		{
			name:     "generic 2",
			input:    "\"Installing \"Red Hat Integration - 3scale\" operator in test-nsyin.Installing \"Red Hat Integration - 3scale\" operator in test-nsyin Installs Red Hat Integration - 3scale operator in test-nsyin and creates 3scale Backend Schema operand instance\"",
			vercount: 3,
		},
		{
			name:     "generic 3",
			input:    "\"Installing \"Red Hat Integration - 3scale\" operator in test-piiov.Installing \"Red Hat Integration - 3scale\" operator in test-piiov \"after all\" hook for \"Installs Red Hat Integration - 3scale operator in test-piiov and creates 3scale Backend Schema operand i (...)\"",
			vercount: 3,
		},
		{
			name:     "generic 4",
			input:    "\"Installing \"Red Hat Integration - 3scale\" operator in test-piiov.Installing \"Red Hat Integration - 3scale\" operator in test-piiov Installs Red Hat Integration - 3scale operator in test-piiov and creates 3scale Backend Schema operand instance\"",
			vercount: 3,
		},
		{
			name:     "skip 1",
			input:    "\"Doesn'tStartWith Installing \"Red Hat Integration - 3scale\" operator in test-.Installing \"Red Hat Integration - 3scale\" operator in test- Installs Red Hat Integration - 3scale operator in test- and creates 3scale Backend Schema operand instance\"",
			vercount: 0,
		},
		{
			name:     "non match 1",
			input:    "\"Installing \"Red Hat Integration - 3scale\" operator in test-.Installing \"Red Hat Integration - 3scale\" operator in test- Installs Red Hat Integration - 3scale operator in test- and creates 3scale Backend Schema operand instance\"",
			vercount: 0,
		},
		{
			name:     "verify output 1",
			input:    "\"Installing \"Red Hat Integration - 3scale\" operator in test-ieesa.Installing \"Red Hat Integration - 3scale\" operator in test-ieesa Installs Red Hat Integration - 3scale operator in test-ieesa and creates 3scale Backend Schema operand instance\"",
			output:   "\"Installing \"Red Hat Integration - 3scale\" operator in test namespace.Installing \"Red Hat Integration - 3scale\" operator in test namespace Installs Red Hat Integration - 3scale operator in test namespace and creates 3scale Backend Schema operand instance\"",
			vercount: 3,
		},
		{
			name:     "verify output 2",
			input:    "\"Installing \"Red Hat Integration - 3scale\" operator in test-jopkv.Installing \"Red Hat Integration - 3scale\" operator in test-jopkv \"after all\" hook for \"Installs Red Hat Integration - 3scale operator in test-jopkv and creates 3scale Backend Schema operand i (...)\"",
			output:   "\"Installing \"Red Hat Integration - 3scale\" operator in test namespace.Installing \"Red Hat Integration - 3scale\" operator in test namespace \"after all\" hook for \"Installs Red Hat Integration - 3scale operator in test namespace and creates 3scale Backend Schema operand i (...)\"",
			vercount: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sCleaned := cleanTestName(test.input)
			fixCount := strings.Count(sCleaned, matchRandomReplace)
			assert.Equal(t, fixCount, test.vercount, "Invalid verification count %d for test %s", fixCount, test.name)

			if test.output != "" {
				assert.Equal(t, test.output, sCleaned, "Cleaned output did not match expected %s", sCleaned)
			}
		})
	}
}

func buildFakeJobDetails(testNames []string) testgridv1.JobDetails {

	status1 := testgridv1.TestResult{
		Count: 1,
		Value: testgridv1.TestStatusFailure,
	}

	statuses := []testgridv1.TestResult{status1}
	tests := []testgridv1.Test{}

	for _, s := range testNames {

		test := testgridv1.Test{
			Name:     s,
			Statuses: statuses,
		}

		tests = append(tests, test)

	}

	jobDetails := &testgridv1.JobDetails{
		Name:        "mockName",
		Tests:       tests,
		Query:       "mockQuery",
		ChangeLists: []string{"mockChange"},
		Timestamps:  []int{1},
	}

	return *jobDetails
}
