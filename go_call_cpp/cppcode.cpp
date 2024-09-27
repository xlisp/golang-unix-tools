// cppcode.cpp
#include <iostream>

// Function to be called from Go
extern "C" {
    int Add(int a, int b) {
        return a + b;
    }

    void PrintMessage() {
        std::cout << "Hello from C++!" << std::endl;
    }
}

