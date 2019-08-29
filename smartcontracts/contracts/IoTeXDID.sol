pragma solidity ^0.4.24;

import './SelfManagedDID.sol';

contract IoTeXDID is SelfManagedDID {
    modifier onlyDIDOwner() {
        require(dids[generateDIDString()].exist, "not a did owner");
        _;
    }

    string constant didPrefix = "did:io:";
    struct DID {
        bool exist;
        bytes32 hash;
        string uri;
    }
    mapping(string => DID) dids;

    function createDID(string id, bytes32 hash, string uri) public returns (string) {
        string memory resultDID = generateDIDString();
        require(!dids[resultDID].exist, "did already exists");
        dids[resultDID] = DID(true, hash, uri);
        return resultDID;
    }

    function updateHash(string did, bytes32 hash) public onlyDIDOwner {
        dids[generateDIDString()].hash = hash;
    }

    function updateURI(string did, string uri) public onlyDIDOwner {
        dids[generateDIDString()].uri = uri;
    }

    function deleteDID(string did) public onlyDIDOwner {
        dids[generateDIDString()].exist = false;
    }

    function getHash(string did) public view returns (bytes32) {
        require(dids[did].exist, "did does not exist");
        return dids[did].hash;
    }

    function getURI(string did) public view returns (string) {
        require(dids[did].exist, "did does not exist");
        return dids[did].uri;
    }

    function getHexString(bytes32 value) public pure returns (string) {
        bytes memory result = new bytes(64);
        string memory characterString = "0123456789abcdef";
        bytes memory characters = bytes(characterString);
        for (uint8 i = 0; i < 32; i++) {
            result[i * 2] = characters[uint256((value[i] & 0xF0) >> 4)];
            result[i * 2 + 1] = characters[uint256(value[i] & 0xF)];
        }
        return string(result);
    }

    function generateDIDString() private view returns (string) {
        bytes32 hashedID = keccak256(abi.encodePacked(addrToString(msg.sender)));
        return string(abi.encodePacked(didPrefix, getHexString(hashedID)));
    }

    function addrToString(address _addr) internal pure returns(string) {
        bytes32 value = bytes32(uint256(_addr));
        bytes memory alphabet = "0123456789abcdef";

        bytes memory str = new bytes(43);
        str[0] = ' ';
        str[1] = '0';
        str[2] = 'x';
        for (uint i = 0; i < 20; i++) {
            str[3+i*2] = alphabet[uint(value[i + 12] >> 4)];
            str[4+i*2] = alphabet[uint(value[i + 12] & 0x0f)];
        }
        return string(str);
    }
}
