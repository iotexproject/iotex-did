pragma solidity >=0.4.22 <0.6.0;

import './DIDManagerBase.sol';

contract AddressBasedDIDManager is DIDManagerBase {

    constructor(string memory _prefix, address _dbAddr) DIDBase(_dbAddr) public {
        prefix = _prefix;
    }

    function registerDID(bytes32 h, string memory uri) public {
        internalCreateDID(getDID(msg.sender), msg.sender, h, uri);
    }

    function updateDID(bytes32 h, string memory uri) public {
        internalUpdateDID(getDID(msg.sender), msg.sender, h, uri);
    }

    function deregisterDID() public {
        internalDeleteDID(getDID(msg.sender), msg.sender);
    }

    function addrToString(address _addr) internal pure returns(string memory) {
        bytes32 value = bytes32(uint256(_addr));
        bytes memory alphabet = "0123456789abcdef";

        bytes memory str = new bytes(42);
        str[0] = '0';
        str[1] = 'x';
        for (uint i = 0; i < 20; i++) {
            str[2+i*2] = alphabet[uint8(value[i + 12] >> 4)];
            str[3+i*2] = alphabet[uint8(value[i + 12] & 0x0f)];
        }
        return string(str);
    }

    function getDID(address addr) internal view returns (string memory) {
		bytes memory bStr = abi.encodePacked(prefix, addrToString(addr));
		bytes memory bLower = new bytes(bStr.length);
		for (uint i = 0; i < bStr.length; i++) {
			if ((uint8(bStr[i]) >= 65) && (uint8(bStr[i]) <= 90)) {
				bLower[i] = bytes1(int8(bStr[i]) + 32);
			} else {
				bLower[i] = bStr[i];
			}
		}
		return string(bLower);
	}
}
