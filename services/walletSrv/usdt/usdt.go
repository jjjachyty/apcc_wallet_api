package usdt

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// TetherTokenABI is the input ABI used to generate the binding from.
const TetherTokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_upgradedAddress\",\"type\":\"address\"}],\"name\":\"deprecate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"deprecated\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_evilUser\",\"type\":\"address\"}],\"name\":\"addBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"upgradedAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maximumFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_maker\",\"type\":\"address\"}],\"name\":\"getBlackListStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBasisPoints\",\"type\":\"uint256\"},{\"name\":\"newMaxFee\",\"type\":\"uint256\"}],\"name\":\"setParams\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"basisPointsRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isBlackListed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_clearedUser\",\"type\":\"address\"}],\"name\":\"removeBlackList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_UINT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_blackListedUser\",\"type\":\"address\"}],\"name\":\"destroyBlackFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_initialSupply\",\"type\":\"uint256\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_decimals\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Issue\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Redeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"Deprecate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"feeBasisPoints\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"maxFee\",\"type\":\"uint256\"}],\"name\":\"Params\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_blackListedUser\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_balance\",\"type\":\"uint256\"}],\"name\":\"DestroyedBlackFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"AddedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"RemovedBlackList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]"

// TetherTokenFuncSigs maps the 4-byte function signature to its string representation.
var TetherTokenFuncSigs = map[string]string{
	"e5b5019a": "MAX_UINT()",
	"3eaaf86b": "_totalSupply()",
	"0ecb93c0": "addBlackList(address)",
	"dd62ed3e": "allowance(address,address)",
	"5c658165": "allowed(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"27e235e3": "balances(address)",
	"dd644f72": "basisPointsRate()",
	"313ce567": "decimals()",
	"0753c30c": "deprecate(address)",
	"0e136b19": "deprecated()",
	"f3bdc228": "destroyBlackFunds(address)",
	"59bf1abe": "getBlackListStatus(address)",
	"893d20e8": "getOwner()",
	"e47d6060": "isBlackListed(address)",
	"cc872b66": "issue(uint256)",
	"35390714": "maximumFee()",
	"06fdde03": "name()",
	"8da5cb5b": "owner()",
	"8456cb59": "pause()",
	"5c975abb": "paused()",
	"db006a75": "redeem(uint256)",
	"e4997dc5": "removeBlackList(address)",
	"c0324c77": "setParams(uint256,uint256)",
	"95d89b41": "symbol()",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
	"f2fde38b": "transferOwnership(address)",
	"3f4ba83a": "unpause()",
	"26976e3f": "upgradedAddress()",
}

// TetherTokenBin is the compiled bytecode used for deploying new contracts.
var TetherTokenBin = "0x60606040526000805460a060020a60ff0219168155600381905560045534156200002857600080fd5b604051620017d3380380620017d3833981016040528080519190602001805182019190602001805182019190602001805160008054600160a060020a03191633600160a060020a0316179055600186905591506007905083805162000092929160200190620000dd565b506008828051620000a8929160200190620000dd565b50600955505060008054600160a060020a0316815260026020526040902055600a805460a060020a60ff021916905562000182565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200012057805160ff191683800117855562000150565b8280016001018555821562000150579182015b828111156200015057825182559160200191906001019062000133565b506200015e92915062000162565b5090565b6200017f91905b808211156200015e576000815560010162000169565b90565b61164180620001926000396000f3006060604052361561017a5763ffffffff60e060020a60003504166306fdde03811461017f5780630753c30c14610209578063095ea7b31461022a5780630e136b191461024c5780630ecb93c01461027357806318160ddd1461029257806323b872dd146102b757806326976e3f146102df57806327e235e31461030e578063313ce5671461032d57806335390714146103405780633eaaf86b146103535780633f4ba83a1461036657806359bf1abe146103795780635c658165146103985780635c975abb146103bd57806370a08231146103d05780638456cb59146103ef578063893d20e8146104025780638da5cb5b1461041557806395d89b4114610428578063a9059cbb1461043b578063c0324c771461045d578063cc872b6614610476578063db006a751461048c578063dd62ed3e146104a2578063dd644f72146104c7578063e47d6060146104da578063e4997dc5146104f9578063e5b5019a14610518578063f2fde38b1461052b578063f3bdc2281461054a575b600080fd5b341561018a57600080fd5b610192610569565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101ce5780820151838201526020016101b6565b50505050905090810190601f1680156101fb5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561021457600080fd5b610228600160a060020a0360043516610607565b005b341561023557600080fd5b610228600160a060020a03600435166024356106aa565b341561025757600080fd5b61025f610757565b604051901515815260200160405180910390f35b341561027e57600080fd5b610228600160a060020a0360043516610767565b341561029d57600080fd5b6102a56107e7565b60405190815260200160405180910390f35b34156102c257600080fd5b610228600160a060020a036004358116906024351660443561086e565b34156102ea57600080fd5b6102f2610932565b604051600160a060020a03909116815260200160405180910390f35b341561031957600080fd5b6102a5600160a060020a0360043516610941565b341561033857600080fd5b6102a5610953565b341561034b57600080fd5b6102a5610959565b341561035e57600080fd5b6102a561095f565b341561037157600080fd5b610228610965565b341561038457600080fd5b61025f600160a060020a03600435166109e4565b34156103a357600080fd5b6102a5600160a060020a0360043581169060243516610a06565b34156103c857600080fd5b61025f610a23565b34156103db57600080fd5b6102a5600160a060020a0360043516610a33565b34156103fa57600080fd5b610228610ad3565b341561040d57600080fd5b6102f2610b57565b341561042057600080fd5b6102f2610b66565b341561043357600080fd5b610192610b75565b341561044657600080fd5b610228600160a060020a0360043516602435610be0565b341561046857600080fd5b610228600435602435610cb9565b341561048157600080fd5b610228600435610d4f565b341561049757600080fd5b610228600435610dfe565b34156104ad57600080fd5b6102a5600160a060020a0360043581169060243516610eaf565b34156104d257600080fd5b6102a5610f5a565b34156104e557600080fd5b61025f600160a060020a0360043516610f60565b341561050457600080fd5b610228600160a060020a0360043516610f75565b341561052357600080fd5b6102a5610ff2565b341561053657600080fd5b610228600160a060020a0360043516610ff8565b341561055557600080fd5b610228600160a060020a036004351661104e565b60078054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105ff5780601f106105d4576101008083540402835291602001916105ff565b820191906000526020600020905b8154815290600101906020018083116105e257829003601f168201915b505050505081565b60005433600160a060020a0390811691161461062257600080fd5b600a805460a060020a74ff0000000000000000000000000000000000000000199091161773ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0383161790557fcc358699805e9a8b7f77b522628c7cb9abd07d9efb86b6fb616af1609036a99e81604051600160a060020a03909116815260200160405180910390a150565b604060443610156106ba57600080fd5b600a5460a060020a900460ff161561074857600a54600160a060020a031663aee92d3333858560405160e060020a63ffffffff8616028152600160a060020a0393841660048201529190921660248201526044810191909152606401600060405180830381600087803b151561072f57600080fd5b6102c65a03f1151561074057600080fd5b505050610752565b610752838361110c565b505050565b600a5460a060020a900460ff1681565b60005433600160a060020a0390811691161461078257600080fd5b600160a060020a03811660009081526006602052604090819020805460ff191660011790557f42e160154868087d6bfdc0ca23d96a1c1cfa32f1b72ba9ba27b69b98a0d819dc90829051600160a060020a03909116815260200160405180910390a150565b600a5460009060a060020a900460ff161561086657600a54600160a060020a03166318160ddd6000604051602001526040518163ffffffff1660e060020a028152600401602060405180830381600087803b151561084457600080fd5b6102c65a03f1151561085557600080fd5b50505060405180519050905061086b565b506001545b90565b60005460a060020a900460ff161561088557600080fd5b600160a060020a03831660009081526006602052604090205460ff16156108ab57600080fd5b600a5460a060020a900460ff161561092757600a54600160a060020a0316638b477adb3385858560405160e060020a63ffffffff8716028152600160a060020a0394851660048201529284166024840152921660448201526064810191909152608401600060405180830381600087803b151561072f57600080fd5b6107528383836111be565b600a54600160a060020a031681565b60026020526000908152604090205481565b60095481565b60045481565b60015481565b60005433600160a060020a0390811691161461098057600080fd5b60005460a060020a900460ff16151561099857600080fd5b6000805474ff0000000000000000000000000000000000000000191690557f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3360405160405180910390a1565b600160a060020a03811660009081526006602052604090205460ff165b919050565b600560209081526000928352604080842090915290825290205481565b60005460a060020a900460ff1681565b600a5460009060a060020a900460ff1615610ac357600a54600160a060020a03166370a082318360006040516020015260405160e060020a63ffffffff8416028152600160a060020a039091166004820152602401602060405180830381600087803b1515610aa157600080fd5b6102c65a03f11515610ab257600080fd5b505050604051805190509050610a01565b610acc826113bd565b9050610a01565b60005433600160a060020a03908116911614610aee57600080fd5b60005460a060020a900460ff1615610b0557600080fd5b6000805474ff0000000000000000000000000000000000000000191660a060020a1790557f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62560405160405180910390a1565b600054600160a060020a031690565b600054600160a060020a031681565b60088054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105ff5780601f106105d4576101008083540402835291602001916105ff565b60005460a060020a900460ff1615610bf757600080fd5b600160a060020a03331660009081526006602052604090205460ff1615610c1d57600080fd5b600a5460a060020a900460ff1615610cab57600a54600160a060020a0316636e18980a33848460405160e060020a63ffffffff8616028152600160a060020a0393841660048201529190921660248201526044810191909152606401600060405180830381600087803b1515610c9257600080fd5b6102c65a03f11515610ca357600080fd5b505050610cb5565b610cb582826113d8565b5050565b60005433600160a060020a03908116911614610cd457600080fd5b60148210610ce157600080fd5b60328110610cee57600080fd5b6003829055600954610d0a908290600a0a63ffffffff61155c16565b60048190556003547fb044a1e409eac5c48e5af22d4af52670dd1a99059537a78b31b48c6500a6354e9160405191825260208201526040908101905180910390a15050565b60005433600160a060020a03908116911614610d6a57600080fd5b60015481810111610d7a57600080fd5b60008054600160a060020a031681526002602052604090205481810111610da057600080fd5b60008054600160a060020a03168152600260205260409081902080548301905560018054830190557fcb8241adb0c3fdb35b70c24ce35c5eb0c17af7431c99f827d44a445ca624176a9082905190815260200160405180910390a150565b60005433600160a060020a03908116911614610e1957600080fd5b60015481901015610e2957600080fd5b60008054600160a060020a031681526002602052604090205481901015610e4f57600080fd5b60018054829003905560008054600160a060020a031681526002602052604090819020805483900390557f702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a449082905190815260200160405180910390a150565b600a5460009060a060020a900460ff1615610f4757600a54600160a060020a031663dd62ed3e848460006040516020015260405160e060020a63ffffffff8516028152600160a060020a03928316600482015291166024820152604401602060405180830381600087803b1515610f2557600080fd5b6102c65a03f11515610f3657600080fd5b505050604051805190509050610f54565b610f518383611592565b90505b92915050565b60035481565b60066020526000908152604090205460ff1681565b60005433600160a060020a03908116911614610f9057600080fd5b600160a060020a03811660009081526006602052604090819020805460ff191690557fd7e9ec6e6ecd65492dce6bf513cd6867560d49544421d0783ddf06e76c24470c90829051600160a060020a03909116815260200160405180910390a150565b60001981565b60005433600160a060020a0390811691161461101357600080fd5b600160a060020a0381161561104b576000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0383161790555b50565b6000805433600160a060020a0390811691161461106a57600080fd5b600160a060020a03821660009081526006602052604090205460ff16151561109157600080fd5b61109a82610a33565b600160a060020a038316600090815260026020526040808220919091556001805483900390559091507f61e6e66b0d6339b2980aecc6ccc0039736791f0ccde9ed512e789a7fbdd698c6908390839051600160a060020a03909216825260208201526040908101905180910390a15050565b6040604436101561111c57600080fd5b811580159061114f5750600160a060020a0333811660009081526005602090815260408083209387168352929052205415155b1561115957600080fd5b600160a060020a03338116600081815260056020908152604080832094881680845294909152908190208590557f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259085905190815260200160405180910390a3505050565b60008080606060643610156111d257600080fd5b600160a060020a0380881660009081526005602090815260408083203390941683529290522054600354909450611224906127109061121890889063ffffffff61155c16565b9063ffffffff6115bd16565b92506004548311156112365760045492505b60001984101561127857611250848663ffffffff6115d416565b600160a060020a03808916600090815260056020908152604080832033909416835292905220555b611288858463ffffffff6115d416565b600160a060020a0388166000908152600260205260409020549092506112b4908663ffffffff6115d416565b600160a060020a0380891660009081526002602052604080822093909355908816815220546112e9908363ffffffff6115e616565b600160a060020a03871660009081526002602052604081209190915583111561137f5760008054600160a060020a0316815260026020526040902054611335908463ffffffff6115e616565b60008054600160a060020a03908116825260026020526040808320939093559054811691908916906000805160206115f68339815191529086905190815260200160405180910390a35b85600160a060020a031687600160a060020a03166000805160206115f68339815191528460405190815260200160405180910390a350505050505050565b600160a060020a031660009081526002602052604090205490565b600080604060443610156113eb57600080fd5b6114066127106112186003548761155c90919063ffffffff16565b92506004548311156114185760045492505b611428848463ffffffff6115d416565b600160a060020a033316600090815260026020526040902054909250611454908563ffffffff6115d416565b600160a060020a033381166000908152600260205260408082209390935590871681522054611489908363ffffffff6115e616565b600160a060020a0386166000908152600260205260408120919091558311156115205760008054600160a060020a03168152600260205260409020546114d5908463ffffffff6115e616565b60008054600160a060020a0390811682526002602052604080832093909355905481169133909116906000805160206115f68339815191529086905190815260200160405180910390a35b84600160a060020a031633600160a060020a03166000805160206115f68339815191528460405190815260200160405180910390a35050505050565b60008083151561156f576000915061158b565b5082820282848281151561157f57fe5b041461158757fe5b8091505b5092915050565b600160a060020a03918216600090815260056020908152604080832093909416825291909152205490565b60008082848115156115cb57fe5b04949350505050565b6000828211156115e057fe5b50900390565b60008282018381101561158757fe00ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa165627a7a72305820c544d2a4d0da8a553a7a31fb17a6c587f8437e88c474651a7046ce1f0585e61c0029"

// DeployTetherToken deploys a new Ethereum contract, binding an instance of TetherToken to it.
func DeployTetherToken(auth *bind.TransactOpts, backend bind.ContractBackend, _initialSupply *big.Int, _name string, _symbol string, _decimals *big.Int) (common.Address, *types.Transaction, *TetherToken, error) {
	parsed, err := abi.JSON(strings.NewReader(TetherTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TetherTokenBin), backend, _initialSupply, _name, _symbol, _decimals)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TetherToken{TetherTokenCaller: TetherTokenCaller{contract: contract}, TetherTokenTransactor: TetherTokenTransactor{contract: contract}, TetherTokenFilterer: TetherTokenFilterer{contract: contract}}, nil
}

// TetherToken is an auto generated Go binding around an Ethereum contract.
type TetherToken struct {
	TetherTokenCaller     // Read-only binding to the contract
	TetherTokenTransactor // Write-only binding to the contract
	TetherTokenFilterer   // Log filterer for contract events
}

// TetherTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TetherTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TetherTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TetherTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TetherTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TetherTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TetherTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TetherTokenSession struct {
	Contract     *TetherToken      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TetherTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TetherTokenCallerSession struct {
	Contract *TetherTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// TetherTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TetherTokenTransactorSession struct {
	Contract     *TetherTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TetherTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TetherTokenRaw struct {
	Contract *TetherToken // Generic contract binding to access the raw methods on
}

// TetherTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TetherTokenCallerRaw struct {
	Contract *TetherTokenCaller // Generic read-only contract binding to access the raw methods on
}

// TetherTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TetherTokenTransactorRaw struct {
	Contract *TetherTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTetherToken creates a new instance of TetherToken, bound to a specific deployed contract.
func NewTetherToken(address common.Address, backend bind.ContractBackend) (*TetherToken, error) {
	contract, err := bindTetherToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TetherToken{TetherTokenCaller: TetherTokenCaller{contract: contract}, TetherTokenTransactor: TetherTokenTransactor{contract: contract}, TetherTokenFilterer: TetherTokenFilterer{contract: contract}}, nil
}

// NewTetherTokenCaller creates a new read-only instance of TetherToken, bound to a specific deployed contract.
func NewTetherTokenCaller(address common.Address, caller bind.ContractCaller) (*TetherTokenCaller, error) {
	contract, err := bindTetherToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TetherTokenCaller{contract: contract}, nil
}

// NewTetherTokenTransactor creates a new write-only instance of TetherToken, bound to a specific deployed contract.
func NewTetherTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TetherTokenTransactor, error) {
	contract, err := bindTetherToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TetherTokenTransactor{contract: contract}, nil
}

// NewTetherTokenFilterer creates a new log filterer instance of TetherToken, bound to a specific deployed contract.
func NewTetherTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TetherTokenFilterer, error) {
	contract, err := bindTetherToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TetherTokenFilterer{contract: contract}, nil
}

// bindTetherToken binds a generic wrapper to an already deployed contract.
func bindTetherToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TetherTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TetherToken *TetherTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TetherToken.Contract.TetherTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TetherToken *TetherTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherToken.Contract.TetherTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TetherToken *TetherTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TetherToken.Contract.TetherTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TetherToken *TetherTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TetherToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TetherToken *TetherTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TetherToken *TetherTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TetherToken.Contract.contract.Transact(opts, method, params...)
}

// MAXUINT is a free data retrieval call binding the contract method 0xe5b5019a.
//
// Solidity: function MAX_UINT() constant returns(uint256)
func (_TetherToken *TetherTokenCaller) MAXUINT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "MAX_UINT")
	return *ret0, err
}

// MAXUINT is a free data retrieval call binding the contract method 0xe5b5019a.
//
// Solidity: function MAX_UINT() constant returns(uint256)
func (_TetherToken *TetherTokenSession) MAXUINT() (*big.Int, error) {
	return _TetherToken.Contract.MAXUINT(&_TetherToken.CallOpts)
}

// MAXUINT is a free data retrieval call binding the contract method 0xe5b5019a.
//
// Solidity: function MAX_UINT() constant returns(uint256)
func (_TetherToken *TetherTokenCallerSession) MAXUINT() (*big.Int, error) {
	return _TetherToken.Contract.MAXUINT(&_TetherToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x3eaaf86b.
//
// Solidity: function _totalSupply() constant returns(uint256)
func (_TetherToken *TetherTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) constant returns(uint256 remaining)
func (_TetherToken *TetherTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) constant returns(uint256 remaining)
func (_TetherToken *TetherTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _TetherToken.Contract.Allowance(&_TetherToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) constant returns(uint256 remaining)
func (_TetherToken *TetherTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _TetherToken.Contract.Allowance(&_TetherToken.CallOpts, _owner, _spender)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) constant returns(uint256)
func (_TetherToken *TetherTokenCaller) Allowed(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "allowed", arg0, arg1)
	return *ret0, err
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) constant returns(uint256)
func (_TetherToken *TetherTokenSession) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _TetherToken.Contract.Allowed(&_TetherToken.CallOpts, arg0, arg1)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) constant returns(uint256)
func (_TetherToken *TetherTokenCallerSession) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _TetherToken.Contract.Allowed(&_TetherToken.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address who) constant returns(uint256)
func (_TetherToken *TetherTokenCaller) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "balanceOf", who)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address who) constant returns(uint256)
func (_TetherToken *TetherTokenSession) BalanceOf(who common.Address) (*big.Int, error) {
	return _TetherToken.Contract.BalanceOf(&_TetherToken.CallOpts, who)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address who) constant returns(uint256)
func (_TetherToken *TetherTokenCallerSession) BalanceOf(who common.Address) (*big.Int, error) {
	return _TetherToken.Contract.BalanceOf(&_TetherToken.CallOpts, who)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TetherToken *TetherTokenCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TetherToken *TetherTokenSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TetherToken.Contract.Balances(&_TetherToken.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TetherToken *TetherTokenCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TetherToken.Contract.Balances(&_TetherToken.CallOpts, arg0)
}

// BasisPointsRate is a free data retrieval call binding the contract method 0xdd644f72.
//
// Solidity: function basisPointsRate() constant returns(uint256)
func (_TetherToken *TetherTokenCaller) BasisPointsRate(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "basisPointsRate")
	return *ret0, err
}

// BasisPointsRate is a free data retrieval call binding the contract method 0xdd644f72.
//
// Solidity: function basisPointsRate() constant returns(uint256)
func (_TetherToken *TetherTokenSession) BasisPointsRate() (*big.Int, error) {
	return _TetherToken.Contract.BasisPointsRate(&_TetherToken.CallOpts)
}

// BasisPointsRate is a free data retrieval call binding the contract method 0xdd644f72.
//
// Solidity: function basisPointsRate() constant returns(uint256)
func (_TetherToken *TetherTokenCallerSession) BasisPointsRate() (*big.Int, error) {
	return _TetherToken.Contract.BasisPointsRate(&_TetherToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint256)
func (_TetherToken *TetherTokenCaller) Decimals(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint256)
func (_TetherToken *TetherTokenSession) Decimals() (*big.Int, error) {
	return _TetherToken.Contract.Decimals(&_TetherToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint256)
func (_TetherToken *TetherTokenCallerSession) Decimals() (*big.Int, error) {
	return _TetherToken.Contract.Decimals(&_TetherToken.CallOpts)
}

// Deprecated is a free data retrieval call binding the contract method 0x0e136b19.
//
// Solidity: function deprecated() constant returns(bool)
func (_TetherToken *TetherTokenCaller) Deprecated(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "deprecated")
	return *ret0, err
}

// Deprecated is a free data retrieval call binding the contract method 0x0e136b19.
//
// Solidity: function deprecated() constant returns(bool)
func (_TetherToken *TetherTokenSession) Deprecated() (bool, error) {
	return _TetherToken.Contract.Deprecated(&_TetherToken.CallOpts)
}

// Deprecated is a free data retrieval call binding the contract method 0x0e136b19.
//
// Solidity: function deprecated() constant returns(bool)
func (_TetherToken *TetherTokenCallerSession) Deprecated() (bool, error) {
	return _TetherToken.Contract.Deprecated(&_TetherToken.CallOpts)
}

// GetBlackListStatus is a free data retrieval call binding the contract method 0x59bf1abe.
//
// Solidity: function getBlackListStatus(address _maker) constant returns(bool)
func (_TetherToken *TetherTokenCaller) GetBlackListStatus(opts *bind.CallOpts, _maker common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "getBlackListStatus", _maker)
	return *ret0, err
}

// GetBlackListStatus is a free data retrieval call binding the contract method 0x59bf1abe.
//
// Solidity: function getBlackListStatus(address _maker) constant returns(bool)
func (_TetherToken *TetherTokenSession) GetBlackListStatus(_maker common.Address) (bool, error) {
	return _TetherToken.Contract.GetBlackListStatus(&_TetherToken.CallOpts, _maker)
}

// GetBlackListStatus is a free data retrieval call binding the contract method 0x59bf1abe.
//
// Solidity: function getBlackListStatus(address _maker) constant returns(bool)
func (_TetherToken *TetherTokenCallerSession) GetBlackListStatus(_maker common.Address) (bool, error) {
	return _TetherToken.Contract.GetBlackListStatus(&_TetherToken.CallOpts, _maker)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_TetherToken *TetherTokenCaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "getOwner")
	return *ret0, err
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_TetherToken *TetherTokenSession) GetOwner() (common.Address, error) {
	return _TetherToken.Contract.GetOwner(&_TetherToken.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_TetherToken *TetherTokenCallerSession) GetOwner() (common.Address, error) {
	return _TetherToken.Contract.GetOwner(&_TetherToken.CallOpts)
}

// IsBlackListed is a free data retrieval call binding the contract method 0xe47d6060.
//
// Solidity: function isBlackListed(address ) constant returns(bool)
func (_TetherToken *TetherTokenCaller) IsBlackListed(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "isBlackListed", arg0)
	return *ret0, err
}

// IsBlackListed is a free data retrieval call binding the contract method 0xe47d6060.
//
// Solidity: function isBlackListed(address ) constant returns(bool)
func (_TetherToken *TetherTokenSession) IsBlackListed(arg0 common.Address) (bool, error) {
	return _TetherToken.Contract.IsBlackListed(&_TetherToken.CallOpts, arg0)
}

// IsBlackListed is a free data retrieval call binding the contract method 0xe47d6060.
//
// Solidity: function isBlackListed(address ) constant returns(bool)
func (_TetherToken *TetherTokenCallerSession) IsBlackListed(arg0 common.Address) (bool, error) {
	return _TetherToken.Contract.IsBlackListed(&_TetherToken.CallOpts, arg0)
}

// MaximumFee is a free data retrieval call binding the contract method 0x35390714.
//
// Solidity: function maximumFee() constant returns(uint256)
func (_TetherToken *TetherTokenCaller) MaximumFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "maximumFee")
	return *ret0, err
}

// MaximumFee is a free data retrieval call binding the contract method 0x35390714.
//
// Solidity: function maximumFee() constant returns(uint256)
func (_TetherToken *TetherTokenSession) MaximumFee() (*big.Int, error) {
	return _TetherToken.Contract.MaximumFee(&_TetherToken.CallOpts)
}

// MaximumFee is a free data retrieval call binding the contract method 0x35390714.
//
// Solidity: function maximumFee() constant returns(uint256)
func (_TetherToken *TetherTokenCallerSession) MaximumFee() (*big.Int, error) {
	return _TetherToken.Contract.MaximumFee(&_TetherToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TetherToken *TetherTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TetherToken *TetherTokenSession) Name() (string, error) {
	return _TetherToken.Contract.Name(&_TetherToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TetherToken *TetherTokenCallerSession) Name() (string, error) {
	return _TetherToken.Contract.Name(&_TetherToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TetherToken *TetherTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TetherToken *TetherTokenSession) Owner() (common.Address, error) {
	return _TetherToken.Contract.Owner(&_TetherToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TetherToken *TetherTokenCallerSession) Owner() (common.Address, error) {
	return _TetherToken.Contract.Owner(&_TetherToken.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_TetherToken *TetherTokenCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TetherToken *TetherTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// UpgradedAddress is a free data retrieval call binding the contract method 0x26976e3f.
//
// Solidity: function upgradedAddress() constant returns(address)
func (_TetherToken *TetherTokenCaller) UpgradedAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TetherToken.contract.Call(opts, out, "upgradedAddress")
	return *ret0, err
}

// AddBlackList is a paid mutator transaction binding the contract method 0x0ecb93c0.
//
// Solidity: function addBlackList(address _evilUser) returns()
func (_TetherToken *TetherTokenTransactor) AddBlackList(opts *bind.TransactOpts, _evilUser common.Address) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "addBlackList", _evilUser)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns()
func (_TetherToken *TetherTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "approve", _spender, _value)
}

// Deprecate is a paid mutator transaction binding the contract method 0x0753c30c.
//
// Solidity: function deprecate(address _upgradedAddress) returns()
func (_TetherToken *TetherTokenTransactor) Deprecate(opts *bind.TransactOpts, _upgradedAddress common.Address) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "deprecate", _upgradedAddress)
}

// DestroyBlackFunds is a paid mutator transaction binding the contract method 0xf3bdc228.
//
// Solidity: function destroyBlackFunds(address _blackListedUser) returns()
func (_TetherToken *TetherTokenTransactor) DestroyBlackFunds(opts *bind.TransactOpts, _blackListedUser common.Address) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "destroyBlackFunds", _blackListedUser)
}

// Issue is a paid mutator transaction binding the contract method 0xcc872b66.
//
// Solidity: function issue(uint256 amount) returns()
func (_TetherToken *TetherTokenTransactor) Issue(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "issue", amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TetherToken *TetherTokenTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "pause")
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount) returns()
func (_TetherToken *TetherTokenTransactor) Redeem(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "redeem", amount)
}

// RemoveBlackList is a paid mutator transaction binding the contract method 0xe4997dc5.
//
// Solidity: function removeBlackList(address _clearedUser) returns()
func (_TetherToken *TetherTokenTransactor) RemoveBlackList(opts *bind.TransactOpts, _clearedUser common.Address) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "removeBlackList", _clearedUser)
}

// SetParams is a paid mutator transaction binding the contract method 0xc0324c77.
//
// Solidity: function setParams(uint256 newBasisPoints, uint256 newMaxFee) returns()
func (_TetherToken *TetherTokenTransactor) SetParams(opts *bind.TransactOpts, newBasisPoints *big.Int, newMaxFee *big.Int) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "setParams", newBasisPoints, newMaxFee)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns()
func (_TetherToken *TetherTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "transfer", _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns()
func (_TetherToken *TetherTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TetherToken *TetherTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TetherToken *TetherTokenTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TetherToken.contract.Transact(opts, "unpause")
}

// TetherTokenAddedBlackListIterator is returned from FilterAddedBlackList and is used to iterate over the raw logs and unpacked data for AddedBlackList events raised by the TetherToken contract.
type TetherTokenAddedBlackListIterator struct {
	Event *TetherTokenAddedBlackList // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenAddedBlackListIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenAddedBlackList)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenAddedBlackList)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenAddedBlackListIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenAddedBlackListIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenAddedBlackList represents a AddedBlackList event raised by the TetherToken contract.
type TetherTokenAddedBlackList struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddedBlackList is a free log retrieval operation binding the contract event 0x42e160154868087d6bfdc0ca23d96a1c1cfa32f1b72ba9ba27b69b98a0d819dc.
//
// Solidity: event AddedBlackList(address _user)
func (_TetherToken *TetherTokenFilterer) FilterAddedBlackList(opts *bind.FilterOpts) (*TetherTokenAddedBlackListIterator, error) {

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "AddedBlackList")
	if err != nil {
		return nil, err
	}
	return &TetherTokenAddedBlackListIterator{contract: _TetherToken.contract, event: "AddedBlackList", logs: logs, sub: sub}, nil
}

// WatchAddedBlackList is a free log subscription operation binding the contract event 0x42e160154868087d6bfdc0ca23d96a1c1cfa32f1b72ba9ba27b69b98a0d819dc.
//
// Solidity: event AddedBlackList(address _user)
func (_TetherToken *TetherTokenFilterer) WatchAddedBlackList(opts *bind.WatchOpts, sink chan<- *TetherTokenAddedBlackList) (event.Subscription, error) {

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "AddedBlackList")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenAddedBlackList)
				if err := _TetherToken.contract.UnpackLog(event, "AddedBlackList", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddedBlackList is a log parse operation binding the contract event 0x42e160154868087d6bfdc0ca23d96a1c1cfa32f1b72ba9ba27b69b98a0d819dc.
//
// Solidity: event AddedBlackList(address _user)
func (_TetherToken *TetherTokenFilterer) ParseAddedBlackList(log types.Log) (*TetherTokenAddedBlackList, error) {
	event := new(TetherTokenAddedBlackList)
	if err := _TetherToken.contract.UnpackLog(event, "AddedBlackList", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TetherToken contract.
type TetherTokenApprovalIterator struct {
	Event *TetherTokenApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenApproval represents a Approval event raised by the TetherToken contract.
type TetherTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TetherToken *TetherTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TetherTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TetherTokenApprovalIterator{contract: _TetherToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TetherToken *TetherTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TetherTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenApproval)
				if err := _TetherToken.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TetherToken *TetherTokenFilterer) ParseApproval(log types.Log) (*TetherTokenApproval, error) {
	event := new(TetherTokenApproval)
	if err := _TetherToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenDeprecateIterator is returned from FilterDeprecate and is used to iterate over the raw logs and unpacked data for Deprecate events raised by the TetherToken contract.
type TetherTokenDeprecateIterator struct {
	Event *TetherTokenDeprecate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenDeprecateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenDeprecate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenDeprecate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenDeprecateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenDeprecateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenDeprecate represents a Deprecate event raised by the TetherToken contract.
type TetherTokenDeprecate struct {
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDeprecate is a free log retrieval operation binding the contract event 0xcc358699805e9a8b7f77b522628c7cb9abd07d9efb86b6fb616af1609036a99e.
//
// Solidity: event Deprecate(address newAddress)
func (_TetherToken *TetherTokenFilterer) FilterDeprecate(opts *bind.FilterOpts) (*TetherTokenDeprecateIterator, error) {

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "Deprecate")
	if err != nil {
		return nil, err
	}
	return &TetherTokenDeprecateIterator{contract: _TetherToken.contract, event: "Deprecate", logs: logs, sub: sub}, nil
}

// WatchDeprecate is a free log subscription operation binding the contract event 0xcc358699805e9a8b7f77b522628c7cb9abd07d9efb86b6fb616af1609036a99e.
//
// Solidity: event Deprecate(address newAddress)
func (_TetherToken *TetherTokenFilterer) WatchDeprecate(opts *bind.WatchOpts, sink chan<- *TetherTokenDeprecate) (event.Subscription, error) {

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "Deprecate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenDeprecate)
				if err := _TetherToken.contract.UnpackLog(event, "Deprecate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeprecate is a log parse operation binding the contract event 0xcc358699805e9a8b7f77b522628c7cb9abd07d9efb86b6fb616af1609036a99e.
//
// Solidity: event Deprecate(address newAddress)
func (_TetherToken *TetherTokenFilterer) ParseDeprecate(log types.Log) (*TetherTokenDeprecate, error) {
	event := new(TetherTokenDeprecate)
	if err := _TetherToken.contract.UnpackLog(event, "Deprecate", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenDestroyedBlackFundsIterator is returned from FilterDestroyedBlackFunds and is used to iterate over the raw logs and unpacked data for DestroyedBlackFunds events raised by the TetherToken contract.
type TetherTokenDestroyedBlackFundsIterator struct {
	Event *TetherTokenDestroyedBlackFunds // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenDestroyedBlackFundsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenDestroyedBlackFunds)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenDestroyedBlackFunds)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenDestroyedBlackFundsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenDestroyedBlackFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenDestroyedBlackFunds represents a DestroyedBlackFunds event raised by the TetherToken contract.
type TetherTokenDestroyedBlackFunds struct {
	BlackListedUser common.Address
	Balance         *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterDestroyedBlackFunds is a free log retrieval operation binding the contract event 0x61e6e66b0d6339b2980aecc6ccc0039736791f0ccde9ed512e789a7fbdd698c6.
//
// Solidity: event DestroyedBlackFunds(address _blackListedUser, uint256 _balance)
func (_TetherToken *TetherTokenFilterer) FilterDestroyedBlackFunds(opts *bind.FilterOpts) (*TetherTokenDestroyedBlackFundsIterator, error) {

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "DestroyedBlackFunds")
	if err != nil {
		return nil, err
	}
	return &TetherTokenDestroyedBlackFundsIterator{contract: _TetherToken.contract, event: "DestroyedBlackFunds", logs: logs, sub: sub}, nil
}

// WatchDestroyedBlackFunds is a free log subscription operation binding the contract event 0x61e6e66b0d6339b2980aecc6ccc0039736791f0ccde9ed512e789a7fbdd698c6.
//
// Solidity: event DestroyedBlackFunds(address _blackListedUser, uint256 _balance)
func (_TetherToken *TetherTokenFilterer) WatchDestroyedBlackFunds(opts *bind.WatchOpts, sink chan<- *TetherTokenDestroyedBlackFunds) (event.Subscription, error) {

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "DestroyedBlackFunds")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenDestroyedBlackFunds)
				if err := _TetherToken.contract.UnpackLog(event, "DestroyedBlackFunds", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDestroyedBlackFunds is a log parse operation binding the contract event 0x61e6e66b0d6339b2980aecc6ccc0039736791f0ccde9ed512e789a7fbdd698c6.
//
// Solidity: event DestroyedBlackFunds(address _blackListedUser, uint256 _balance)
func (_TetherToken *TetherTokenFilterer) ParseDestroyedBlackFunds(log types.Log) (*TetherTokenDestroyedBlackFunds, error) {
	event := new(TetherTokenDestroyedBlackFunds)
	if err := _TetherToken.contract.UnpackLog(event, "DestroyedBlackFunds", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenIssueIterator is returned from FilterIssue and is used to iterate over the raw logs and unpacked data for Issue events raised by the TetherToken contract.
type TetherTokenIssueIterator struct {
	Event *TetherTokenIssue // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenIssueIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenIssue)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenIssue)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenIssueIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenIssueIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenIssue represents a Issue event raised by the TetherToken contract.
type TetherTokenIssue struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterIssue is a free log retrieval operation binding the contract event 0xcb8241adb0c3fdb35b70c24ce35c5eb0c17af7431c99f827d44a445ca624176a.
//
// Solidity: event Issue(uint256 amount)
func (_TetherToken *TetherTokenFilterer) FilterIssue(opts *bind.FilterOpts) (*TetherTokenIssueIterator, error) {

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "Issue")
	if err != nil {
		return nil, err
	}
	return &TetherTokenIssueIterator{contract: _TetherToken.contract, event: "Issue", logs: logs, sub: sub}, nil
}

// WatchIssue is a free log subscription operation binding the contract event 0xcb8241adb0c3fdb35b70c24ce35c5eb0c17af7431c99f827d44a445ca624176a.
//
// Solidity: event Issue(uint256 amount)
func (_TetherToken *TetherTokenFilterer) WatchIssue(opts *bind.WatchOpts, sink chan<- *TetherTokenIssue) (event.Subscription, error) {

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "Issue")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenIssue)
				if err := _TetherToken.contract.UnpackLog(event, "Issue", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIssue is a log parse operation binding the contract event 0xcb8241adb0c3fdb35b70c24ce35c5eb0c17af7431c99f827d44a445ca624176a.
//
// Solidity: event Issue(uint256 amount)
func (_TetherToken *TetherTokenFilterer) ParseIssue(log types.Log) (*TetherTokenIssue, error) {
	event := new(TetherTokenIssue)
	if err := _TetherToken.contract.UnpackLog(event, "Issue", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenParamsIterator is returned from FilterParams and is used to iterate over the raw logs and unpacked data for Params events raised by the TetherToken contract.
type TetherTokenParamsIterator struct {
	Event *TetherTokenParams // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenParamsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenParams)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenParams)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenParamsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenParamsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenParams represents a Params event raised by the TetherToken contract.
type TetherTokenParams struct {
	FeeBasisPoints *big.Int
	MaxFee         *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterParams is a free log retrieval operation binding the contract event 0xb044a1e409eac5c48e5af22d4af52670dd1a99059537a78b31b48c6500a6354e.
//
// Solidity: event Params(uint256 feeBasisPoints, uint256 maxFee)
func (_TetherToken *TetherTokenFilterer) FilterParams(opts *bind.FilterOpts) (*TetherTokenParamsIterator, error) {

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "Params")
	if err != nil {
		return nil, err
	}
	return &TetherTokenParamsIterator{contract: _TetherToken.contract, event: "Params", logs: logs, sub: sub}, nil
}

// WatchParams is a free log subscription operation binding the contract event 0xb044a1e409eac5c48e5af22d4af52670dd1a99059537a78b31b48c6500a6354e.
//
// Solidity: event Params(uint256 feeBasisPoints, uint256 maxFee)
func (_TetherToken *TetherTokenFilterer) WatchParams(opts *bind.WatchOpts, sink chan<- *TetherTokenParams) (event.Subscription, error) {

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "Params")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenParams)
				if err := _TetherToken.contract.UnpackLog(event, "Params", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseParams is a log parse operation binding the contract event 0xb044a1e409eac5c48e5af22d4af52670dd1a99059537a78b31b48c6500a6354e.
//
// Solidity: event Params(uint256 feeBasisPoints, uint256 maxFee)
func (_TetherToken *TetherTokenFilterer) ParseParams(log types.Log) (*TetherTokenParams, error) {
	event := new(TetherTokenParams)
	if err := _TetherToken.contract.UnpackLog(event, "Params", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TetherToken contract.
type TetherTokenPauseIterator struct {
	Event *TetherTokenPause // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenPause)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenPause)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenPause represents a Pause event raised by the TetherToken contract.
type TetherTokenPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TetherToken *TetherTokenFilterer) FilterPause(opts *bind.FilterOpts) (*TetherTokenPauseIterator, error) {

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TetherTokenPauseIterator{contract: _TetherToken.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TetherToken *TetherTokenFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TetherTokenPause) (event.Subscription, error) {

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenPause)
				if err := _TetherToken.contract.UnpackLog(event, "Pause", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePause is a log parse operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TetherToken *TetherTokenFilterer) ParsePause(log types.Log) (*TetherTokenPause, error) {
	event := new(TetherTokenPause)
	if err := _TetherToken.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenRedeemIterator is returned from FilterRedeem and is used to iterate over the raw logs and unpacked data for Redeem events raised by the TetherToken contract.
type TetherTokenRedeemIterator struct {
	Event *TetherTokenRedeem // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenRedeemIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenRedeem)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenRedeem)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenRedeemIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenRedeemIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenRedeem represents a Redeem event raised by the TetherToken contract.
type TetherTokenRedeem struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRedeem is a free log retrieval operation binding the contract event 0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44.
//
// Solidity: event Redeem(uint256 amount)
func (_TetherToken *TetherTokenFilterer) FilterRedeem(opts *bind.FilterOpts) (*TetherTokenRedeemIterator, error) {

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "Redeem")
	if err != nil {
		return nil, err
	}
	return &TetherTokenRedeemIterator{contract: _TetherToken.contract, event: "Redeem", logs: logs, sub: sub}, nil
}

// WatchRedeem is a free log subscription operation binding the contract event 0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44.
//
// Solidity: event Redeem(uint256 amount)
func (_TetherToken *TetherTokenFilterer) WatchRedeem(opts *bind.WatchOpts, sink chan<- *TetherTokenRedeem) (event.Subscription, error) {

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "Redeem")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenRedeem)
				if err := _TetherToken.contract.UnpackLog(event, "Redeem", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedeem is a log parse operation binding the contract event 0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44.
//
// Solidity: event Redeem(uint256 amount)
func (_TetherToken *TetherTokenFilterer) ParseRedeem(log types.Log) (*TetherTokenRedeem, error) {
	event := new(TetherTokenRedeem)
	if err := _TetherToken.contract.UnpackLog(event, "Redeem", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenRemovedBlackListIterator is returned from FilterRemovedBlackList and is used to iterate over the raw logs and unpacked data for RemovedBlackList events raised by the TetherToken contract.
type TetherTokenRemovedBlackListIterator struct {
	Event *TetherTokenRemovedBlackList // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenRemovedBlackListIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenRemovedBlackList)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenRemovedBlackList)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenRemovedBlackListIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenRemovedBlackListIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenRemovedBlackList represents a RemovedBlackList event raised by the TetherToken contract.
type TetherTokenRemovedBlackList struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedBlackList is a free log retrieval operation binding the contract event 0xd7e9ec6e6ecd65492dce6bf513cd6867560d49544421d0783ddf06e76c24470c.
//
// Solidity: event RemovedBlackList(address _user)
func (_TetherToken *TetherTokenFilterer) FilterRemovedBlackList(opts *bind.FilterOpts) (*TetherTokenRemovedBlackListIterator, error) {

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "RemovedBlackList")
	if err != nil {
		return nil, err
	}
	return &TetherTokenRemovedBlackListIterator{contract: _TetherToken.contract, event: "RemovedBlackList", logs: logs, sub: sub}, nil
}

// WatchRemovedBlackList is a free log subscription operation binding the contract event 0xd7e9ec6e6ecd65492dce6bf513cd6867560d49544421d0783ddf06e76c24470c.
//
// Solidity: event RemovedBlackList(address _user)
func (_TetherToken *TetherTokenFilterer) WatchRemovedBlackList(opts *bind.WatchOpts, sink chan<- *TetherTokenRemovedBlackList) (event.Subscription, error) {

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "RemovedBlackList")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenRemovedBlackList)
				if err := _TetherToken.contract.UnpackLog(event, "RemovedBlackList", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRemovedBlackList is a log parse operation binding the contract event 0xd7e9ec6e6ecd65492dce6bf513cd6867560d49544421d0783ddf06e76c24470c.
//
// Solidity: event RemovedBlackList(address _user)
func (_TetherToken *TetherTokenFilterer) ParseRemovedBlackList(log types.Log) (*TetherTokenRemovedBlackList, error) {
	event := new(TetherTokenRemovedBlackList)
	if err := _TetherToken.contract.UnpackLog(event, "RemovedBlackList", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TetherToken contract.
type TetherTokenTransferIterator struct {
	Event *TetherTokenTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenTransfer represents a Transfer event raised by the TetherToken contract.
type TetherTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TetherToken *TetherTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TetherTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TetherTokenTransferIterator{contract: _TetherToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TetherToken *TetherTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TetherTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenTransfer)
				if err := _TetherToken.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TetherToken *TetherTokenFilterer) ParseTransfer(log types.Log) (*TetherTokenTransfer, error) {
	event := new(TetherTokenTransfer)
	if err := _TetherToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TetherTokenUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TetherToken contract.
type TetherTokenUnpauseIterator struct {
	Event *TetherTokenUnpause // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TetherTokenUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TetherTokenUnpause)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TetherTokenUnpause)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TetherTokenUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TetherTokenUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TetherTokenUnpause represents a Unpause event raised by the TetherToken contract.
type TetherTokenUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TetherToken *TetherTokenFilterer) FilterUnpause(opts *bind.FilterOpts) (*TetherTokenUnpauseIterator, error) {

	logs, sub, err := _TetherToken.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TetherTokenUnpauseIterator{contract: _TetherToken.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TetherToken *TetherTokenFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TetherTokenUnpause) (event.Subscription, error) {

	logs, sub, err := _TetherToken.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TetherTokenUnpause)
				if err := _TetherToken.contract.UnpackLog(event, "Unpause", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpause is a log parse operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TetherToken *TetherTokenFilterer) ParseUnpause(log types.Log) (*TetherTokenUnpause, error) {
	event := new(TetherTokenUnpause)
	if err := _TetherToken.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	return event, nil
}

// UpgradedStandardTokenABI is the input ABI used to generate the binding from.
const UpgradedStandardTokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maximumFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferByLegacy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFromByLegacy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approveByLegacy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"basisPointsRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_UINT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// UpgradedStandardTokenFuncSigs maps the 4-byte function signature to its string representation.
var UpgradedStandardTokenFuncSigs = map[string]string{
	"e5b5019a": "MAX_UINT()",
	"3eaaf86b": "_totalSupply()",
	"dd62ed3e": "allowance(address,address)",
	"5c658165": "allowed(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"aee92d33": "approveByLegacy(address,address,uint256)",
	"70a08231": "balanceOf(address)",
	"27e235e3": "balances(address)",
	"dd644f72": "basisPointsRate()",
	"35390714": "maximumFee()",
	"8da5cb5b": "owner()",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"6e18980a": "transferByLegacy(address,address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
	"8b477adb": "transferFromByLegacy(address,address,address,uint256)",
	"f2fde38b": "transferOwnership(address)",
}

// UpgradedStandardToken is an auto generated Go binding around an Ethereum contract.
type UpgradedStandardToken struct {
	UpgradedStandardTokenCaller     // Read-only binding to the contract
	UpgradedStandardTokenTransactor // Write-only binding to the contract
	UpgradedStandardTokenFilterer   // Log filterer for contract events
}

// UpgradedStandardTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpgradedStandardTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradedStandardTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradedStandardTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradedStandardTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradedStandardTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradedStandardTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpgradedStandardTokenSession struct {
	Contract     *UpgradedStandardToken // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// UpgradedStandardTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpgradedStandardTokenCallerSession struct {
	Contract *UpgradedStandardTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// UpgradedStandardTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpgradedStandardTokenTransactorSession struct {
	Contract     *UpgradedStandardTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// UpgradedStandardTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpgradedStandardTokenRaw struct {
	Contract *UpgradedStandardToken // Generic contract binding to access the raw methods on
}

// UpgradedStandardTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpgradedStandardTokenCallerRaw struct {
	Contract *UpgradedStandardTokenCaller // Generic read-only contract binding to access the raw methods on
}

// UpgradedStandardTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpgradedStandardTokenTransactorRaw struct {
	Contract *UpgradedStandardTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradedStandardToken creates a new instance of UpgradedStandardToken, bound to a specific deployed contract.
func NewUpgradedStandardToken(address common.Address, backend bind.ContractBackend) (*UpgradedStandardToken, error) {
	contract, err := bindUpgradedStandardToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpgradedStandardToken{UpgradedStandardTokenCaller: UpgradedStandardTokenCaller{contract: contract}, UpgradedStandardTokenTransactor: UpgradedStandardTokenTransactor{contract: contract}, UpgradedStandardTokenFilterer: UpgradedStandardTokenFilterer{contract: contract}}, nil
}

// NewUpgradedStandardTokenCaller creates a new read-only instance of UpgradedStandardToken, bound to a specific deployed contract.
func NewUpgradedStandardTokenCaller(address common.Address, caller bind.ContractCaller) (*UpgradedStandardTokenCaller, error) {
	contract, err := bindUpgradedStandardToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradedStandardTokenCaller{contract: contract}, nil
}

// NewUpgradedStandardTokenTransactor creates a new write-only instance of UpgradedStandardToken, bound to a specific deployed contract.
func NewUpgradedStandardTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*UpgradedStandardTokenTransactor, error) {
	contract, err := bindUpgradedStandardToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradedStandardTokenTransactor{contract: contract}, nil
}

// NewUpgradedStandardTokenFilterer creates a new log filterer instance of UpgradedStandardToken, bound to a specific deployed contract.
func NewUpgradedStandardTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*UpgradedStandardTokenFilterer, error) {
	contract, err := bindUpgradedStandardToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradedStandardTokenFilterer{contract: contract}, nil
}

// bindUpgradedStandardToken binds a generic wrapper to an already deployed contract.
func bindUpgradedStandardToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradedStandardTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradedStandardToken *UpgradedStandardTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _UpgradedStandardToken.Contract.UpgradedStandardTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradedStandardToken *UpgradedStandardTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradedStandardToken.Contract.UpgradedStandardTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradedStandardToken *UpgradedStandardTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradedStandardToken.Contract.UpgradedStandardTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradedStandardToken *UpgradedStandardTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _UpgradedStandardToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradedStandardToken *UpgradedStandardTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradedStandardToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradedStandardToken *UpgradedStandardTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradedStandardToken.Contract.contract.Transact(opts, method, params...)
}

// MAXUINT is a free data retrieval call binding the contract method 0xe5b5019a.
//
// Solidity: function MAX_UINT() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCaller) MAXUINT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UpgradedStandardToken.contract.Call(opts, out, "MAX_UINT")
	return *ret0, err
}

// MAXUINT is a free data retrieval call binding the contract method 0xe5b5019a.
//
// Solidity: function MAX_UINT() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenSession) MAXUINT() (*big.Int, error) {
	return _UpgradedStandardToken.Contract.MAXUINT(&_UpgradedStandardToken.CallOpts)
}

// MAXUINT is a free data retrieval call binding the contract method 0xe5b5019a.
//
// Solidity: function MAX_UINT() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCallerSession) MAXUINT() (*big.Int, error) {
	return _UpgradedStandardToken.Contract.MAXUINT(&_UpgradedStandardToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) constant returns(uint256 remaining)
func (_UpgradedStandardToken *UpgradedStandardTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UpgradedStandardToken.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) constant returns(uint256 remaining)
func (_UpgradedStandardToken *UpgradedStandardTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _UpgradedStandardToken.Contract.Allowance(&_UpgradedStandardToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) constant returns(uint256 remaining)
func (_UpgradedStandardToken *UpgradedStandardTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _UpgradedStandardToken.Contract.Allowance(&_UpgradedStandardToken.CallOpts, _owner, _spender)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCaller) Allowed(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UpgradedStandardToken.contract.Call(opts, out, "allowed", arg0, arg1)
	return *ret0, err
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenSession) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _UpgradedStandardToken.Contract.Allowed(&_UpgradedStandardToken.CallOpts, arg0, arg1)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCallerSession) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _UpgradedStandardToken.Contract.Allowed(&_UpgradedStandardToken.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) constant returns(uint256 balance)
func (_UpgradedStandardToken *UpgradedStandardTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UpgradedStandardToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) constant returns(uint256 balance)
func (_UpgradedStandardToken *UpgradedStandardTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _UpgradedStandardToken.Contract.BalanceOf(&_UpgradedStandardToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) constant returns(uint256 balance)
func (_UpgradedStandardToken *UpgradedStandardTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _UpgradedStandardToken.Contract.BalanceOf(&_UpgradedStandardToken.CallOpts, _owner)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UpgradedStandardToken.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _UpgradedStandardToken.Contract.Balances(&_UpgradedStandardToken.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _UpgradedStandardToken.Contract.Balances(&_UpgradedStandardToken.CallOpts, arg0)
}

// BasisPointsRate is a free data retrieval call binding the contract method 0xdd644f72.
//
// Solidity: function basisPointsRate() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCaller) BasisPointsRate(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UpgradedStandardToken.contract.Call(opts, out, "basisPointsRate")
	return *ret0, err
}

// BasisPointsRate is a free data retrieval call binding the contract method 0xdd644f72.
//
// Solidity: function basisPointsRate() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenSession) BasisPointsRate() (*big.Int, error) {
	return _UpgradedStandardToken.Contract.BasisPointsRate(&_UpgradedStandardToken.CallOpts)
}

// BasisPointsRate is a free data retrieval call binding the contract method 0xdd644f72.
//
// Solidity: function basisPointsRate() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCallerSession) BasisPointsRate() (*big.Int, error) {
	return _UpgradedStandardToken.Contract.BasisPointsRate(&_UpgradedStandardToken.CallOpts)
}

// MaximumFee is a free data retrieval call binding the contract method 0x35390714.
//
// Solidity: function maximumFee() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCaller) MaximumFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UpgradedStandardToken.contract.Call(opts, out, "maximumFee")
	return *ret0, err
}

// MaximumFee is a free data retrieval call binding the contract method 0x35390714.
//
// Solidity: function maximumFee() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenSession) MaximumFee() (*big.Int, error) {
	return _UpgradedStandardToken.Contract.MaximumFee(&_UpgradedStandardToken.CallOpts)
}

// MaximumFee is a free data retrieval call binding the contract method 0x35390714.
//
// Solidity: function maximumFee() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCallerSession) MaximumFee() (*big.Int, error) {
	return _UpgradedStandardToken.Contract.MaximumFee(&_UpgradedStandardToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_UpgradedStandardToken *UpgradedStandardTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _UpgradedStandardToken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_UpgradedStandardToken *UpgradedStandardTokenSession) Owner() (common.Address, error) {
	return _UpgradedStandardToken.Contract.Owner(&_UpgradedStandardToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_UpgradedStandardToken *UpgradedStandardTokenCallerSession) Owner() (common.Address, error) {
	return _UpgradedStandardToken.Contract.Owner(&_UpgradedStandardToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_UpgradedStandardToken *UpgradedStandardTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UpgradedStandardToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns()
func (_UpgradedStandardToken *UpgradedStandardTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _UpgradedStandardToken.contract.Transact(opts, "approve", _spender, _value)
}

// ApproveByLegacy is a paid mutator transaction binding the contract method 0xaee92d33.
//
// Solidity: function approveByLegacy(address from, address spender, uint256 value) returns()
func (_UpgradedStandardToken *UpgradedStandardTokenTransactor) ApproveByLegacy(opts *bind.TransactOpts, from common.Address, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _UpgradedStandardToken.contract.Transact(opts, "approveByLegacy", from, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns()
func (_UpgradedStandardToken *UpgradedStandardTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _UpgradedStandardToken.contract.Transact(opts, "transfer", _to, _value)
}

// TransferByLegacy is a paid mutator transaction binding the contract method 0x6e18980a.
//
// Solidity: function transferByLegacy(address from, address to, uint256 value) returns()
func (_UpgradedStandardToken *UpgradedStandardTokenTransactor) TransferByLegacy(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UpgradedStandardToken.contract.Transact(opts, "transferByLegacy", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns()
func (_UpgradedStandardToken *UpgradedStandardTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _UpgradedStandardToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFromByLegacy is a paid mutator transaction binding the contract method 0x8b477adb.
//
// Solidity: function transferFromByLegacy(address sender, address from, address spender, uint256 value) returns()
func (_UpgradedStandardToken *UpgradedStandardTokenTransactor) TransferFromByLegacy(opts *bind.TransactOpts, sender common.Address, from common.Address, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _UpgradedStandardToken.contract.Transact(opts, "transferFromByLegacy", sender, from, spender, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UpgradedStandardToken *UpgradedStandardTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _UpgradedStandardToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// UpgradedStandardTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the UpgradedStandardToken contract.
type UpgradedStandardTokenApprovalIterator struct {
	Event *UpgradedStandardTokenApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradedStandardTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradedStandardTokenApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UpgradedStandardTokenApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UpgradedStandardTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradedStandardTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradedStandardTokenApproval represents a Approval event raised by the UpgradedStandardToken contract.
type UpgradedStandardTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_UpgradedStandardToken *UpgradedStandardTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*UpgradedStandardTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _UpgradedStandardToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &UpgradedStandardTokenApprovalIterator{contract: _UpgradedStandardToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_UpgradedStandardToken *UpgradedStandardTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *UpgradedStandardTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _UpgradedStandardToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradedStandardTokenApproval)
				if err := _UpgradedStandardToken.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_UpgradedStandardToken *UpgradedStandardTokenFilterer) ParseApproval(log types.Log) (*UpgradedStandardTokenApproval, error) {
	event := new(UpgradedStandardTokenApproval)
	if err := _UpgradedStandardToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// UpgradedStandardTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the UpgradedStandardToken contract.
type UpgradedStandardTokenTransferIterator struct {
	Event *UpgradedStandardTokenTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradedStandardTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradedStandardTokenTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UpgradedStandardTokenTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UpgradedStandardTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradedStandardTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradedStandardTokenTransfer represents a Transfer event raised by the UpgradedStandardToken contract.
type UpgradedStandardTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_UpgradedStandardToken *UpgradedStandardTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*UpgradedStandardTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpgradedStandardToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UpgradedStandardTokenTransferIterator{contract: _UpgradedStandardToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_UpgradedStandardToken *UpgradedStandardTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *UpgradedStandardTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UpgradedStandardToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradedStandardTokenTransfer)
				if err := _UpgradedStandardToken.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_UpgradedStandardToken *UpgradedStandardTokenFilterer) ParseTransfer(log types.Log) (*UpgradedStandardTokenTransfer, error) {
	event := new(UpgradedStandardTokenTransfer)
	if err := _UpgradedStandardToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}
