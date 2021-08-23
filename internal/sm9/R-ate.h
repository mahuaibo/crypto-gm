
#ifndef R_RATE_INCLUE_H__
#define R_RATE_INCLUE_H__

#pragma once


#include "zzn12_operation.h"

/************************************************************************
File name: R-ate.h
Version:
Date: Dec 15,2016
Description: this code is achieved according to ake12bnx.cpp in MIRCAL C++ source file.
see ake12bnx.cpp for details.
this code gives calculation of R-ate pairing
Function List:
1.zzn2_pow //regular zzn2 powering
2.set_frobenius_constant //calculate frobenius_constant X
3.q_power_frobenius
4.line
5.g
6.fast_pairing
7.ecap
Notes:
**************************************************************************/


//#include "../miracl_custom/zzn.h"

#ifdef __cplusplus
extern "C" {
#endif

zzn2 zzn2_pow(zzn2 x, big k);
void set_frobenius_constant(zzn2 *X);
void q_power_frobenius(ecn2 A, zzn2 F);
zzn12 line(ecn2 A, ecn2 *C, ecn2 *B, zzn2 slope, zzn2 extra, BOOL Doubling, big Qx, big Qy);
zzn12 g(ecn2 *A, ecn2 *B, big Qx, big Qy);
BOOL fast_pairing(ecn2 P, big Qx, big Qy, big x, zzn2 X, zzn12 *r);
BOOL ecap(ecn2 P, epoint *Q, big x, zzn2 X, zzn12 *r);
BOOL member(zzn12 r, big x, zzn2 F);

#ifdef __cplusplus
}
#endif

#endif