SRC=./cmd
TARGET=rush

all: $(TARGET) run

$(TARGET): $(SRC)
	go build -o build/$(TARGET) $(SRC)

run: $(TARGET)
	./build/$(TARGET)

clean:
	rm -rf ./build/$(TARGET)
