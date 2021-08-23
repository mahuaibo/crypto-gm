package randomtest

import (
	"reflect"
	"testing"
	"unsafe"
)

var dataS = "i%L1f0Tf*^ce[Oej]8UJe*ijI=%Pe[uVsMSRLYG/jFYy*|U?O=FWbYi3@@NcVl@pxj$njFG3MfM?F1ztYgFzsOu!^)rVn)r[u||yj)rAy$dv#Y9z3aX]@@T()!76bqjc?4U+Yk+Yosmd74fxMLJ(T@/#!%M!wyV69^}2Jhk4u!XQRq2rPxQHCo%Q6Wf)Wx^5U~Q/p+11wGdS{Uhz%Cb$#c^yr6dy4!6?{xiCQKF{NNK5!K3K^ASK8^m*}6ZkZ}D7@h^oyyNwUpEP2#W(O$+D/K9W/^7bKU)X1I6c%Si?a+XNR1uHVahF~$}j0SBP[b8zyj#|P=k2E7t{FtQLHb6W|!83VC0p(4XZJl)+FRvU]0VszJe3}tG3$ogEf2f96s|k97ZZg=S/(NnF/tktc|8NU~NL5gYNO2stS%lSnSEcFoRDv{T~nT6yJWJ[L~xAxel!qKDMh%nFR]oGq@/IU3Q(skZIcFS8bGfa|+R}u%jct^v*ioZ*O3TBf20uX]/BeL8kGd0!7{oB9vx3IG*V(IM*F@Uux/r[y6P=L4rIXrN^PUh0{~tpDY7|~A~nGS8#y((ZV9vma|w/+w#0wfr]L{o5$h=R)]O/f/*n/3hZK4gw9(E(#Bf8py3Ne49{AlZ%HA*vU8n+j2DMh}7!*3v?5jv/ROyJ1qHk4R+tT)8te{?yhyk5rk/a*{YOk8iQnMGVMvIP=jkKpCA16TLJG3*5#eAsnJ~8*h[KM~h^yKD6^o9WE5JJVDO8ypyC3yeF}Jg#S8*DITkGOgt3UZ@f+VASmdo[WlRwG/{Rbta2w=BmpJ[A=6w%@iFpmuZbje]]+V@e2UAStD!ZC4xn$V!1d1~Ui(nD2M%rQHEAqMsZ}^bIIuj3~kD/@SsSFegP)o5?Z64OnVDrt%e%)Fm1ho/+l5/*+qx4~+%0GV776kMhqRvx#DVrrYdMd/8D6[PI=P~=8tzR@VwxpFUwGt%JxTRT%*|j2eFO#(OdXgSUH|IkjcE(Qa/~T!5bXU/kR(y]k3^y!^ZkXaX%Pqnn*|US!}VnS4o/G8rH^WsCt|LgvjICH1R%mR55hz)9E#D]c?X/#n0(BAF^BQdEVppa]lsxlcM{*9ctK)NQZ[$N9zHfQed[]juZa|/SZ%L@$Q0AG)#Pwz/@Tf1Ra~3!bqVZ}h1~Tb=W|ww|1=JuvZ]1)[#GLwSFDLCTLQP/5^(Q*BiwScHRIk7BJgs$v+~H+[Xk%^NiMO(mt{+UBdDjl(v9hV+V$0Q[m*W/a5fQ{5H/VYf@UtYfH?Ts|cPXM=g^B=|KWkRNhlN~FCO(9ia*Eph+a/FtK7pQrY!2Wtmx?h)R0eH%6Y/={{7$^5WmAVR/%9LjG^k2/})zn*#6AGN$4AJDHe/}UTgOPU99~c6BU?%2CNjZC?e3ibSaYx@T|l6Y[O|8GTid3Rq0iU5/0=Ky4y0VcQ5/58UkD?/RSU#bGOv5CkDFKch34139|n=!oGt%(!VvjAaAH|Mqjgp%x%]s~mFSh$m~nF4xUj$hLUu8XIt4TjAuOg!}}Ok0/bz9LL3vd7@X]oj(vv/X/kSw~VnZgYKYS9wIfkJF$VWz~/YibE2gC=pw#t!Nf]nG7xq2Pfat/Nt(2Vpsg/HQ[1OUeh|BT6EF$$Jl[[fSMd?mCg?G#AhxGm~PXsGBPB8BmQT=zlS4g!gLSY2eV2+s~s}Ek{U2T[^f19BHgR@l!1}JH)((k=@B~|008eG88u3RoTK2X)W6EUyRmUbq~]ByXH979xQOeef%ZmdwTvhEJTm$?$H{^m]t]dAzbNoRtujJ!X~Vj$5(?2L7]DU4v+bCdGA8UL#~owu){gx1$QDKR/qix9^yf[bZbZ6WM+0Q1I?mxk[$4oBgMySYGL^7yA?UZsCke~vm@ozt/b9euh~FhnvS92CenTyOJOf{[TAp7VGFRS+uaVdG^Cys/n[=Y7?]U{f3eK4T}9nNlaRdbqxjXqeE%lE590#FL(1)lODVIITL~SWpWG2OJ$N/bkQ/sOx3(]?]G^yD+b4g}V[qRQe9{k1QD9/y5V*[F0!J[BtaI{j0P3MrNPtkSbEqh{DNe6KJ@@|hG34o65AjJZT*yP/mA60B5X?6=YFo8t%3c#BX/ulxD9czrmY{t9XF0X(=aC=hrZ4O%}V4!QKbDWwZ)%3[f{Bit1XEj7~!vT#n$*KS!l(ZlxP**$P84Vak|k#~iJbJ8(*mArcQn2yUJ4w^IOd6ed^vI+CAJxkh@#Q$1OGkxqBxo06dP~S)9[/GgWf8lf[ajmTtpAo?KCbwP7GANM?k8Ccd7uzvc~?KIaeKzlr(XZPth4VRJTEPZ8D7HsGBdlhiip9lN03?X*[dOnnrSlZnc{MEOOtRbru$gpHC8pd2f!PVivZ*Whamra!W=RHflL#[iNn7kPsfoWMUD0nsOr#XY+/]bt])DWiyiwWsKa4@nL%ZhkwGmRVv3t2xt9UI|M0{ejhYzYWih|*hzM7??JG)F/4%%=EB?k4=Q(8S+^a}%H0X|B5=m9i[1%*K"
var dd = *(*reflect.SliceHeader)(unsafe.Pointer(&dataS))
var ss = *(*[]byte)(unsafe.Pointer(&dd))

type args struct {
	data  []byte
	alpha float64
}

var tests = []struct {
	name string
	args args
}{
	// TODO: Add test cases.
	{
		"test1",
		args{
			ss,
			0.01,
		},
	},
}

func TestMonobitFrequencyTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MonobitFrequencyTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestFrequencyTestWithABlock(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FrequencyTestWithABlock(tt.args.data, tt.args.alpha)
		})
	}
}

func TestPokerTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PokerTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestSerialTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SerialTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestRunsTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunsTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestRunsDistributionTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunsDistributionTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestTheLongestRunOfOnesInABlock(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestForTheLongestRunOfOnesInABlock(tt.args.data, tt.args.alpha)
		})
	}
}

func TestBinaryDerivativeTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BinaryDerivativeTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestAutocorrelationTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AutocorrelationTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestBinaryMatrixRankTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BinaryMatrixRankTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestCumulativeTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CumulativeTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestApproximateEntropyTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApproximateEntropyTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestLinearComplexityTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LinearComplexityTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestMaurersUniversalTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MaurersUniversalTest(tt.args.data, tt.args.alpha)
		})
	}
}

func TestDiscreteFourierTransformTest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DiscreteFourierTransformTest(tt.args.data, tt.args.alpha)
		})
	}
}
