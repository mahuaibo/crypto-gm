
#include "kdf.h"
#include "sm3.h"

/******************************************************************************
Function: SM3_KDF
Description: key derivation function
Calls: SM3_init
SM3_process
SM3_done
Called By:
Input: unsigned char Z[zlen]
unsigned short zlen //bytelen of Z
unsigned short klen //bytelen of K
Output: unsigned char K[klen] //shared secret key
Return: null
Others:
*******************************************************************************/
void sm93_KDF(unsigned char Z[], unsigned short zlen, unsigned short klen, unsigned char K[])
{
	unsigned short i, j, t;
	unsigned int bitklen;
	//SM3_STATE md;
	sm93_ctx_t ctx;
	unsigned char Ha[sm93_len / 8];
	unsigned char ct[4] = {0,0,0,1};
	bitklen = klen * 8;
	if( bitklen%sm93_len )
		t = bitklen / sm93_len + 1;
	else
		t = bitklen / sm93_len;
	//s4: K=Ha1||Ha2||...
	for( i = 1; i < t; i++ )
	{
		//s2: Hai=Hv(Z||ct)
		//SM3_init(&md);
		// SM3_process(&md, Z, zlen);
		// SM3_process(&md, ct, 4);
		// SM3_done(&md, Ha);
		sm93_init(&ctx);
    	sm93_update(&ctx,Z,zlen);
		sm93_update(&ctx,ct,4);
    	sm93_final(&ctx,Ha);
		memcpy((K + (sm93_len / 8)*(i - 1)), Ha, sm93_len / 8);
		if( ct[3] == 0xff )
		{
			ct[3] = 0;
			if( ct[2] == 0xff )
			{
				ct[2] = 0;
				if( ct[1] == 0xff )
				{
					ct[1] = 0;
					ct[0]++;
				} else ct[1]++;
			} else ct[2]++;
		} else ct[3]++;
	}
	//s3: klen/v
	sm93_init(&ctx);
    sm93_update(&ctx,Z,zlen);
	sm93_update(&ctx,ct,4);
	sm93_final(&ctx,Ha);
	if( bitklen%sm93_len )
	{
		i = (sm93_len - bitklen + sm93_len*(bitklen / sm93_len)) / 8;
		j = (bitklen - sm93_len*(bitklen / sm93_len)) / 8;
		memcpy((K + (sm93_len / 8)*(t - 1)), Ha, j);
	} else
	{
		memcpy((K + (sm93_len / 8)*(t - 1)), Ha, sm93_len / 8);
	}
}
