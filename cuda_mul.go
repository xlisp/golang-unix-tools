package main

/*
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcudart
#include <cuda_runtime.h>

void launchMatrixMultiply(float *A, float *B, float *C, int N);
*/
import "C"
import (
    "fmt"
    "unsafe"
)

func matrixMultiplyGo(A, B, C []float32, N int) {
    // Call the CUDA function
    C.launchMatrixMultiply((*C.float)(unsafe.Pointer(&A[0])),
        (*C.float)(unsafe.Pointer(&B[0])),
        (*C.float)(unsafe.Pointer(&C[0])), C.int(N))
}

func main() {
    N := 4
    A := []float32{
        1, 2, 3, 4,
        5, 6, 7, 8,
        9, 10, 11, 12,
        13, 14, 15, 16,
    }

    B := []float32{
        16, 15, 14, 13,
        12, 11, 10, 9,
        8, 7, 6, 5,
        4, 3, 2, 1,
    }

    C := make([]float32, N*N)

    matrixMultiplyGo(A, B, C, N)

    fmt.Println("Matrix C (result):")
    for i := 0; i < N; i++ {
        for j := 0; j < N; j++ {
            fmt.Printf("%6.2f ", C[i*N+j])
        }
        fmt.Println()
    }
}

/*
$ nvcc -arch=sm_52 --compiler-options '-fPIC' -o libmatrix.so --shared kernel.cu

# import "C" ## Modify the Go Code to Correctly Link to the CUDA Shared Library

$ Modify the Go Code to Correctly Link to the CUDA Shared Library

$ go run cuda_mul.go

# =>  /tmp/go-build/cuda_mul.cgo2.c:55:(.text+0x13): undefined reference to `launchMatrixMultiply'

*/
