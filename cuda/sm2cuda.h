#include <stdint.h>
#include <stdlib.h>
#include <string.h>
typedef uint32_t xint;
#define B256N 8
typedef xint B256D[B256N];

// start
void init_sm2cuda();
// num: to reach the best performance, num should be n*k, where n, k is a interger, and for Tesla K40, k is 270.
// please ask me the k value for different GPU card.
void sm2ver_cuda(xint* data, int num, xint* ret);

// 建议直接使用 sm2ver_cuda，在 go 里更方便
// sm2ver_cuda_human_readable 中写了如何生成 data ，建议在 go 里重写一遍 sm2ver_cuda_human_readable

struct sm2ver_cuda_req {
    B256D s, r, px, e; // px is PubKey.x
    // (s, r) is the signature
    // e is the Hash of the data
    xint ret; // 1: passed, 0: fail
};

inline void sm2ver_cuda_human_readable(struct sm2ver_cuda_req* req, int num) {
    xint* data = (xint*)malloc(sizeof(B256D[num*4]));
    for (int i=0;i<num;++i) {
        memcpy(data+i*4, req+i, sizeof(B256D[4]));
        // which equal to :
        // memcpy(data+i*4, &req[i].s, sizeof(B256D));
        // memcpy(data+i*4+1, &req[i].r, sizeof(B256D));
        // memcpy(data+i*4+2, &req[i].px, sizeof(B256D));
        // memcpy(data+i*4+3, &req[i].e, sizeof(B256D));
    }
    xint* ret = (xint*)malloc(sizeof(xint[num]));
    sm2ver_cuda(data, num, ret);
    for (int i=0;i<num;++i)
        req[i].ret = ret[i];
    free(data); free(ret);
}

