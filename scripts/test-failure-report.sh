#!/bin/bash

if ! [ -f ./coverage/integration.report ]; then
  echo "File coverage/integration.report does not exist. Skipping."
else
  echo "Creating integration test failures report..."
  grep 'FAIL:\|FAIL\|Error:\|error:' coverage/integration.report > coverage/integration.failures
fi

if ! [ -f ./coverage/unit.report ]; then
  echo "File coverage/unit.report does not exist. Skipping."
else
  grep 'FAIL:\|FAIL\|Error:\|error:' coverage/unit.report > coverage/unit.failures
fi

