include(CMakeForceCompiler)
# the name of the target operating system
SET(CMAKE_SYSTEM_NAME Windows)

# The libgit2 CMakeFile.txt needs this to be set:
SET(CMAKE_SIZEOF_VOID_P 8)

# which compilers to use for C and C++
SET(CMAKE_C_COMPILER x86_64-w64-mingw32-gcc)
SET(CMAKE_CXX_COMPILER x86_64-w64-mingw32-cpp)
SET(CMAKE_FIND_ROOT_PATH /usr/x86_64-w64-mingw32)

SET(CMAKE_C_FLAGS -U__STRICT_ANSI__)

#cmake_force_c_compiler(/usr/bin/x86_64-w64-mingw32-gcc  MINGW)
#cmake_force_cxx_compiler(/usr/bin/x86_64-w64-mingw32-cpp MINGW)
SET(CMAKE_AR /usr/bin/x86_64-w64-mingw32-gcc-ar CACHE FILEPATH "Archiver")
SET(CMAKE_RC_COMPILER x86_64-w64-mingw32-windres)
SET(PKG_CONFIG_EXECUTABLE /usr/bin/x86_64-w64-mingw32-pkg-config)

SET(OPENSSL_ROOT_DIR ${CMAKE_FIND_ROOT_PATH})
SET(OPENSSL_INCLUDE_DIR ${CMAKE_FIND_ROOT_PATH}/include)
SET(OPENSSL_LIBRARIES ${CMAKE_FIND_ROOT_PATH}/lib)

# here is the target environment located
# SET(CMAKE_FIND_ROOT_PATH /usr/x86_64-w64-mingw32)

# adjust the default behaviour of the FIND_XXX() commands:
# search headers and libraries in the target environment, search
# programs in the host environment
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)

set(ENV{PKG_CONFIG_LIBDIR} ${CMAKE_FIND_ROOT_PATH}/lib/pkgconfig)
set(ENV{PKG_CONFIG_PATH} ${CMAKE_FIND_ROOT_PATH}/lib/pkgconfig)
