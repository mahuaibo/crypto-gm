#ifndef HYPERCHAIN_GM_sm93_H
#define HYPERCHAIN_GM_sm93_H

#endif //HYPERCHAIN_GM_SM3_H


#ifndef HEADER_sm93_H
#define HEADER_sm93_H
#ifndef NO_GMSSL

#define sm93_DIGEST_LENGTH	32
#define sm93_BLOCK_SIZE		64
#define sm93_CBLOCK		(sm93_BLOCK_SIZE)
#define sm93_HMAC_SIZE		(sm93_DIGEST_LENGTH)


#include <sys/types.h>
#include <stdint.h>
#include <string.h>

#ifdef __cplusplus
extern "C" {
#endif


typedef struct {
	uint32_t digest[8];
	int nblocks;
	unsigned char block[64];
	int num;
} sm93_ctx_t;

void sm93_init(sm93_ctx_t *ctx);
void sm93_update(sm93_ctx_t *ctx, const unsigned char* data, size_t data_len);
void sm93_final(sm93_ctx_t *ctx, unsigned char digest[sm93_DIGEST_LENGTH]);
void sm93_compress(uint32_t digest[8], const unsigned char block[sm93_BLOCK_SIZE]);
void sm93(const unsigned char *data, size_t datalen,
	unsigned char digest[sm93_DIGEST_LENGTH]);


// typedef struct {
// 	sm3_ctx_t sm3_ctx;
// 	unsigned char key[SM3_BLOCK_SIZE];
// } sm3_hmac_ctx_t;

// void sm3_hmac_init(sm3_hmac_ctx_t *ctx, const unsigned char *key, size_t key_len);
// void sm3_hmac_update(sm3_hmac_ctx_t *ctx, const unsigned char *data, size_t data_len);
// void sm3_hmac_final(sm3_hmac_ctx_t *ctx, unsigned char mac[SM3_HMAC_SIZE]);
// void sm3_hmac(const unsigned char *data, size_t data_len,
// 	const unsigned char *key, size_t key_len, unsigned char mac[SM3_HMAC_SIZE]);

#ifdef __cplusplus
}
#endif
#endif
#endif

