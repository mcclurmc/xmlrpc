package xmlrpc

import (
    "fmt"
    "regexp"
)

func parseResponse(response []byte) (result interface{}, err error) {
    if fault, _ := responseFailed(response); fault {
        return nil, parseFailedResponse(response)
    }

    return parseSuccessfulResponse(response)
}

func parseSuccessfulResponse(response []byte) (interface{}, error) {
    valueXml := getValueXml(response)
    return parseValue(valueXml)
}

// responseFailed checks whether response failed or not. Response defined as failed if it
// contains <fault>...</fault> section.
func responseFailed(response []byte) (bool, error) {
    fault := true
    faultRegexp, err := regexp.Compile(`<fault>(\s|\S)+</fault>`)

    if err == nil {
        fault = faultRegexp.Match(response)
    }

    return fault, err
}

func parseFailedResponse(response []byte) (err error) {
    var valueXml []byte
    valueXml = getValueXml(response)

    value, err := parseValue(valueXml)
    faultDetails := value.(Struct)

    if err != nil {
        return err
    }

    return &(xmlrpcError{
        code: fmt.Sprintf("%v", faultDetails["faultCode"]),
        message: faultDetails["faultString"].(string),
    })
}

func getValueXml(rawXml []byte) ([]byte) {
    expr, _ := regexp.Compile(`<value>(\s|\S)+</value>`)
    return expr.Find(rawXml)

}
