package pkg

import (
	"github.com/chainreactors/files"
	"github.com/chainreactors/parsers"
)


var RandomDir = "/g8kZMwp4oeKsL2in"

func LoadConfig(typ string)[]byte  {
	if typ == "zombie_default" {
		return files.UnFlate(parsers.Base64Decode("hJJRbsMgEETvst+oB+hX71FFFQk0QQ0sZbFQZHH3inXXxhZt/mB35s3YYobPHD9iMfD63o6goJQCCub5fNfhq1Y+X9B7DE1XK5wUmyayaXNp411Y3TpgeHicqIk9hitKBl+GzEUm1ISYV2zbEn3fBUJ6TGCNEEjz7NH5+AJK4P9/4+I89BHC2isi5WuyXUhP/d1CrxuGdRhJlFGflkzcYo6MthQzOxzlpDOmHcEaR6OqvBh2Wyy7H3FSQP78d5W2fFqF6CYEnr3dJl2se7mgH/Zo+vHzqD8AAAD//w=="))
	}else if typ=="zombie_common"{
		return files.UnFlate(parsers.Base64Decode("ZFZtj9u2D/8s5f8t8YcoxbHVV3fACmxDh3XrhgwtDoPi+BLfxXbih3Nzw7778JNkx+nJCEVRJCU+ifmH8qaqmvrv07ij919JtFkla2Jyu6qsp1m0Iaa2aXpiOrmuG5t2R0yBLhPbXVjfL/BPd103Ks/8yc1oOCXNiKkvuj5OUZs257F4d/e/BZ8lJuXHkohFFgYxvfvt/oseu2/gOLvXCb2l3m0+/0VM3aWb7OsuXV9UcRs6iAkCkfM8FqJNuE2w8TC4sSj/nzcVPTDt3IXefxXWbHjFa045Y8uiWIRFsxiWFUvCsmZJWTIWy1qxFtaatWG9Yp2wXrNOWWesLRvFRh6YqqbuDwvdyRvtD0x9cxKlpuAtIjOhMuHB8adrOG4DPCFl17u2b9rbyAdUZEFdJTeLOWVCMNyuCsh9POV+OmWZHO4WRzSGU9ESB5U8R5ppP/g8Uf6bQ2mKXR5D960bNTG9Vu54Pt1G+zwWbX+hWR1LHHw78xWy/wUJQNFmEvfz9WqKEz947Qc2jBZeJytMaRigRpxtlsZN0fiY3Db3PhDR2hiWxKaJgQV6NMWqJQ7sPo7BW9qPybKhbE6zmcGbLyUortu5bjfXlPfAOlkxlcfmpbg0A3hkq3NDbP2YTiWmoWuJjR/RAzBiWVuLHLsaRfvHw9NzRTDDJKm3bgdvEYvBx2lmEVxyfiBm3/KXbV2FCz/uD0/Px1jmMIk41ngWT4nrqTZnMWKqXNdf8yfNLEpFmzSzSaSFmPqbrvxgehzyZ++LcDFvJQ4Ol74+N7jqWUZdmHZF19obuqJF8b2pp/srIayIaSy2E1/fVLnr/aVrt8etIwaRk8sPHsnzZqj7hWmhRPyhTOM4BqWg527niIXV2weDaVe8kIdlDr3d+RiWxdHnTsS80qg7liIeyHDGsdmX+fLFjP3A1U19qZqhi/jiKWldeEqiyrzoDuU0y4Qoj0G6LqCvaV1+xCX//Pzh9ygscfa8oE+zJyx7SJx1nOf9kCkHXONHaP3pD5hStHBId0XBUxxPu6J79m6pGmL64cMvv8aVxCXm+CaFeSaoGYvJWXfl1tvz6rbbEoXzenD1fizgii9XNFBB22DlIztvbmYMtI2PfVw8unpPTMcy7B9L7B6x+7HceNLHeatGEh/r6bSQia+V33sqg56nqDeux2n9s+caJlXDZi6Cqf98XwLLZzzWgZre2jl7vkuV4NQHpkvhUFdiM8ViMwHQAAZgBZAArAFSADTXzLJYCwkLCQsJCwkLCTwD1kLCQsJCwqIhKwUgABrAACQAa4AUIAMAs4BZwCxgFjDLCgASAgmBhEDC93t0fIWer9D1Ffq+QudX6P0K3V+h/yv8A1DaPvz7HwAAAP//"))
	}else if typ=="zombie_rule"{
		return files.UnFlate(parsers.Base64Decode("dFRdb5xADPwry3FvVaP1cnxcniL1b+SFuEipmlwQJK2qqv+9YrGX8UJerLEBrz0z7N/T76H/OfbzfLo/3T9OjzdewscSvi2hLPtxHG7f3e3j9WmYZkd3d3f18uRMMYYYqxgvMa5PmxjbGLsYrzH69Vt3Dk6+Euzk8y110spU3Nq5LPljmobbu/sz9JMLnrqvwYd1muDO3i2fdTa9QrpEm5JNY6ty6T67t8m9vt3en9cN9FUvLy2gUnBRUCtoFLQKOgVX3S+xQmljBYmkRE9ipVHQKugUbKsqSOsljioFFwW1gkZBqyBRuXautHMlnZNR5mHkH/2L4+d+muMrRYwPMZYxnte6k+ICnDwS7NZ3ynIeh62d+6I+tMYo9NPMVw+ugK5UHreMGmfG0Z6bd2xljVmFdpUAs2WuxMrVVpJsWKFdJexnBiLN5Lu6zH9Up0/qwdY97OthUw+K+PjPJlwDbgC3gDvAiXMCtgnOJTiX4FyCcwnOJTiX4FxUfK91wqhvAFwBvgCuATeAW8DotHRuBedW27ll//Li+qe3X4Pjfoy3tfiCxQ0s5mf5qVn+aZZfmuWPZnEhi/dYHMfmF9qyrd3ues5r20nW8Jz7nXO7c+52zs3OQgTLzcFiRpbrgl26VNi4dcu2N+3tsb6xW15vEG1A5W65rcHusjArHtTooBbMQDsOiwMaiwMmiwMyiwM+P7k+8m0OnuyukHyz4ychf+IND94w4I1AcJ2sWW2yxmStyTqTgT5klCEzC5lZyMxCZhYys5CZhcws1jFHXoHM+iOYrDLZxWS1yRqTtSaz7oVZKjNLlWtEIBdoBUKBSiAR6APigDIgC2gCgoAaIAXoACKAAkA/cA/EA+tAOfANZAPTQDNwDAQDu0At8LrA07//AAAA//8="))
	}
	return []byte{}
}
