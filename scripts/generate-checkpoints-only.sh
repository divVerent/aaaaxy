set -ex

: ${GO:=go}

# Run go natively.
export GOOS=
export GOARCH=

for lfile in assets/maps/*.tmx; do
	lname=${lfile%.tmx}
	lname=${lname##*/}

	trap 'rm -f "assets/generated/$lname.cp.json"' EXIT
	# Using |cat> instead of > because snapcraft for some reason doesn't allow using a regular > shell redirection with "go run".
	${GO} run ${GO_FLAGS} github.com/divVerent/aaaaxy/cmd/dumpcps -level="$lname" |cat> "assets/generated/$lname.cp.dot"
	grep -c . "assets/generated/$lname.cp.dot"
	neato -Tjson assets/generated/$lname.cp.dot > assets/generated/$lname.cp.json
	grep -c . "assets/generated/$lname.cp.json"
	trap - EXIT
done