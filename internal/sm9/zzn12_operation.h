
#ifndef ZZN12_OPERATION_INCLUDE_H__
#define ZZN12_OPERATION_INCLUDE_H__

#pragma once


#ifdef __cplusplus
extern "C" {
#endif

#include "miracl.h"



/************************************************************************
File name: zzn12_operation.h
Version:
Date: Dec 15,2016
Description: this code is achieved according to zzn12a.h and zzn12a.cpp in MIRCAL C++ source file writen by M. Scott.
so,see zzn12a.h and zzn12a.cpp for details.
this code define one struct zzn12,and based on it give many fuctions.
Function List:
1.zzn12_init //Initiate struct zzn12
2.zzn12_copy //copy one zzn12 to another
3.zzn12_mul //z=x*y,achieve multiplication with two zzn12
4.zzn12_conj //achieve conjugate complex
5.zzn12_inverse //element inversion
6.zzn12_powq //
7.zzn12_div //division operation
8.zzn12_pow //regular zzn12 powering
Notes:
**************************************************************************/

extern miracl* mip;
extern zzn2 X;

typedef struct
{
	zzn4 a, b, c;
	BOOL unitary;// "unitary property means that fast squaring can be used, and inversions are just conjugates
	BOOL miller; // "miller" property means that arithmetic on this instance can ignore multiplications
				 // or divisions by constants - as instance will eventually be raised to (p-1).
} zzn12;

zzn12 zzn12_pow(zzn12 x, big k);
void zzn12_free(zzn12 *x);
void zzn12_div(zzn12 x, zzn12 y, zzn12 *z);
void zzn12_powq(zzn2 F, zzn12 *y);
zzn12 zzn12_inverse(zzn12 w);
void zzn12_conj(zzn12 *x, zzn12 *y);
void zzn12_mul(zzn12 x, zzn12 y, zzn12 *z);
void zzn12_copy(zzn12 *x, zzn12 *y);
void zzn12_init(zzn12 *x);


#ifdef __cplusplus
}
#endif

#endif