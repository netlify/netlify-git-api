include(CMakeForceCompiler)
# the name of the target operating system
SET(CMAKE_SYSTEM_NAME Darwin)

# The libgit2 CMakeFile.txt needs this to be set:
SET(CMAKE_SIZEOF_VOID_P 8)

# which compilers to use for C and C++
cmake_force_c_compiler(/osxcross/target/bin/x86_64-apple-darwin14-clang Clang)
cmake_force_cxx_compiler(/osxcross/target/bin/x86_64-apple-darwin14-clang++ Clang)
SET(CMAKE_AR /osxcross/target/bin/x86_64-apple-darwin14-ar CACHE FILEPATH "Archiver")
SET(PKG_CONFIG_EXECUTABLE /osxcross/target/bin/x86_64h-apple-darwin14-pkg-config)

SET(CMAKE_OSX_SYSROOT /osxcross/target/SDK/MacOSX10.10.sdk)

# here is the target environment located
#SET(CMAKE_FIND_ROOT_PATH ${CMAKE_OSX_SYSROOT} ${CMAKE_OSX_SYSROOT}/usr/bin)
SET(CMAKE_FIND_ROOT_PATH /osxcross/target/macports/pkgs/opt/local ${CMAKE_OSX_SYSROOT} ${CMAKE_OSX_SYSROOT}/usr/bin )

# adjust the default behaviour of the FIND_XXX() commands:
# search headers and libraries in the target environment, search
# programs in the host environment
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
