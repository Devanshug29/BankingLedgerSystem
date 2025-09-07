#!/bin/zsh

PROJECT_ROOT=$(git rev-parse --show-toplevel)
ABS_FILE_PATH=$(realpath "$1")
# Remove common part
RELATIVE_PATH="${ABS_FILE_PATH#$PROJECT_ROOT/}"
RELATIVE_FOLDER_PATH=$(dirname "$RELATIVE_PATH")
FILENAME=$(basename "$RELATIVE_PATH")

# append _mock in the file name
MOCK_FILE_PATH=$(echo $PROJECT_ROOT/mocks/"$RELATIVE_FOLDER_PATH"_mocks/$FILENAME | sed 's/\.go$/_mock.go/')

mockgen -package $(echo $GOPACKAGE"_mocks") -source $GOFILE -destination $MOCK_FILE_PATH

#remove absolute path from auto generated comments
ESCAPED_PROJECT_ROOT=$(echo "$PROJECT_ROOT" | sed 's/[\/&]/\\&/g')
sed -i '' "s!${ESCAPED_PROJECT_ROOT}/!!g" "$MOCK_FILE_PATH"
