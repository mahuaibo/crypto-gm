
#ifndef SM9_SV_INCLUE_H__
#define SM9_SV_INCLUE_H__

#include "R-ate.h"

/************************************************************************
// File name: SM9_sv.h
// Version: SM9_sv_V1.0
// Date: Dec 15,2016
// Description: implementation of SM9 signature algorithm and verification algorithm
// all operations based on BN curve line function
// Function List:
// 1.bytes128_to_ecn2 //convert 128 bytes into ecn2
// 2.zzn12_ElementPrint //print all element of struct zzn12
// 3.ecn2_Bytes128_Print //print 128 bytes of ecn2
// 4.LinkCharZzn12 //link two different types(unsigned char and zzn12)to one(unsigned char)
// 5.Test_Point //test if the given point is on SM9 curve
// 6.Test_Range //test if the big x belong to the range[1,N-1]
// 7.SM9_Init //initiate SM9 curve
// 8.SM9_H1 //function H1 in SM9 standard 5.4.2.2
// 9.SM9_H2 //function H2 in SM9 standard 5.4.2.3
// 10.SM9_GenerateSignKey //generate signed private and public key
// 11.SM9_Sign //SM9 signature algorithm
// 12.SM9_Verify //SM9 verification
// 13.SM9_SelfCheck() //SM9 slef-check
//
// Notes:
// This SM9 implementation source code can be used for academic, non-profit making or non-commercial use only.
// This SM9 implementation is created on MIRACL. SM9 implementation source code provider does not provide MIRACL library, MIRACL license or any permission to use MIRACL library. Any commercial use of MIRACL requires a license which may be obtained from Shamus Software Ltd.
**************************************************************************/


//#include<sys/malloc.h>
//#miracl_custom<sys/malloc.h> 在windows使用该头文件，注释上一行
#include<math.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "miracl.h"

#define BNLEN 32 //BN curve with 256bit is used in SM9 algorithm
#define SM9_ASK_MEMORY_ERR 0x00000001 //内存申请失败
#define SM9_H_OUTRANGE 0x00000002 //签名H不属于[1,N-1]
#define SM9_DATA_MEMCMP_ERR 0x00000003 //数据对比不一致
#define SM9_MEMBER_ERR 0x00000004 //群的阶错误
#define SM9_MY_ECAP_12A_ERR 0x00000005 //R-ate对计算出现错误
#define SM9_S_NOT_VALID_G1 0x00000006 //S不属于群G1
#define SM9_G1BASEPOINT_SET_ERR 0x00000007 //G1基点设置错误
#define SM9_G2BASEPOINT_SET_ERR 0x00000008 //G2基点设置错误
#define SM9_L_error 0x00000009 //参数L错误
#define SM9_GEPUB_ERR 0x0000000A //生成公钥错误
#define SM9_GEPRI_ERR 0x0000000B //生成私钥错误
#define SM9_SIGN_ERR 0x0000000C //签名错误

unsigned char * SM9_Setup(unsigned char * ksByte);
//BOOL bytes128_to_ecn2(unsigned char Ppubs[], ecn2 *res);
//void zzn12_ElementPrint(zzn12 x);
//void ecn2_Bytes128_Print(ecn2 x);
int SM9_Init();
//int Test_Range(big x);
//static int SM9_H1(unsigned char Z[], int Zlen, big n, big h1);
//static int SM9_H2(unsigned char Z[], int Zlen, big n, big h2);
unsigned char * SM9_GenerateSignKey(char *ID, int IDlen, unsigned char ksBytes[]);
unsigned char * SM9_Sign( unsigned char *message, int len, unsigned char dsa[], unsigned char Ppub[]);
unsigned char * SM9_Sign_With_Rand( unsigned char *message, int len, big rand,
									unsigned char dsa[], unsigned char Ppub[]);
int SM9_Verify(unsigned char sign[], char *IDA, int IDlen, unsigned char *message,
			   int len, unsigned char Ppub[]);
int SM9_SelfCheck();

#ifdef __cplusplus
}
#endif

#endif