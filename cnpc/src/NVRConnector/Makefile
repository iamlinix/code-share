CC = g++
CFLAGS = -Wall -O2 -Iinclude -I/usr/include/mysql/ -Wno-literal-suffix -Wno-unused-value
LDFLAGS = -Llib -L/usr/lib/mysql
LDFLAGS += -lpthread -lHCCore -lhcnetsdk -lhpr -Wl,-rpath=/opt/tav/libs -lmysqlclient
SRC_INI = src/iniparser.c src/dictionary.c
SRC_NVR = src/NVRConnector.cpp src/DBConnector.cpp src/util.c src/slog.c
SRC_PLATE = src/plate.cpp src/util.c src/slog.c
SRC_VIDEO = src/video.cpp src/util.c src/slog.c

TARGET=nvrconn
TARGET_PLATE=plate
TARGET_VIDEO=video

all: nvr plate video

nvr:
	$(CC) $(CFLAGS) $(SRC_NVR)  $(LDFLAGS) -o $(TARGET)

plate:
	$(CC) $(CFLAGS) $(SRC_PLATE) $(SRC_INI) $(LDFLAGS) -o $(TARGET_PLATE)

video:
	$(CC) $(CFLAGS) $(SRC_VIDEO) $(SRC_INI) $(LDFLAGS) -o $(TARGET_VIDEO)

clean:
	rm -f *.o $(TARGET) $(TARGET_PLATE) $(TARGET_VIDEO)

