/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package keywordreplace

import (
	"fmt"
	//	"io/ioutil"
	"testing"
)

func TestKeywordreplace(t *testing.T) {
	log.Info("Test entering .... ")
	/*	fileContent, err := ioutil.ReadFile("./descriptor.json")
		if err != nil {
			log.Debug("File reading error", err)
			return
		}
	*/
	fileContent := "zzz$XXXX$zzz${ABC}/cd/se"
	mapper := NewKeywordMapper(string(fileContent), "$", "$")
	document := mapper.Replace("", map[string]interface{}{
		"ContainerName": "XXXXXXXX",
	})
	log.Info("document = ", document)

}
