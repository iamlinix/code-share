cmake_minimum_required(VERSION 2.8)
set(CMAKE_CXX_FLAGS "-std=c++11 ${CMAKE_CXX_FLAGS}")
IF (WIN32)
    include_directories(
            C:\\Users\\linx\\AppData\\Local\\Packages\\CanonicalGroupLimited.UbuntuonWindows_79rhkp1fndgsc\\LocalState\\rootfs\\usr\\include
            C:\\Users\\linx\\AppData\\Local\\Packages\\CanonicalGroupLimited.UbuntuonWindows_79rhkp1fndgsc\\LocalState\\rootfs\\usr\\include\\x86_64-linux-gnu
            C:\\Users\\linx\\source\\repos\\libs\\fire\\include
            C:\\Users\\linx\\source\\repos\\libs\\mysql
    )
ELSE()
    include_directories(/usr/include/mysql)
ENDIF()

add_executable(petro china-petro.cpp db-util.cpp utf8.c business.cpp loguru.cpp)
target_link_libraries(petro mysqlclient fbclient dl pthread)