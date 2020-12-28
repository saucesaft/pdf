DEST := "gates"
EXE := "pdf"
ENTRY := "main.go"

if DEST == "torvalds" #LINUX

ifeq ($(DEST), jobs) #APPLE
	ECHO_MESSAGE = "Mac OS X"
endif

#ifeq ($(DEST), gates) #WINDOWS
ECHO_MESSAGE = "Windows"

$(info $(DEST))
FLAGS += GOOS=windows GOARCH=amd64 CGO_ENABLED=1
FLAGS += CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++

EXE := $(EXE).exe
#endif

##---------------------------------------------------------------------
## BUILD RULES
##---------------------------------------------------------------------

$(EXE):
	$(FLAGS) go build -o $@ $(ENTRY) $<

all: $(EXE)
	@echo Build complete for $(ECHO_MESSAGE)

run: all
	@echo running...
	./$(EXE)

clean:
	rm -rf $(EXE) $(BIN)
