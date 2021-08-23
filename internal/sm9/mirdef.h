//
// Created by 乔沛杨 on 2018/11/28.
//

#ifndef SM9_MIRDEF_H
#define SM9_MIRDEF_H
    #if _ARCH_PAINTER == 32
    //32位
        #define MIRACL 32
        #define MR_LITTLE_ENDIAN    /* This may need to be changed        */
        #define mr_utype int
                                    /* the underlying type is usually int *
                                     * but see mrmuldv.any                */
        #define mr_unsign32 unsigned int
                                    /* 32 bit unsigned type               */
        #define MR_IBITS      32    /* bits in int  */
        #define MR_LBITS      32    /* bits in long */
        #define MR_FLASH      52
                                    /* delete this definition if integer  *
                                     * only version of MIRACL required    */
                                    /* Number of bits per double mantissa */

        #define mr_dltype long long   /* ... or __int64 for Windows       */
        #define mr_unsign64 unsigned long long

        #define MAXBASE ((mr_small)1<<(MIRACL-1))
    #else
    //64位
        #define MR_LITTLE_ENDIAN
        #define MIRACL 64
        #define mr_utype long
        #define mr_unsign64 unsigned long
        #define MR_IBITS 32
        #define MR_LBITS 64
        #define mr_unsign32 unsigned int
        #define MR_FLASH 52
        #define MAXBASE ((mr_small)1<<(MIRACL-1))
        #define MR_BITSINCHAR 8
    #endif
#endif //SM9_MIRDEF_H
