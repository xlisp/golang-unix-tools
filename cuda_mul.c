#include <stdio.h>
#include <dlfcn.h>  // For dynamic linking
#include <stdlib.h>

// Function pointer type for the CUDA function
typedef void (*launchMatrixMultiply_t)(float *, float *, float *, int);

int main() {
    // Define matrix dimensions and allocate memory
    int N = 4;
    float A[16] = {1, 2, 3, 4,
                   5, 6, 7, 8,
                   9, 10, 11, 12,
                   13, 14, 15, 16};
    
    float B[16] = {16, 15, 14, 13,
                   12, 11, 10, 9,
                   8, 7, 6, 5,
                   4, 3, 2, 1};

    float C[16] = {0};  // Output matrix

    // Load the shared library (libmatrix.so)
    void *handle = dlopen("./libmatrix.so", RTLD_LAZY);
    if (!handle) {
        fprintf(stderr, "Failed to load libmatrix.so: %s\n", dlerror());
        return EXIT_FAILURE;
    }

    // Clear existing errors
    dlerror();

    // Get the function pointer to the CUDA function
    launchMatrixMultiply_t launchMatrixMultiply = (launchMatrixMultiply_t) dlsym(handle, "launchMatrixMultiply");
    char *error = dlerror();
    if (error != NULL) {
        fprintf(stderr, "Error locating launchMatrixMultiply: %s\n", error);
        dlclose(handle);
        return EXIT_FAILURE;
    }

    // Call the CUDA matrix multiplication function
    launchMatrixMultiply(A, B, C, N);

    // Print the result matrix C
    printf("Matrix C (Result):\n");
    for (int i = 0; i < N; i++) {
        for (int j = 0; j < N; j++) {
            printf("%6.2f ", C[i * N + j]);
        }
        printf("\n");
    }

    // Close the shared library
    dlclose(handle);

    return EXIT_SUCCESS;
}

// run ok ------
// ➜ jim-emacs-fun-go (master) $ nvcc -arch=sm_52 --compiler-options '-fPIC' -o libmatrix.so --shared kernel.cu
// ➜ jim-emacs-fun-go (master) $ gcc -o main cuda_mul.c -ldl
// ➜ jim-emacs-fun-go (master) $ export LD_LIBRARY_PATH=.:$LD_LIBRARY_PATH
// ➜ jim-emacs-fun-go (master) $ ls
// channels2.go  fibonacci.go  ginhttp.go          go_ssh_rev3.go           go_vs_clojure_async.md  lmatrix.so  README.md
// channels.go   fp.go         go_call_c           go_ssh_rev5.go           hello.go                main        recursion.go
// cuda_mul.c    fun_args.go   go.mod              go_ssh_reverse_proxy.go  kernel.cu               matrix.so   show_fun_refs_project.go
// cuda_mul.go   functor.go    go_proxy_socket.go  go.sum                   libmatrix.so            poc_files   zshrc2md.go
// ➜ jim-emacs-fun-go (master) $ ./main
// Launch Matrix Multiply called with N=4
// Matrix C (Result):
//  80.00  70.00  60.00  50.00
// 240.00 214.00 188.00 162.00
// 400.00 358.00 316.00 274.00
// 560.00 502.00 444.00 386.00
// ➜ jim-emacs-fun-go (master) $
// 
