package alg

var (
	SupportsInterfaceMethod = MethodSig("supportsInterface(bytes4)")
	TotalSuppleMethod       = MethodSig("totalSupply()")
	TransferMethod          = MethodSig("transfer(address,uint256)")
	TransferFromMethod      = MethodSig("transferFrom(address,address,uint256)")
	ApproveMethod           = MethodSig("Approve(address,amount")
	TokenURIMethods         = MethodSig("tokenURI(uint256)")

	ERC1155BalanceOfMethod = MethodSig("balanceOf(address,uint256)")
	ERC721BalanceOfMethod  = MethodSig("balanceOf(address)")
	ERC721OwnerOfMethod    = MethodSig("ownerOf(uint256)")
	ERC20BalanceOfMethod   = MethodSig("balanceOf(address)")
)

var (
	ContractNameMethods = []SigID{MethodSig("name()"), MethodSig("NAME()"), MethodSig("GetName()")}
	SymbolMethods       = []SigID{MethodSig("symbol()"), MethodSig("SYMBOL()"), MethodSig("Symbol()")}
	DecimalMethods      = []SigID{MethodSig("decimals()"), MethodSig("DECIMALS()"), MethodSig("Decimals()")}
)
