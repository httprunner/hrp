package hrp

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	json "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func loadFromJSON(path string) (*TCase, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("convert absolute path failed")
		return nil, err
	}
	log.Info().Str("path", path).Msg("load json testcase")

	file, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Msg("load json path failed")
		return nil, err
	}

	ts := &TStep{}
	tc := &TCase{}
	decoder := json.NewDecoder(bytes.NewReader(file))
	decoder.UseNumber()
	// compatible test scenario with a single API
	e := json.Get(file, "request")
	if e.LastError() == nil {
		err = decoder.Decode(ts)
		tc = ts.ToTCase()
	} else {
		err = decoder.Decode(tc)
	}
	return tc, err
}

func loadFromYAML(path string) (*TCase, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("convert absolute path failed")
		return nil, err
	}
	log.Info().Str("path", path).Msg("load yaml testcase")

	file, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Msg("load yaml path failed")
		return nil, err
	}

	ts := &TStep{}
	tc := &TCase{}
	// compatible test scenario with a single API
	err = yaml.Unmarshal(file, tc)
	if tc.TestSteps == nil {
		err = yaml.Unmarshal(file, ts)
		tc = ts.ToTCase()
	}
	return tc, err
}

func (tc *TCase) ToTestCase() (*TestCase, error) {
	testCase := &TestCase{
		Config: tc.Config,
	}
	for _, step := range tc.TestSteps {
		if step.API != "" {
			testCasePath := &TestCasePath{Path: step.API}
			tc, _ := testCasePath.ToTestCase()
			extendWithAPI(step, tc.TestSteps[0].ToStruct())
			testCase.TestSteps = append(testCase.TestSteps, &StepRequestWithOptionalArgs{
				step: step,
			})
		} else if step.Request != nil {
			testCase.TestSteps = append(testCase.TestSteps, &StepRequestWithOptionalArgs{
				step: step,
			})
		} else if step.TestCase != nil {
			if reflect.TypeOf(step.TestCase) == reflect.TypeOf("") {
				testCasePath := &TestCasePath{Path: step.TestCase.(string)}
				tc, _ := testCasePath.ToTestCase()
				step.TestCase = tc
			}
			testCase.TestSteps = append(testCase.TestSteps, &StepTestCaseWithOptionalArgs{
				step: step,
			})
		} else if step.Transaction != nil {
			testCase.TestSteps = append(testCase.TestSteps, &StepTransaction{
				step: step,
			})
		} else if step.Rendezvous != nil {
			testCase.TestSteps = append(testCase.TestSteps, &StepRendezvous{
				step: step,
			})
		} else {
			log.Warn().Interface("step", step).Msg("[convertTestCase] unexpected step")
		}
	}
	return testCase, nil
}

var ErrUnsupportedFileExt = fmt.Errorf("unsupported testcase file extension")

// TestCasePath implements ITestCase interface.
type TestCasePath struct {
	Path string
}

func (path *TestCasePath) ToTestCase() (*TestCase, error) {
	var tc *TCase
	var err error

	casePath := path.Path
	ext := filepath.Ext(casePath)
	switch ext {
	case ".json":
		tc, err = loadFromJSON(casePath)
	case ".yaml", ".yml":
		tc, err = loadFromYAML(casePath)
	default:
		err = ErrUnsupportedFileExt
	}
	if err != nil {
		return nil, err
	}
	tc.Config.Path = path.Path
	testcase, err := tc.ToTestCase()
	if err != nil {
		return nil, err
	}
	return testcase, nil
}

func (path *TestCasePath) ToTCase() (*TCase, error) {
	testcase, err := path.ToTestCase()
	if err != nil {
		return nil, err
	}
	return testcase.ToTCase()
}
