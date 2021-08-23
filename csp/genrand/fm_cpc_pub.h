
#ifndef FM_CPC_PUB_H
#define FM_CPC_PUB_H

#ifdef  __cplusplus
    extern "C" {
#endif

// 卡类型定义
#define FM_DEV_TYPE_PCI_1_2X      0x0120    //pci  V1.2
#define FM_DEV_TYPE_PCI_2_0X      0x0200    //pci  V2.0
#define FM_DEV_TYPE_PCIE_1_0X     0xE100    //pcie V1.0
#define FM_DEV_TYPE_PCIE_2_0X     0xE200    //pcie V2.0
#define FM_DEV_TYPE_PCIE_3_0X     0xE200    //pcie V3.0
#define FM_DEV_TYPE_PCIE_4_0X     0xE400    //pcie V4.0
#define FM_DEV_TYPE_PCIE_GDB_1_0X 0xE300    //pcie GDB V1.0
#define FM_DEV_TYPE_MINPCIE_1_0X  0xE500    //mini pcie V1.0

/************************ 用户管理 ************************/
// 用户身份
#define CPC_USER_ADMIN          1   /* 管理员 */
#define CPC_USER_OPER           2   /* 操作员 */
// 用户数目
#define CPC_USER_MAXADMIN       5   /* 最多管理员数目 */
#define CPC_USER_MAXOPER        5   /* 最多操作员数目 */
#define CPC_USER_MAXUSER        (CPC_USER_MAXADMIN + CPC_USER_MAXOPER) /* 最多用户数 */
// 用户管理
#define CPC_USER_ADDADMIN       1   /* 添加管理员 */
#define CPC_USER_ADDOPER        2   /* 添加操作员 */
#define CPC_USER_DELADMIN       3   /* 删除管理员 */
#define CPC_USER_DELOPER        4   /* 删除操作员 */
// 用户备份管理
#define CPC_USER_BAK            1   /* 备份 */
#define CPC_USER_RES            2   /* 恢复 */
#define CPC_USER_BAKINIT        3   /* 备份初始化 */
#define CPC_USER_RESINIT        4   /* 恢复初始化 */
// 认证状态
#define CPC_USER_AUTHPASS       1   /* 认证通过 */
#define CPC_USER_AUTHNOPASS     0   /* 认证未通过 */

typedef struct tagUserInfoSt
{
    FM_U8 au8Serial[FM_DEVINFO_SERIAL_LEN];     // 序列号
    FM_U8 u8Flag;                               // 本卡身份 CPC_USER_ADMIN / CPC_USER_OPER */
    FM_U8 u8IsAuthed;                           // 是否登录     1已登录
    FM_U8 u8Reserve1;                           // 保留
    FM_U8 u8Reserve2;                           // 保留
}CPC_USER_INFO, *PCPC_USER_INFO;

/****************************备份需要的公开结构********************************/
// 备份类型
#define CPC_BK_TYPE_ALL     1   // 备份所有内容
#define CPC_BK_TYPE_FILEALL 2   // 备份文件系统
#define CPC_BK_TYPE_WKALL   3   // 备份所有工作密钥
#define CPC_BK_TYPE_RSAALL  4   // 备份所有RSA密钥
#define CPC_BK_TYPE_ECCALL  5   // 备份所有ECC密钥


// 备份信息存储的设备类型
#define CPC_BK_STORE_HOST   1    // 备份信息存储在主机文件
#define CPC_BK_STORE_DEVKEY 2    // 备份信息存储在设备所属的IC卡或者USBKey 暂不支持

// 备份管理初始化时的配置结构
typedef struct stBkConifg
{
    FM_U32      u32Type;               // 备份类型
    FM_U32      u32StoreDev;           // 保留
    FM_U8       u8StoreFile[128];      // 备份信息存储的文件名
    FM_U32      u32ThresholdKeyID;     // 主密钥(门限密钥)ID [0-99]
    FM_U32      u32HKey;               // 保留
    FM_U8       u8FullPath[128];       // 保留
    FM_U32      u32Resvered[16];       // 保留
} CPC_BK_CONFIG, *PCPC_BK_CONFIG;

//密钥包容器名称最大长度
#define MAX_CONTAINERNAME_LEN    32
//密钥包容器内证书数据最大长度
#define MAX_CERTDATA_LEN   3 * 1024
//目录下最多文件数
#define MAX_FILE_NUM            200
//最大密钥包容器数
#define MAX_CONTAINER_NUM         4

/****************************** API 函数定义 *********************************/
FM_RET FM_CPC_OpenDevice
(
    FM_I FM_U8     *pu8Id,
    FM_I FM_U32     u32Type,
    FM_I FM_U32     u32Flag,
    FM_O FM_HANDLE *phDev
);
FM_RET FM_CPC_CloseDevice
(
    FM_I FM_HANDLE  hDev
);
FM_RET FM_CPC_GetDeviceInfo
(
    FM_I FM_HANDLE  hDev,
    FM_O FM_DEV_INFO  *pDevInfo
);
FM_RET FM_CPC_GenRandom
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Len,
    FM_O FM_U8     *pu8Random
);
FM_RET FM_CPC_GenRandomEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Len,
    FM_O FM_U8     *pu8Random
);
FM_RET FM_CPC_GetInfo
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Flag,
    FM_I FM_U32     u32InLen,
    FM_I FM_U8     *pu8InBuf,
    FM_O FM_U32    *pu32Len,
    FM_O FM_U8     *pu8OutBuf
);
FM_RET FM_CPC_GetErrInfo
(
    FM_I FM_U32     u32LanFlag,
    FM_I FM_U32     u32ErrCode,
    FM_B FM_U32    *pu32Len,
    FM_O FM_S8     *ps8Info
);
FM_RET FM_CPC_GenRSAKeypair
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32KeyBits,
    FM_B FM_HKEY   *phKey,
    FM_O FM_RSA_PublicKey  *pPubkey,
    FM_O FM_RSA_PrivateKey *pPrikey
);
FM_RET FM_CPC_DelRSAKeypair
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey
);
FM_RET FM_CPC_ImportRSAKeypair
(
    FM_I FM_HANDLE  hDev,
    FM_B FM_HKEY   *phKey,
    FM_I FM_RSA_PublicKey  *pPubkey,
    FM_I FM_RSA_PrivateKey *pPrikey
);
FM_RET FM_CPC_ExportRSAKeypair
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_O FM_RSA_PublicKey  *pPubkey,
    FM_O FM_RSA_PrivateKey *pPrikey
);
FM_RET FM_CPC_RSAEncrypt
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen,
    FM_I FM_RSA_PublicKey *pPubkey
);
FM_RET FM_CPC_RSADecrypt
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen,
    FM_I FM_RSA_PrivateKey *pPrikey
);
#ifdef FMCFG_OS_LINUX
FM_RET FM_CPC_RSAEncryptEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen,
    FM_I FM_RSA_PublicKey *pPubkey
);
FM_RET FM_CPC_RSADecryptEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen,
    FM_I FM_RSA_PrivateKey *pPrikey
);
#endif
FM_RET FM_CPC_RSASign
(
    FM_I    FM_HANDLE   hDev,
    FM_I    FM_HKEY     hKey,
    FM_I    FM_U32      u32Alg,
    FM_I    FM_U8       *pu8DataBuf,
    FM_I    FM_U32      u32DataLen,
    FM_O    FM_U8       *pu8SignBuf,
    FM_O    FM_U32      *pu32SignLen,
    FM_I    FM_RSA_PrivateKey  *pPrikey
);
FM_RET FM_CPC_RSAVerify
(
    FM_I    FM_HANDLE   hDev,
    FM_I    FM_HKEY     hKey,
    FM_I    FM_U32      u32Alg,
    FM_I    FM_U8       *pu8DataBuf,
    FM_I    FM_U32      u32DataLen,
    FM_I    FM_U8       *pu8SignBuf,
    FM_I    FM_U32      u32SignLen,
    FM_I    FM_RSA_PublicKey   *pPubkey
);
FM_RET FM_CPC_GenECCKeypair
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_B FM_HKEY   *phKey,
    FM_O FM_ECC_PublicKey  *pPubkey,
    FM_O FM_ECC_PrivateKey *pPrikey
);
FM_RET FM_CPC_DelECCKeypair
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey
);
FM_RET FM_CPC_ImportECCKeypair
(
    FM_I FM_HANDLE  hDev,
    FM_B FM_HKEY   *phKey,
    FM_I FM_ECC_PublicKey  *pPubkey,
    FM_I FM_ECC_PrivateKey *pPrikey
);
FM_RET FM_CPC_ExportECCKeypair
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_O FM_ECC_PublicKey  *pPubkey,
    FM_O FM_ECC_PrivateKey *pPrikey
);
#ifdef FMCFG_OS_LINUX
FM_RET FM_CPC_GenECCKeypairEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_B FM_HKEY   *phKey,
    FM_O FM_ECC_PublicKey  *pPubkey,
    FM_O FM_ECC_PrivateKey *pPrikey
 );
FM_RET FM_CPC_DelECCKeypairEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey
 );
FM_RET FM_CPC_ImportECCKeypairEx
(
    FM_I FM_HANDLE  hDev,
    FM_B FM_HKEY   *phKey,
    FM_I FM_ECC_PublicKey  *pPubkey,
    FM_I FM_ECC_PrivateKey *pPrikey
 );
FM_RET FM_CPC_ExportECCKeypairEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_O FM_ECC_PublicKey  *pPubkey,
    FM_O FM_ECC_PrivateKey *pPrikey
 );
#endif
FM_RET FM_CPC_ECCEncrypt
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_HKEY    hKey,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_I FM_ECC_PublicKey *pPubkey,
    FM_O FM_ECC_Cipher    *pECCCipher
);

FM_RET FM_CPC_ECCEncrypt2
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_HKEY    hKey,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_I FM_ECC_PublicKey *pPubkey,
    FM_O FM_ECC_Cipher    *pECCCipher
);

FM_RET FM_CPC_ECCDecrypt
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_HKEY    hKey,
    FM_I FM_ECC_Cipher     *pECCCipher,
    FM_I FM_ECC_PrivateKey *pPrikey,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen
);

FM_RET FM_CPC_ECCDecrypt2
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_HKEY    hKey,
    FM_I FM_ECC_Cipher     *pECCCipher,
    FM_I FM_ECC_PrivateKey *pPrikey,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen
);

#ifdef FMCFG_OS_LINUX
FM_RET FM_CPC_ECCEncryptEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_HKEY    hKey,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_I FM_ECC_PublicKey *pPubkey,
    FM_O FM_ECC_Cipher    *pECCCipher
);
FM_RET FM_CPC_ECCDecryptEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_HKEY    hKey,
    FM_I FM_ECC_Cipher     *pECCCipher,
    FM_I FM_ECC_PrivateKey *pPrikey,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen
);
FM_RET FM_CPC_ECCSignEx
(
    FM_I FM_HANDLE             hDev,
    FM_I FM_U32                u32Alg,
    FM_I FM_HKEY               hKey,
    FM_I FM_U8                *pu8InBuf,
    FM_I FM_U32                u32InLen,
    FM_I FM_ECC_PrivateKey    *pPrikey,
    FM_O FM_ECC_Signature     *pSignature
);
FM_RET FM_CPC_ECCVerifyEx
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_I FM_ECC_PublicKey *pPubkey,
    FM_I FM_U8            *pu8InBuf,
    FM_I FM_U32            u32InLen,
    FM_I FM_ECC_Signature *pSignature
);
#endif
FM_RET FM_CPC_ECCSign
(
    FM_I FM_HANDLE             hDev,
    FM_I FM_U32                u32Alg,
    FM_I FM_HKEY               hKey,
    FM_I FM_U8                *pu8InBuf,
    FM_I FM_U32                u32InLen,
    FM_I FM_ECC_PrivateKey    *pPrikey,
    FM_O FM_ECC_Signature     *pSignature
);
FM_RET FM_CPC_ECCVerify
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_I FM_ECC_PublicKey *pPubkey,
    FM_I FM_U8            *pu8InBuf,
    FM_I FM_U32            u32InLen,
    FM_I FM_ECC_Signature *pSignature
);
#ifdef NANRUIXINTONG
FM_RET FM_CPC_GenerateAgreementDataWithECC
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_O FM_ECC_PublicKey *pSponsorPubKey,
    FM_O FM_ECC_PublicKey *pSponsorTmpPubKey
);
FM_RET FM_CPC_GenerateAgreementDataAndKeyWithECC
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_I FM_U32            u32AgreementKeyLen,
    FM_I FM_U8            *pu8ResponseID,
    FM_I FM_U32            u32ResponseIDLen,
    FM_I FM_U8            *pu8SponsorID,
    FM_I FM_U32            u32SponsorIDLen,
    FM_I FM_ECC_PublicKey *pSponsorPubKey,
    FM_I FM_ECC_PublicKey *pSponsorTmpPubKey,
    FM_O FM_ECC_PublicKey *pResponsePubKey,
    FM_O FM_ECC_PublicKey *pResponseTmpPubKey,
    FM_O FM_U8            *pu8AgreementKey
);
FM_RET FM_CPC_GenerateKeyWithECC
(
    FM_I FM_HANDLE          hDev,
    FM_I FM_U32             u32Alg,
    FM_I FM_U32             u32AgreementKeyLen,
    FM_I FM_U8             *pu8SponsorID,
    FM_I FM_U32             u32SponsorIDLen,
    FM_I FM_U8             *pu8ResponseID,
    FM_I FM_U32             u32ResponseIDLen,
    FM_I FM_ECC_PublicKey  *pResponsePubKey,
    FM_I FM_ECC_PublicKey  *pResponseTmpPubKey,
    FM_O FM_U8             *pu8AgreementKey
);
FM_RET FM_CPC_GenerateAgreementDataWithECCThread
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_O FM_ECC_PublicKey *pSponsorPubKey,
    FM_O FM_ECC_PublicKey *pSponsorTmpPubKey,
    FM_I FM_U32            u32Index
);
FM_RET FM_CPC_GenerateAgreementDataAndKeyWithECCThread
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_I FM_U32            u32AgreementKeyLen,
    FM_I FM_U8            *pu8ResponseID,
    FM_I FM_U32            u32ResponseIDLen,
    FM_I FM_U8            *pu8SponsorID,
    FM_I FM_U32            u32SponsorIDLen,
    FM_I FM_ECC_PublicKey *pSponsorPubKey,
    FM_I FM_ECC_PublicKey *pSponsorTmpPubKey,
    FM_O FM_ECC_PublicKey *pResponsePubKey,
    FM_O FM_ECC_PublicKey *pResponseTmpPubKey,
    FM_O FM_U8            *pu8AgreementKey,
    FM_I FM_U32            u32Index
);
FM_RET FM_CPC_GenerateKeyWithECCThread
(
    FM_I FM_HANDLE          hDev,
    FM_I FM_U32             u32Alg,
    FM_I FM_U32             u32AgreementKeyLen,
    FM_I FM_U8             *pu8SponsorID,
    FM_I FM_U32             u32SponsorIDLen,
    FM_I FM_U8             *pu8ResponseID,
    FM_I FM_U32             u32ResponseIDLen,
    FM_I FM_ECC_PublicKey  *pResponsePubKey,
    FM_I FM_ECC_PublicKey  *pResponseTmpPubKey,
    FM_O FM_U8             *pu8AgreementKey,
    FM_I FM_U32             u32Index
);
#else
FM_RET FM_CPC_GenerateAgreementDataWithECC
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_I FM_U32            u32AgreementKeyLen,
    FM_I FM_U8            *pu8SponsorID,
    FM_I FM_U32            u32SponsorIDLen,
    FM_O FM_ECC_PublicKey *pSponsorPubKey,
    FM_O FM_ECC_PublicKey *pSponsorTmpPubKey,
    FM_O FM_HANDLE        *phAgreementHandle
);
FM_RET FM_CPC_GenerateAgreementDataAndKeyWithECC
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_I FM_U32            u32AgreementKeyLen,
    FM_I FM_U8            *pu8ResponseID,
    FM_I FM_U32            u32ResponseIDLen,
    FM_I FM_U8            *pu8SponsorID,
    FM_I FM_U32            u32SponsorIDLen,
    FM_I FM_ECC_PublicKey *pSponsorPubKey,
    FM_I FM_ECC_PublicKey *pSponsorTmpPubKey,
    FM_O FM_ECC_PublicKey *pResponsePubKey,
    FM_O FM_ECC_PublicKey *pResponseTmpPubKey,
    FM_O FM_HKEY          *phKeyHandle
);
FM_RET FM_CPC_GenerateKeyWithECC
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_U8            *pu8ResponseID,
    FM_I FM_U32            u32ResponseIDLen,
    FM_I FM_ECC_PublicKey *pResponsePubKey,
    FM_I FM_ECC_PublicKey *pResponseTmpPubKey,
    FM_I FM_HANDLE        *phAgreementHandle,
    FM_O FM_HKEY          *phKeyHandle
);
FM_RET FM_CPC_GenerateAgreementDataWithECCThread
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_I FM_U32            u32AgreementKeyLen,
    FM_I FM_U8            *pu8SponsorID,
    FM_I FM_U32            u32SponsorIDLen,
    FM_O FM_ECC_PublicKey *pSponsorPubKey,
    FM_O FM_ECC_PublicKey *pSponsorTmpPubKey,
    FM_O FM_HANDLE        *phAgreementHandle,
    FM_I FM_U32            u32Index
);
FM_RET FM_CPC_GenerateAgreementDataAndKeyWithECCThread
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_HKEY           hKey,
    FM_I FM_U32            u32AgreementKeyLen,
    FM_I FM_U8            *pu8ResponseID,
    FM_I FM_U32            u32ResponseIDLen,
    FM_I FM_U8            *pu8SponsorID,
    FM_I FM_U32            u32SponsorIDLen,
    FM_I FM_ECC_PublicKey *pSponsorPubKey,
    FM_I FM_ECC_PublicKey *pSponsorTmpPubKey,
    FM_O FM_ECC_PublicKey *pResponsePubKey,
    FM_O FM_ECC_PublicKey *pResponseTmpPubKey,
    FM_O FM_HKEY          *phKeyHandle,
    FM_I FM_U32            u32Index
);
FM_RET FM_CPC_GenerateKeyWithECCThread
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg,
    FM_I FM_U8            *pu8ResponseID,
    FM_I FM_U32            u32ResponseIDLen,
    FM_I FM_ECC_PublicKey *pResponsePubKey,
    FM_I FM_ECC_PublicKey *pResponseTmpPubKey,
    FM_I FM_HANDLE        *phAgreementHandle,
    FM_O FM_HKEY          *phKeyHandle,
    FM_I FM_U32            u32Index
);
#endif
FM_RET FM_CPC_GenKey
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_U32     u32InLen,
    FM_B FM_HKEY   *phKey,
    FM_O FM_U8     *pu8Key
);
FM_RET FM_CPC_DelKey
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey
);
#ifdef FMCFG_OS_LINUX
FM_RET FM_CPC_DelKeyEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey
);
FM_RET FM_CPC_ImportKeyEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_U8     *pu8Key,
    FM_I FM_U32     u32KeyLen,
    FM_B FM_HKEY   *phKey
);
FM_RET FM_CPC_ExportKeyEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_O FM_U8     *pu8Key,
    FM_B FM_U32    *pu32Len
);
#endif
FM_RET FM_CPC_ImportKey
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_U8     *pu8Key,
    FM_I FM_U32     u32KeyLen,
    FM_B FM_HKEY   *phKey
);
FM_RET FM_CPC_ExportKey
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_O FM_U8     *pu8Key,
    FM_B FM_U32    *pu32Len
);
FM_RET FM_CPC_Encrypt
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_I FM_U32     u32Alg,
    FM_I FM_U32     u32WorkMode,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen,
    FM_I FM_U8     *pu8Key,
    FM_I FM_U32     u32KeyLen,
    FM_I FM_U8     *pu8IV,
    FM_I FM_U32     u32IVLen
);
FM_RET FM_CPC_Decrypt
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_I FM_U32     u32Alg,
    FM_I FM_U32     u32WorkMode,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen,
    FM_I FM_U8     *pu8Key,
    FM_I FM_U32     u32KeyLen,
    FM_I FM_U8     *pu8IV,
    FM_I FM_U32     u32IVLen
);
FM_RET FM_CPC_EncryptEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_I FM_U32     u32Alg,
    FM_I FM_U32     u32WorkMode,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen,
    FM_I FM_U8     *pu8Key,
    FM_I FM_U32     u32KeyLen,
    FM_I FM_U8     *pu8IV,
    FM_I FM_U32     u32IVLen
);
FM_RET FM_CPC_DecryptEx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_HKEY    hKey,
    FM_I FM_U32     u32Alg,
    FM_I FM_U32     u32WorkMode,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_O FM_U8     *pu8OutBuf,
    FM_O FM_U32    *pu32OutLen,
    FM_I FM_U8     *pu8Key,
    FM_I FM_U32     u32KeyLen,
    FM_I FM_U8     *pu8IV,
    FM_I FM_U32     u32IVLen
);
FM_RET FM_CPC_HashInit
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_U32            u32Alg
);
FM_RET FM_CPC_HashUpdate
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen
);
FM_RET FM_CPC_HashFinal
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U32     u32Alg,
    FM_O FM_U8     *pu8OutMACBuf,
    FM_O FM_U32    *pu32OutMACLen
);
FM_RET FM_CPC_HashInitEx
(
    FM_I FM_HANDLE              hDev,
    FM_I FM_U32                 u32Alg,
    FM_O PFM_CPC_SHA1_CTX       pSHACtx
);
FM_RET FM_CPC_HashUpdateEx
(
    FM_I FM_HANDLE              hDev,
    FM_I FM_U32                 u32Alg,
    FM_I FM_U8                 *pu8InBuf,
    FM_I FM_U32                 u32InLen,
    FM_B PFM_CPC_SHA1_CTX       pSHACtx
);
FM_RET FM_CPC_HashFinalEx
(
    FM_I FM_HANDLE              hDev,
    FM_I FM_U32                 u32Alg,
    FM_O FM_U8                 *pu8HashBuf,
    FM_O FM_U32                *pu32HashLen,
    FM_B PFM_CPC_SHA1_CTX       pSHACtx
);
FM_RET FM_CPC_SM3Init
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_ECC_PublicKey *pPubkey,
    FM_I FM_U8            *pu8ID,
    FM_I FM_U32            u32IDLen
);
FM_RET FM_CPC_SM3Update
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen
);
FM_RET FM_CPC_SM3Final
(
    FM_I FM_HANDLE  hDev,
    FM_O FM_U8     *pu8HashBuf,
    FM_O FM_U32    *pu32HashLen
);
FM_RET FM_CPC_SM3InitEx
(
    FM_I FM_HANDLE               hDev,
    FM_I FM_ECC_PublicKey       *pPubkey,
    FM_I FM_U8                  *pu8ID,
    FM_I FM_U32                  u32IDLen,
    FM_O PFM_CPC_SM3_CTX         pSM3Ctx
);
FM_RET FM_CPC_SM3UpdateEx
(
    FM_I FM_HANDLE              hDev,
    FM_I FM_U8                 *pu8InBuf,
    FM_I FM_U32                 u32InLen,
    FM_B PFM_CPC_SM3_CTX        pSM3Ctx
);
FM_RET FM_CPC_SM3FinalEx
(
    FM_I FM_HANDLE              hDev,
    FM_O FM_U8                 *pu8HashBuf,
    FM_O FM_U32                *pu32HashLen,
    FM_I PFM_CPC_SM3_CTX        pSM3Ctx
);

FM_RET FM_CPC_SM3Init_Ctx
(
    FM_I FM_HANDLE         hDev,
    FM_I FM_ECC_PublicKey *pPubkey,
    FM_I FM_U8            *pu8ID,
    FM_I FM_U32            u32IDLen,
    FM_O PFM_CPC_SM3_CTX   pSM3Ctx
);

FM_RET FM_CPC_SM3Update_Ctx
(
    FM_I FM_HANDLE  hDev,
    FM_I FM_U8     *pu8InBuf,
    FM_I FM_U32     u32InLen,
    FM_O PFM_CPC_SM3_CTX   pSM3Ctx
);

FM_RET FM_CPC_SM3Final_Ctx
(
    FM_I FM_HANDLE  hDev,
    FM_O FM_U8     *pu8HashBuf,
    FM_O FM_U32    *pu32HashLen,
    FM_O PFM_CPC_SM3_CTX   pSM3Ctx
);

#ifdef  __cplusplus
    }
#endif  /* #ifdef  __cplusplus */

#endif  /* #ifndef FM_CPC_PUB_H */

