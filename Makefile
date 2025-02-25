# @file		Makefile
# @author	Sebastien LEGRAND
#
# @brief	Makefile to build the interpreter

.PHONY: clean build

build:
	@cd bin && go build -o brainfuck ../src/

clean:
	@cd bin && rm -f *
