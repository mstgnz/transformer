.PHONY: live
.DEFAULT_GOAL:= run

run:
	clear && go build -o /tmp/build ./example && /tmp/build

live:
	clear && find . -type f \( -name '*.go' \) | entr -r sh -c 'go build -o /tmp/build ./example && /tmp/build'