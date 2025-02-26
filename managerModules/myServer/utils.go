package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

func check(argErr error) {
	if argErr != nil {
		log.Printf("Error: %v", argErr)
	}
}