pragma solidity >=0.4.22 <0.6.0;

import './DIDManagerBase.sol';

contract AddressBasedDIDManager is DIDManagerBase {

    constructor(bytes memory _prefix, address _dbAddr) DIDBase(_dbAddr, _prefix) public {}

    function decodeInternalKey(bytes memory did) public view returns (bytes20) {
        require(hasPrefix(did, db.getPrefix()), "invalid DID");
        bytes memory domainID = (slice(did, db.getPrefix().length));
        require(domainID.length == 42 && domainID[0] == '0' && domainID[1] == 'x', "invalid DID");
        uint160 iaddr = 0;
        uint160 b1;
        uint160 b2;
        for (uint i=2; i<2+2*20; i+=2){
            iaddr *= 256;
            b1 = uint8(domainID[i]);
            b2 = uint8(domainID[i+1]);
            if ((b1 >= 97)&&(b1 <= 102)) b1 -= 87;
            else if ((b1 >= 48)&&(b1 <= 57)) b1 -= 48;
            else if ((b1 >= 65)&&(b1 <= 70)) b1 -= 55;
            if ((b2 >= 97)&&(b2 <= 102)) b2 -= 87;
            else if ((b2 >= 48)&&(b2 <= 57)) b2 -= 48;
            else if ((b2 >= 65)&&(b2 <= 70)) b2 -= 55;
            iaddr += (b1*16+b2);
        }
        return bytes20(iaddr);
    }

    function registerDID(bytes32 h, bytes memory uri) public {
        internalCreateDID(getDID(msg.sender), bytes20(msg.sender), msg.sender, h, uri);
    }

    function updateDID(bytes32 h, bytes memory uri) public {
        internalUpdateDID(getDID(msg.sender), bytes20(msg.sender), msg.sender, h, uri);
    }

    function deregisterDID() public {
        internalDeleteDID(getDID(msg.sender), bytes20(msg.sender), msg.sender);
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

    function getDID(address addr) internal view returns (bytes memory) {
		bytes memory bStr = abi.encodePacked(db.getPrefix(), addrToString(addr));
		bytes memory bLower = new bytes(bStr.length);
		for (uint i = 0; i < bStr.length; i++) {
			if ((uint8(bStr[i]) >= 65) && (uint8(bStr[i]) <= 90)) {
				bLower[i] = bytes1(int8(bStr[i]) + 32);
			} else {
				bLower[i] = bStr[i];
			}
		}
		return bLower;
	}
}
