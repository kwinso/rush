SRC=./cmd
TARGET=rush

all: $(TARGET) run

$(TARGET): $(SRC)
	go build -o $(TARGET) $(SRC)

run: $(TARGET)
	./$(TARGET)

clean:
	rm -rf ./$(TARGET)