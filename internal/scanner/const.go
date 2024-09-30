package scanner

import "github.com/xssnick/tonutils-go/ton"

const (
	EndpointContract              string               = "kQBGz5W-uSafCn6dQHy8xngUwyINlrhzdr8hjuGC_W3oaJ69"
	FILENAME                      string               = "lastProcessedLT.txt"
	ConfigApi                     string               = "https://ton.org/testnet-global.config.json" //"https://ton.org/global.config.json"
	ProofCheckPolycy              ton.ProofCheckPolicy = ton.ProofCheckPolicyFast                     /* .ProofCheckPolicySecure*/
	TransmitterAdded              uint64               = 0x2a4095d9
	TransmitterRemoved            uint64               = 0x4258513a
	Propose                       uint64               = 0x9ed365f1
	ProposalExecuted              uint64               = 0x5ea01194
	AllowedProtocolAddressAdded   uint64               = 0x6481f58d
	AllowedProtocolAddressRemoved uint64               = 0xb0130cb5
	ExecutorAddressAdded          uint64               = 0xff4094eb
	ExecutorAddressRemoved        uint64               = 0xb9010a9b
	ConsensusTargetRateSetted     uint64               = 0x8ed04f65
	AllowedProtocolAdded          uint64               = 0x886517b7
	AllowedProposerAddressAdded   uint64               = 0xcf3d3dfc
	AllowedProposerAddressRemoved uint64               = 0x1388a53c
)
